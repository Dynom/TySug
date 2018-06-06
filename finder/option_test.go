package finder

import "testing"

func TestOptExampleAlgorithm(t *testing.T) {
	alg := exampleAlgorithm

	if s := alg("", "apple juice"); s != 0 {
		t.Errorf("Expected the example algorithm to return 0 when an argument is empty.")
	}

	if s := alg("apple juice", ""); s != 0 {
		t.Errorf("Expected the example algorithm to return 0 when an argument is empty.")
	}

	if s := alg("apple", "juice"); s != 0 {
		t.Errorf("Expected the example algorithm to return 0 when the values don't match.")
	}

	if s := alg("tree", "trie"); s != 1 {
		t.Errorf("Expected the example algorithm to return 1 when the first letters match.")
	}
}
