package service

// Interface is the type any service must implement
type Interface interface {
	Rank(input string) (string, float64)
}
