package service

type Interface interface {
	Rank(input string) (string, float64)
}
