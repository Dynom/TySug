package server

import (
	"fmt"
	"os"
)

// Logger is the interface TySugServer uses for logging. Conveniently,
// github.com/sirupsen/logrus.Logger implements this interface.
type Logger interface {
	Errorf(format string, args ...interface{})
}

type defaultLogger struct{}

// Errorf tries to write the (formated) error message to stderr. It will
// print to stdout if the writting to stderr fails.
func (defaultLogger) Errorf(format string, args ...interface{}) {
	_, err := fmt.Fprintf(os.Stderr, format, args...)
	// We do not expect to error here, but we should print it somewhere reasonable.
	if err != nil {
		fmt.Println(append(
			[]interface{}{
				"Error while logging error: " + err.Error() + " (original message: +" + format + "+) Extra:",
			}, args...,
		)...)
	}
}
