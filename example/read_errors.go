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

func readErrors() {
	// create write client
	cli := lindb.NewClient("http://localhost:9000")
	w := cli.Write("_internal")
	// get error chan
	errCh := w.Errors()
	go func() {
		for err := range errCh {
			fmt.Printf("got err:%s\n", err)
		}
	}()

	// write data
	for i := 0; i < 10; i++ {
		p := api.NewPoint("cpu").
			AddTag("host", "host1").
			AddField(api.NewSum("load", 10.0)).
			AddField(api.NewLast("usage", 24.0))
		w.AddPoint(context.TODO(), p)
	}

	// close write client
	w.Close()
}
