package check

import (
	"fmt"
	"slices"

	"github.com/YakDriver/tfproviderdocs/check/contents"
)

type ContentsCheck struct {
	Options *ContentsOptions
}

// ContentsOptions represents configuration options for Contents.
type ContentsOptions struct {
	*FileOptions

	Enable                                 bool
	EnhancedRegionChecks                   bool
	ProviderName                           string
	RequireAttributesSection               contents.SectionRequirement
	RequireTimeoutsSection                 contents.SectionRequirement
	RequireImportSection                   contents.SectionRequirement
	RequireSchemaOrdering                  bool
	IgnoreContentsCheck                    []string
	IgnoreEnhancedRegionCheck              []string
	IgnoreEnhancedRegionCheckSubcategories []string

	DisableRegionArgumentCheck         bool
	DisallowAttributesSection          bool
	AttributesSectionDisallowedMessage string
	DisallowImportSection              bool
	ImportSectionDisallowedMessage     string
	ArgumentsBylineTexts               []string
}

func NewContentsCheck(opts *ContentsOptions) *ContentsCheck {
	check := &ContentsCheck{
		Options: opts,
	}

	if check.Options == nil {
		check.Options = &ContentsOptions{}
	}

	if check.Options.FileOptions == nil {
		check.Options.FileOptions = &FileOptions{}
	}

	return check
}

func (check *ContentsCheck) Run(path string, exampleLanguage string, subcategory *string) error {
	if !check.Options.Enable {
		return nil
	}

	checkOpts := &contents.CheckOptions{
		ArgumentsSection: &contents.CheckArgumentsSectionOptions{
			EnhancedRegionChecks:  check.Options.EnhancedRegionChecks,
			RegionAware:           true,
			RequireSchemaOrdering: check.Options.RequireSchemaOrdering,
			ExpectedBylineTexts:   check.Options.ArgumentsBylineTexts,
		},
		AttributesSection: &contents.CheckAttributesSectionOptions{
			RequireSchemaOrdering: check.Options.RequireSchemaOrdering,
			RequireSection:        check.Options.RequireAttributesSection,
		},
		ExamplesSection: &contents.CheckExamplesSectionOptions{
			ExpectedCodeBlockLanguage: exampleLanguage,
		},
		TimeoutsSection: &contents.CheckTimeoutsSectionOptions{
			RequireSection: check.Options.RequireTimeoutsSection,
		},
		ImportSection: &contents.CheckImportSectionOptions{
			RequireSection: check.Options.RequireImportSection,
		},
		DisallowAttributesSection:          check.Options.DisallowAttributesSection,
		AttributesSectionDisallowedMessage: check.Options.AttributesSectionDisallowedMessage,
		DisallowImportSection:              check.Options.DisallowImportSection,
		ImportSectionDisallowedMessage:     check.Options.ImportSectionDisallowedMessage,
	}

	if check.Options.DisableRegionArgumentCheck {
		checkOpts.ArgumentsSection.RegionAware = false
	}

	doc := contents.NewDocument(path, check.Options.ProviderName)

	if err := doc.Parse(); err != nil {
		return fmt.Errorf("error parsing file: %w", err)
	}

	if len(check.Options.IgnoreContentsCheck) > 0 && slices.Contains(check.Options.IgnoreContentsCheck, doc.ResourceName) {
		return nil
	}

	if len(check.Options.IgnoreEnhancedRegionCheck) > 0 && slices.Contains(check.Options.IgnoreEnhancedRegionCheck, doc.ResourceName) {
		checkOpts.ArgumentsSection.RegionAware = false
	}

	if len(check.Options.IgnoreEnhancedRegionCheckSubcategories) > 0 && subcategory != nil && slices.Contains(check.Options.IgnoreEnhancedRegionCheckSubcategories, *subcategory) {
		checkOpts.ArgumentsSection.RegionAware = false
	}

	if err := doc.Check(checkOpts); err != nil {
		return err
	}

	return nil
}
