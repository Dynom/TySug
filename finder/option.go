package finder

type Option func(sug *TySug)

func OptExampleAlgorithm(sug *TySug) {
	sug.Alg = exampleAlgorithm
}

func exampleAlgorithm(a, b string) float64 {
	if len(a) == 0 || len(b) == 0 {
		return 0
	}

	if a[0] == b[0] {
		return 1
	}

	return 0
}
