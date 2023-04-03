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

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/client_go/api"
)

func TestClient_Write(t *testing.T) {
	c := NewClient("http://localhost:8080")
	assert.NotNil(t, c.Write("test"))

	c = NewClientWithOptions("http://localhost:8080", nil)
	assert.NotNil(t, c.Write("test"))
}

func TestClient_DataQuery(t *testing.T) {
	c := NewClient("http://localhost:8080")
	assert.NotNil(t, c.DataQuery())
}

func TestWrite(t *testing.T) {
	cli := NewClient("http://localhost:9000")
	w := cli.Write("_internal")
	p := api.NewPoint("cpu")
	p.AddTag("host", "host1")
	p.AddField(api.NewSum("mem", 10.0))
	w.AddPoint(context.TODO(), p)

	go func() {
		for err := range w.Errors() {
			fmt.Println("****")
			fmt.Println(err)
			fmt.Println("****")
		}
	}()
	w.Close()
	time.Sleep(time.Second)
}
