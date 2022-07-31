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

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteOptions(t *testing.T) {
	assert.Equal(t, 1_000, DefaultWriteOptions().BatchSize())
	assert.Equal(t, int64(1_000), DefaultWriteOptions().FlushInterval())
	assert.Equal(t, 3, DefaultWriteOptions().MaxRetries())
	assert.Equal(t, 100, DefaultWriteOptions().RetryBufferLimit())
	assert.True(t, DefaultWriteOptions().UseGZip())
	assert.Nil(t, DefaultWriteOptions().DefaultTags())

	opt := DefaultWriteOptions().SetUseGZip(false).
		SetFlushInterval(3_000).
		SetBatchSize(2_000).
		SetMaxRetries(10).
		SetRetryBufferLimit(1_000).
		AddDefaultTag("k1", "v1").
		AddDefaultTag("k2", "v2")
	assert.Equal(t, 2_000, opt.BatchSize())
	assert.Equal(t, int64(3_000), opt.FlushInterval())
	assert.Equal(t, 10, opt.MaxRetries())
	assert.Equal(t, 1_000, opt.RetryBufferLimit())
	assert.False(t, opt.UseGZip())
	assert.Equal(t, map[string]string{"k1": "v1", "k2": "v2"}, opt.DefaultTags())
}
