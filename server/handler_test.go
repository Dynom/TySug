package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"bytes"
	"encoding/json"

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
