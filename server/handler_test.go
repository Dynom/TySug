package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/sirupsen/logrus/hooks/test"
)

func TestServiceHandlerHappyFlow(t *testing.T) {
	logger, _ := test.NewNullLogger()
	recorder := httptest.NewRecorder()

	{ // setup
		payload, _ := json.Marshal(tySugRequest{Input: "baz"})
		req := httptest.NewRequest(http.MethodPost, "/list/foo", bytes.NewReader(payload))
		sr := NewServiceRegistry()
		sr.Register("foo", stubSvc{FindResult: "bar"})

		hf := http.HandlerFunc(serviceHandler(logger, sr, nil))
		hf.ServeHTTP(recorder, req)
	}

	t.Run("test status code", func(t *testing.T) {
		if status := recorder.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})

	t.Run("test response result", func(t *testing.T) {
		expect := tySugResponse{Result: "bar"}
		var result tySugResponse
		err := json.Unmarshal(recorder.Body.Bytes(), &result)
		if err != nil {
			t.Errorf("unexpected error while unmarshalling the response type %s", err)
		}

		if result.Result != expect.Result {
			t.Errorf("expected the input to be %s, instead I got %s", expect.Result, result.Result)
		}
	})
}

func TestServiceHandlerInvalidListName(t *testing.T) {
	logger, _ := test.NewNullLogger()

	t.Run("status code, nonexistent list name", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		{ // setup
			payload, _ := json.Marshal(tySugRequest{Input: "baz"})
			req := httptest.NewRequest(http.MethodPost, "/list/not-existing", bytes.NewReader(payload))

			sr := NewServiceRegistry()
			hf := http.HandlerFunc(serviceHandler(logger, sr, nil))
			hf.ServeHTTP(recorder, req)
		}

		if status := recorder.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
		}
	})

	t.Run("status code, unspecified list name", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		{ // setup
			payload, _ := json.Marshal(tySugRequest{Input: "baz"})
			req := httptest.NewRequest(http.MethodPost, "/list/", bytes.NewReader(payload))

			sr := NewServiceRegistry()
			hf := http.HandlerFunc(serviceHandler(logger, sr, nil))
			hf.ServeHTTP(recorder, req)
		}

		if status := recorder.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})
}

func TestRequestID(t *testing.T) {
	recorder1 := httptest.NewRecorder()
	recorder2 := httptest.NewRecorder()

	{ // setup
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		h := createRequestIDHandler(http.NewServeMux())
		hf := http.HandlerFunc(h)

		// Recording two requests
		hf.ServeHTTP(recorder1, req)
		hf.ServeHTTP(recorder2, req)
	}

	rid1 := recorder1.Result().Header.Get(HeaderRequestID)
	rid2 := recorder2.Result().Header.Get(HeaderRequestID)
	if rid1 == rid2 {
		t.Errorf("Did not expect the request ID's to be identical: %s vs %s", rid1, rid2)
	}

	if rid1 == "" {
		t.Errorf("Did not expect the request ID to be empty.")
	}
}

func TestCORS(t *testing.T) {
	recorder := httptest.NewRecorder()

	reqOrigin := "https://example.org"
	reqMethod := http.MethodPost
	reqHeader := "X-Foo-Bar"

	{ // setup
		req := httptest.NewRequest(http.MethodOptions, "/", nil)
		req.Header.Set("Origin", reqOrigin)
		req.Header.Set("Access-Control-Request-Method", reqMethod)
		req.Header.Set("Access-Control-Request-Headers", reqHeader)

		c := createCORSType(nil)
		c.Handler(http.NotFoundHandler()).ServeHTTP(recorder, req)
	}

	resultOrigin := recorder.Result().Header.Get("Access-Control-Allow-Origin")
	resultMethods := recorder.Result().Header.Get("Access-Control-Allow-Methods")
	resultHeaders := recorder.Result().Header.Get("Access-Control-Allow-Headers")

	if resultOrigin != "*" {
		t.Errorf("Expected the origin to be a wildcard, since none were allow-listed.")
		t.Logf("%+v", recorder.Result())
	}

	if !strings.Contains(resultHeaders, reqHeader) {
		t.Errorf("Expected the headers to be present")
	}

	if !strings.Contains(resultMethods, reqMethod) {
		t.Errorf("Expected the methods to be present")
	}
}

func BenchmarkRequestIDGeneration(b *testing.B) {
	h := createRequestIDHandler(noopHandler{})
	request := httptest.NewRequest(http.MethodPost, "/list/foo", nil)
	recorder := httptest.NewRecorder()

	b.Run("implementation", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			h(recorder, request)
		}

		b.Logf("Last request ID is %s", recorder.Header().Get(HeaderRequestID))
	})

	b.Run("buffer", func(b *testing.B) {
		buf := strings.Builder{}
		buf.Grow(256)
		var result string
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buf.Reset()
			buf.WriteString("1228958371")
			buf.WriteString("-")
			buf.WriteString(strconv.Itoa(i))
			result = buf.String()
		}

		b.Logf("Last request ID is: %s", result)
	})

	b.Run("buffer custom", func(b *testing.B) {
		buf := make([]byte, 0, 3)
		var result string
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buf = buf[0:0]
			buf = append(buf, "1228958371"...)
			buf = append(buf, "-"...)
			buf = append(buf, strconv.Itoa(i)...)
			result = string(buf)
		}

		b.Logf("Last request ID is: %s", result)
	})

	// The fastest (for short strings) and simplest solution
	b.Run("string", func(b *testing.B) {
		foo := "1228958371"
		var result string

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result = foo + "-" + strconv.Itoa(i)
		}

		b.Logf("Last request ID is: %s", result)
	})
}

type noopHandler struct{}

func (noopHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}
