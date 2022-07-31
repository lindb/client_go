// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package http

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"runtime"
	"time"
)

var (
	UserAgent = fmt.Sprintf("lindb-client-go/%s  (%s; %s)", "0.0.1", runtime.GOOS, runtime.GOARCH)
)

// Options represents http configuration options for communicating with LinDB server.
type Options struct {
	// Request timeout(s), default 30.
	reqTimeout int64
	// TLS configuration for secure connection, default nil.
	tlsConfig *tls.Config
}

// HTTPClient returns the HTTP client with setting.
func (o *Options) HTTPClient() *http.Client {
	return &http.Client{
		Timeout: time.Second * time.Duration(o.reqTimeout),
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout: 5 * time.Second,
			TLSClientConfig:     o.TLSConfig(),
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     90 * time.Second,
		},
	}
}

// SetReqTimeout sets the request timeout.
func (o *Options) SetReqTimeout(timeout int64) *Options {
	o.reqTimeout = timeout
	return o
}

// ReqTimeout returns the request timeout.
func (o *Options) ReqTimeout() int64 {
	return o.reqTimeout
}

// SetTLSConfig sets TLS configuration for secure connection.
func (o *Options) SetTLSConfig(tlsConfig *tls.Config) *Options {
	o.tlsConfig = tlsConfig
	return o
}

// TLSConfig returns TLS configuration.
func (o *Options) TLSConfig() *tls.Config {
	return o.tlsConfig
}

// DefaultOptions returns an Options with default.
func DefaultOptions() *Options {
	return &Options{
		reqTimeout: 30, // set default request timeout
	}
}
