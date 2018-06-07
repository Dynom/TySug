# TySug
TySug is both a library and a webservice for suggesting alternatives.

## As a webservice
@todo

## As a library
You can use the various components that make TySug individually or as a whole.

### Example
```go
// Note: The arguments are case-sensitive. Normalize the data to avoid possible problems 
referenceList := []string{"example", "amplifier", "ample"}
ts := tysug.New(referenceList)

alt, score := ts.Find("exampel")
// alt   = example
// score = 0.9714285714285714 
```
if you want to use a different algorithm, simply wrap your algorithm as an `Option` and pass it as argument to the Finder. You can find your inspiration in unit-tests / examples.

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
 - https://www.joyofdata.de/blog/comparison-of-string-distance-algorithms/


# Examples
## Finding common e-mail domain typos
To help people prevent submitting an incorrect e-mail address, one could try the following

```go
func SuggestAlternative(email string, domains []string) (string, float64) {

	i := strings.LastIndex(email, "@")
	if i <= 0 || i >= len(email) {
		return email, 0
	}

	// Extracting the local and domain parts
	localPart := email[:i]
	hostname := email[i+1:]

	sug, _ := tysug.New(domains)
	alternative, score := sug.Find(hostname)

	if score > 0.9 {
		combined := localPart + "@" + alternative
		return combined, score
	}

	return email, score
}

```