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
	"github.com/lindb/client_go/api"
)

func writeData() {
	// create write client with options
	cli := lindb.NewClientWithOptions(
		"http://localhost:9000",
		lindb.DefaultOptions().SetBatchSize(200).
			SetReqTimeout(60).
			SetRetryBufferLimit(100).
			SetMaxRetries(3),
	)
	// get write client
	w := cli.Write("_internal")
	// get error chan
	errCh := w.Errors()
	go func() {
		for err := range errCh {
			fmt.Printf("got err:%s\n", err)
		}
	}()

	// write some metric data
	for i := 0; i < 10; i++ {
		// write cpu data
		w.AddPoint(context.TODO(), api.NewPoint("cpu").
			AddTag("host", "host1").
			AddField(api.NewSum("load", 10.0)).
			AddField(api.NewLast("usage", 24.0)))
		// write memory data
		w.AddPoint(context.TODO(), api.NewPoint("memory").
			AddTag("host", "host1").
			AddField(api.NewLast("used", 10.0)).
			AddField(api.NewLast("total", 24.0)))
	}

	// close write client
	w.Close()
}
