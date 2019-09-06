package client

import (
	"net"
	"time"

	"github.com/lindb/client_go/config"
	"github.com/lindb/lindb/pkg/logger"
)

// Sender defines a component to send binary data to remote receiver.
type Sender interface {
	// Send sends byte array data to remote receiver.
	// Retry limited times when encountering error in the process.
	// Concurrent unsafe.
	Send(data []byte)
}

// sender implements Sender.
type sender struct {
	// manager providers remote addresses
	manager config.AddressManager
	// conn is the current net connection
	conn net.Conn
	// retry times limit
	retryLimit int
	// dial timeout
	dialTimeout int64
	logger      *logger.Logger
}

// Send sends byte array data to remote receiver.
// Retry limited times when encountering error in the process.
// Concurrent unsafe.
func (s *sender) Send(data []byte) {
	for i := 0; i < s.retryLimit; i++ {
		if s.conn == nil {
			s.connect()
			continue
		}
		n, err := s.conn.Write(data)
		if err != nil {
			s.logger.Error("write", logger.Error(err), logger.Int32("writeSize", int32(n)))
			s.closeConn()
			s.connect()
		}

		// success send data, break loop
		break
	}
}

// closeConn close current connection.
func (s *sender) closeConn() {
	if s.conn != nil {
		if err := s.conn.Close(); err != nil {
			s.logger.Error("close conn", logger.Error(err), logger.String("address", s.conn.RemoteAddr().String()))
		}
		s.logger.Info("close conn", logger.String("address", s.conn.RemoteAddr().String()))
		s.conn = nil
	}
}

// connect try to create a new connection.
func (s *sender) connect() {
	node, err := s.manager.RandomNext()
	if err != nil {
		s.logger.Error("get next node err", logger.Error(err))
		return
	}
	dialer := net.Dialer{Timeout: time.Duration(s.dialTimeout) * time.Second}
	conn, err := dialer.Dial("tcp", config.BuildAddress(node))
	if err != nil {
		s.logger.Error("dial err", logger.Error(err), logger.String("address", node.Indicator()))
		return
	}
	s.logger.Info("connect success", logger.String("address", node.Indicator()))
	s.conn = conn
}

// NewSender creates a new Sender with clientConfig, address manager.
func NewSender(clientConfig *config.ClientConfig, manager config.AddressManager) Sender {
	s := &sender{
		manager:     manager,
		dialTimeout: clientConfig.DialTimeout,
		retryLimit:  clientConfig.RetryLimit,
		logger:      logger.GetLogger("client", "Sender"),
	}
	s.connect()
	return s
}
