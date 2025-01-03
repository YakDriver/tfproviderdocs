package check

import (
	"testing"
)

func TestLegacyEphemeralFileCheck(t *testing.T) {
	testCases := []struct {
		Name            string
		BasePath        string
		Path            string
		ExampleLanguage string
		Options         *LegacyEphemeralFileOptions
		ExpectError     bool
	}{
		{
			Name:            "valid",
			BasePath:        "testdata/valid-legacy-files",
			Path:            "ephemeral.html.markdown",
			ExampleLanguage: "terraform",
		},
		{
			Name:            "invalid extension",
			BasePath:        "testdata/invalid-legacy-files",
			Path:            "ephemeral_invalid_extension.txt",
			ExampleLanguage: "terraform",
			ExpectError:     true,
		},
		{
			Name:            "invalid frontmatter",
			BasePath:        "testdata/invalid-legacy-files",
			Path:            "ephemeral_invalid_frontmatter.html.markdown",
			ExampleLanguage: "terraform",
			ExpectError:     true,
		},
		{
			Name:            "invalid frontmatter with sidebar_current",
			BasePath:        "testdata/invalid-legacy-files",
			Path:            "ephemeral_with_sidebar_current.html.markdown",
			ExampleLanguage: "terraform",
			ExpectError:     true,
		},
		{
			Name:            "invalid frontmatter without layout",
			BasePath:        "testdata/invalid-legacy-files",
			Path:            "ephemeral_without_layout.html.markdown",
			ExampleLanguage: "terraform",
			ExpectError:     true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			if testCase.Options == nil {
				testCase.Options = &LegacyEphemeralFileOptions{}
			}

			if testCase.Options.FileOptions == nil {
				testCase.Options.FileOptions = &FileOptions{
					BasePath: testCase.BasePath,
				}
			}

			got := NewLegacyEphemeralFileCheck(testCase.Options).Run(testCase.Path, testCase.ExampleLanguage)

			if got == nil && testCase.ExpectError {
				t.Errorf("expected error, got no error")
			}

			if got != nil && !testCase.ExpectError {
				t.Errorf("expected no error, got error: %s", got)
			}
		})
	}
}
