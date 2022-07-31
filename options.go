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

package lindb

import (
	"crypto/tls"

	"github.com/lindb/client_go/api"
	"github.com/lindb/client_go/internal/http"
)

// Options represents configuration options for client.
type Options struct {
	httpOptions  *http.Options     // HTTP options
	writeOptions *api.WriteOptions // Write options
}

// SetTLSConfig sets TLS configuration for secure connection.
func (o *Options) SetTLSConfig(tlsConfig *tls.Config) *Options {
	o.HTTPOptions().SetTLSConfig(tlsConfig)
	return o
}

// SetReqTimeout sets HTTP request timeout(sec)
func (o *Options) SetReqTimeout(timeout int64) *Options {
	o.HTTPOptions().SetReqTimeout(timeout)
	return o
}

// HTTPOptions returns the HTTP options, if not set return default options.
func (o *Options) HTTPOptions() *http.Options {
	if o.httpOptions == nil {
		o.httpOptions = http.DefaultOptions()
	}
	return o.httpOptions
}

// SetBatchSize sets batch size in single write request.
func (o *Options) SetBatchSize(batchSize int) *Options {
	o.WriteOptions().SetBatchSize(batchSize)
	return o
}

// SetFlushInterval sets flush interval(ms)
func (o *Options) SetFlushInterval(interval int64) *Options {
	o.WriteOptions().SetFlushInterval(interval)
	return o
}

// SetUseGZip sets whether to use GZip compress write data.
func (o *Options) SetUseGZip(useGZip bool) *Options {
	o.WriteOptions().SetUseGZip(useGZip)
	return o
}

// AddDefaultTag adds default tag for all metrics.
func (o *Options) AddDefaultTag(key, value string) *Options {
	o.WriteOptions().AddDefaultTag(key, value)
	return o
}

// SetMaxRetries sets maximum count of retry attempts of failed write.
func (o *Options) SetMaxRetries(maxRetries int) *Options {
	o.WriteOptions().SetMaxRetries(maxRetries)
	return o
}

// SetRetryBufferLimit sets maximum number of write request to keep for retry.
func (o *Options) SetRetryBufferLimit(retryBufferLimit int) *Options {
	o.WriteOptions().SetRetryBufferLimit(retryBufferLimit)
	return o
}

// WriteOptions returns the write options, if not set return default options.
func (o *Options) WriteOptions() *api.WriteOptions {
	if o.writeOptions == nil {
		o.writeOptions = api.DefaultWriteOptions()
	}
	return o.writeOptions
}

// DefaultOptions creates an Options with default.
func DefaultOptions() *Options {
	return &Options{
		httpOptions:  http.DefaultOptions(),
		writeOptions: api.DefaultWriteOptions(),
	}
}
