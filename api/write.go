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
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/klauspost/compress/gzip"

	"github.com/lindb/common/series"

	"github.com/lindb/client_go/internal"
	httppkg "github.com/lindb/client_go/internal/http"
)

const (
	// ContentTypeFlat represents flat buffer content type.
	ContentTypeFlat = "application/flatbuffer"
)

var (
	errTooManyRetryRequests = errors.New("too many retry requests, drop current request")
	errTooManyRetry         = errors.New("max retry attempt")
)

// retryReq represents request need to retry.
type retryReq struct {
	data     io.Reader
	attempts int
}

// Write represents write client for writing time series data asynchronously.
type Write interface {
	// AddPoint adds a time series point into buffer.
	AddPoint(ctx context.Context, point *Point)
	// Errors watches error in background goroutine.
	Errors() <-chan error
	// Close closes write client, before close try to send pending points.
	Close()
}

// write implements Write interface.
type write struct {
	endpoint     string
	database     string
	writeOptions *WriteOptions
	client       *http.Client

	bufferCh    chan *Point
	sendCh      chan []byte
	errCh       chan error
	stopBatchCh chan struct{}
	stopSendCh  chan struct{}
	doneCh      chan struct{}

	builder     *series.RowBuilder
	buf         *bytes.Buffer
	gzipBuf     *bytes.Buffer
	gzipWriter  *gzip.Writer
	batchedSize int

	closed bool
	mutex  sync.Mutex
}

// NewWrite creates an asynchronously write client.
func NewWrite(endpoint, database string, writeOptions *WriteOptions, httpOptions *httppkg.Options) Write {
	w := &write{
		endpoint:     fmt.Sprintf("%s/api/v1/write?db=%s", endpoint, database),
		database:     database,
		client:       httpOptions.HTTPClient(),
		writeOptions: writeOptions,
		bufferCh:     make(chan *Point, writeOptions.BatchSize()+1),
		sendCh:       make(chan []byte),
		errCh:        make(chan error),
		stopBatchCh:  make(chan struct{}),
		stopSendCh:   make(chan struct{}),
		doneCh:       make(chan struct{}),
		builder:      series.CreateRowBuilder(),
		buf:          &bytes.Buffer{},
	}
	if writeOptions.UseGZip() {
		w.gzipBuf = &bytes.Buffer{}
		w.gzipWriter = gzip.NewWriter(w.gzipBuf)
	}
	go w.bufferProc() // process point->data([]byte)
	go w.sendProc()   // send data to server
	return w
}

// AddPoint adds a time series point into buffer.
func (w *write) AddPoint(ctx context.Context, point *Point) {
	if point == nil || !point.Valid() {
		return
	}
	select {
	case <-ctx.Done():
	case w.bufferCh <- point:
	}
}

// Errors watches error in background goroutine.
func (w *write) Errors() <-chan error {
	return w.errCh
}

// Close closes write client, before close try to send pending points.
func (w *write) Close() {
	if w.closed {
		return
	}

	w.mutex.Lock()
	defer w.mutex.Unlock()

	// double check
	if w.closed {
		return
	}

	close(w.stopBatchCh)
	close(w.bufferCh)
	<-w.doneCh // wait buffer process completed

	close(w.stopSendCh)
	close(w.sendCh)
	<-w.doneCh // wait send process completed

	close(w.errCh)
	w.closed = true
}

// bufferProc consumes time series point from buffer chan, marshals point then put data into send buffer.
func (w *write) bufferProc() {
	batchSize := w.writeOptions.BatchSize()
	ticker := time.NewTicker(time.Duration(w.writeOptions.flushInterval) * time.Millisecond)

	defer func() {
		ticker.Stop()
		w.doneCh <- struct{}{}
	}()

	for {
		select {
		case point := <-w.bufferCh:
			if err := w.batchPoint(point); err != nil {
				w.emitErr(err)
				continue
			}
			// check batch buffer is full
			if w.batchedSize >= batchSize {
				w.flushBuffer()
			}
		case <-ticker.C:
			w.flushBuffer()
		case <-w.stopBatchCh:
			// try to batch pending points
			for point := range w.bufferCh {
				if err := w.batchPoint(point); err != nil {
					w.emitErr(err)
				}
			}
			w.flushBuffer()
			return
		}
	}
}

// flushBuffer flushes buffer data, put data into send chan, then clear buffer.
func (w *write) flushBuffer() {
	if w.batchedSize == 0 {
		return
	}
	data := w.buf.Bytes()
	w.buf.Reset() // reset batch buf
	w.batchedSize = 0

	// put data into send chan
	w.sendCh <- data
}

// batchPoint marshals point, if success put data into buffer.
func (w *write) batchPoint(point *Point) error {
	if point == nil {
		return nil
	}
	defer w.builder.Reset()

	builder := w.builder

	builder.AddNameSpace(internal.String2ByteSlice(point.namespace))
	builder.AddMetricName(internal.String2ByteSlice(point.MetricName()))
	builder.AddTimestamp(point.Timestamp().UnixMilli())

	addTag := func(tags map[string]string) error {
		for k, v := range tags {
			if err := builder.AddTag(internal.String2ByteSlice(k), internal.String2ByteSlice(v)); err != nil {
				return err
			}
		}
		return nil
	}
	// add default tags
	if err := addTag(w.writeOptions.DefaultTags()); err != nil {
		return err
	}
	// add tags of current point
	if err := addTag(point.Tags()); err != nil {
		return err
	}

	// write field
	fields := point.Fields()
	for _, f := range fields {
		if err := f.write(builder); err != nil {
			return err
		}
	}

	// put point into buffer
	data, err := builder.Build()
	if err != nil {
		return err
	}
	_, err = w.buf.Write(data)
	if err != nil {
		return err
	}
	w.batchedSize++
	return nil
}

// sendProc consumes batched write data, then send it to broker.
func (w *write) sendProc() {
	defer func() {
		// invoke when send goroutine exit.
		w.doneCh <- struct{}{}
	}()
	retryBufferLimit := w.writeOptions.RetryBufferLimit()
	maxRetries := w.writeOptions.MaxRetries()
	retryBuffers := make([]*retryReq, 0)
	retry := func(data io.Reader, attempt int) {
		if attempt >= maxRetries {
			w.emitErr(errTooManyRetry)
			return
		}
		if len(retryBuffers) > retryBufferLimit {
			w.emitErr(errTooManyRetryRequests)
			return
		}
		retryBuffers = append(retryBuffers, &retryReq{
			data:     data,
			attempts: attempt + 1,
		})
	}
	// send write data
	send := func(data []byte) bool {
		if len(data) == 0 {
			return false
		}
		// try compress data
		reqData, err := w.compress(data)
		if err != nil {
			w.emitErr(err)
			return true
		}
		if err := w.send(reqData); err != nil {
			w.emitErr(err)
			retry(reqData, 0)
			return false
		}
		return true
	}
	// send failed write request
	sendRetryReq := func(needRetry bool) {
		if len(retryBuffers) > 0 {
			messages := retryBuffers
			retryBuffers = make([]*retryReq, 0)
			for _, msg := range messages {
				if err := w.send(msg.data); err != nil {
					w.emitErr(err)
					if needRetry {
						retry(msg.data, msg.attempts)
					}
				}
			}
		}
	}
	for {
		select {
		case data := <-w.sendCh:
			if send(data) {
				// if send ok, retry pending failed request
				sendRetryReq(true)
			}
		case <-w.stopSendCh:
			// try to send pending messages
			for data := range w.sendCh {
				_ = send(data)
			}
			sendRetryReq(false)
			return
		}
	}
}

// send write data to broker.
func (w *write) send(data io.Reader) error {
	req, _ := http.NewRequestWithContext(context.TODO(), http.MethodPut, w.endpoint, data)
	if w.gzipWriter != nil {
		req.Header.Set("Content-Encoding", "gzip")
	}
	req.Header.Set("User-Agent", httppkg.UserAgent)
	req.Header.Set("Content-Type", ContentTypeFlat)

	resp, err := w.client.Do(req)
	defer func() {
		// need close resp body by defer, maybe resp is not nil when throw some err
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		// get error msg, return it as error
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return errors.New(string(b))
	}
	// send data success
	return nil
}

// compress request body if it needs.
func (w *write) compress(data []byte) (*bytes.Buffer, error) {
	if w.gzipWriter != nil {
		w.gzipWriter.Reset(w.gzipBuf)
		if _, err := w.gzipWriter.Write(data); err != nil {
			return nil, err
		}
		if err := w.gzipWriter.Close(); err != nil {
			return nil, err
		}
		return w.gzipBuf, nil
	}
	return bytes.NewBuffer(data), nil
}

// emitErr emits error into chan.
func (w *write) emitErr(err error) {
	select {
	case w.errCh <- err:
	default:
		// no err read, cannot put err into chan
	}
}
