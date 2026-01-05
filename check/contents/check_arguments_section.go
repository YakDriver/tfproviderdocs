// Copyright (c) IBM Corp. 2019-2026
// SPDX-License-Identifier: MPL-2.0

// Copyright (c) HashiCorp, Inc. 2019-2026
package contents

import (
	"fmt"
	"slices"
	"sort"
	"strings"
)

type CheckArgumentsSectionOptions struct {
	EnhancedRegionChecks  bool
	RegionAware           bool // The resource is Region-aware and has a top-level region argument.
	RequireSchemaOrdering bool
	ExpectedBylineTexts   []string
	AllowedHeadingTexts   []string
	AllowMissingByline    bool
}

func (d *Document) checkArgumentsSection() error {
	checkOpts := &CheckArgumentsSectionOptions{}

	if d.CheckOptions != nil && d.CheckOptions.ArgumentsSection != nil {
		checkOpts = d.CheckOptions.ArgumentsSection
	}

	section := d.Sections.Arguments

	if section == nil {
		return fmt.Errorf("missing arguments section: ## Argument Reference")
	}

	heading := section.Heading

	if heading.Level != 2 {
		return fmt.Errorf("arguments section heading level (%d) should be: 2", heading.Level)
	}

	headingText := string(heading.Text(d.source))
	allowedHeadingTexts := []string{"Argument Reference"}

	if len(checkOpts.AllowedHeadingTexts) > 0 {
		allowedHeadingTexts = checkOpts.AllowedHeadingTexts
	}

	foundHeading := slices.Contains(allowedHeadingTexts, headingText)

	if !foundHeading {
		formatted := make([]string, len(allowedHeadingTexts))
		for i, v := range allowedHeadingTexts {
			formatted[i] = fmt.Sprintf("%q", v)
		}
		return fmt.Errorf("arguments section heading (%s) should be one of: %s", headingText, strings.Join(formatted, ", "))
	}

	paragraphs := section.Paragraphs
	expectedBylineTexts := []string{
		"This resource supports the following arguments:",
		"This ephemeral resource supports the following arguments:",
		"This list resource supports the following arguments:",
		"This action supports the following arguments:",
		"The following arguments are required:",
		"The following arguments are optional:",
		"This resource does not support any arguments.",
		"This ephemeral resource does not support any arguments.",
		"This list resource does not support any arguments.",
		"This action does not support any arguments.",
		"This data source does not support any arguments.",
		"This data source supports the following arguments:",
	}

	if len(checkOpts.ExpectedBylineTexts) > 0 {
		expectedBylineTexts = checkOpts.ExpectedBylineTexts
	}

	allowedTexts := make([]string, len(expectedBylineTexts))
	for i, v := range expectedBylineTexts {
		allowedTexts[i] = fmt.Sprintf("%q", v)
	}
	allowedTextsMessage := strings.Join(allowedTexts, ", ")

	switch len(paragraphs) {
	case 0:
		if !checkOpts.AllowMissingByline {
			return fmt.Errorf("argument section byline should be one of: %s", allowedTextsMessage)
		}
	default:
		if len(expectedBylineTexts) == 0 {
			break
		}

		paragraphText := string(paragraphs[0].Text(d.source))

		found := slices.Contains(expectedBylineTexts, paragraphText)

		if !found {
			return fmt.Errorf("argument section byline (%s) should be one of: %s", paragraphText, allowedTextsMessage)
		}

		if paragraphText == "The following arguments are required:" {
			// Check for Optionals.
			if n := len(section.SchemaAttributeLists); n > 0 {
				if slices.ContainsFunc(section.SchemaAttributeLists[0].Items, func(item *SchemaAttributeListItem) bool {
					return item.Optional
				}) {
					return fmt.Errorf("required arguments section contains an Optional argument")
				}
			}

			if n := len(paragraphs); n > 1 {
				// A following paragraph must be "The following arguments are optional:"
				want := "The following arguments are optional:"
				idx := -1
				for i := 1; i < n; i++ {
					if string(paragraphs[i].Text(d.source)) == want {
						idx = i
						break
					}
				}

				if idx < 0 {
					return fmt.Errorf("argument section byline (%s) should be: %q", paragraphText, want)
				}

				// Check for Required.
				if n := len(section.SchemaAttributeLists); n > idx {
					if slices.ContainsFunc(section.SchemaAttributeLists[idx].Items, func(item *SchemaAttributeListItem) bool {
						return item.Required
					}) {
						return fmt.Errorf("optional arguments section contains a Required argument")
					}
				}
			}
		}
	}

	if checkOpts.EnhancedRegionChecks && checkOpts.RegionAware {
		found := false
		for _, list := range section.SchemaAttributeLists {
			if slices.ContainsFunc(list.Items, func(item *SchemaAttributeListItem) bool {
				return item.Name == "region" && item.Optional
			}) {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("arguments section does not contain an Optional region argument")
		}
	}

	if checkOpts.RequireSchemaOrdering {
		for _, list := range section.SchemaAttributeLists {
			if !sort.IsSorted(SchemaAttributeListItemByName(list.Items)) {
				return fmt.Errorf("arguments section is not sorted by name")
			}
		}
	}

	return nil
}
