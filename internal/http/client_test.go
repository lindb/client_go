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

package http

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func TestClient_DoPut(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/error":
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("500 - Something bad happened!"))
		case "/gzip":
			gz := gzip.NewWriter(w)
			defer gz.Close()
			gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
			gzr.Header().Set("Content-Encoding", "gzip")
			gzr.WriteHeader(http.StatusOK)
			_, _ = gzr.Write([]byte("gzip"))
		default:
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("Good!"))
		}
	}))
	defer ts.Close()

	cases := []struct {
		name    string
		prepare func()
		assert  func(resp []byte, err error)
	}{
		{
			name: "new request failure",
			prepare: func() {
				newRequestFn = func(_ context.Context, _, _ string, _ io.Reader) (*http.Request, error) {
					return nil, fmt.Errorf("err")
				}
			},
			assert: func(resp []byte, err error) {
				assert.Nil(t, resp)
				assert.Error(t, err)
			},
		},
		{
			name: "do request failure",
			prepare: func() {
				newRequestFn = func(ctx context.Context, method, _ string, body io.Reader) (*http.Request, error) {
					return http.NewRequestWithContext(context.TODO(), method, "", body)
				}
			},
			assert: func(resp []byte, err error) {
				assert.Nil(t, resp)
				assert.Error(t, err)
			},
		},
		{
			name: "get internal error",
			prepare: func() {
				newRequestFn = func(ctx context.Context, method, endpoint string, body io.Reader) (*http.Request, error) {
					return http.NewRequestWithContext(context.TODO(), method, endpoint+"/error", body)
				}
			},
			assert: func(resp []byte, err error) {
				assert.Nil(t, resp)
				assert.Error(t, err)
			},
		},
		{
			name: "get internal error, read error msg failure",
			prepare: func() {
				newRequestFn = func(ctx context.Context, method, endpoint string, body io.Reader) (*http.Request, error) {
					return http.NewRequestWithContext(context.TODO(), method, endpoint+"/error", body)
				}
				readAllFn = func(_ io.Reader) ([]byte, error) {
					return nil, fmt.Errorf("err")
				}
			},
			assert: func(resp []byte, err error) {
				assert.Nil(t, resp)
				assert.Error(t, err)
			},
		},
		{
			name: "read resp failure",
			prepare: func() {
				readAllFn = func(_ io.Reader) ([]byte, error) {
					return nil, fmt.Errorf("err")
				}
			},
			assert: func(resp []byte, err error) {
				assert.Nil(t, resp)
				assert.Error(t, err)
			},
		},
		{
			name:    "read resp successfully",
			prepare: func() {},
			assert: func(resp []byte, err error) {
				assert.Equal(t, []byte("Good!"), resp)
				assert.NoError(t, err)
			},
		},
		{
			name: "read gzip resp failure",
			prepare: func() {
				newRequestFn = func(ctx context.Context, method, endpoint string, body io.Reader) (*http.Request, error) {
					return http.NewRequestWithContext(context.TODO(), method, endpoint+"/gzip", body)
				}
				newGzipReaderFn = func(_ io.Reader) (*gzip.Reader, error) {
					return nil, fmt.Errorf("err")
				}
			},
			assert: func(resp []byte, err error) {
				assert.Nil(t, resp)
				assert.Error(t, err)
			},
		},
		{
			name: "read gzip resp successfully",
			prepare: func() {
				newRequestFn = func(ctx context.Context, method, endpoint string, body io.Reader) (*http.Request, error) {
					return http.NewRequestWithContext(context.TODO(), method, endpoint+"/gzip", body)
				}
			},
			assert: func(resp []byte, err error) {
				assert.Equal(t, []byte("gzip"), resp)
				assert.NoError(t, err)
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(_ *testing.T) {
			defer func() {
				newRequestFn = http.NewRequestWithContext
				readAllFn = io.ReadAll
				newGzipReaderFn = gzip.NewReader
			}()
			tt.prepare()
			resp, err := DoPut(context.TODO(), &http.Client{}, ts.URL, []byte{1})
			tt.assert(resp, err)
		})
	}
}
