# LinDB Client Go

[![LICENSE](https://img.shields.io/github/license/lindb/client_go)](https://github.com/lindb/client_go/blob/main/LICENSE)
[![Language](https://img.shields.io/badge/Language-Go-blue.svg)](https://golang.org/)
[![Go Report Card](https://goreportcard.com/badge/github.com/lindb/client_go)](https://goreportcard.com/report/github.com/lindb/client_go)
[![codecov](https://codecov.io/gh/lindb/client_go/branch/main/graph/badge.svg)](https://codecov.io/gh/lindb/client_go)
[![Github Actions Status](https://github.com/lindb/client_go/workflows/LinDB%20CI/badge.svg)](https://github.com/lindb/client_go/actions?query=workflow%3A%22LinDB+CI%22)
[![GoDoc](https://img.shields.io/badge/Godoc-reference-blue.svg)](https://godoc.org/github.com/lindb/client_go)

This repository contains the reference Go client for LinDB.

- [Features](#features)
- [How To Use](#how-to-use)
  - [Installation](#installation)
  - [Write Data](#write-data)
  - [Reading Background Process Errors](#reading-background-process-errors)
  - [Write Options](#options)

## Features

- Write data
  - Write data use asynchronous
  - Support field type(sum/min/max/last/first/histogram)
  - [FlatBuf Protocol](https://github.com/lindb/common/blob/main/proto/v1/metrics.fbs)

## How To Use

### Installation

Go 1.18 or later is required.

- Add the client package to your project dependencies (go.mod).
   ```sh
   go get github.com/lindb/client_go
   ```
  
- Add import `github.com/lindb/client_go` to your source code.

### Write data

```go
package main

import (
	"context"
	"fmt"
	"time"

	lindb "github.com/lindb/client_go"
	"github.com/lindb/client_go/api"
)

func main() {
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
```

### Reading background process errors

Write client doesn't log any error. Can use [Errors()](https://pkg.go.dev/github.com/lindb/client_go/api#Write) method, which returns the channel for reading errors occurring
during async writes.

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/lindb/client_go"
	"github.com/lindb/client_go/api"
)

func main() {
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
```

### Options

```go
package http

import "crypto/tls"

type Options struct {
	// Request timeout(s), default 30.
	reqTimeout int64
	// TLS configuration for secure connection, default nil.
	tlsConfig *tls.Config
}
```

See [tls.Config](https://pkg.go.dev/crypto/tls#Config) for detail.

```go
package api

type WriteOptions struct {
	// Number of series sent in single write request, default 1000.
	batchSize int
	// Flush interval(ms) which is buffer flushed if it has not been already written, default 1000.
	flushInterval int64
	// Whether to use GZip compress write data, default true.
	useGZip bool
	// Default tags are added to each written series.
	defaultTags map[string]string
	// Maximum count of retry attempts of failed writes, default 3.
	maxRetries int
	// Maximum number of write request to keep for retry, default 100.
	retryBufferLimit int
}
```