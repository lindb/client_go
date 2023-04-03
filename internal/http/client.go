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
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"io"
	"net/http"
)

// for testing
var (
	newRequestFn    = http.NewRequestWithContext
	readAllFn       = io.ReadAll
	newGzipReaderFn = gzip.NewReader
)

// DoPut sends put request based on given client/endpoint/request body.
func DoPut(ctx context.Context, cli *http.Client, endpoint string, body []byte) ([]byte, error) {
	req, err := newRequestFn(ctx, http.MethodPut, endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Content-Type", "application/json")

	resp, err := cli.Do(req)
	defer func() {
		// need close resp body by defer, maybe resp is not nil when throw some err
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		// get error msg, return it as error
		b, err0 := readAllFn(resp.Body)
		if err0 != nil {
			return nil, err0
		}
		return nil, errors.New(string(b))
	}
	if resp.Header.Get("Content-Encoding") == "gzip" {
		// read by gzip
		respBody, err0 := newGzipReaderFn(resp.Body)
		if err0 != nil {
			return nil, err0
		}
		resp.Body = respBody
	}

	respData, err := readAllFn(resp.Body)
	if err != nil {
		return nil, err
	}
	return respData, nil
}
