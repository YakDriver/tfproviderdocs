package check

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-multierror"
)

type RegistryActionFileOptions struct {
	*FileOptions

	Contents     *ContentsOptions
	FrontMatter  *FrontMatterOptions
	ProviderName string
}

type RegistryActionFileCheck struct {
	FileCheck

	Options *RegistryActionFileOptions
}

func NewRegistryActionFileCheck(opts *RegistryActionFileOptions) *RegistryActionFileCheck {
	check := &RegistryActionFileCheck{
		Options: opts,
	}

	if check.Options == nil {
		check.Options = &RegistryActionFileOptions{}
	}

	if check.Options.Contents == nil {
		check.Options.Contents = &ContentsOptions{}
	}

	check.Options.Contents.Enable = true

	if check.Options.Contents.ProviderName == "" {
		check.Options.Contents.ProviderName = check.Options.ProviderName
	}

	check.Options.Contents.DisableRegionArgumentCheck = true
	check.Options.Contents.DisallowAttributesSection = true
	if check.Options.Contents.AttributesSectionDisallowedMessage == "" {
		check.Options.Contents.AttributesSectionDisallowedMessage = "actions documentation cannot include an attributes section"
	}
	check.Options.Contents.DisallowImportSection = true
	if check.Options.Contents.ImportSectionDisallowedMessage == "" {
		check.Options.Contents.ImportSectionDisallowedMessage = "actions documentation cannot include an import section"
	}
	if len(check.Options.Contents.ArgumentsBylineTexts) == 0 {
		check.Options.Contents.ArgumentsBylineTexts = []string{
			"This action supports the following arguments:",
			"The following arguments are required:",
			"The following arguments are optional:",
			"This action does not support any arguments.",
		}
	}

	if check.Options.FileOptions == nil {
		check.Options.FileOptions = &FileOptions{}
	}

	if check.Options.FrontMatter == nil {
		check.Options.FrontMatter = &FrontMatterOptions{}
	}

	check.Options.FrontMatter.NoLayout = true
	check.Options.FrontMatter.NoSidebarCurrent = true
	check.Options.FrontMatter.RequireDescription = true
	check.Options.FrontMatter.RequirePageTitle = true
	check.Options.FrontMatter.RequireSubcategory = true

	return check
}

func (check *RegistryActionFileCheck) Run(path string, exampleLanguage string) error {
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

	if !IsValidCdktfDirectory(filepath.Dir(fullpath)) {
		if err := NewContentsCheck(check.Options.Contents).Run(fullpath, exampleLanguage, subcategory); err != nil {
			return fmt.Errorf("%s: error checking file contents: %w", path, err)
		}
	}
	return nil
}

func (check *RegistryActionFileCheck) RunAll(files []string, exampleLanguage string) error {
	var result *multierror.Error

	for _, file := range files {
		if err := check.Run(file, exampleLanguage); err != nil {
			result = multierror.Append(result, err)
		}
	}

	return result.ErrorOrNil()
}
