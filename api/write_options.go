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

package api

// WriteOptions represents write configuration.
type WriteOptions struct {
	// Number of series sent in single write request, default 1000.
	batchSize int
	// Flush interval(ms) which is buffer flushed if it has not been already written, default 1000.
	flushInterval int64
	// Whether to use GZip compress write data, default true.
	useGZip bool
	// Default tags are added to each written series.
	defaultTags map[string]string
	// Maximum count of retry attempts of failed writes, default 3.
	maxRetries int
	// Maximum number of write request to keep for retry, default 100.
	retryBufferLimit int
}

// SetBatchSize sets batch size in single write request.
func (opt *WriteOptions) SetBatchSize(batchSize int) *WriteOptions {
	opt.batchSize = batchSize
	return opt
}

// BatchSize returns the number of batch size in single write request.
func (opt *WriteOptions) BatchSize() int {
	return opt.batchSize
}

// SetFlushInterval sets flush interval(ms).
func (opt *WriteOptions) SetFlushInterval(interval int64) *WriteOptions {
	opt.flushInterval = interval
	return opt
}

// FlushInterval returns the flush interval(ms).
func (opt *WriteOptions) FlushInterval() int64 {
	return opt.flushInterval
}

// SetUseGZip sets whether to use GZip compress write data.
func (opt *WriteOptions) SetUseGZip(useGZip bool) *WriteOptions {
	opt.useGZip = useGZip
	return opt
}

// UseGZip returns whether to use GZip compress write data.
func (opt *WriteOptions) UseGZip() bool {
	return opt.useGZip
}

// AddDefaultTag adds default tag.
func (opt *WriteOptions) AddDefaultTag(key, value string) *WriteOptions {
	if opt.defaultTags == nil {
		opt.defaultTags = make(map[string]string)
	}
	opt.defaultTags[key] = value
	return opt
}

// DefaultTags returns the default tags for all metrics.
func (opt *WriteOptions) DefaultTags() map[string]string {
	return opt.defaultTags
}

// SetMaxRetries sets maximum count of retry attempts of failed write.
func (opt *WriteOptions) SetMaxRetries(maxRetries int) *WriteOptions {
	opt.maxRetries = maxRetries
	return opt
}

// MaxRetries returns maximum count of retry attempts of failed write.
func (opt *WriteOptions) MaxRetries() int {
	return opt.maxRetries
}

// SetRetryBufferLimit sets maximum number of write request to keep for retry.
func (opt *WriteOptions) SetRetryBufferLimit(retryBufferLimit int) *WriteOptions {
	opt.retryBufferLimit = retryBufferLimit
	return opt
}

// RetryBufferLimit returns maximum number of write request to keep for retry.
func (opt *WriteOptions) RetryBufferLimit() int {
	return opt.retryBufferLimit
}

// DefaultWriteOptions creates a WriteOptions with default.
func DefaultWriteOptions() *WriteOptions {
	return &WriteOptions{
		batchSize:        1_000,
		flushInterval:    1_000, // 1s
		useGZip:          true,
		maxRetries:       3,
		retryBufferLimit: 1_00,
	}
}
