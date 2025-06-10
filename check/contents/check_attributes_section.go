package contents

import (
	"fmt"
	"sort"
)

type CheckAttributesSectionOptions struct {
	RequireSchemaOrdering bool
}

func (d *Document) checkAttributesSection() error {
	checkOpts := &CheckAttributesSectionOptions{}

	if d.CheckOptions != nil && d.CheckOptions.AttributesSection != nil {
		checkOpts = d.CheckOptions.AttributesSection
	}

	section := d.Sections.Attributes

	if section == nil {
		return fmt.Errorf("missing attribute section: ## Attribute Reference")
	}

	heading := section.Heading

	if heading.Level != 2 {
		return fmt.Errorf("attribute section heading level (%d) should be: 2", heading.Level)
	}

	headingText := string(heading.Text(d.source))
	expectedHeadingTexts := []string{
		"Attribute Reference",
	}

	if headingText != expectedHeadingTexts[0] {
		return fmt.Errorf("attribute section heading (%s) should be: %q", headingText, expectedHeadingTexts[0])
	}

	paragraphs := section.Paragraphs
	expectedBylineTexts := []string{
		"This resource exports the following attributes in addition to the arguments above:",
		"This ephemeral resource exports the following attributes in addition to the arguments above:",
		"This data source exports the following attributes in addition to the arguments above:",
		"This resource exports no additional attributes.",
		"This ephemeral resource exports no additional attributes.",
		"This data source exports no additional attributes.",
	}

	switch len(paragraphs) {
	case 0:
		return fmt.Errorf("attribute section byline should be: %q, %q, %q, or %q", expectedBylineTexts[0], expectedBylineTexts[1], expectedBylineTexts[2], expectedBylineTexts[3])
	case 1:
		paragraphText := string(paragraphs[0].Text(d.source))

		found := false

		for _, v := range expectedBylineTexts {
			if paragraphText == v {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("attribute section byline (%s) should be: %q, %q, %q, or %q", paragraphText, expectedBylineTexts[0], expectedBylineTexts[1], expectedBylineTexts[2], expectedBylineTexts[3])
		}
	}

	if checkOpts.RequireSchemaOrdering {
		for _, list := range section.SchemaAttributeLists {
			if !sort.IsSorted(SchemaAttributeListItemByName(list.Items)) {
				return fmt.Errorf("attribute section is not sorted by name")
			}
		}
	}

	return nil
}
