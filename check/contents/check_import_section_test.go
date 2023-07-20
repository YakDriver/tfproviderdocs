package contents

import (
	"testing"
)

func TestCheckImportSection(t *testing.T) {
	testCases := []struct {
		Name         string
		Path         string
		ProviderName string
		ExpectError  bool
	}{
		{
			Name:         "passing",
			Path:         "testdata/import/passing.md",
			ProviderName: "test",
		},
		{
			Name:         "passing cannot import",
			Path:         "testdata/import/passing_cannot_import.md",
			ProviderName: "test",
		},
		{
			Name:         "passing code block after",
			Path:         "testdata/import/passing_code_block_after.md",
			ProviderName: "test",
		},
		{
			Name:         "wrong code block resource type",
			Path:         "testdata/import/wrong_code_block_resource_type.md",
			ProviderName: "test",
			ExpectError:  true,
		},
		{
			Name:         "wrong heading level",
			Path:         "testdata/import/wrong_heading_level.md",
			ProviderName: "test",
			ExpectError:  true,
		},
		{
			Name:         "wrong heading text",
			Path:         "testdata/import/wrong_heading_text.md",
			ProviderName: "test",
			ExpectError:  true,
		},
		{
			Name:         "wrong passive voice",
			Path:         "testdata/import/wrong_passive_voice.md",
			ProviderName: "test",
			ExpectError:  true,
		},
		{
			Name:         "wrong eg",
			Path:         "testdata/import/wrong_eg.md",
			ProviderName: "test",
			ExpectError:  true,
		},
		{
			Name:         "wrong without codeblock",
			Path:         "testdata/import/wrong_no_codeblock.md",
			ProviderName: "test",
			ExpectError:  true,
		},
		{
			Name:         "wrong console prefix",
			Path:         "testdata/import/wrong_console_prefix.md",
			ProviderName: "test",
			ExpectError:  true,
		},
		{
			Name:         "wrong code block order",
			Path:         "testdata/import/wrong_code_block_order.md",
			ProviderName: "test",
			ExpectError:  true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			doc := NewDocument(testCase.Path, testCase.ProviderName)

			if err := doc.Parse(); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			got := doc.checkImportSection()

			if got == nil && testCase.ExpectError {
				t.Errorf("expected error, got no error")
			}

			if got != nil && !testCase.ExpectError {
				t.Errorf("expected no error, got error: %s", got)
			}
		})
	}
}
