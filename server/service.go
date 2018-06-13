package server

// Service is the type any service must implement
type Service interface {
	Rank(input string) (string, float64)
}
