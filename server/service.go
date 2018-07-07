package server

import "context"

// Service is the type any service must implement
type Service interface {
	Find(ctx context.Context, input string) (string, float64, bool)
}
