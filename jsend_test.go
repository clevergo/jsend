package jsend

import "testing"
import "net/http/httptest"
import "encoding/json"
import "reflect"
import "net/http"

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

	contentType := "application/json"
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
		if contentType != actualContentType {
			t.Errorf("expected content type %q, got %q", contentType, actualContentType)
		}
		if test.status != 0 && test.status != response.Result().StatusCode {
			t.Errorf("expected status code %d, got %d", test.status, response.Result().StatusCode)
		}
		var actualBody Body
		if err = json.Unmarshal(response.Body.Bytes(), &actualBody); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(test.body, actualBody) {
			t.Errorf("expected body %v, got %v", test.body, actualBody)
		}
	}
}
