// Copyright (c) IBM Corp. 2019-2026
// SPDX-License-Identifier: MPL-2.0

// Copyright (c) HashiCorp, Inc. 2019-2026
package contents

import (
	"fmt"
	"strings"
)

func (d *Document) checkTitleSection() error {
	section := d.Sections.Title

	if section == nil {
		return fmt.Errorf("missing title section: # Resource: %s", d.ResourceName)
	}

	heading := section.Heading

	if heading.Level != 1 {
		return fmt.Errorf("title section heading level (%d) should be: 1", heading.Level)
	}

	headingText := string(heading.Text(d.source))

	validPrefixes := []string{"Action", "Data Source", "Ephemeral", "List Resource", "Resource"}

	if d.CheckOptions != nil && d.CheckOptions.TitleSection != nil && len(d.CheckOptions.TitleSection.AllowedPrefixes) > 0 {
		validPrefixes = d.CheckOptions.TitleSection.AllowedPrefixes
	}
	isValidPrefix := false
	for _, prefix := range validPrefixes {
		if strings.HasPrefix(headingText, fmt.Sprintf("%s: ", prefix)) {
			isValidPrefix = true
			break
		}
	}

	if !isValidPrefix {
		return fmt.Errorf("title section heading (%s) should have one of these prefixes: %v", headingText, validPrefixes)
	}

	if len(section.FencedCodeBlocks) > 0 {
		return fmt.Errorf("title section code examples should be in Example Usage section")
	}

	return nil
}
