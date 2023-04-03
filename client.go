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

import "github.com/lindb/client_go/api"

// Client represents the api to communicate with LinDB backend server.
// Ref InfluxDB client: https://github.com/influxdata/influxdb-client-go
type Client interface {
	// Write returns an asynchronous write client.
	Write(database string) api.Write
	// DataQuery returns a metric data query client.
	DataQuery() api.DataQuery
}

// client implements the Client interface.
type client struct {
	brokerEndpoint string
	options        *Options
}

// NewClientWithOptions creates a Client with backend endpoint and options.
func NewClientWithOptions(brokerEndpoint string, options *Options) Client {
	if options == nil {
		options = DefaultOptions()
	}
	return &client{
		brokerEndpoint: brokerEndpoint,
		options:        options,
	}
}

// NewClient creates a Client with backend endpoint and default options.
func NewClient(brokerEndpoint string) Client {
	return &client{
		brokerEndpoint: brokerEndpoint,
		options:        DefaultOptions(),
	}
}

// Write returns an asynchronous write client.
func (c *client) Write(database string) api.Write {
	return api.NewWrite(c.brokerEndpoint, database, c.options.WriteOptions(), c.options.HTTPOptions())
}

// DataQuery returns a metric data query client.
func (c *client) DataQuery() api.DataQuery {
	return api.NewDataQuery(c.brokerEndpoint, c.options.HTTPOptions())
}
