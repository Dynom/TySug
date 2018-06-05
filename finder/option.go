package finder

type Option func(sug *TySug)

func OptExampleAlgorithm(sug *TySug) {
	sug.Alg = func(a, b string) float64 {
		if a == b {
			return 1
		}

		return 0
	}
}
