package client

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/lindb/client_go/config"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/rpc/proto/field"
)

var (
	// errClientClosed is the error returned when writing to a closed client.
	errClientClosed = errors.New("client is closed")

	// errTooManyDatabases is the error returned when recorded database num exceeds config limit.
	errTooManyDatabases = errors.New("database num exceeds the config limit")
)

// Client defines the basic write, close methods for a lindb client.
type Client interface {
	// WriteList writes a metricList to client buffer, not guarantees to write metrics to lindb storage.
	WriteList(metricList *field.MetricList) error
	// Write writes a metric with the specific database to storage.
	Write(database string, metric *field.Metric) error
	// Close closes the client, all writes after close will result in errClientClosed.
	Close() error
}

// client implements Client.
type client struct {
	// close status mark
	closed uint32
	// config for client
	clientConfig *config.ClientConfig
	// sender handles all sending things
	sender Sender
	// database -> metric buffer map
	bufferMap sync.Map
	// lock for bufferMap
	locker4map sync.Mutex
	// num for database
	databaseNum int
	logger      *logger.Logger
}

// Write writes a metric with the specific database to storage.
func (c *client) Write(database string, metric *field.Metric) error {
	return c.WriteList(&field.MetricList{
		Database: database,
		Metrics:  []*field.Metric{metric},
	})
}

// Close closes the client, all writes after close will result in errClientClosed.
func (c *client) Close() error {
	atomic.StoreUint32(&c.closed, 1)
	c.flush()
	c.logger.Info("client closed")
	return nil
}

// WriteList writes a metricList to storage.
func (c *client) WriteList(metricList *field.MetricList) error {
	if atomic.LoadUint32(&c.closed) == 1 {
		return errClientClosed
	}

	if metricList == nil || len(metricList.Metrics) == 0 {
		return nil
	}

	database := metricList.Database

	val, ok := c.bufferMap.Load(database)
	if !ok {
		// double check
		c.locker4map.Lock()
		val, ok = c.bufferMap.Load(database)
		if !ok {
			if c.databaseNum > c.clientConfig.DatabaseLimit {
				return errTooManyDatabases
			}
			val = NewBuffer(database, c.clientConfig.BufferSize)
			c.bufferMap.Store(database, val)
		}
		c.locker4map.Unlock()
	}

	buffer, _ := val.(*buffer)

	var (
		err error
	)
	for _, metric := range metricList.Metrics {
		err = c.write(buffer, metric)
		if err != nil {
			return err
		}

	}
	return nil
}

// write validates the metric, all invalidate metric will be dropped silently.
// then put the metric to buffer for batching.
func (c *client) write(buffer *buffer, metric *field.Metric) error {
	if metric == nil {
		return nil
	}

	if metric.Name == "" {
		return nil
	}

	if len(metric.Fields) == 0 {
		return nil
	}

	if metric.Timestamp == 0 {
		metric.Timestamp = timeutil.Now()
	}

	return buffer.In(metric)
}

// initFlushTask inits a flush task to flush data and send to sender.
func (c *client) initFlushTask() {
	go func() {
		for {
			if atomic.LoadUint32(&c.closed) == 1 {
				break
			}
			count := c.flush()
			if count == 0 {
				//. todo config, time unit
				time.Sleep(20 * time.Millisecond)
			}
		}
		c.logger.Info("stop flush task")
	}()

	c.logger.Info("start flush task")
}

// flush traverse all the buffers and sends batched metrics.
func (c *client) flush() int {
	count := 0
	c.bufferMap.Range(func(key, value interface{}) bool {
		buffer, _ := value.(*buffer)
		metrics := buffer.Out()
		size := len(metrics)
		if size == 0 {
			return true
		}

		count += size

		database, _ := key.(string)

		metricList := &field.MetricList{
			Database: database,
			Metrics:  metrics,
		}

		metricListSizeInBytes := metricList.Size()
		bytes := make([]byte, 4+metricListSizeInBytes)

		writer := stream.NewSliceWriter(bytes[:4])
		writer.PutInt32(int32(metricListSizeInBytes))

		_, err := metricList.MarshalTo(bytes[4:])
		if err != nil {
			c.logger.Error("metricList Marshal", logger.Error(err), logger.String("database", database))
			return true
		}

		c.sender.Send(bytes)
		return true
	})

	c.logger.Info("flush", logger.Int32("count", int32(count)))

	return count
}

// NewClient creates a client from config file located by configPath, broker url must exist in config file.
func NewClientFromConfigPath(configPath string) (Client, error) {
	c, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	return newClient(c), nil
}

// NewClient creates a client with broker url.
func NewClient(brokerURL string) (Client, error) {
	if brokerURL == "" {
		return nil, fmt.Errorf("broker url is not provided")
	}
	c := config.NewDefaultConfig()
	c.BrokerURL = brokerURL
	return newClient(c), nil
}

// newClient returns a client with provided ClientConfig.
func newClient(c *config.ClientConfig) Client {
	m := config.NewAddressManager(c)
	s := NewSender(c, m)

	cli := &client{
		clientConfig: c,
		sender:       s,
		logger:       logger.GetLogger("client", "Client"),
	}

	cli.initFlushTask()
	return cli
}
