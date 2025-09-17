package check

import "testing"

func TestRegistryActionFileCheck(t *testing.T) {
	testCases := []struct {
		Name            string
		BasePath        string
		Path            string
		ExampleLanguage string
		Options         *RegistryActionFileOptions
		ExpectError     bool
	}{
		{
			Name:            "valid",
			BasePath:        "testdata/valid-registry-files",
			Path:            "action.md",
			ExampleLanguage: "terraform",
		},
		{
			Name:            "invalid extension",
			BasePath:        "testdata/invalid-registry-files",
			Path:            "action_invalid_extension.markdown",
			ExampleLanguage: "terraform",
			ExpectError:     true,
		},
		{
			Name:            "invalid frontmatter",
			BasePath:        "testdata/invalid-registry-files",
			Path:            "action_invalid_frontmatter.md",
			ExampleLanguage: "terraform",
			ExpectError:     true,
		},
		{
			Name:            "invalid layout",
			BasePath:        "testdata/invalid-registry-files",
			Path:            "action_with_layout.md",
			ExampleLanguage: "terraform",
			ExpectError:     true,
		},
		{
			Name:            "with attributes",
			BasePath:        "testdata/invalid-registry-files",
			Path:            "action_with_attributes.md",
			ExampleLanguage: "terraform",
			ExpectError:     true,
		},
		{
			Name:            "with import",
			BasePath:        "testdata/invalid-registry-files",
			Path:            "action_with_import.md",
			ExampleLanguage: "terraform",
			ExpectError:     true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			if testCase.Options == nil {
				testCase.Options = &RegistryActionFileOptions{}
			}

			if testCase.Options.FileOptions == nil {
				testCase.Options.FileOptions = &FileOptions{
					BasePath: testCase.BasePath,
				}
			}

			got := NewRegistryActionFileCheck(testCase.Options).Run(testCase.Path, testCase.ExampleLanguage)

			if got == nil && testCase.ExpectError {
				t.Errorf("expected error, got no error")
			}

			if got != nil && !testCase.ExpectError {
				t.Errorf("expected no error, got error: %s", got)
			}
		})
	}
}
