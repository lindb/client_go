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

	"github.com/lindb/common/models"
	"github.com/lindb/common/pkg/encoding"

	httppkg "github.com/lindb/client_go/internal/http"
)

// For testing
var (
	doPutFn = httppkg.DoPut
)

// DataQuery represents data query client for query time series data from database.
type DataQuery interface {
	// DataQuery queries time series data from database by given ql.
	// LinQL ref: https://lindb.io/guide/lin-ql.html#metric-query
	// Example: select heap_objects from lindb.runtime.mem where 'role' in ('Broker') group by node
	DataQuery(ctx context.Context, database, ql string) (*models.ResultSet, error)
	// MetadataQuery queries metric metadata from database by given ql.
	// LinQL ref: https://lindb.io/guide/lin-ql.html#metric-meta-query
	// Example: show fields from lindb.runtime.mem
	MetadataQuery(ctx context.Context, database, ql string) (*models.Metadata, error)
}

// dataQuery implements DataQuery interface.
type dataQuery struct {
	endpoint string
	client   *http.Client
}

// NewDataQuery creates a data query client.
func NewDataQuery(endpoint string, httpOptions *httppkg.Options) DataQuery {
	return &dataQuery{
		endpoint: fmt.Sprintf("%s/api/v1/exec", endpoint),
		client:   httpOptions.HTTPClient(),
	}
}

// DataQuery queries time series data from database by given ql.
// LinQL ref: https://lindb.io/guide/lin-ql.html#metric-query
// Example: select heap_objects from lindb.runtime.mem where 'role' in ('Broker') group by node
func (q *dataQuery) DataQuery(ctx context.Context, database, ql string) (*models.ResultSet, error) {
	resp, err := q.sendRequest(ctx, database, ql)
	if err != nil {
		return nil, err
	}
	rs := &models.ResultSet{}
	err = encoding.JSONUnmarshal(resp, rs)
	if err != nil {
		return nil, err
	}
	return rs, nil
}

// MetadataQuery queries metric metadata from database by given ql.
// LinQL ref: https://lindb.io/guide/lin-ql.html#metric-meta-query
// Example: show fields from lindb.runtime.mem
func (q *dataQuery) MetadataQuery(ctx context.Context, database, ql string) (*models.Metadata, error) {
	resp, err := q.sendRequest(ctx, database, ql)
	if err != nil {
		return nil, err
	}
	rs := &models.Metadata{}
	err = encoding.JSONUnmarshal(resp, rs)
	if err != nil {
		return nil, err
	}
	return rs, nil
}

// sendRequest sends query request.
func (q *dataQuery) sendRequest(ctx context.Context, database, ql string) ([]byte, error) {
	param := struct {
		Database string `json:"db"`
		QL       string `json:"sql"`
	}{
		Database: database,
		QL:       ql,
	}
	return doPutFn(ctx, q.client, q.endpoint, encoding.JSONMarshal(&param))
}
