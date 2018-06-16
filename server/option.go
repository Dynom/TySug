package server

import (
	"net/http"

	"fmt"

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

// WithCORS adds the CORS handler to the request handling
func WithCORS(allowedOrigins []string) Option {

	if cap(allowedOrigins) == 0 {
		allowedOrigins = []string{"*"}
	}

	c := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{http.MethodPost},
		AllowedOrigins:   allowedOrigins,
	})

	return func(server *TySugServer) {
		server.handlers = append(server.handlers, c.Handler)
	}
}

// WithInputLimitValidator specifies a max input-value limit validator
func WithInputLimitValidator(inputMax int) Option {
	return func(server *TySugServer) {
		server.validators = append(server.validators, func(TSRequest tySugRequest) error {
			if len(TSRequest.Input) > inputMax {
				return fmt.Errorf("WithInputLimitValidator input exceeds server specified maximum of %d bytes", inputMax)
			}

			return nil
		})
	}
}
