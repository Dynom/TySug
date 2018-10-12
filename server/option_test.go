package server

import (
	"testing"

	"github.com/sirupsen/logrus/hooks/test"
)

func TestWithInputLimitValidator(t *testing.T) {
	s := TySugServer{}

	before := len(s.validators)
	WithInputLimitValidator(12)(&s)

	after := len(s.validators)

	if before >= after {
		t.Errorf("Expected a validator to have been added (Before: %d, after: %d)", before, after)
	}
}

func TestCreateInputLimitValidator(t *testing.T) {
	v := createInputLimitValidator(12)
	req := tySugRequest{Input: "more than twelve"}
	if err := v(req); err == nil {
		t.Errorf("Expected the request to be invalid, since the input is more than 12 bytes in size %+v", err)
	}

	req = tySugRequest{Input: "lt 12"}
	if err := v(req); err != nil {
		t.Errorf("Expected the request to be valid, since the input is less than 12 bytes in size %+v", err)
	}

	req = tySugRequest{Input: ""}
	if err := v(req); err == nil {
		t.Errorf("Expected the request to be invalid, since the input is empty %+v", err)
	}
}

func TestWithCORS(t *testing.T) {
	t.Run("OneOrigin", func(t *testing.T) {
		s := TySugServer{}

		hc := len(s.handlers)
		WithCORS([]string{"localhost"})(&s)

		if l := len(s.handlers); hc > l || l != 1 {
			t.Errorf("Expected exactly one handler, instead I got %d.", l)
		}
	})

	t.Run("NoOrigin", func(t *testing.T) {
		s := TySugServer{}
		log, hook := test.NewNullLogger()
		s.Logger = log

		hc := len(s.handlers)
		lmc := len(hook.Entries)
		WithCORS([]string{})(&s)

		if l := len(s.handlers); hc > l || l != 1 {
			t.Errorf("Expected exactly one handler, instead I got %d.", l)
		}

		// We're expecting a log message warning about unsafe usage since no origin was specified.
		if l := len(hook.Entries); lmc >= l || l != 1 {
			t.Errorf("Expected exactly one log messages, instead I got %d.", l)
		}
	})
}

func TestWithGzipHandler(t *testing.T) {
	s := TySugServer{}

	hc := len(s.handlers)
	WithGzipHandler()(&s)

	if l := len(s.handlers); hc > l || l != 1 {
		t.Errorf("Expected exactly one handler, instead I got %d.", l)
	}
}

func TestWithLogger(t *testing.T) {
	s := TySugServer{}

	WithLogger(nil)(&s)
	if s.Logger == nil {
		t.Errorf("Expected a default logger to be defined, if none was specified")
	}
}
