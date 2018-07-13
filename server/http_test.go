package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetRequestFromHTTPRequest(t *testing.T) {

	t.Run("empty payload", func(t *testing.T) {
		request, err := createStubbedTySugRequest(strings.NewReader("{}"))
		if err != nil {
			t.Errorf("Expected error was thrown '%s'\n%+v", err, request)
		}
	})

	t.Run("payload too large", func(t *testing.T) {
		i := strings.Repeat("a", maxBodySize+1)
		payload, _ := json.Marshal(tySugRequest{Input: i})
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(payload))

		request, err := getRequestFromHTTPRequest(req)
		if err != ErrBodyTooLarge {
			t.Errorf("Expected an error to be thrown %+v", request)
		}
	})

	t.Run("invalid request body", func(t *testing.T) {
		request, err := createStubbedTySugRequest(strings.NewReader("{Input: nil}"))
		if err != ErrInvalidRequestBody {
			t.Errorf("Expected an error to be thrown %+v", request)
		}
	})
}

func createStubbedTySugRequest(r io.Reader) (tySugRequest, error) {
	req := httptest.NewRequest(http.MethodPost, "/", r)
	return getRequestFromHTTPRequest(req)
}
