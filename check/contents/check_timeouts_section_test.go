package contents

import (
	"testing"
)

func TestCheckTimeoutsSection(t *testing.T) {
	testCases := []struct {
		Name         string
		Path         string
		ProviderName string
		CheckOptions *CheckOptions
		ExpectError  bool
	}{
		{
			Name:         "passing",
			Path:         "testdata/timeouts/passing.md",
			ProviderName: "test",
		},
		{
			Name:         "forbidden",
			Path:         "testdata/timeouts/passing.md",
			ProviderName: "test",
			CheckOptions: &CheckOptions{
				TimeoutsSection: &CheckTimeoutsSectionOptions{
					RequireSection: Forbidden,
				},
			},
			ExpectError: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			doc := NewDocument(testCase.Path, testCase.ProviderName)

			if err := doc.Parse(); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			doc.CheckOptions = testCase.CheckOptions

			got := doc.checkTimeoutsSection()

			if got == nil && testCase.ExpectError {
				t.Errorf("expected error, got no error")
			}

			if got != nil && !testCase.ExpectError {
				t.Errorf("expected no error, got error: %s", got)
			}
		})
	}
}
