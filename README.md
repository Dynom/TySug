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