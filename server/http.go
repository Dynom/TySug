package server

import (
	"encoding/json"
	"net/http"

	"time"

	"io"
	"io/ioutil"

	"errors"

	"github.com/rs/cors"
)

const maxBodySize = 1 << 20 // 1Mb

// Errors
var (
	ErrMissingBody        = errors.New("missing body")
	ErrInvalidRequest     = errors.New("invalid request")
	ErrInvalidRequestBody = errors.New("invalid request body")
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

	Logger Logger
}

// Option is a handy type used for configuration purposes
type Option func(*TySugServer)

// WithLogger sets the logger to be used when encountering http-related errors.
// Errors are written to the standard error output in most cases. Printing on the
// standard output is reserved to extreme case where writing on stderr failed.
func WithLogger(logger Logger) Option {
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
func NewHTTP(svc Service, mux *http.ServeMux, options ...Option) TySugServer {
	tySug := TySugServer{
		Logger: defaultLogger{},
	}

	for _, opt := range options {
		opt(&tySug)
	}

	c := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut},
	})

	mux.HandleFunc("/", createRequestHandler(tySug.Logger, svc))

	tySug.server = &http.Server{
		ReadHeaderTimeout: 2 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 19, // 512 kb
		Handler:           defaultHeaderHandler(c.Handler(mux)),
	}

	return tySug
}

func createRequestHandler(logger Logger, svc Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := getRequestFromHTTPRequest(r)
		if err != nil {
			w.WriteHeader(400)
			_, writeErr := w.Write([]byte(err.Error()))
			if writeErr != nil {
				logger.Errorf("Errored while writing 400 error: %s (original error: %q)", writeErr, err)
			}
			return
		}

		var res response
		res.Result, res.Score = svc.Rank(req.Input)

		response, err := json.Marshal(res)
		if err != nil {
			w.WriteHeader(500)
			_, writeErr := w.Write([]byte("unable to marshal result, b00m"))
			logger.Errorf("Errored while writing 400 error: %s (original marshaling error: %q)", writeErr, err)
			return
		}

		_, err = w.Write(response)
		if err != nil {
			logger.Errorf("Errored while writing respones: %s", err)
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
