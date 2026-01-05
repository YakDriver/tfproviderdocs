// Copyright IBM Corp. 2019, 2026
// SPDX-License-Identifier: MPL-2.0

package check

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-multierror"
)

type RegistryEphemeralFileOptions struct {
	*FileOptions

	Contents     *ContentsOptions
	FrontMatter  *FrontMatterOptions
	ProviderName string
}

type RegistryEphemeralFileCheck struct {
	FileCheck

	Options *RegistryEphemeralFileOptions
}

func NewRegistryEphemeralFileCheck(opts *RegistryEphemeralFileOptions) *RegistryEphemeralFileCheck {
	check := &RegistryEphemeralFileCheck{
		Options: opts,
	}

	if check.Options == nil {
		check.Options = &RegistryEphemeralFileOptions{}
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

	check.Options.FrontMatter.NoLayout = true
	check.Options.FrontMatter.NoSidebarCurrent = true

	return check
}

func (check *RegistryEphemeralFileCheck) Run(path string, exampleLanguage string) error {
	fullpath := check.Options.FullPath(path)

	log.Printf("[DEBUG] Checking file: %s", fullpath)

	if err := RegistryFileExtensionCheck(path); err != nil {
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

func (check *RegistryEphemeralFileCheck) RunAll(files []string, exampleLanguage string) error {
	var result *multierror.Error

	for _, file := range files {
		if err := check.Run(file, exampleLanguage); err != nil {
			result = multierror.Append(result, err)
		}
	}

	return result.ErrorOrNil()
}
