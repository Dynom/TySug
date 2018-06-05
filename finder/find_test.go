package finder

import "testing"

func TestNewWithCustomAlgorithm(t *testing.T) {
	sug, _ := New([]string{"b"}, OptExampleAlgorithm)

	var score float64

	_, score = sug.Find("a")
	if score != 0 {
		t.Errorf("Expected the score to be 0, instead I got %f.", score)
	}

	_, score = sug.Find("b")
	if score != 1 {
		t.Errorf("Expected the score to be 1, instead I got %f.", score)
	}
}

func TestNoAlgorithm(t *testing.T) {
	_, err := New([]string{})

	if err != ErrNoAlgorithmDefined {
		t.Errorf("Expected an error to be returned when no algorithm was specified.")
	}
}