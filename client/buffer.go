package client

import (
	"errors"

	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/rpc/proto/field"
)

var errBufferFull = errors.New("buffer is full, consider increase buffer size")

type Buffer interface {
	In(metric *field.Metric) error
	Out() []*field.Metric
}

// buffer batches the incoming metrics.
type buffer struct {
	database string
	ch       chan *field.Metric
	logger   *logger.Logger
}

// In puts the metric into the batch buffer, metric will be dropped if the buffer is full.
func (b *buffer) In(metric *field.Metric) error {
	select {
	case b.ch <- metric:
	default:
		// drop is buffer is full
		b.logger.Error("buffer full, drop metric", logger.String("metricName", metric.Name), logger.String("database", b.database))
		return errBufferFull
	}
	return nil
}

// Out gets a batch of metrics from the buffer.
func (b *buffer) Out() []*field.Metric {
	// todo determine size
	batch := make([]*field.Metric, 0, len(b.ch))
	for {
		select {
		case metric := <-b.ch:
			batch = append(batch, metric)
		default:
			return batch
		}
	}
}

// NewBuffer creates a buffer with provided size.
func NewBuffer(database string, size int) Buffer {
	return &buffer{
		database: database,
		ch:       make(chan *field.Metric, size),
		logger:   logger.GetLogger("client", "buffer"),
	}
}
