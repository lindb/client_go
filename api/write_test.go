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
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	httppkg "github.com/lindb/client_go/internal/http"
)

func TestNewWrite(t *testing.T) {
	w := NewWrite("http://localhost:9000", "test",
		DefaultWriteOptions(), httppkg.DefaultOptions())
	assert.NotNil(t, w)
	errCh := w.Errors()
	assert.NotNil(t, errCh)
	w.Close()
	w.Close() // ignore it
}

func TestWriteData(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`ok`))
	}))
	defer svr.Close()

	w := NewWrite(svr.URL, "test",
		DefaultWriteOptions().AddDefaultTag("key", "value"), httppkg.DefaultOptions())
	for i := 0; i < 10; i++ {
		w.AddPoint(context.TODO(), NewPoint("cpu").
			AddTag("key1", "value1").AddField(NewLast("load", 10.0)))
	}
	w.Close()
}

func TestWriteData_Failure(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`error`))
	}))
	defer svr.Close()

	w := NewWrite(svr.URL, "test",
		DefaultWriteOptions().SetMaxRetries(2).SetBatchSize(1).
			SetRetryBufferLimit(50), httppkg.DefaultOptions())
	for i := 0; i < 100; i++ {
		w.AddPoint(context.TODO(), NewPoint("cpu").
			AddTag("key1", "value1").AddField(NewLast("load", 10.0)))
	}
	w.Close()
}

func TestAddPoint(t *testing.T) {
	t.Run("invalid point", func(t *testing.T) {
		w := write{}
		w.AddPoint(context.TODO(), NewPoint("cpu"))
	})
	t.Run("add point timeout", func(t *testing.T) {
		w := write{bufferCh: make(chan *Point)}
		ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*10)
		defer cancel()
		w.AddPoint(ctx, NewPoint("cpu").AddField(NewLast("load", 10.0)))
	})
}

func TestAddWrongPoint(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`ok`))
	}))
	defer svr.Close()

	t.Run("wrong common tags", func(t *testing.T) {
		w := NewWrite(svr.URL, "test",
			DefaultWriteOptions().AddDefaultTag("key", ""), httppkg.DefaultOptions())
		w.AddPoint(context.TODO(), NewPoint("cpu").AddField(NewLast("load", 10.0)))
		w.Close()
	})
	t.Run("wrong field data", func(t *testing.T) {
		w := NewWrite(svr.URL, "test",
			DefaultWriteOptions(), httppkg.DefaultOptions())
		w.AddPoint(context.TODO(), NewPoint("cpu").AddField(NewSum("load", math.Inf(0))))
		w.Close()
	})
	t.Run("wrong point tags", func(t *testing.T) {
		w := NewWrite(svr.URL, "test",
			DefaultWriteOptions(), httppkg.DefaultOptions())
		w.AddPoint(context.TODO(), NewPoint("cpu").AddTag("key", "").AddField(NewLast("load", 10.0)))
		errCh := w.Errors()
		go func() {
			for err := range errCh {
				fmt.Println(err)
			}
		}()
		w.Close()
	})
}
