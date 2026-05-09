package diff

// Status describes the relationship of a key between two env files.
type Status string

const (
	StatusOK       Status = "ok"
	StatusMissing  Status = "missing"
	StatusConflict Status = "conflict"
)

// Result holds the comparison outcome for a single key.
type Result struct {
	Key     string `json:"key"`
	Status  Status `json:"status"`
	ValueA  string `json:"value_a,omitempty"`
	ValueB  string `json:"value_b,omitempty"`
	FileA   string `json:"file_a,omitempty"`
	FileB   string `json:"file_b,omitempty"`
}

// Compare compares two parsed env maps (key→value) and returns a slice of
// Result entries describing every key found in either map.
func Compare(a, b map[string]string, fileA, fileB string) []Result {
	seen := make(map[string]bool)
	var results []Result

	for k, va := range a {
		seen[k] = true
		vb, ok := b[k]
		switch {
		case !ok:
			results = append(results, Result{Key: k, Status: StatusMissing, ValueA: va, FileA: fileA, FileB: fileB})
		case va != vb:
			results = append(results, Result{Key: k, Status: StatusConflict, ValueA: va, ValueB: vb, FileA: fileA, FileB: fileB})
		default:
			results = append(results, Result{Key: k, Status: StatusOK, ValueA: va, ValueB: vb, FileA: fileA, FileB: fileB})
		}
	}

	for k, vb := range b {
		if seen[k] {
			continue
		}
		results = append(results, Result{Key: k, Status: StatusMissing, ValueB: vb, FileA: fileA, FileB: fileB})
	}

	return results
}
