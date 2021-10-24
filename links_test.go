package main

import "testing"

func TestGetIssueOrPr(t *testing.T) {
	tests := []struct {
		r *linkParseResult
		o normalizedIssuePr
	}{
		{&linkParseResult{
			issueNumber:   0,
			prNumber:      1,
			commentNumber: 0,
			links:         make(map[string]string),
		},
			normalizedIssuePr{1, "prs"},
		},
		{&linkParseResult{
			issueNumber:   1,
			prNumber:      0,
			commentNumber: 0,
			links:         make(map[string]string),
		},
			normalizedIssuePr{1, "issues"},
		},
	}

	for _, test := range tests {
		if g := getIssueOrPr(test.r); g != test.o {
			t.Errorf("Expected %v but got %v", test.o, g)
		}
	}
}
