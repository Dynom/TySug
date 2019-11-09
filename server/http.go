package server

import (
	"encoding/json"
	"net/http"

	"time"

	"io"
	"io/ioutil"

	"errors"

	"net/http/pprof"

	"github.com/sirupsen/logrus"
	tomb "gopkg.in/tomb.v2"
)

const maxBodySize = 1 << 20 // 1Mb

// Errors
var (
	ErrMissingBody        = errors.New("missing body")
	ErrInvalidRequest     = errors.New("invalid request")
	ErrInvalidRequestBody = errors.New("invalid request body")
	ErrBodyTooLarge       = errors.New("body too large")
)

type contextKey int

// Context value parameters
const (
	CtxRequestID contextKey = iota
)

// Header constants
const (
	HeaderRequestID = "X-Request-ID"
)

type tySugRequest struct {
	Input string `json:"input"`
}

type tySugResponse struct {
	Result string  `json:"result"`
	Score  float64 `json:"score"`
	Exact  bool    `json:"exact_match"`
}

type pprofConfig struct {
	Enable bool
	Prefix string
}

// Validator is a type to validate a client request, returning a nil error means all went well.
type Validator func(TSRequest tySugRequest) error

// TySugServer the HTTP server
type TySugServer struct {
	server     *http.Server
	handlers   []func(h http.Handler) http.Handler
	validators []Validator
	profConfig pprofConfig

	Logger *logrus.Logger
}

// ListenOnAndServe allows to set the host:port URL late. It calls ListenAndServe()
func (tss *TySugServer) ListenOnAndServe(addr string) error {
	tss.server.Addr = addr

	return tss.server.ListenAndServe()
}

// NewHTTP constructs a new TySugServer
func NewHTTP(sr ServiceRegistry, mux *http.ServeMux, options ...Option) TySugServer {
	tySug := TySugServer{
		Logger: logrus.StandardLogger(),
	}

	for _, opt := range options {
		opt(&tySug)
	}

	var handler http.Handler = createRequestIDHandler(mux)
	for _, h := range tySug.handlers {
		handler = h(handler)
	}

	mux.HandleFunc("/", http.NotFound)
	mux.HandleFunc("/list/", serviceHandler(tySug.Logger, sr, tySug.validators))

	tySug.server = &http.Server{
		ReadHeaderTimeout: 2 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second, // Is overridden, when the profiler is enabled.
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 19, // 512 kb
		Handler:           handler,
	}

	if tySug.profConfig.Enable {
		configureProfiler(tySug, mux)
	}

	return tySug
}

func configureProfiler(s TySugServer, mux *http.ServeMux) {
	var prefix string
	if s.profConfig.Prefix != "" {
		prefix = s.profConfig.Prefix
	} else {
		prefix = "debug"
	}

	mux.HandleFunc(`/`+prefix+`/pprof/`, pprof.Index)
	mux.HandleFunc(`/`+prefix+`/pprof/cmdline`, pprof.Cmdline)
	mux.HandleFunc(`/`+prefix+`/pprof/profile`, pprof.Profile)
	mux.HandleFunc(`/`+prefix+`/pprof/symbol`, pprof.Symbol)
	mux.HandleFunc(`/`+prefix+`/pprof/trace`, pprof.Trace)

	// The profiler needs at least 30 seconds to use /prefix/pprof/profile
	s.server.WriteTimeout = 31 * time.Second
}

func createRequestHandler(logger *logrus.Logger, svc Service, validators []Validator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, ctx := tomb.WithContext(r.Context())

		ctxLogger := logger.WithFields(logrus.Fields{
			"request_id": ctx.Value(CtxRequestID),
		})

		req, reqErr := getRequestFromHTTPRequest(r)
		if reqErr != nil {
			if reqErr == ErrInvalidRequestBody {
				ctxLogger.Errorf("Missing or invalid request body.")
			} else {
				ctxLogger.Errorf("Unable to process HTTP request, %s.", reqErr)
			}

			w.WriteHeader(400)
			_, writeErr := w.Write([]byte(reqErr.Error()))
			if writeErr != nil {
				ctxLogger.Errorf("Error while writing 400 error: %s (original error: %q)", writeErr, reqErr)
			}
			return
		}

		for _, v := range validators {
			if vErr := v(req); vErr != nil {
				ctxLogger.WithFields(logrus.Fields{
					"error": vErr,
				}).Error("Request validation failed")

				w.WriteHeader(400)
				_, writeErr := w.Write([]byte("Validation failed."))
				if writeErr != nil {
					ctxLogger.Errorf("Error while writing 400 error: %s", writeErr)
				}
				return
			}
		}

		var res tySugResponse

		start := time.Now()
		res.Result, res.Score, res.Exact = svc.Find(ctx, req.Input)

		ctxLogger.WithFields(logrus.Fields{
			"input":       req.Input,
			"suggestion":  res.Result,
			"score":       res.Score,
			"duration_µs": time.Since(start) / time.Microsecond,
		}).Debug("Completed new ranking request")

		if !t.Alive() {
			ctxLogger.Info("Request got canceled")
		}

		response, err := json.Marshal(res)
		if err != nil {
			w.WriteHeader(500)
			_, writeErr := w.Write([]byte("unable to marshal result, b00m"))
			ctxLogger.Errorf("Error while writing 500 error: %s (original marshaling error: %q)", writeErr, err)
			return
		}

		_, err = w.Write(response)
		if err != nil {
			ctxLogger.Errorf("Error while writing response: %s", err)
		}
	}
}

func getRequestFromHTTPRequest(r *http.Request) (tySugRequest, error) {
	var req tySugRequest

	var maxSizePlusOne int64 = maxBodySize + 1

	if r.Body == nil {
		return req, ErrMissingBody
	}

	b, err := ioutil.ReadAll(io.LimitReader(r.Body, maxSizePlusOne))
	if err != nil {
		if err == io.EOF {
			return req, ErrMissingBody
		}
		return req, ErrInvalidRequest
	}

	if int64(len(b)) == maxSizePlusOne {
		return req, ErrBodyTooLarge
	}

	err = json.Unmarshal(b, &req)
	if err != nil {
		return req, ErrInvalidRequestBody
	}

	return req, nil
}
