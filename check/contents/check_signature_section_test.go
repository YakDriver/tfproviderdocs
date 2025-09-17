package contents

import "testing"

func TestCheckSignatureSection(t *testing.T) {
	testCases := []struct {
		Name        string
		Path        string
		ExpectError bool
	}{
		{
			Name: "passing",
			Path: "testdata/signature/passing.md",
		},
		{
			Name:        "missing section",
			Path:        "testdata/signature/missing.md",
			ExpectError: true,
		},
		{
			Name:        "wrong heading",
			Path:        "testdata/signature/wrong_heading.md",
			ExpectError: true,
		},
		{
			Name:        "missing code block",
			Path:        "testdata/signature/missing_code_block.md",
			ExpectError: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			doc := NewDocument(testCase.Path, "test")

			if err := doc.Parse(); err != nil {
				t.Fatalf("unexpected parse error: %s", err)
			}

			doc.CheckOptions = &CheckOptions{
				TitleSection: &CheckTitleSectionOptions{AllowedPrefixes: []string{"Function"}},
				ExamplesSection: &CheckExamplesSectionOptions{
					ExpectedCodeBlockLanguage: "terraform",
				},
				ArgumentsSection: &CheckArgumentsSectionOptions{
					AllowedHeadingTexts: []string{"Arguments"},
					AllowMissingByline:  true,
				},
				SignatureSection: &CheckSignatureSectionOptions{
					RequireSection:      Required,
					AllowedHeadingTexts: []string{"Signature"},
					RequireCodeBlock:    true,
				},
			}

			err := doc.checkSignatureSection()

			if testCase.ExpectError {
				if err == nil {
					t.Fatalf("expected error, got none")
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %s", err)
				}
			}
		})
	}
}
