# slogctx

[![Go Reference](https://pkg.go.dev/badge/github.com/phsym/slogctx.svg)](https://pkg.go.dev/github.com/phsym/slogctx) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/phsym/slogctx/master/LICENSE) [![Build](https://github.com/phsym/slogctx/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/phsym/slogctx/actions/workflows/go.yml) [![codecov](https://codecov.io/gh/phsym/slogctx/graph/badge.svg?token=ZIJT9L79QP)](https://codecov.io/gh/phsym/slogctx) [![Go Report Card](https://goreportcard.com/badge/github.com/phsym/slogctx)](https://goreportcard.com/report/github.com/phsym/slogctx)

> WORK IN PROGRESS


## Installation
```bash
go get github.com/phsym/slogctx@latest
```

## Example
```go
package main

import (
	"errors"
	"log/slog"
	"os"

	"github.com/phsym/slogctx"
)

func main() {
	logger := slog.New(
		console.NewHandler(os.Stderr, &console.HandlerOptions{Level: slog.LevelDebug}),
	)
	slog.SetDefault(logger)
	// TBD ...
}
```


## Performances
See [benchmark file](./bench_test.go) for details.

TBD ...