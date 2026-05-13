package grouper_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/grouper"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "DB_HOST", Status: diff.StatusConflict},
		{Key: "DB_PORT", Status: diff.StatusMissing},
		{Key: "AWS_REGION", Status: diff.StatusMatch},
		{Key: "AWS_SECRET", Status: diff.StatusConflict},
		{Key: "APP_NAME", Status: diff.StatusMatch},
		{Key: "UNRELATED", Status: diff.StatusMissing},
	}
}

func TestByPrefix_GroupsCorrectly(t *testing.T) {
	results := sampleResults()
	opts := grouper.DefaultOptions()
	groups := grouper.ByPrefix(results, []string{"DB_", "AWS_"}, opts)

	if len(groups) != 3 {
		t.Fatalf("expected 3 groups (DB_, AWS_, OTHER), got %d", len(groups))
	}
	if groups[0].Prefix != "DB_" || len(groups[0].Results) != 2 {
		t.Errorf("DB_ group wrong: %+v", groups[0])
	}
	if groups[1].Prefix != "AWS_" || len(groups[1].Results) != 2 {
		t.Errorf("AWS_ group wrong: %+v", groups[1])
	}
	// OTHER should contain APP_NAME and UNRELATED
	if len(groups[2].Results) != 2 {
		t.Errorf("OTHER group should have 2 results, got %d", len(groups[2].Results))
	}
}

func TestByPrefix_ResultsSortedByKey(t *testing.T) {
	results := []diff.Result{
		{Key: "DB_PORT", Status: diff.StatusMissing},
		{Key: "DB_HOST", Status: diff.StatusConflict},
	}
	groups := grouper.ByPrefix(results, []string{"DB_"}, grouper.DefaultOptions())
	if len(groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(groups))
	}
	if groups[0].Results[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST first, got %s", groups[0].Results[0].Key)
	}
}

func TestByPrefix_MinGroupSizeFilters(t *testing.T) {
	results := sampleResults()
	opts := grouper.DefaultOptions()
	opts.MinGroupSize = 3
	groups := grouper.ByPrefix(results, []string{"DB_", "AWS_"}, opts)
	// DB_ has 2, AWS_ has 2, OTHER has 2 — all below threshold
	if len(groups) != 0 {
		t.Errorf("expected 0 groups with MinGroupSize=3, got %d", len(groups))
	}
}

func TestByPrefix_EmptyPrefixList(t *testing.T) {
	results := sampleResults()
	groups := grouper.ByPrefix(results, []string{}, grouper.DefaultOptions())
	if len(groups) != 1 {
		t.Fatalf("expected 1 group (OTHER only), got %d", len(groups))
	}
	if len(groups[0].Results) != len(results) {
		t.Errorf("OTHER should contain all %d results, got %d", len(results), len(groups[0].Results))
	}
}

func TestByPrefix_EmptyResults(t *testing.T) {
	groups := grouper.ByPrefix([]diff.Result{}, []string{"DB_"}, grouper.DefaultOptions())
	if len(groups) != 0 {
		t.Errorf("expected 0 groups for empty input, got %d", len(groups))
	}
}

func TestByPrefix_CustomOtherLabel(t *testing.T) {
	results := []diff.Result{
		{Key: "MISC_KEY", Status: diff.StatusMatch},
	}
	opts := grouper.DefaultOptions()
	opts.OtherLabel = "UNGROUPED"
	groups := grouper.ByPrefix(results, []string{"DB_"}, opts)
	if len(groups) != 1 || groups[0].Prefix != "UNGROUPED" {
		t.Errorf("expected UNGROUPED label, got %+v", groups)
	}
}
