// Copyright (c) IBM Corp. 2019-2026
// SPDX-License-Identifier: MPL-2.0

// Copyright (c) HashiCorp, Inc. 2019-2026
package check

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-multierror"
)

type LegacyEphemeralFileOptions struct {
	*FileOptions

	Contents     *ContentsOptions
	FrontMatter  *FrontMatterOptions
	ProviderName string
}

type LegacyEphemeralFileCheck struct {
	FileCheck

	Options *LegacyEphemeralFileOptions
}

func NewLegacyEphemeralFileCheck(opts *LegacyEphemeralFileOptions) *LegacyEphemeralFileCheck {
	check := &LegacyEphemeralFileCheck{
		Options: opts,
	}

	if check.Options == nil {
		check.Options = &LegacyEphemeralFileOptions{}
	}

	if check.Options.Contents == nil {
		check.Options.Contents = &ContentsOptions{}
	}

	if check.Options.Contents.ProviderName == "" {
		check.Options.Contents.ProviderName = check.Options.ProviderName
	}

	if check.Options.FileOptions == nil {
		check.Options.FileOptions = &FileOptions{}
	}

	if check.Options.FrontMatter == nil {
		check.Options.FrontMatter = &FrontMatterOptions{}
	}

	check.Options.FrontMatter.NoSidebarCurrent = true
	check.Options.FrontMatter.RequireDescription = true
	check.Options.FrontMatter.RequireLayout = true
	check.Options.FrontMatter.RequirePageTitle = true

	return check
}

func (check *LegacyEphemeralFileCheck) Run(path string, exampleLanguage string) error {
	fullpath := check.Options.FullPath(path)

	log.Printf("[DEBUG] Checking file: %s", fullpath)

	// skip cdktf directories
	if IsValidCdktfDirectory(path) {
		log.Printf("[DEBUG] Skipping: %s", path)
		return nil
	}

	if FileIgnoreCheck(path) {
		log.Printf("[DEBUG] Skipping: %s", path)
		return nil
	}

	if err := LegacyFileExtensionCheck(path); err != nil {
		return fmt.Errorf("%s: error checking file extension: %w", path, err)
	}

	if err := FileSizeCheck(fullpath); err != nil {
		return fmt.Errorf("%s: error checking file size: %w", path, err)
	}

	content, err := os.ReadFile(fullpath)

	if err != nil {
		return fmt.Errorf("%s: error reading file: %w", path, err)
	}

	subcategory, err := NewFrontMatterCheck(check.Options.FrontMatter).Run(content)

	if err != nil {
		return fmt.Errorf("%s: error checking file frontmatter: %w", path, err)
	}

	// We don't want to check the content for CDKTF files since they are converted
	if !IsValidCdktfDirectory(filepath.Dir(fullpath)) {
		if err := NewContentsCheck(check.Options.Contents).Run(fullpath, exampleLanguage, subcategory); err != nil {
			return fmt.Errorf("%s: error checking file contents: %w", path, err)
		}
	}
	return nil
}

func (check *LegacyEphemeralFileCheck) RunAll(files []string, exampleLanguage string) error {
	var result *multierror.Error

	for _, file := range files {
		if err := check.Run(file, exampleLanguage); err != nil {
			result = multierror.Append(result, err)
		}
	}

	return result.ErrorOrNil()
}
