# TySug
[![CircleCI](https://circleci.com/gh/Dynom/TySug.svg?style=svg)](https://circleci.com/gh/Dynom/TySug)
[![Go Report Card](https://goreportcard.com/badge/github.com/Dynom/TySug)](https://goreportcard.com/report/github.com/Dynom/TySug)
[![GoDoc](https://godoc.org/github.com/Dynom/TySug?status.svg)](https://godoc.org/github.com/Dynom/TySug)
[![codecov](https://codecov.io/gh/Dynom/TySug/branch/master/graph/badge.svg)](https://codecov.io/gh/Dynom/TySug)

TySug is both a library and a webservice for suggesting alternatives.

As Webservice
_`curl -s "http://host:port" --data-binary '{"input": "gmail.co"}' | jq .`_
```json
{
  "result": "gmail.com",
  "score": 0.9777777777777777
}
```

or as a library
```go
referenceList := []string{"example", "amplifier", "ample"}
ts := finder.New(referenceList, finder.OptSetAlgorithm(myAlgorithm))

alt, score := ts.Find("exampel")
// alt   = example
// score = 0.9714285714285714 
```

The goal is to provide an extensible application that helps with finding possible spelling errors. You can use it 
out-of-the-box as a library, a webservice or as a set of packages to build your own application.

By default it uses [Jaro-Winkler](https://en.wikipedia.org/wiki/Jaro%E2%80%93Winkler_distance) to calculate similarity.

## As a webservice
@todo

## As a library
You can use the various components that make up TySug individually or as a whole.

### Example


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
 - https://www.joyofdata.de/blog/comparison-of-string-distance-algorithms/

### Dealing with confidence
When adding your own algorithm, you'll need to handle the "confidence" element yourself. By default TySug's finder will 
handle it just fine, but depending on the scale the algorithm uses you'll need to either normalize the scale or deal 
with the score. 

_Note: Be careful not to introduce bias when converting scale_
```go
var someAlgorithm finder.AlgWrapper = func(a, b string) float64 {

    // Result is, in this instance, the amount of steps taken to achieve equality
    // Algorithms like Jaro produce a value between 0.0 and 1.0
    score := someAlgorithm.CalculateDistance(a, b)
    
    // Finding the longest string
    var ml int
    if len(a) >= len(b) {
        ml = len(a)
    } else {
        ml = len(b)
    }
    
    // This introduces a bias. Inputs of longer lengths get a slight favour over shorter ones, causing deletions to weigh less.
    return 1 - (score / float64(ml))
}

sug := finder.New([]list, finder.OptSetAlgorithm(someAlgorithm))
bestMatch, score := sug.Find(input)
// Here score might be 0.8 for a string of length 10, with 2 mutations
```

Without converting the scale, you'll have no bias, however you need to deal with a range where closer to 0 means less changes:
```go
// typically produces a range from (-1 * maximumInputLength) to 0
return -1 * score
```

# Examples
## Finding common e-mail domain typos
To help people avoid submitting an incorrect e-mail address, one could try the following:

```go
func SuggestAlternative(email string, domains []string) (string, float64) {

    i := strings.LastIndex(email, "@")
    if i <= 0 || i >= len(email) {
        return email, 0
    }
    
    // Extracting the local and domain parts
    localPart := email[:i]
    hostname := email[i+1:]
    
    sug, _ := finder.New(domains)
    alternative, score := sug.Find(hostname)
    
    if score > 0.9 {
        combined := localPart + "@" + alternative
        return combined, score
    }
    
    return email, score
}
```

Do note that:
 - The comparisons are done in a case-sensitive manner, so you probably want to normalize the input and the
   reference list.
 - The reference list order is significant, the first of an equal score wins the election. Put more common words first.
 - Score is very dependant on the algorithm used, you'll want to tweak it for your use-case.
 
 
# Wish list

- Become keyboard aware -- _The keyboard layout could help with identifying more "logical" typing mistakes 
  "beer" versus "beek" or "bee5". They might result with the same score, but "beer" might be more suitable with "bee5" 
  as input_.
