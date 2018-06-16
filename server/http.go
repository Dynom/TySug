package server

import (
	"encoding/json"
	"net/http"

	"time"

	"io"
	"io/ioutil"

	"errors"

	"context"

	"math/rand"

	"sync"

	"strings"

	"strconv"

	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"gopkg.in/tomb.v2"
)

const maxBodySize = 1 << 20 // 1Mb

// Errors
var (
	ErrMissingBody        = errors.New("missing body")
	ErrInvalidRequest     = errors.New("invalid request")
	ErrInvalidRequestBody = errors.New("invalid request body")
)

const (
	CtxRequestID = iota
)

type request struct {
	Input string `json:"input"`
}

type response struct {
	Result string  `json:"result"`
	Score  float64 `json:"score"`
}

// TySugServer the HTTP server
type TySugServer struct {
	server *http.Server

	Logger *logrus.Logger
}

// Option is a handy type used for configuration purposes
type Option func(*TySugServer)

// WithLogger sets the logger to be used when encountering http-related errors.
// Errors are written to the standard error output in most cases. Printing on the
// standard output is reserved to extreme case where writing on stderr failed.
func WithLogger(logger *logrus.Logger) Option {
	if logger == nil {
		logger = logrus.StandardLogger()
	}

	return func(server *TySugServer) {
		server.Logger = logger
	}
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

	c := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{http.MethodPost},
	})

	mux.HandleFunc("/", http.NotFound)
	mux.HandleFunc("/list/", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Path[6:]
		if name == "" {
			tySug.Logger.Info("no list name defined")
			return
		}

		if !sr.HasServiceForList(name) {
			tySug.Logger.Infof("list '%s' not defined", name)
			return
		}

		svc := sr.GetServiceForList(name)
		hf := createRequestHandler(
			tySug.Logger,
			svc,
		)

		hf(w, r)
	})

	tySug.server = &http.Server{
		ReadHeaderTimeout: 2 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 19, // 512 kb
		Handler:           defaultHeaderHandler(c.Handler(createRequestIDHandler(mux))),
	}

	return tySug
}

func createRequestIDHandler(h http.Handler) http.HandlerFunc {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	instanceID := rnd.Int31()

	var requestCounter int
	var lock = sync.Mutex{}
	var buf = strings.Builder{}
	return func(w http.ResponseWriter, r *http.Request) {
		lock.Lock()
		requestCounter++
		buf.Reset()
		buf.WriteString(strconv.Itoa(int(instanceID)))
		buf.WriteString("-")
		buf.WriteString(strconv.Itoa(requestCounter))
		requestID := buf.String()
		lock.Unlock()

		ctx := context.WithValue(r.Context(), CtxRequestID, requestID)
		h.ServeHTTP(w, r.WithContext(ctx))
	}
}

func createRequestHandler(logger *logrus.Logger, svc Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, ctx := tomb.WithContext(r.Context())

		req, err := getRequestFromHTTPRequest(r)
		if err != nil {
			w.WriteHeader(400)
			_, writeErr := w.Write([]byte(err.Error()))
			if writeErr != nil {
				logger.Errorf("Error while writing 400 error: %s (original error: %q)", writeErr, err)
			}
			return
		}

		var res response

		start := time.Now()
		res.Result, res.Score = svc.Find(ctx, req.Input)

		logger.WithFields(logrus.Fields{
			"input":       req.Input,
			"suggestion":  res.Result,
			"score":       res.Score,
			"duration_µs": time.Since(start) / time.Microsecond,
			"request_id":  ctx.Value(CtxRequestID),
		}).Debug("Completed new ranking request")

		if !t.Alive() {
			logger.WithFields(logrus.Fields{
				"request_id": ctx.Value(CtxRequestID),
			}).Info("Request got cancelled")
		}

		response, err := json.Marshal(res)
		if err != nil {
			w.WriteHeader(500)
			_, writeErr := w.Write([]byte("unable to marshal result, b00m"))
			logger.Errorf("Error while writing 500 error: %s (original marshaling error: %q)", writeErr, err)
			return
		}

		_, err = w.Write(response)
		if err != nil {
			logger.Errorf("Error while writing response: %s", err)
		}
	}
}

func defaultHeaderHandler(h http.Handler) http.HandlerFunc {

	type kv struct {
		Key   string
		Value string
	}

	return func(w http.ResponseWriter, req *http.Request) {
		for _, h := range []kv{
			{Key: "Strict-Transport-Security", Value: "max-age=31536000; includeSubDomains"},
			{Key: "Content-Security-Policy", Value: "default-src 'none'"},
			{Key: "X-Frame-Options", Value: "DENY"},
			{Key: "X-XSS-Protection", Value: "1; mode=block"},
			{Key: "X-Content-Type-Options", Value: "nosniff"},
			{Key: "Referrer-Policy", Value: "strict-origin"},
		} {
			w.Header().Set(h.Key, h.Value)
		}

		h.ServeHTTP(w, req)
	}
}

func getRequestFromHTTPRequest(r *http.Request) (request, error) {
	var req request

	b, err := ioutil.ReadAll(io.LimitReader(r.Body, maxBodySize))
	if err != nil {
		if err == io.EOF {
			return req, ErrMissingBody
		}
		return req, ErrInvalidRequest
	}

	err = json.Unmarshal(b, &req)
	if err != nil {
		return req, ErrInvalidRequestBody
	}

	return req, nil
}
