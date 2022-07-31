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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/client_go/api"
	"github.com/lindb/client_go/internal/http"
)

func TestOptions(t *testing.T) {
	assert.Equal(t, http.DefaultOptions(), (&Options{}).HTTPOptions())
	assert.Equal(t, api.DefaultWriteOptions(), (&Options{}).WriteOptions())

	opt := DefaultOptions()
	assert.Equal(t, http.DefaultOptions(), opt.HTTPOptions())
	assert.Equal(t, api.DefaultWriteOptions(), opt.WriteOptions())

	opt.AddDefaultTag("k1", "v1").SetUseGZip(false).SetBatchSize(2_000).
		SetMaxRetries(10).SetRetryBufferLimit(3_000).
		SetFlushInterval(1_000).SetReqTimeout(60).SetTLSConfig(&tls.Config{})
	assert.False(t, opt.WriteOptions().UseGZip())
	assert.Equal(t, 2_000, opt.WriteOptions().BatchSize())
	assert.Equal(t, int64(1_000), opt.WriteOptions().FlushInterval())
	assert.Equal(t, int64(60), opt.HTTPOptions().ReqTimeout())
	assert.Equal(t, 10, opt.WriteOptions().MaxRetries())
	assert.Equal(t, 3_000, opt.WriteOptions().RetryBufferLimit())
	assert.NotNil(t, opt.HTTPOptions().TLSConfig())
}
