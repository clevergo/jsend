// Copyright 2020 CleverGo. All rights reserved.
// Use of this source code is governed by a MIT style license that can be found
// in the LICENSE file.

package jsend

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrite(t *testing.T) {
	tests := []struct {
		status int
		body   Body
	}{
		{
			body: Body{Status: StatusError, Message: "error"},
		},
		{
			body: Body{Status: StatusError, Message: "error", Code: 10000},
		},
		{
			body: Body{Status: StatusError, Message: "error", Code: 10001, Data: "error data"},
		},
		{
			status: http.StatusInternalServerError,
			body:   Body{Status: StatusError, Message: "error"},
		},
		{
			body: Body{Status: StatusFail, Data: "fail"},
		},
		{
			status: http.StatusForbidden,
			body:   Body{Status: StatusFail, Data: "fail"},
		},
		{
			body: Body{Status: StatusSuccess, Data: "success"},
		},
	}

	contentType := "application/json; charset=utf-8"
	for _, test := range tests {
		response := httptest.NewRecorder()
		var err error
		var statuses []int
		if test.status != 0 {
			statuses = append(statuses, test.status)
		}
		if test.body.Status == StatusSuccess {
			err = Success(response, test.body.Data, statuses...)
		} else if test.body.Status == StatusFail {
			err = Fail(response, test.body.Data, statuses...)
		} else {
			if test.body.Data != nil {
				err = ErrorCodeData(response, test.body.Message, test.body.Code, test.body.Data, statuses...)
			} else if test.body.Code != 0 {
				err = ErrorCode(response, test.body.Message, test.body.Code, statuses...)
			} else {
				err = Error(response, test.body.Message, statuses...)
			}
		}
		if err != nil {
			t.Fatal(err)
		}

		actualContentType := response.Header().Get("Content-Type")
		assert.Equal(t, contentType, actualContentType)
		if test.status != 0 {
			assert.Equal(t, test.status, response.Result().StatusCode)
		}
		var actualBody Body
		if err = json.Unmarshal(response.Body.Bytes(), &actualBody); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, test.body, actualBody)
	}
}

func TestWriteError(t *testing.T) {
	w := httptest.NewRecorder()
	err := Write(w, NewError("msg", 0, make(chan int)))
	assert.NotNil(t, err)
}
