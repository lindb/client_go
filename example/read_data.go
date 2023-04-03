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

package main

import (
	"context"
	"fmt"

	lindb "github.com/lindb/client_go"
)

var (
	cli   = lindb.NewClient("http://localhost:9000")
	query = cli.DataQuery()
)

func ReadMetricData() {
	// LinQL ref: https://lindb.io/guide/lin-ql.html#metric-query
	data, err := query.DataQuery(context.TODO(),
		"_internal",
		"select heap_objects from lindb.runtime.mem where time>now()-2m and 'role' in ('Broker') group by node")
	if err != nil {
		fmt.Println(err)
		return
	}
	// print table
	_, table := data.ToTable()
	fmt.Println(table)
}

func ReadMetricMetadata() {
	// LinQL ref: https://lindb.io/guide/lin-ql.html#metric-meta-query
	qls := []string{
		"show namespaces",
		"show metrics",
		"show tag keys from lindb.runtime.mem",
		"show tag values from lindb.runtime.mem with key=namespace",
		"show fields from lindb.runtime.mem",
	}
	for _, ql := range qls {
		data, err := query.MetadataQuery(context.TODO(), "_internal", ql)
		if err != nil {
			fmt.Println(err)
			return
		}
		// print table
		_, table := data.ToTable()
		fmt.Println(table)
	}
}
