# JSend's implementation writen in Go(golang)
[![Build Status](https://img.shields.io/travis/clevergo/jsend?style=flat-square)](https://travis-ci.org/clevergo/jsend)
[![Coverage Status](https://img.shields.io/coveralls/github/clevergo/jsend?style=flat-square)](https://coveralls.io/github/clevergo/jsend?branch=master)
[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/clevergo.tech/jsend?tab=doc)
[![Go Report Card](https://goreportcard.com/badge/github.com/clevergo/jsend?style=flat-square)](https://goreportcard.com/report/github.com/clevergo/jsend)
[![Release](https://img.shields.io/github/release/clevergo/jsend.svg?style=flat-square)](https://github.com/clevergo/jsend/releases)
[![Downloads](https://img.shields.io/endpoint?url=https://pkg.clevergo.tech/api/badges/downloads/total/clevergo.tech/jsend&style=flat-square)](https://pkg.clevergo.tech/clevergo.tech/jsend)
[![Chat](https://img.shields.io/badge/chat-telegram-blue?style=flat-square)](https://t.me/clevergotech)
[![Community](https://img.shields.io/badge/community-forum-blue?style=flat-square&color=orange)](https://forum.clevergo.tech)

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
