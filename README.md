# JSend's implementation writen in Go(golang)
[![Build Status](https://travis-ci.org/clevergo/jsend.svg?branch=master)](https://travis-ci.org/clevergo/jsend)
[![Coverage Status](https://coveralls.io/repos/github/clevergo/jsend/badge.svg?branch=master)](https://coveralls.io/github/clevergo/jsend?branch=master)
[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/clevergo.tech/jsend?tab=doc)
[![Go Report Card](https://goreportcard.com/badge/github.com/clevergo/jsend)](https://goreportcard.com/report/github.com/clevergo/jsend)
[![Release](https://img.shields.io/github/release/clevergo/jsend.svg?style=flat-square)](https://github.com/clevergo/jsend/releases)

This package is an implementation of [JSend](https://github.com/omniti-labs/jsend) specification written in Go(golang).

## Installation

```shell
go get clevergo.tech/jsend
```

## Usage

Usage is pretty simple.

```go
// success response
jsend.Success(w, data)
// fail response
jsend.Fail(w, data)
// error response
jsend.Error(w, message)
// error response with extra code
jsend.ErrorCode(w, message, code)
// error response with extra code and data
jsend.ErrorCodeData(w, message, code, data)
```

It can also be integrated with web framework, such as Gin, Echo, CleverGo:

```go
// success response
ctx.JSON(http.StatusOK, jsend.New(data))
// fail response
ctx.JSON(http.StatusOK, jsend.NewFail(data))
// error response
ctx.JSON(http.StatusOK, jsend.NewError(message, code, data))
```

Checkout [example](https://github.com/clevergo/examples/tree/master/jsend) for details.

### Error Handling

It is application responsibility to handle error.

### Status Code

By default status code `http.StatusOK` was used implicitly,
it can also be specified by the last parameter if necessary.

```go
jsend.Success(w, data, http.StatusOK)
jsend.Fail(w, data, http.StatusForbidden)
jsend.Error(w, message, http.StatusInternalServerError)
```