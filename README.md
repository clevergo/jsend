# JSend's implementation writen in Go(golang)
[![Build Status](https://travis-ci.org/clevergo/jsend.svg?branch=master)](https://travis-ci.org/clevergo/jsend) [![Coverage Status](https://coveralls.io/repos/github/clevergo/jsend/badge.svg?branch=master)](https://coveralls.io/github/clevergo/jsend?branch=master)  [![GoDoc](https://img.shields.io/badge/godoc-reference-blue)](https://pkg.go.dev/github.com/clevergo/jsend) [![Go Report Card](https://goreportcard.com/badge/github.com/clevergo/jsend)](https://goreportcard.com/report/github.com/clevergo/jsend) [![Release](https://img.shields.io/github/release/clevergo/jsend.svg?style=flat-square)](https://github.com/clevergo/jsend/releases)

This package is an implementation of [JSend](https://github.com/omniti-labs/jsend) specification written in Go(golang).

## Installation

```shell
go get github.com/clevergo/jsend
```

## Usage

Usage is pretty simple.

```go
// success response
jsend.Success(w, "any type of data")
// fail response
jsend.Fail(w, "any type of data")
// error response
jsend.Error(w, "error message")
// error response with extra code
jsend.ErrorCode(w, "error message", errorCode)
// error response with extra code and data
jsend.ErrorCodeData(w, "error message", errorCode, "any type of data")
```

See [Example](example).

### Error Handling

It is application responsibility to handle error, see [Example](example) for further more detail.

### Status Code

By default status code `http.StatusOK` was used implicitly,
it can also be specified by the last parameter if necessary.

```go
jsend.Success(w, "any type of data", http.StatusOK)
jsend.Fail(w, "any type of data", http.StatusForbidden)
jsend.Error(w, "error message", http.StatusInternalServerError)
```