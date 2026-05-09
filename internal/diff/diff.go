package diff

// Result holds the comparison result between two env files.
type Result struct {
	// MissingInSecond contains keys present in first but absent in second.
	MissingInSecond []string
	// MissingInFirst contains keys present in second but absent in first.
	MissingInFirst []string
	// Conflicts contains keys present in both files but with different values.
	Conflicts []Conflict
}

// Conflict describes a key whose value differs between two env files.
type Conflict struct {
	Key        string
	FirstValue  string
	SecondValue string
}

// Compare takes two parsed env maps and returns a Result describing
// missing keys and conflicting values between them.
func Compare(first, second map[string]string) Result {
	var result Result

	for key, val1 := range first {
		val2, ok := second[key]
		if !ok {
			result.MissingInSecond = append(result.MissingInSecond, key)
			continue
		}
		if val1 != val2 {
			result.Conflicts = append(result.Conflicts, Conflict{
				Key:         key,
				FirstValue:  val1,
				SecondValue: val2,
			})
		}
	}

	for key := range second {
		if _, ok := first[key]; !ok {
			result.MissingInFirst = append(result.MissingInFirst, key)
		}
	}

	return result
}

// HasDifferences returns true if the Result contains any discrepancies.
func (r Result) HasDifferences() bool {
	return len(r.MissingInFirst) > 0 ||
		len(r.MissingInSecond) > 0 ||
		len(r.Conflicts) > 0
}
