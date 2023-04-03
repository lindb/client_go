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
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/lindb/common/models"
	"github.com/stretchr/testify/assert"

	httppkg "github.com/lindb/client_go/internal/http"
)

func TestDataQuery_Data(t *testing.T) {
	cases := []struct {
		name    string
		prepare func()
		assert  func(rs *models.ResultSet, err error)
	}{
		{
			name: "send request failure",
			prepare: func() {
				doPutFn = func(_ context.Context, _ *http.Client, _ string, _ []byte) ([]byte, error) {
					return nil, fmt.Errorf("err")
				}
			},
			assert: func(rs *models.ResultSet, err error) {
				assert.Nil(t, rs)
				assert.Error(t, err)
			},
		},
		{
			name: "unmarshal failure",
			prepare: func() {
				doPutFn = func(_ context.Context, _ *http.Client, _ string, _ []byte) ([]byte, error) {
					return []byte("abc"), nil
				}
			},
			assert: func(rs *models.ResultSet, err error) {
				assert.Nil(t, rs)
				assert.Error(t, err)
			},
		},
		{
			name: "successfully",
			prepare: func() {
				doPutFn = func(_ context.Context, _ *http.Client, _ string, _ []byte) ([]byte, error) {
					return []byte("{}"), nil
				}
			},
			assert: func(rs *models.ResultSet, err error) {
				assert.NotNil(t, rs)
				assert.NoError(t, err)
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(_ *testing.T) {
			api := NewDataQuery("test", httppkg.DefaultOptions())
			tt.prepare()
			rs, err := api.DataQuery(context.TODO(), "test", "select load from cpu")
			tt.assert(rs, err)
		})
	}
}

func TestDataQuery_MetaData(t *testing.T) {
	cases := []struct {
		name    string
		prepare func()
		assert  func(rs *models.Metadata, err error)
	}{
		{
			name: "send request failure",
			prepare: func() {
				doPutFn = func(_ context.Context, _ *http.Client, _ string, _ []byte) ([]byte, error) {
					return nil, fmt.Errorf("err")
				}
			},
			assert: func(rs *models.Metadata, err error) {
				assert.Nil(t, rs)
				assert.Error(t, err)
			},
		},
		{
			name: "unmarshal failure",
			prepare: func() {
				doPutFn = func(_ context.Context, _ *http.Client, _ string, _ []byte) ([]byte, error) {
					return []byte("abc"), nil
				}
			},
			assert: func(rs *models.Metadata, err error) {
				assert.Nil(t, rs)
				assert.Error(t, err)
			},
		},
		{
			name: "successfully",
			prepare: func() {
				doPutFn = func(_ context.Context, _ *http.Client, _ string, _ []byte) ([]byte, error) {
					return []byte("{}"), nil
				}
			},
			assert: func(rs *models.Metadata, err error) {
				assert.NotNil(t, rs)
				assert.NoError(t, err)
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(_ *testing.T) {
			api := NewDataQuery("test", httppkg.DefaultOptions())
			tt.prepare()
			rs, err := api.MetadataQuery(context.TODO(), "test", "show fields from cpu")
			tt.assert(rs, err)
		})
	}
}
