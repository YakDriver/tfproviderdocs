package contents

import (
	"fmt"
	"strings"

	"github.com/YakDriver/tfproviderdocs/markdown"
)

func (d *Document) checkImportSection() error {
	section := d.Sections.Import

	if section == nil {
		return nil
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

	for _, fencedCodeBlock := range section.FencedCodeBlocks {
		text := markdown.FencedCodeBlockText(fencedCodeBlock, d.source)

		if !strings.Contains(text, d.ResourceName) {
			return fmt.Errorf("import section code block text should contain resource name: %s", d.ResourceName)
		}
	}

	return nil
}
