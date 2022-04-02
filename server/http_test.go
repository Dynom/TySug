package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Dynom/TySug/server/service"

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

func TestConfigureProfiler(t *testing.T) {
	t.Run("Testing the WriteTimeout", func(t *testing.T) {
		s := TySugServer{server: &http.Server{}}
		configureProfiler(s, http.NewServeMux())
		if s.server.WriteTimeout < 30 {
			t.Errorf("Expected a write timeout of at least 30 seconds when debugging is enabled.\n%+v", s)
		}
	})

	t.Run("Testing if we fall back to /debug/", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/debug/pprof/", nil)

		mux := http.NewServeMux()
		configureProfiler(TySugServer{server: &http.Server{}}, mux)
		mux.ServeHTTP(w, r)

		if w.Code != 200 {
			t.Error("Expected the default route to be available under /debug/")
		}
	})

	t.Run("Testing if listen on the custom prefix", func(t *testing.T) {
		const prefix = "213c06c51ed"

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/"+prefix+"/pprof/", nil)

		s := TySugServer{server: &http.Server{}}
		s.profConfig.Prefix = prefix

		mux := http.NewServeMux()
		configureProfiler(s, mux)
		mux.ServeHTTP(w, r)

		if w.Code != 200 {
			t.Errorf("Expected the custom route to be available under /%s/.\n%+v", prefix, w)
		}
	})
}

func TestHeaders(t *testing.T) {
	// Setting up the service registry
	sr := NewServiceRegistry()
	log, _ := test.NewNullLogger()
	svc, err := service.NewDomain([]string{"foo"}, log)
	if err != nil {
		t.Error(err)
	}

	sr.Register("test", svc)

	t.Run("No extra headers defined", func(t *testing.T) {
		h := NewHTTP(sr, http.NewServeMux(),
			WithDefaultHeaders(nil),
		)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/list/test", strings.NewReader(`{"input": "foo"}`))

		h.server.Handler.ServeHTTP(w, r)

		if w.Code != 200 {
			t.Errorf("Sanity check failed, expecting a successful request. Instead I received: %+v", w.Body)
		}

		if h := w.Header().Get("X-Test-Header"); h != "" {
			t.Errorf("Expecting the test header to not be defined. Instead I got: %+v", h)
		}
	})

	t.Run("With extra headers defined", func(t *testing.T) {
		h := NewHTTP(sr, http.NewServeMux(),
			WithDefaultHeaders(http.Header{
				"X-Test-Header": {"beep"},
				"X-Beep":        {"42"},
			}),
		)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/list/test", strings.NewReader(`{"input": "foo"}`))

		h.server.Handler.ServeHTTP(w, r)

		if w.Code != 200 {
			t.Errorf("Sanity check failed, expecting a successful request. Instead I received: %+v", w.Body)
		}

		if h := w.Header().Get("X-Test-Header"); h != "beep" {
			t.Errorf("Expecting the test header 'X-Test-Header' to be defined.")
		}

		if h := w.Header().Get("X-Beep"); h != "42" {
			t.Errorf("Expecting the test header 'X-Beep' to be defined.")
		}
	})
}

func createStubbedTySugRequest(r io.Reader) (tySugRequest, error) {
	req := httptest.NewRequest(http.MethodPost, "/", r)
	return getRequestFromHTTPRequest(req)
}
