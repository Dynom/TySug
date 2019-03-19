# TySug

[![CircleCI](https://circleci.com/gh/Dynom/TySug.svg?style=svg)](https://circleci.com/gh/Dynom/TySug)
[![Go Report Card](https://goreportcard.com/badge/github.com/Dynom/TySug)](https://goreportcard.com/report/github.com/Dynom/TySug)
[![GoDoc](https://godoc.org/github.com/Dynom/TySug?status.svg)](https://godoc.org/github.com/Dynom/TySug)
[![codecov](https://codecov.io/gh/Dynom/TySug/branch/master/graph/badge.svg)](https://codecov.io/gh/Dynom/TySug)
[![Microbadger](https://images.microbadger.com/badges/image/dynom/tysug.svg)](https://microbadger.com/images/dynom/tysug)

TySug is collection of packages, together they form a keyboard layout aware alternative word suggester. It can be used as both a library and a webservice.

![shcool](https://raw.githubusercontent.com/Dynom/TySug/master/docs/shcool.jpg)

The primary supported use-case is to help with spelling mistakes against short popular word lists (e.g. domain names). 
Which is useful in helping to prevent typos in e.g. e-mail addresses, detect spam, phishing ([Typosquatting](https://en.m.wikipedia.org/wiki/Typosquatting)), etc. 

The goal is to provide an extensible library that helps with finding possible spelling errors. You can use it 
out-of-the-box as a library, a webservice or as a set of packages to build your own application.

Currently it's a fairly naive approach and not (yet) backed by ML.


# Using TySug

You can use TySug as stand-alone webservice to match against a known-list. If you have Docker you'll have it up and running in a few minutes. 

## TL;DR

If you have Docker, just run:

```bash
docker run --rm -it dynom/tysug:latest
```

_If you don't have Docker, you can download the binary from the [releases](https://github.com/Dynom/TySug/releases) page._

In a different terminal, run:

```bash
curl -s "http://127.0.0.1:1337/list/domains" --data-binary '{"input": "gmail.co"}'
```

## As Webservice

_`curl -s "http://host:port/list/domains" --data-binary '{"input": "gmail.co"}' | jq .`_
```json
{
  "result": "gmail.com",
  "score": 0.9777777777777777,
  "exact_match": false
}
```

- The webservice uses [Jaro-Winkler](https://en.wikipedia.org/wiki/Jaro%E2%80%93Winkler_distance) to calculate similarity.
- The example uses [jq](https://stedolan.github.io/jq/), just omit it if you don't have it installed.


### The path /list/< name >

The name corresponds with a list definition in the [config.toml](https://github.com/Dynom/TySug/blob/master/config.toml). Using this approach the service can be used for various 
types of data. This is both for efficiency (shorter lists to iterate over) and to be more opinionated. when no list by 
that name is found, a 404 is returned.


## As a library
TySug is a collection of stand-alone packages. In each library you can find a README covering the details.
```go
import "github.com/Dynom/TySug/finder"
```
```go
referenceList := []string{"example", "amplifier", "ample"}
ts := finder.New(referenceList, finder.WithAlgorithm(myAlgorithm))

alt, score, exact := ts.Find("exampel")
// alt   = example
// score = 0.9714285714285714
// exact = false (not an exact match in our reference list)
```

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

### Dealing with confidence

When adding your own algorithm, you'll need to handle the "confidence" element yourself. By default TySug's finder will 
handle it just fine, but depending on the scale the algorithm uses you'll need to either normalize the scale or deal 
with the score. 

_Note: Be careful not to introduce bias when converting scale._
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

sug := finder.New([]list, finder.WithAlgorithm(someAlgorithm))
bestMatch, score := sug.Find(input)
// Here score might be 0.8 for a string of length 10, with 2 mutations
```

Without converting the scale, you'll have no bias, however you need to deal with a range where closer to 0 means less changes:
```go
// typically produces a range from (-1 * maximumInputLength) to 0
return -1 * score
```
# Details

## Reference lists

The reference list is a list with known/approved words. TySug's webservice is not optimised to deal with large lists, 
instead it aims for "opinionated" lists. This way you can have a list of domain names or country names. This keeps the 
service snappy and less prone to false-positives.

Large is relative. The size is strongly related to the processing time, longer lists take more time 
[O(N)](http://bigocheatsheet.com/). Test and keep the list within your response-time limits :-). 

### Case-sensitivity

TySug does not normalise words. This means that words are treated in a case-sensitive matter. This is done mostly to
avoid doing unnecessary work in the hot-path. Typically you'll want to make sure both your lists and your input uses the
same casing.

### Ordering

The reference list order is significant. The first of an equal score wins the election. So you'll want to put more 
common, popular, etc. words first in the list. 

## Keyboard layout awareness

Tysug's webservice is keyboard layout aware. This means that when the input is 'bee5' and the reference list contains the 
words 'beer' and 'beek', the word 'beer' is favoured on a Query-US keyboard.

This happens because of a two-pass approach. In the first pass a list of words is collected with 1 or more words with the
same score. If more than 1 word is found with the same score, the keyboard algorithm is applied. Most string-distance
algorithms factor in the "cost" of reaching equality. The amount of "cost" it takes with one letter difference, in the 
same location within a word (E.g.: bee5 versus beer or beek) is typically the same. Making in the assumption that a 
word is typed by a human on a keyboard and that fingers need to travel a distance to reach certain buttons. Factoring in
this assumption could produce better results in the right context.

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
    alternative, score, exact := sug.Find(hostname)

    if exact || score > 0.9 {
        combined := localPart + "@" + alternative
        return combined, score
    }

    return email, score
}
```

# Typos
Dealing with typos is complicated and heavily context dependent.

- Atomic typos -- Typing a (contextual) incorrect, but correctly spelled word (e.g.: _beer_ where you meant: _beet_).
- Intentional typos -- Typing "[teh](https://en.m.wikipedia.org/wiki/Teh)" instead of "the".
- Marking Typos -- Intentional "typos" (e.g.: Bee5^Hr -> _Beer_ or "World Wide Mess^WWeb" -> _World Wide Web_.) 

# Resources

- [https://www.digitalcoding.com/tools/typo-generator.html](https://www.digitalcoding.com/tools/typo-generator.html)
- [http://aspell.net](http://aspell.net) 

# Further reading

- How Difficult is it to Develop a Perfect Spell-checker? A Cross-linguistic Analysis through Complex Network Approach - [http://citeseerx.ist.psu.edu/viewdoc/summary?doi=10.1.1.146.4390](http://citeseerx.ist.psu.edu/viewdoc/summary?doi=10.1.1.146.4390)
- Typographical and Orthographical Spelling Error Correction - [http://citeseerx.ist.psu.edu/viewdoc/summary?doi=10.1.1.70.9592](http://citeseerx.ist.psu.edu/viewdoc/summary?doi=10.1.1.70.9592)
- How to Write a Spelling Corrector - [https://norvig.com/spell-correct.html](https://norvig.com/spell-correct.html)
- Using the Web for Language Independent Spellchecking and Autocorrection - [http://static.googleusercontent.com/media/research.google.com/en/us/pubs/archive/36180.pdf](http://static.googleusercontent.com/media/research.google.com/en/us/pubs/archive/36180.pdf)
- Spellchecking by computer - [https://www.dcs.bbk.ac.uk/..roger/spellchecking.html](https://www.dcs.bbk.ac.uk/~roger/spellchecking.html)

# Wishlist
- Support for [Marking Typos](https://en.wikipedia.org/wiki/Typographical_error#Marking_typos). -- Probably not particularly useful, but seems fun to implement.