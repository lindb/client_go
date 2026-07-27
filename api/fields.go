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
	"github.com/lindb/common/metric"
	"github.com/lindb/common/proto/gen/v1/flatMetricsV1"

	"github.com/lindb/client_go/internal"
)

// Field represents field interface.
type Field interface {
	// write field data into broker row builder.
	write(builder *metric.RowBuilder) error
}

// exemplar represents an exemplar, which is a sample input measurement.
// Exemplars also hold information about the environment when the measurement
// was recorded, for example the span and trace ID of the active span when the
// exemplar was recorded.
type exemplar struct {
	name     string
	traceID  string
	spanID   string
	duration int64
}

// NewExemplar creates an exemplar.
func NewExemplar(name, traceID, spanID string, duration int64) Field {
	return &exemplar{
		name:     name,
		traceID:  traceID,
		spanID:   spanID,
		duration: duration,
	}
}

// write exemplar data into broker row builder(flat protocol)
func (e *exemplar) write(builder *metric.RowBuilder) error {
	return builder.AddExemplar(internal.String2ByteSlice(e.name),
		internal.String2ByteSlice(e.traceID),
		internal.String2ByteSlice(e.spanID),
		e.duration,
	)
}

// simpleField represents simple field like(sum/min/max/last etc.)
type simpleField struct {
	name      string                        // field name
	v         float64                       // field value
	fieldType flatMetricsV1.SimpleFieldType // field type
}

// write field data into broker row builder(flat protocol)
func (s *simpleField) write(builder *metric.RowBuilder) error {
	return builder.AddSimpleField(internal.String2ByteSlice(s.name), s.fieldType, s.v)
}

// Sum represents sum field, implements Field interface.
type Sum struct {
	simpleField
}

// NewSum creates a Sum filed.
func NewSum(name string, v float64) Field {
	return &Sum{
		simpleField: simpleField{
			fieldType: flatMetricsV1.SimpleFieldTypeDeltaSum,
			name:      name,
			v:         v,
		},
	}
}

// Min represents min field, implements Field interface.
type Min struct {
	simpleField
}

// NewMin creates a Min field.
func NewMin(name string, v float64) Field {
	return &Min{
		simpleField: simpleField{
			fieldType: flatMetricsV1.SimpleFieldTypeMin,
			name:      name,
			v:         v,
		},
	}
}

// Max represents max field, implements Field interface.
type Max struct {
	simpleField
}

// NewMax creates a Max field.
func NewMax(name string, v float64) Field {
	return &Max{
		simpleField: simpleField{
			fieldType: flatMetricsV1.SimpleFieldTypeMax,
			name:      name,
			v:         v,
		},
	}
}

// First represents first field, implements Field interface.
type First struct {
	simpleField
}

// NewFirst creates a First field.
func NewFirst(name string, v float64) Field {
	return &First{
		simpleField: simpleField{
			fieldType: flatMetricsV1.SimpleFieldTypeFirst,
			name:      name,
			v:         v,
		},
	}
}

// Last represents last field, implements Field interface.
type Last struct {
	simpleField
}

// NewLast creates a Last field.
func NewLast(name string, v float64) Field {
	return &Last{
		simpleField: simpleField{
			fieldType: flatMetricsV1.SimpleFieldTypeLast,
			name:      name,
			v:         v,
		},
	}
}

// Histogram represents histogram field(compound field).
type Histogram struct {
	min, max, sum, count float64
	values, bounds       []float64
}

// NewHistogram creates a Histogram field.
func NewHistogram(minVal, maxVal, sum, count float64, values, bounds []float64) Field {
	return &Histogram{
		min:    minVal,
		max:    maxVal,
		sum:    sum,
		count:  count,
		values: values,
		bounds: bounds,
	}
}

// write histogram data into broker row builder.
func (h *Histogram) write(builder *metric.RowBuilder) error {
	if err := builder.AddCompoundFieldMMSC(h.min, h.max, h.sum, h.count); err != nil {
		return err
	}
	return builder.AddCompoundFieldData(h.values, h.bounds)
}
