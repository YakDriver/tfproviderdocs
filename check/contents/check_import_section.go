// Copyright (c) HashiCorp, Inc. 2019-2026
// SPDX-License-Identifier: MPL-2.0

package contents

import (
	"fmt"
	"strings"

	"github.com/YakDriver/tfproviderdocs/markdown"
)

type CheckImportSectionOptions struct {
	RequireSection SectionRequirement
}

func (d *Document) checkImportSection() error {
	checkOpts := &CheckImportSectionOptions{}

	if d.CheckOptions != nil && d.CheckOptions.ImportSection != nil {
		checkOpts = d.CheckOptions.ImportSection
	}

	section := d.Sections.Import

	if section == nil {
		if checkOpts.RequireSection == Required {
			return fmt.Errorf("missing import section: ## Import")
		}

		return nil
	} else {
		if checkOpts.RequireSection == Forbidden {
			return fmt.Errorf("import section should not be present")
		}
	}

	heading := section.Heading

	if heading.Level != 2 {
		return fmt.Errorf("import section heading level (%d) should be: 2", heading.Level)
	}

	headingText := string(heading.Text(d.source))
	expectedHeadingText := "Import"

	if headingText != expectedHeadingText {
		return fmt.Errorf("import section heading (%s) should be: %s", headingText, expectedHeadingText)
	}

	paragraphs := section.Paragraphs
	problems := [][]string{
		{
			"can be imported", // problem
			"use active voice instead: Import X using A, B, C.", // fix message
		},
		{
			"e.g",                          // problem
			"instead use \"For example:\"", // fix message
		},
		{
			"E.g",                          // problem
			"instead use \"For example:\"", // fix message
		},
	}

	if len(paragraphs) > 0 {
		text := string(paragraphs[0].Text(d.source))

		for _, v := range problems {
			problem := v[0]
			msg := v[1]
			if strings.Contains(text, problem) {
				return fmt.Errorf("import section should not include %q, %s", problem, msg)
			}
		}

		suffix := ". For example:"
		if !strings.HasSuffix(text, suffix) && !strings.Contains(text, "cannot import") {
			return fmt.Errorf("import section should conclude with %q (or state \"You cannot import ...\")", suffix)
		}
	}

	if len(paragraphs) > 0 && !strings.Contains(string(paragraphs[0].Text(d.source)), "cannot import") && len(section.FencedCodeBlocks) < 1 {
		return fmt.Errorf("import section should have a code block (or state \"You cannot import ...\")")
	}

	hitConsole := false
	for i, fencedCodeBlock := range section.FencedCodeBlocks {
		text := markdown.FencedCodeBlockText(fencedCodeBlock, d.source)

		if !strings.Contains(text, d.ResourceName) {
			return fmt.Errorf("import section code block text should contain resource name: %s", d.ResourceName)
		}

		if i == 0 && (!strings.Contains(markdown.FencedCodeBlockLanguage(fencedCodeBlock, d.source), "terraform") || !strings.HasPrefix(text, "import {")) {
			return fmt.Errorf("the first import section code block should have an import block using type 'terraform' (i.e., ```terraform\nimport {)")
		}

		if strings.Contains(markdown.FencedCodeBlockLanguage(fencedCodeBlock, d.source), "console") && !strings.HasPrefix(text, "% ") {
			return fmt.Errorf("import section code block type 'console' should begin with '%% '")
		}

		if !strings.Contains(markdown.FencedCodeBlockLanguage(fencedCodeBlock, d.source), "console") && !strings.Contains(markdown.FencedCodeBlockLanguage(fencedCodeBlock, d.source), "terraform") {
			return fmt.Errorf("import section code block type should be 'console' or 'terraform' (i.e., ```console or ```terraform)")
		}

		if strings.Contains(markdown.FencedCodeBlockLanguage(fencedCodeBlock, d.source), "console") {
			hitConsole = true
		}

		if hitConsole && strings.Contains(markdown.FencedCodeBlockLanguage(fencedCodeBlock, d.source), "terraform") && strings.HasPrefix(text, "import ") {
			return fmt.Errorf("import section: all code blocks of type 'terraform' should be before code blocks of type 'console'")
		}
	}

	return nil
}
