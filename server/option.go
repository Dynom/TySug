package server

import (
	"net/http"

	"fmt"

	"errors"

	"github.com/NYTimes/gziphandler"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

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

type Header struct {
	Name  string
	Value string
}

func WithDefaultHeaders(headers []Header) Option {
	return func(server *TySugServer) {
		server.handlers = append(server.handlers, func(handler http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				for _, h := range headers {
					w.Header().Set(h.Name, h.Value)
				}

				handler.ServeHTTP(w, req)
			})
		})
	}
}

// WithCORS adds the CORS handler to the request handling
func WithCORS(allowedOrigins []string) Option {
	c := createCORSType(allowedOrigins)

	return func(server *TySugServer) {
		if len(allowedOrigins) == 0 {
			server.Logger.Warn("Allowing any Origin. This is an insecure CORS setup, this is NOT recommended for real world usage!")
		}

		server.handlers = append(server.handlers, c.Handler)
	}
}

func createCORSType(allowedOrigins []string) *cors.Cors {
	c := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{http.MethodPost},
		AllowedOrigins:   allowedOrigins,
	})

	return c
}

// WithInputLimitValidator specifies a max input-value limit validator
func WithInputLimitValidator(inputMax int) Option {
	return func(server *TySugServer) {
		server.validators = append(server.validators, createInputLimitValidator(inputMax))
	}
}

func createInputLimitValidator(inputMax int) Validator {
	return func(TSRequest tySugRequest) error {
		if len(TSRequest.Input) > inputMax {
			return fmt.Errorf("WithInputLimitValidator input exceeds server specified maximum of %d bytes", inputMax)
		}

		if len(TSRequest.Input) == 0 {
			return errors.New("WithInputLimitValidator input may not be empty")
		}

		return nil
	}
}

// WithGzipHandler adds a gzip handler to the server's handlers
func WithGzipHandler() Option {
	return func(server *TySugServer) {
		server.handlers = append(server.handlers, gziphandler.GzipHandler)
	}
}

// WithPProf enables pprof
func WithPProf(prefix string) Option {
	return func(server *TySugServer) {
		server.profConfig.Enable = true
		server.profConfig.Prefix = prefix
	}
}
