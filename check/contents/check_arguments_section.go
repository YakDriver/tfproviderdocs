package contents

import (
	"fmt"
	"slices"
	"sort"
)

type CheckArgumentsSectionOptions struct {
	EnhancedRegionChecks  bool
	RegionAware           bool // The resource is Region-aware and has a top-level region argument.
	RequireSchemaOrdering bool
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
	expectedHeadingText := "Argument Reference"

	if headingText != expectedHeadingText {
		return fmt.Errorf("arguments section heading (%s) should be: %s", headingText, expectedHeadingText)
	}

	paragraphs := section.Paragraphs
	expectedBylineTexts := []string{
		"This resource supports the following arguments:",
		"The following arguments are required:",
		"The following arguments are optional:",
		"This resource does not support any arguments.",
		"This data source does not support any arguments.",
		"This data source supports the following arguments:",
	}

	switch len(paragraphs) {
	case 0:
		return fmt.Errorf("argument section byline should be: %q, %q, %q, or %q", expectedBylineTexts[0], expectedBylineTexts[1], expectedBylineTexts[2], expectedBylineTexts[3])
	default:
		paragraphText := string(paragraphs[0].Text(d.source))

		found := false

		for _, v := range expectedBylineTexts {
			if paragraphText == v {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("argument section byline (%s) should be: %q, %q, %q, or %q", paragraphText, expectedBylineTexts[0], expectedBylineTexts[1], expectedBylineTexts[2], expectedBylineTexts[3])
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
		for _, list := range section.SchemaAttributeLists {
			if !slices.ContainsFunc(list.Items, func(item *SchemaAttributeListItem) bool {
				return item.Name == "region" && item.Optional
			}) {
				return fmt.Errorf("arguments section does not contain an Optional region argument")
			}
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
