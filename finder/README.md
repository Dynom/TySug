# Finder
[![godoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/Dynom/TySug/finder) [![Go Report Card](https://goreportcard.com/badge/github.com/Dynom/TySug)](https://goreportcard.com/report/github.com/Dynom/TySug)

Finder is a library that finds the best match against a list of strings, using an algorithm of choice.

## Example 
```go
import "github.com/Dynom/TySug/finder"
```
```go
referenceList := []string{"example", "amplifier", "ample"}
ts := finder.New(referenceList, finder.WithAlgorithm(finder.NewJaroWinklerDefault()))

alt, score, exact := ts.Find("exampel")
// alt   = example
// score = 0.9714285714285714
// exact = false (not an exact match in our reference list)
```

## Algorithms
You're free to specify your own algorithm. By default Jaro Winkler is available, this gives you freedom around different input lengths (in contrast to Levenshtein).


### Using a different algorithm

if you want to use a different algorithm, simply wrap your algorithm in a `finder.Algorithm` compatible type and pass 
it as argument to the Finder. You can find inspiration in the unit-tests / examples.

Possible considerations:
 - [Levenshtein](https://en.wikipedia.org/wiki/Levenshtein_distance)
 - [Damerau-Levenshtein](https://en.wikipedia.org/wiki/Damerau%E2%80%93Levenshtein_distance)
 - [LCS](https://en.wikipedia.org/wiki/Longest_common_subsequence_problem)
 - [q-gram](https://en.wikipedia.org/wiki/N-gram)
 - [Cosine](https://en.wikipedia.org/wiki/Cosine_similarity)
 - [Jaccard](https://en.wikipedia.org/wiki/Jaccard_index)
 - [Jaro / Jaro-Winkler](https://en.wikipedia.org/wiki/Jaro%E2%80%93Winkler_distance)
 - [Smith-Waterman](https://en.wikipedia.org/wiki/Smith%E2%80%93Waterman_algorithm)
 - [Sift4](https://siderite.blogspot.com/2014/11/super-fast-and-accurate-string-distance.html) (used in [mailcheck.js](https://github.com/mailcheck/mailcheck))
 
Sources:
 - [joyofdata.de/blog/comparison-of-string-distance-algorithms/](https://www.joyofdata.de/blog/comparison-of-string-distance-algorithms/)

## Case-sensitivity

Finder does not normalise words. This means that words are treated in a case-sensitive matter. This is done mostly to
avoid doing unnecessary work in the hot-path. Typically you'll want to make sure both your lists and your input uses the
same casing.

## Ordering

The reference list order is significant. The first of an equal score wins the election. So you'll want to put more 
common, popular, etc. words first in the list. 