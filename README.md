TySug is a library for suggesting alternatives.

Example:
```go
referenceList := []string{"example", "amplifier", "ample"}
ts := tysug.New(referenceList)

alt, score := ts.Find("exampel")
// alt   - example
// score - 0.9714285714285714 
```

If you want to use your own algorithms, reference the finder package directly.

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