package classification

import (
	"testing"

	"github.com/vllm-project/semantic-router/src/semantic-router/pkg/config"
)

func TestKeywordClassifier_Classify_FindsAll(t *testing.T) {
	// Define two rules that should both match the input
	rules := []config.KeywordRule{
		{
			Name:          "rule1",
			Operator:      "OR",
			Keywords:      []string{"foo"},
			CaseSensitive: false,
		},
		{
			Name:          "rule2",
			Operator:      "OR",
			Keywords:      []string{"bar"},
			CaseSensitive: false,
		},
	}

	classifier, err := NewKeywordClassifier(rules)
	if err != nil {
		t.Fatalf("Failed to create classifier: %v", err)
	}

	text := "foo and bar"

	// Current expected behavior (bug): Classify returns only first match
	// Desired behavior: Classify (or new method) returns both "rule1" and "rule2"
	// For now, let's verify what Classify returns.

	// Verify FindAllMatches finds BOTH rules
	matches, err := classifier.FindAllMatches(text)
	if err != nil {
		t.Fatalf("FindAllMatches failed: %v", err)
	}

	if len(matches) != 2 {
		t.Fatalf("Expected 2 matched rules, got %d", len(matches))
	}

	foundRule1 := false
	foundRule2 := false

	for _, m := range matches {
		if m.RuleName == "rule1" {
			foundRule1 = true
		} else if m.RuleName == "rule2" {
			foundRule2 = true
		}
	}

	if !foundRule1 || !foundRule2 {
		t.Errorf("Did not find both rules. Found rule1: %v, found rule2: %v", foundRule1, foundRule2)
	}
}
