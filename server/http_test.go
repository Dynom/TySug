package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sirupsen/logrus/hooks/test"
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

	t.Run("empty request body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodHead, "/", nil)
		req.Body = nil

		request, err := getRequestFromHTTPRequest(req)
		if err != ErrMissingBody {
			t.Errorf("Expected an error to be thrown %+v", request)
		}
	})
}

func TestCreateRequestHandler(t *testing.T) {
	logger, _ := test.NewNullLogger()

	t.Run("missing body", func(t *testing.T) {
		h := createRequestHandler(logger, nil, nil)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/", nil)
		h.ServeHTTP(w, r)

		if w.Code != 400 {
			t.Errorf("Expected error was thrown\n%+v", w)
		}
	})

	t.Run("broken JSON", func(t *testing.T) {
		h := createRequestHandler(logger, nil, nil)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("{B0roken JSON"))
		h.ServeHTTP(w, r)

		if w.Code != 400 {
			t.Errorf("Expected error was thrown\n%+v", w)
		}
	})

	t.Run("validation is performed", func(t *testing.T) {
		h := createRequestHandler(logger, nil, []Validator{createInputLimitValidator(3)})

		payload, _ := json.Marshal(tySugRequest{Input: "input too large"})
		r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(payload))
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)

		t.Logf("Body: %s", w.Body)
		if w.Code != 400 {
			t.Errorf("Expected error was thrown\n%+v", w)
		}
	})
}

func TestNewHTTP(t *testing.T) {
	const pprofPrefix = "blup"
	h := NewHTTP(ServiceRegistry{}, http.DefaultServeMux,
		WithGzipHandler(),
		WithPProf(pprofPrefix),
	)

	r := httptest.NewRequest(http.MethodGet, "/"+pprofPrefix+"/pprof/heap", nil)
	w := httptest.NewRecorder()
	h.server.Handler.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Errorf("Expected a server with pprof enabled at the correct prefix %+v", h)
	}
}

func createStubbedTySugRequest(r io.Reader) (tySugRequest, error) {
	req := httptest.NewRequest(http.MethodPost, "/", r)
	return getRequestFromHTTPRequest(req)
}
