// Copyright (c) IBM Corp. 2019-2026
// SPDX-License-Identifier: MPL-2.0

// Copyright (c) HashiCorp, Inc. 2019-2026
package contents

import (
	"fmt"
	"slices"
	"strings"
)

type CheckSignatureSectionOptions struct {
	RequireSection      SectionRequirement
	AllowedHeadingTexts []string
	RequireCodeBlock    bool
}

func (d *Document) checkSignatureSection() error {
	opts := d.CheckOptions.SignatureSection
	if opts == nil {
		return nil
	}

	section := d.Sections.Signature

	if section == nil {
		if opts.RequireSection == Required {
			return fmt.Errorf("missing signature section: ## Signature")
		}
		return nil
	}

	if opts.RequireSection == Forbidden {
		return fmt.Errorf("signature section should not be present")
	}

	heading := section.Heading
	if heading.Level != 2 {
		return fmt.Errorf("signature section heading level (%d) should be: 2", heading.Level)
	}

	headingText := string(heading.Text(d.source))
	allowedHeadingTexts := []string{"Signature"}

	if len(opts.AllowedHeadingTexts) > 0 {
		allowedHeadingTexts = opts.AllowedHeadingTexts
	}

	foundHeading := slices.Contains(allowedHeadingTexts, headingText)

	if !foundHeading {
		formatted := make([]string, len(allowedHeadingTexts))
		for i, v := range allowedHeadingTexts {
			formatted[i] = fmt.Sprintf("%q", v)
		}
		return fmt.Errorf("signature section heading (%s) should be one of: %s", headingText, strings.Join(formatted, ", "))
	}

	if opts.RequireCodeBlock && len(section.FencedCodeBlocks) == 0 {
		return fmt.Errorf("signature section must include a code block")
	}

	return nil
}
