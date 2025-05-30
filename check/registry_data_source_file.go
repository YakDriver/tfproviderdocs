package check

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-multierror"
)

type RegistryDataSourceFileOptions struct {
	*FileOptions

	Contents    *ContentsOptions
	FrontMatter *FrontMatterOptions
}

type RegistryDataSourceFileCheck struct {
	FileCheck

	Options *RegistryDataSourceFileOptions
}

func NewRegistryDataSourceFileCheck(opts *RegistryDataSourceFileOptions) *RegistryDataSourceFileCheck {
	check := &RegistryDataSourceFileCheck{
		Options: opts,
	}

	if check.Options == nil {
		check.Options = &RegistryDataSourceFileOptions{}
	}

	if check.Options.Contents == nil {
		check.Options.Contents = &ContentsOptions{}
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

func (check *RegistryDataSourceFileCheck) Run(path string, exampleLanguage string) error {
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

func (check *RegistryDataSourceFileCheck) RunAll(files []string, exampleLanguage string) error {
	var result *multierror.Error

	for _, file := range files {
		if err := check.Run(file, exampleLanguage); err != nil {
			result = multierror.Append(result, err)
		}
	}

	return result.ErrorOrNil()
}
