package server

// Service is the type any service must implement
type Service interface {
	Find(input string) (string, float64)
}
