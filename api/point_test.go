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
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPoint(t *testing.T) {
	now := time.Now()
	sum := NewSum("s", 10.0)
	p := NewPoint("test").AddTag("k1", "v1").
		SetNamespace("ns").SetTimestamp(now).AddField(sum)
	assert.Equal(t, "ns", p.Namespace())
	assert.Equal(t, "test", p.MetricName())
	assert.Equal(t, map[string]string{"k1": "v1"}, p.Tags())
	assert.Len(t, p.Fields(), 1)
	assert.Equal(t, sum, p.Fields()[0])
	assert.Equal(t, now, p.Timestamp())
	assert.True(t, p.Valid())

	p = NewPoint("")
	assert.False(t, NewPoint("").Valid())
	assert.False(t, NewPoint("xx").Valid())
}
