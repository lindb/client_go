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
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/common/series"
)

func TestSimpleField(t *testing.T) {
	builder, releaseFunc := series.NewRowBuilder()
	defer releaseFunc(builder)

	assert.NoError(t, NewSum("sum", 10.0).write(builder))
	assert.NoError(t, NewMax("max", 10.0).write(builder))
	assert.NoError(t, NewMin("min", 10.0).write(builder))
	assert.NoError(t, NewFirst("first", 10.0).write(builder))
	assert.NoError(t, NewLast("last", 10.0).write(builder))
}

func TestHistogramField(t *testing.T) {
	builder, releaseFunc := series.NewRowBuilder()
	defer releaseFunc(builder)

	histogram := NewHistogram(1.0, 10.0, 100.0, 20.0, []float64{1, 2, 3}, []float64{1, 2, math.Inf(0)})
	assert.NoError(t, histogram.write(builder))

	histogram = NewHistogram(1.0, 10.0, 100.0, 20.0, []float64{1, 2, 3, 4}, []float64{1, 2, 3})
	assert.Error(t, histogram.write(builder))

	histogram = NewHistogram(-1.0, 10.0, 100.0, 20.0, []float64{1, 2, 3}, []float64{1, 2, math.Inf(0)})
	assert.Error(t, histogram.write(builder))
}
