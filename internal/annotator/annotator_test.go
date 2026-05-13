package annotator_test

import (
	"testing"

	"github.com/user/envdiff/internal/annotator"
	"github.com/user/envdiff/internal/diff"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "DB_PASSWORD", Status: diff.StatusMissing},
		{Key: "APP_NAME", Status: diff.StatusConflict, Values: []string{"app", "myapp"}},
		{Key: "API_TOKEN", Status: diff.StatusConflict, Values: []string{"tok1", "tok2"}},
		{Key: "PORT", Status: diff.StatusMatch},
	}
}

func TestAnnotate_LengthMatchesInput(t *testing.T) {
	rs := sampleResults()
	anns := annotator.Annotate(rs, annotator.DefaultOptions())
	if len(anns) != len(rs) {
		t.Fatalf("expected %d annotations, got %d", len(rs), len(anns))
	}
}

func TestAnnotate_SensitiveDetected(t *testing.T) {
	rs := []diff.Result{{Key: "DB_PASSWORD", Status: diff.StatusMissing}}
	anns := annotator.Annotate(rs, annotator.DefaultOptions())
	if !anns[0].IsSensitive {
		t.Error("expected DB_PASSWORD to be marked sensitive")
	}
}

func TestAnnotate_NonSensitive(t *testing.T) {
	rs := []diff.Result{{Key: "APP_NAME", Status: diff.StatusConflict}}
	anns := annotator.Annotate(rs, annotator.DefaultOptions())
	if anns[0].IsSensitive {
		t.Error("expected APP_NAME to not be sensitive")
	}
}

func TestAnnotate_DefaultValueApplied(t *testing.T) {
	rs := []diff.Result{{Key: "PORT", Status: diff.StatusMissing}}
	opts := annotator.Options{
		Defaults:          map[string]string{"PORT": "8080"},
		SensitivePatterns: []string{},
	}
	anns := annotator.Annotate(rs, opts)
	if !anns[0].HasDefault {
		t.Error("expected HasDefault=true for PORT")
	}
	if anns[0].DefaultVal != "8080" {
		t.Errorf("expected default '8080', got %q", anns[0].DefaultVal)
	}
}

func TestAnnotate_HintMissingSensitive(t *testing.T) {
	rs := []diff.Result{{Key: "API_TOKEN", Status: diff.StatusMissing}}
	anns := annotator.Annotate(rs, annotator.DefaultOptions())
	if anns[0].Hint == "" {
		t.Error("expected non-empty hint for missing sensitive key")
	}
}

func TestAnnotate_HintMatchIsEmpty(t *testing.T) {
	rs := []diff.Result{{Key: "PORT", Status: diff.StatusMatch}}
	anns := annotator.Annotate(rs, annotator.DefaultOptions())
	if anns[0].Hint != "" {
		t.Errorf("expected empty hint for matching key, got %q", anns[0].Hint)
	}
}

func TestAnnotate_DefaultOptions_HasPatterns(t *testing.T) {
	opts := annotator.DefaultOptions()
	if len(opts.SensitivePatterns) == 0 {
		t.Error("DefaultOptions should provide sensitive patterns")
	}
}
