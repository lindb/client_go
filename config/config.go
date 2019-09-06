package config

import (
	"fmt"

	"github.com/lindb/lindb/pkg/fileutil"
)

// default config
const (
	// http
	syncAddressInterval int64 = 60
	syncAddressTimeout  int64 = 5

	// tcp
	dialTimeout int64 = 1
	retryLimit  int   = 2

	// buffer
	bufferSize    int = 1024
	databaseLimit int = 30
)

// ClientConfig defines config for client.
type ClientConfig struct {
	// http
	BrokerURL           string `toml:"brokerURL"`
	SyncAddressInterval int64  `toml:"syncAddressInterval"`
	SyncAddressTimeout  int64  `toml:"syncAddressTimeout"`

	//tcp
	DialTimeout int64 `toml:"dialTimeout"`
	RetryLimit  int   `toml:"retryLimit"`

	// buffer
	BufferSize    int `toml:"bufferSize"`
	DatabaseLimit int `toml:"databaseLimit"`
}

// LoadConfig loads Client from file path, return error when file not exists or config parse error.
func LoadConfig(path string) (*ClientConfig, error) {
	if !fileutil.Exist(path) {
		return nil, fmt.Errorf("config path:%s not exists", path)
	}
	clientConfig := NewDefaultConfig()
	err := fileutil.DecodeToml(path, clientConfig)
	if err != nil {
		return nil, fmt.Errorf("load config from path:%s err:%v", path, err)
	}

	if clientConfig.BrokerURL == "" {
		return nil, fmt.Errorf("broker url is not provided")
	}

	return clientConfig, nil
}

func NewDefaultConfig() *ClientConfig {
	return &ClientConfig{
		SyncAddressInterval: syncAddressInterval,
		SyncAddressTimeout:  syncAddressTimeout,
		DialTimeout:         dialTimeout,
		RetryLimit:          retryLimit,
		BufferSize:          bufferSize,
		DatabaseLimit:       databaseLimit,
	}
}
