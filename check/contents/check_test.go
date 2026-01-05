// Copyright (c) HashiCorp, Inc. 2019-2026
// SPDX-License-Identifier: MPL-2.0

package contents

import (
	"testing"
)

func TestCheck(t *testing.T) {
	testCases := []struct {
		Name         string
		Path         string
		ProviderName string
		CheckOptions *CheckOptions
		ExpectError  bool
	}{
		{
			Name:         "passing",
			Path:         "testdata/full.md",
			ProviderName: "test",
		},
		{
			Name:         "action disallow sections",
			Path:         "testdata/action.md",
			ProviderName: "test",
			CheckOptions: &CheckOptions{
				DisallowAttributesSection:          true,
				AttributesSectionDisallowedMessage: "actions documentation cannot include an attributes section",
				DisallowImportSection:              true,
				ImportSectionDisallowedMessage:     "actions documentation cannot include an import section",
			},
		},
		{
			Name:         "disallow attributes error",
			Path:         "testdata/full.md",
			ProviderName: "test",
			CheckOptions: &CheckOptions{
				DisallowAttributesSection: true,
			},
			ExpectError: true,
		},
		{
			Name:         "disallow import error",
			Path:         "testdata/full.md",
			ProviderName: "test",
			CheckOptions: &CheckOptions{
				DisallowImportSection: true,
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

			got := doc.Check(testCase.CheckOptions)

			if got == nil && testCase.ExpectError {
				t.Errorf("expected error, got no error")
			}

			if got != nil && !testCase.ExpectError {
				t.Errorf("expected no error, got error: %s", got)
			}
		})
	}
}
