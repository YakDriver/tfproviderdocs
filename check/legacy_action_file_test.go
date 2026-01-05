// Copyright (c) HashiCorp, Inc. 2019-2026
// SPDX-License-Identifier: MPL-2.0

package check

import "testing"

func TestLegacyActionFileCheck(t *testing.T) {
	testCases := []struct {
		Name            string
		BasePath        string
		Path            string
		ExampleLanguage string
		Options         *LegacyActionFileOptions
		ExpectError     bool
	}{
		{
			Name:            "valid",
			BasePath:        "testdata/valid-legacy-files",
			Path:            "action.html.markdown",
			ExampleLanguage: "terraform",
		},
		{
			Name:            "invalid extension",
			BasePath:        "testdata/invalid-legacy-files",
			Path:            "action_invalid_extension.txt",
			ExampleLanguage: "terraform",
			ExpectError:     true,
		},
		{
			Name:            "invalid frontmatter",
			BasePath:        "testdata/invalid-legacy-files",
			Path:            "action_invalid_frontmatter.html.markdown",
			ExampleLanguage: "terraform",
			ExpectError:     true,
		},
		{
			Name:            "with sidebar current",
			BasePath:        "testdata/invalid-legacy-files",
			Path:            "action_with_sidebar_current.html.markdown",
			ExampleLanguage: "terraform",
			ExpectError:     true,
		},
		{
			Name:            "without layout",
			BasePath:        "testdata/invalid-legacy-files",
			Path:            "action_without_layout.html.markdown",
			ExampleLanguage: "terraform",
			ExpectError:     true,
		},
		{
			Name:            "with attributes",
			BasePath:        "testdata/invalid-legacy-files",
			Path:            "action_with_attributes.html.markdown",
			ExampleLanguage: "terraform",
			ExpectError:     true,
		},
		{
			Name:            "with import",
			BasePath:        "testdata/invalid-legacy-files",
			Path:            "action_with_import.html.markdown",
			ExampleLanguage: "terraform",
			ExpectError:     true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			if testCase.Options == nil {
				testCase.Options = &LegacyActionFileOptions{}
			}

			if testCase.Options.FileOptions == nil {
				testCase.Options.FileOptions = &FileOptions{
					BasePath: testCase.BasePath,
				}
			}

			got := NewLegacyActionFileCheck(testCase.Options).Run(testCase.Path, testCase.ExampleLanguage)

			if got == nil && testCase.ExpectError {
				t.Errorf("expected error, got no error")
			}

			if got != nil && !testCase.ExpectError {
				t.Errorf("expected no error, got error: %s", got)
			}
		})
	}
}
