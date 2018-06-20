package server

import (
	"context"
	"testing"
)

type stubSvc struct {
}

func (stubSvc) Find(ctx context.Context, input string) (string, float64) {
	return "", 0
}

func TestHasServiceForList(t *testing.T) {
	sr := NewServiceRegistry()
	sr.Register("existing", stubSvc{})

	if sr.HasServiceForList("non-existing") {
		t.Errorf("Expected false when checking for a non-existing service")
	}

	if !sr.HasServiceForList("existing") {
		t.Errorf("Expected true when checking for an existing service")
	}
}

func TestGetServiceForList(t *testing.T) {
	sr := NewServiceRegistry()
	sr.Register("existing", stubSvc{})

	{
		result := sr.GetServiceForList("non-existing")

		if svc, ok := result.(stubSvc); ok {
			t.Errorf("Expected the non-existing service not to exist! %#v", svc)
		}
	}

	{
		result := sr.GetServiceForList("existing")

		if svc, ok := result.(stubSvc); !ok {
			t.Errorf("Expected the existing service to exist! %#v", svc)
		}
	}
}
