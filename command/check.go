// Copyright (c) HashiCorp, Inc. 2019-2026
// SPDX-License-Identifier: MPL-2.0

package command

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/YakDriver/tfproviderdocs/check"
	"github.com/YakDriver/tfproviderdocs/check/contents"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/mitchellh/cli"
)

type CheckCommandConfig struct {
	AllowedGuideSubcategories                  string
	AllowedGuideSubcategoriesFile              string
	AllowedResourceSubcategories               string
	AllowedResourceSubcategoriesFile           string
	EnableContentsCheck                        bool
	EnableEnhancedRegionCheck                  bool
	IgnoreCdktfMissingFiles                    bool
	IgnoreContentsCheckDataSources             string
	IgnoreContentsCheckActions                 string
	IgnoreContentsCheckEphemerals              string
	IgnoreContentsCheckFunctions               string
	IgnoreContentsCheckResources               string
	IgnoreEnhancedRegionCheckDataSources       string
	IgnoreEnhancedRegionCheckDataSourcesFile   string
	IgnoreEnhancedRegionCheckEphemerals        string
	IgnoreEnhancedRegionCheckEphemeralsFile    string
	IgnoreEnhancedRegionCheckResources         string
	IgnoreEnhancedRegionCheckResourcesFile     string
	IgnoreEnhancedRegionCheckSubcategories     string
	IgnoreEnhancedRegionCheckSubcategoriesFile string
	IgnoreFileMismatchDataSources              string
	IgnoreFileMismatchActions                  string
	IgnoreFileMismatchEphemerals               string
	IgnoreFileMismatchFunctions                string
	IgnoreFileMismatchListResources            string
	IgnoreFileMismatchResources                string
	IgnoreFileMissingDataSources               string
	IgnoreFileMissingActions                   string
	IgnoreFileMissingEphemerals                string
	IgnoreFileMissingFunctions                 string
	IgnoreFileMissingListResources             string
	IgnoreFileMissingResources                 string
	LogLevel                                   string
	Path                                       string
	ProviderName                               string
	ProviderSource                             string
	ProvidersSchemaJson                        string
	RequireGuideSubcategory                    bool
	RequireResourceSubcategory                 bool
	RequireSchemaOrdering                      bool
}

// CheckCommand is a Command implementation
type CheckCommand struct {
	Ui cli.Ui
}

func (*CheckCommand) Help() string {
	optsBuffer := bytes.NewBuffer([]byte{})
	opts := tabwriter.NewWriter(optsBuffer, 0, 0, 1, ' ', 0)
	LogLevelFlagHelp(opts)
	fmt.Fprintf(opts, CommandHelpOptionFormat, "-allowed-guide-subcategories", "Comma separated list of allowed guide frontmatter subcategories.")
	fmt.Fprintf(opts, CommandHelpOptionFormat, "-allowed-guide-subcategories-file", "Path to newline separated file of allowed guide frontmatter subcategories.")
	fmt.Fprintf(opts, CommandHelpOptionFormat, "-allowed-resource-subcategories", "Comma separated list of allowed data source and resource frontmatter subcategories.")
	fmt.Fprintf(opts, CommandHelpOptionFormat, "-allowed-resource-subcategories-file", "Path to newline separated file of allowed data source and resource frontmatter subcategories.")
	fmt.Fprintf(opts, CommandHelpOptionFormat, "-enable-contents-check", "(Experimental) Enable contents checking.")
	fmt.Fprintf(opts, CommandHelpOptionFormat, "-enable-enhanced-region-check", "Enable enhanced Region functionality checks (requires -enable-contents-check).")
	fmt.Fprintf(opts, CommandHelpOptionFormat, "-ignore-cdktf-missing-files", "Ignore checks for missing CDK for Terraform documentation files when iteratively introducing them in large providers.")
	fmt.Fprintf(opts, CommandHelpOptionFormat, "-ignore-contents-check-data-sources", "Comma separated list of data sources to ignore contents checking.")
	fmt.Fprintf(opts, CommandHelpOptionFormat, "-ignore-contents-check-actions", "Comma separated list of actions to ignore contents checking.")
	fmt.Fprintf(opts, CommandHelpOptionFormat, "-ignore-contents-check-ephemerals", "Comma separated list of ephemerals to ignore contents checking.")
	fmt.Fprintf(opts, CommandHelpOptionFormat, "-ignore-contents-check-functions", "Comma separated list of functions to ignore contents checking.")
	fmt.Fprintf(opts, CommandHelpOptionFormat, "-ignore-contents-check-resources", "Comma separated list of resources to ignore contents checking.")
	fmt.Fprintf(opts, CommandHelpOptionFormat, "-ignore-enhanced-region-check-data-sources", "Comma separated list of data sources to ignore enhanced Region functionality checks.")
	fmt.Fprintf(opts, CommandHelpOptionFormat, "-ignore-enhanced-region-check-data-sources-file", "Path to newline separated file of data sources to ignore enhanced Region functionality checks.")
	fmt.Fprintf(opts, CommandHelpOptionFormat, "-ignore-enhanced-region-check-ephemerals", "Comma separated list of ephemerals to ignore enhanced Region functionality checks.")
	fmt.Fprintf(opts, CommandHelpOptionFormat, "-ignore-enhanced-region-check-ephemerals-file", "Path to newline separated file of ephemerals to ignore enhanced Region functionality checks.")
	fmt.Fprintf(opts, CommandHelpOptionFormat, "-ignore-enhanced-region-check-resources", "Comma separated list of resources to ignore enhanced Region functionality checks.")
	fmt.Fprintf(opts, CommandHelpOptionFormat, "-ignore-enhanced-region-check-resources-file", "Path to newline separated file of resources to ignore enhanced Region functionality checks.")
	fmt.Fprintf(opts, CommandHelpOptionFormat, "-ignore-enhanced-region-check-subcategories", "Comma separated list of frontmatter subcategories to ignore enhanced Region functionality checks.")
	fmt.Fprintf(opts, CommandHelpOptionFormat, "-ignore-enhanced-region-check-subcategories-file", "Path to newline separated file of frontmatter subcategories to ignore enhanced Region functionality checks.")
	fmt.Fprintf(opts, CommandHelpOptionFormat, "-ignore-file-mismatch-data-sources", "Comma separated list of data sources to ignore mismatched/extra files.")
	fmt.Fprintf(opts, CommandHelpOptionFormat, "-ignore-file-mismatch-actions", "Comma separated list of actions to ignore mismatched/extra files.")
	fmt.Fprintf(opts, CommandHelpOptionFormat, "-ignore-file-mismatch-ephemerals", "Comma separated list of ephemerals to ignore mismatched/extra files.")
	fmt.Fprintf(opts, CommandHelpOptionFormat, "-ignore-file-mismatch-functions", "Comma separated list of functions to ignore mismatched/extra files.")
	fmt.Fprintf(opts, CommandHelpOptionFormat, "-ignore-file-mismatch-list-resources", "Comma separated list of list resources to ignore mismatched/extra files.")
	fmt.Fprintf(opts, CommandHelpOptionFormat, "-ignore-file-mismatch-resources", "Comma separated list of resources to ignore mismatched/extra files.")
	fmt.Fprintf(opts, CommandHelpOptionFormat, "-ignore-file-missing-data-sources", "Comma separated list of data sources to ignore missing files.")
	fmt.Fprintf(opts, CommandHelpOptionFormat, "-ignore-file-missing-actions", "Comma separated list of actions to ignore missing files.")
	fmt.Fprintf(opts, CommandHelpOptionFormat, "-ignore-file-missing-ephemerals", "Comma separated list of ephemerals to ignore missing files.")
	fmt.Fprintf(opts, CommandHelpOptionFormat, "-ignore-file-missing-functions", "Comma separated list of functions to ignore missing files.")
	fmt.Fprintf(opts, CommandHelpOptionFormat, "-ignore-file-missing-list-resources", "Comma separated list of list resources to ignore missing files.")
	fmt.Fprintf(opts, CommandHelpOptionFormat, "-ignore-file-missing-resources", "Comma separated list of resources to ignore missing files.")
	fmt.Fprintf(opts, CommandHelpOptionFormat, "-provider-name", "Terraform Provider short name (e.g. aws). Automatically determined if -provider-source is given or if current working directory or provided path is prefixed with terraform-provider-*.")
	fmt.Fprintf(opts, CommandHelpOptionFormat, "-provider-source", "Terraform Provider source address (e.g. registry.terraform.io/hashicorp/aws) for Terraform CLI 0.13 and later -providers-schema-json. Automatically sets -provider-name by dropping hostname and namespace prefix.")
	fmt.Fprintf(opts, CommandHelpOptionFormat, "-providers-schema-json", "Path to terraform providers schema -json file. Enables enhanced validations.")
	fmt.Fprintf(opts, CommandHelpOptionFormat, "-require-guide-subcategory", "Require guide frontmatter subcategory.")
	fmt.Fprintf(opts, CommandHelpOptionFormat, "-require-resource-subcategory", "Require data source and resource frontmatter subcategory.")
	fmt.Fprintf(opts, CommandHelpOptionFormat, "-require-schema-ordering", "Require schema attribute lists to be alphabetically ordered (requires -enable-contents-check).")
	opts.Flush()

	helpText := fmt.Sprintf(`
Usage: tfproviderdocs check [options] [PATH]

  Performs documentation directory and file checks against the given Terraform Provider codebase.

Options:

%s
`, optsBuffer.String())

	return strings.TrimSpace(helpText)
}

func (c *CheckCommand) Name() string { return "check" }

func configureCheckCommandFlags(flags *flag.FlagSet, config *CheckCommandConfig) {
	LogLevelFlag(flags, &config.LogLevel)
	flags.StringVar(&config.AllowedGuideSubcategories, "allowed-guide-subcategories", "", "")
	flags.StringVar(&config.AllowedGuideSubcategoriesFile, "allowed-guide-subcategories-file", "", "")
	flags.StringVar(&config.AllowedResourceSubcategories, "allowed-resource-subcategories", "", "")
	flags.StringVar(&config.AllowedResourceSubcategoriesFile, "allowed-resource-subcategories-file", "", "")
	flags.BoolVar(&config.EnableContentsCheck, "enable-contents-check", false, "")
	flags.BoolVar(&config.EnableEnhancedRegionCheck, "enable-enhanced-region-check", false, "")
	flags.BoolVar(&config.IgnoreCdktfMissingFiles, "ignore-cdktf-missing-files", false, "")
	flags.StringVar(&config.IgnoreContentsCheckDataSources, "ignore-contents-check-data-sources", "", "")
	flags.StringVar(&config.IgnoreContentsCheckActions, "ignore-contents-check-actions", "", "")
	flags.StringVar(&config.IgnoreContentsCheckEphemerals, "ignore-contents-check-ephemerals", "", "")
	flags.StringVar(&config.IgnoreContentsCheckFunctions, "ignore-contents-check-functions", "", "")
	flags.StringVar(&config.IgnoreContentsCheckResources, "ignore-contents-check-resources", "", "")
	flags.StringVar(&config.IgnoreEnhancedRegionCheckDataSources, "ignore-enhanced-region-check-data-sources", "", "")
	flags.StringVar(&config.IgnoreEnhancedRegionCheckDataSourcesFile, "ignore-enhanced-region-check-data-sources-file", "", "")
	flags.StringVar(&config.IgnoreEnhancedRegionCheckEphemerals, "ignore-enhanced-region-check-ephemerals", "", "")
	flags.StringVar(&config.IgnoreEnhancedRegionCheckEphemeralsFile, "ignore-enhanced-region-check-ephemerals-file", "", "")
	flags.StringVar(&config.IgnoreEnhancedRegionCheckResources, "ignore-enhanced-region-check-resources", "", "")
	flags.StringVar(&config.IgnoreEnhancedRegionCheckResourcesFile, "ignore-enhanced-region-check-resources-file", "", "")
	flags.StringVar(&config.IgnoreEnhancedRegionCheckSubcategories, "ignore-enhanced-region-check-subcategories", "", "")
	flags.StringVar(&config.IgnoreEnhancedRegionCheckSubcategoriesFile, "ignore-enhanced-region-check-subcategories-file", "", "")
	flags.StringVar(&config.IgnoreFileMismatchDataSources, "ignore-file-mismatch-data-sources", "", "")
	flags.StringVar(&config.IgnoreFileMismatchActions, "ignore-file-mismatch-actions", "", "")
	flags.StringVar(&config.IgnoreFileMismatchEphemerals, "ignore-file-mismatch-ephemerals", "", "")
	flags.StringVar(&config.IgnoreFileMismatchFunctions, "ignore-file-mismatch-functions", "", "")
	flags.StringVar(&config.IgnoreFileMismatchListResources, "ignore-file-mismatch-list-resources", "", "")
	flags.StringVar(&config.IgnoreFileMismatchResources, "ignore-file-mismatch-resources", "", "")
	flags.StringVar(&config.IgnoreFileMissingDataSources, "ignore-file-missing-data-sources", "", "")
	flags.StringVar(&config.IgnoreFileMissingActions, "ignore-file-missing-actions", "", "")
	flags.StringVar(&config.IgnoreFileMissingEphemerals, "ignore-file-missing-ephemerals", "", "")
	flags.StringVar(&config.IgnoreFileMissingFunctions, "ignore-file-missing-functions", "", "")
	flags.StringVar(&config.IgnoreFileMissingListResources, "ignore-file-missing-list-resources", "", "")
	flags.StringVar(&config.IgnoreFileMissingResources, "ignore-file-missing-resources", "", "")
	flags.StringVar(&config.ProviderName, "provider-name", "", "")
	flags.StringVar(&config.ProviderSource, "provider-source", "", "")
	flags.StringVar(&config.ProvidersSchemaJson, "providers-schema-json", "", "")
	flags.BoolVar(&config.RequireGuideSubcategory, "require-guide-subcategory", false, "")
	flags.BoolVar(&config.RequireResourceSubcategory, "require-resource-subcategory", false, "")
	flags.BoolVar(&config.RequireSchemaOrdering, "require-schema-ordering", false, "")
}

func (c *CheckCommand) Run(args []string) int {
	var config CheckCommandConfig

	flags := flag.NewFlagSet(c.Name(), flag.ContinueOnError)
	flags.Usage = func() { c.Ui.Info(c.Help()) }
	configureCheckCommandFlags(flags, &config)

	if err := flags.Parse(args); err != nil {
		flags.Usage()
		return 1
	}

	args = flags.Args()

	if len(args) == 1 {
		config.Path = args[0]
	}

	ConfigureLogging(c.Name(), config.LogLevel)

	if config.ProviderName == "" && config.ProviderSource != "" {
		providerSourceParts := strings.Split(config.ProviderSource, "/")
		config.ProviderName = providerSourceParts[len(providerSourceParts)-1]
	}

	if config.ProviderName == "" {
		if config.Path == "" {
			config.ProviderName = providerNameFromCurrentDirectory()
		} else {
			config.ProviderName = providerNameFromPath(config.Path)
		}
	}

	if config.ProviderName == "" {
		log.Printf("[WARN] Unable to determine provider name. Contents and enhanced validations may fail.")
	} else {
		log.Printf("[DEBUG] Found provider name: %s", config.ProviderName)
	}

	directories, err := check.GetDirectories(config.Path)

	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error getting Terraform Provider documentation directories: %s", err))
		return 1
	}

	if len(directories) == 0 {
		if config.Path == "" {
			c.Ui.Error("No Terraform Provider documentation directories found in current path")
		} else {
			c.Ui.Error(fmt.Sprintf("No Terraform Provider documentation directories found in path: %s", config.Path))
		}

		return 1
	}

	var allowedGuideSubcategories []string
	if v := config.AllowedGuideSubcategories; v != "" {
		allowedGuideSubcategories = strings.Split(v, ",")
	}

	if v := config.AllowedGuideSubcategoriesFile; v != "" {
		var err error
		allowedGuideSubcategories, err = allowedSubcategoriesFile(v)

		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error getting allowed guide subcategories: %s", err))
			return 1
		}
	}

	var allowedResourceSubcategories []string
	if v := config.AllowedResourceSubcategories; v != "" {
		allowedResourceSubcategories = strings.Split(v, ",")
	}

	if v := config.AllowedResourceSubcategoriesFile; v != "" {
		var err error
		allowedResourceSubcategories, err = allowedSubcategoriesFile(v)

		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error getting allowed resource subcategories: %s", err))
			return 1
		}
	}

	var ignoreContentsCheckDataSources []string
	var ignoreContentsCheckActions []string
	if v := config.IgnoreContentsCheckDataSources; v != "" {
		ignoreContentsCheckDataSources = strings.Split(v, ",")
	}

	if v := config.IgnoreContentsCheckActions; v != "" {
		ignoreContentsCheckActions = strings.Split(v, ",")
	}

	var ignoreContentsCheckEphemerals []string
	if v := config.IgnoreContentsCheckEphemerals; v != "" {
		ignoreContentsCheckEphemerals = strings.Split(v, ",")
	}

	var ignoreContentsCheckFunctions []string
	if v := config.IgnoreContentsCheckFunctions; v != "" {
		ignoreContentsCheckFunctions = strings.Split(v, ",")
	}

	var ignoreContentsCheckResources []string
	if v := config.IgnoreContentsCheckResources; v != "" {
		ignoreContentsCheckResources = strings.Split(v, ",")
	}

	var ignoreEnhancedRegionCheckDataSources []string
	if v := config.IgnoreEnhancedRegionCheckDataSources; v != "" {
		ignoreEnhancedRegionCheckDataSources = strings.Split(v, ",")
	}

	if v := config.IgnoreEnhancedRegionCheckDataSourcesFile; v != "" {
		var err error
		ignoreEnhancedRegionCheckDataSources, err = allowedSubcategoriesFile(v)

		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error getting ignore enhanced Region check data sources: %s", err))
			return 1
		}
	}

	var ignoreEnhancedRegionCheckEphemerals []string
	if v := config.IgnoreEnhancedRegionCheckEphemerals; v != "" {
		ignoreEnhancedRegionCheckEphemerals = strings.Split(v, ",")
	}

	if v := config.IgnoreEnhancedRegionCheckEphemeralsFile; v != "" {
		var err error
		ignoreEnhancedRegionCheckEphemerals, err = allowedSubcategoriesFile(v)

		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error getting ignore enhanced Region check ephemerals: %s", err))
			return 1
		}
	}

	var ignoreEnhancedRegionCheckResources []string
	if v := config.IgnoreEnhancedRegionCheckResources; v != "" {
		ignoreEnhancedRegionCheckResources = strings.Split(v, ",")
	}

	if v := config.IgnoreEnhancedRegionCheckResourcesFile; v != "" {
		var err error
		ignoreEnhancedRegionCheckResources, err = allowedSubcategoriesFile(v)

		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error getting ignore enhanced Region check resources: %s", err))
			return 1
		}
	}

	var ignoreEnhancedRegionCheckSubcategories []string
	if v := config.IgnoreEnhancedRegionCheckSubcategories; v != "" {
		ignoreEnhancedRegionCheckSubcategories = strings.Split(v, ",")
	}

	if v := config.IgnoreEnhancedRegionCheckSubcategoriesFile; v != "" {
		var err error
		ignoreEnhancedRegionCheckSubcategories, err = allowedSubcategoriesFile(v)

		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error getting ignore enhanced Region check subcategories: %s", err))
			return 1
		}
	}

	var ignoreFileMismatchDataSources []string
	var ignoreFileMismatchActions []string
	if v := config.IgnoreFileMismatchDataSources; v != "" {
		ignoreFileMismatchDataSources = strings.Split(v, ",")
	}

	if v := config.IgnoreFileMismatchActions; v != "" {
		ignoreFileMismatchActions = strings.Split(v, ",")
	}

	var ignoreFileMismatchEphemerals []string
	if v := config.IgnoreFileMismatchEphemerals; v != "" {
		ignoreFileMismatchEphemerals = strings.Split(v, ",")
	}

	var ignoreFileMismatchFunctions []string
	if v := config.IgnoreFileMismatchFunctions; v != "" {
		ignoreFileMismatchFunctions = strings.Split(v, ",")
	}

	var ignoreFileMismatchListResources []string
	if v := config.IgnoreFileMismatchListResources; v != "" {
		ignoreFileMismatchListResources = strings.Split(v, ",")
	}

	var ignoreFileMismatchResources []string
	if v := config.IgnoreFileMismatchResources; v != "" {
		ignoreFileMismatchResources = strings.Split(v, ",")
	}

	var ignoreFileMissingDataSources []string
	var ignoreFileMissingActions []string
	if v := config.IgnoreFileMissingDataSources; v != "" {
		ignoreFileMissingDataSources = strings.Split(v, ",")
	}

	if v := config.IgnoreFileMissingActions; v != "" {
		ignoreFileMissingActions = strings.Split(v, ",")
	}
	var ignoreFileMissingEphemerals []string
	if v := config.IgnoreFileMissingEphemerals; v != "" {
		ignoreFileMissingEphemerals = strings.Split(v, ",")
	}

	var ignoreFileMissingFunctions []string
	if v := config.IgnoreFileMissingFunctions; v != "" {
		ignoreFileMissingFunctions = strings.Split(v, ",")
	}

	var ignoreFileMissingListResources []string
	if v := config.IgnoreFileMissingListResources; v != "" {
		ignoreFileMissingListResources = strings.Split(v, ",")
	}

	var ignoreFileMissingResources []string
	if v := config.IgnoreFileMissingResources; v != "" {
		ignoreFileMissingResources = strings.Split(v, ",")
	}

	var actionNames, dataSourceNames, ephemeralNames, listResourceNames, resourceNames, functionNames []string
	if config.ProvidersSchemaJson != "" {
		ps, err := providerSchemas(config.ProvidersSchemaJson)

		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error enabling Terraform Provider schema checks: %s", err))
			return 1
		}

		if config.ProviderName == "" {
			msg := `Unknown provider name for enabling Terraform Provider schema checks.

Check that the current working directory or provided path is prefixed with terraform-provider-*.`
			c.Ui.Error(msg)
			return 1
		}

		actionNames = providerSchemasActions(ps, config.ProviderName, config.ProviderSource)
		dataSourceNames = providerSchemasDataSources(ps, config.ProviderName, config.ProviderSource)
		ephemeralNames = providerSchemasEphemerals(ps, config.ProviderName, config.ProviderSource)
		functionNames = providerSchemasFunctions(ps, config.ProviderName, config.ProviderSource)
		listResourceNames = providerSchemasListResources(ps, config.ProviderName, config.ProviderSource)
		resourceNames = providerSchemasResources(ps, config.ProviderName, config.ProviderSource)
	}

	fileOpts := &check.FileOptions{
		BasePath: config.Path,
	}
	checkOpts := &check.CheckOptions{
		// action
		RegistryActionFile: &check.RegistryActionFileOptions{
			Contents: &check.ContentsOptions{
				Enable:                             config.EnableContentsCheck,
				RequireSchemaOrdering:              config.RequireSchemaOrdering,
				IgnoreContentsCheck:                ignoreContentsCheckActions,
				ProviderName:                       config.ProviderName,
				TitleSectionPrefixes:               []string{"Action"},
				DisableRegionArgumentCheck:         true,
				DisallowAttributesSection:          true,
				AttributesSectionDisallowedMessage: "actions documentation cannot include an attributes section",
				DisallowImportSection:              true,
				ImportSectionDisallowedMessage:     "actions documentation cannot include an import section",
				ArgumentsBylineTexts: []string{
					"This action supports the following arguments:",
					"The following arguments are required:",
					"The following arguments are optional:",
					"This action does not support any arguments.",
				},
			},
			FileOptions: fileOpts,
			FrontMatter: &check.FrontMatterOptions{
				AllowedSubcategories: allowedResourceSubcategories,
				RequireSubcategory:   config.RequireResourceSubcategory,
			},
			ProviderName: config.ProviderName,
		},
		LegacyActionFile: &check.LegacyActionFileOptions{
			Contents: &check.ContentsOptions{
				Enable:                             config.EnableContentsCheck,
				RequireSchemaOrdering:              config.RequireSchemaOrdering,
				IgnoreContentsCheck:                ignoreContentsCheckActions,
				ProviderName:                       config.ProviderName,
				TitleSectionPrefixes:               []string{"Action"},
				DisableRegionArgumentCheck:         true,
				DisallowAttributesSection:          true,
				AttributesSectionDisallowedMessage: "actions documentation cannot include an attributes section",
				DisallowImportSection:              true,
				ImportSectionDisallowedMessage:     "actions documentation cannot include an import section",
				ArgumentsBylineTexts: []string{
					"This action supports the following arguments:",
					"The following arguments are required:",
					"The following arguments are optional:",
					"This action does not support any arguments.",
				},
			},
			FileOptions: fileOpts,
			FrontMatter: &check.FrontMatterOptions{
				AllowedSubcategories: allowedResourceSubcategories,
				RequireSubcategory:   config.RequireResourceSubcategory,
			},
			ProviderName: config.ProviderName,
		},
		ActionFileMismatch: &check.FileMismatchOptions{
			IgnoreFileMismatch: ignoreFileMismatchActions,
			IgnoreFileMissing:  ignoreFileMissingActions,
			ProviderName:       config.ProviderName,
			ResourceType:       check.ResourceTypeAction,
			ResourceNames:      actionNames,
		},

		// data source
		RegistryDataSourceFile: &check.RegistryDataSourceFileOptions{
			Contents: &check.ContentsOptions{
				Enable:                                 config.EnableContentsCheck,
				EnhancedRegionChecks:                   config.EnableEnhancedRegionCheck,
				RequireAttributesSection:               contents.Required,
				RequireImportSection:                   contents.Forbidden,
				RequireSchemaOrdering:                  config.RequireSchemaOrdering,
				IgnoreContentsCheck:                    ignoreContentsCheckDataSources,
				IgnoreEnhancedRegionCheck:              ignoreEnhancedRegionCheckDataSources,
				IgnoreEnhancedRegionCheckSubcategories: ignoreEnhancedRegionCheckSubcategories,
				ProviderName:                           config.ProviderName,
				TitleSectionPrefixes:                   []string{"Data Source"},
			},
			FileOptions: fileOpts,
			FrontMatter: &check.FrontMatterOptions{
				AllowedSubcategories: allowedResourceSubcategories,
				RequireSubcategory:   config.RequireResourceSubcategory,
			},
		},
		LegacyDataSourceFile: &check.LegacyDataSourceFileOptions{
			Contents: &check.ContentsOptions{
				Enable:                                 config.EnableContentsCheck,
				EnhancedRegionChecks:                   config.EnableEnhancedRegionCheck,
				RequireAttributesSection:               contents.Required,
				RequireImportSection:                   contents.Forbidden,
				RequireSchemaOrdering:                  config.RequireSchemaOrdering,
				IgnoreContentsCheck:                    ignoreContentsCheckDataSources,
				IgnoreEnhancedRegionCheck:              ignoreEnhancedRegionCheckDataSources,
				IgnoreEnhancedRegionCheckSubcategories: ignoreEnhancedRegionCheckSubcategories,
				ProviderName:                           config.ProviderName,
				TitleSectionPrefixes:                   []string{"Data Source"},
			},
			FileOptions: fileOpts,
			FrontMatter: &check.FrontMatterOptions{
				AllowedSubcategories: allowedResourceSubcategories,
				RequireSubcategory:   config.RequireResourceSubcategory,
			},
		},
		DataSourceFileMismatch: &check.FileMismatchOptions{
			IgnoreFileMismatch: ignoreFileMismatchDataSources,
			IgnoreFileMissing:  ignoreFileMissingDataSources,
			ProviderName:       config.ProviderName,
			ResourceType:       check.ResourceTypeDataSource,
			ResourceNames:      dataSourceNames,
		},

		// ephemeral
		RegistryEphemeralFile: &check.RegistryEphemeralFileOptions{
			Contents: &check.ContentsOptions{
				Enable:                                 config.EnableContentsCheck,
				EnhancedRegionChecks:                   config.EnableEnhancedRegionCheck,
				RequireAttributesSection:               contents.Required,
				RequireImportSection:                   contents.Forbidden,
				RequireSchemaOrdering:                  config.RequireSchemaOrdering,
				IgnoreContentsCheck:                    ignoreContentsCheckEphemerals,
				IgnoreEnhancedRegionCheck:              ignoreEnhancedRegionCheckEphemerals,
				IgnoreEnhancedRegionCheckSubcategories: ignoreEnhancedRegionCheckSubcategories,
				ProviderName:                           config.ProviderName,
				TitleSectionPrefixes:                   []string{"Ephemeral"},
			},
			FileOptions: fileOpts,
			FrontMatter: &check.FrontMatterOptions{
				AllowedSubcategories: allowedResourceSubcategories,
				RequireSubcategory:   config.RequireResourceSubcategory,
			},
			ProviderName: config.ProviderName,
		},
		LegacyEphemeralFile: &check.LegacyEphemeralFileOptions{
			Contents: &check.ContentsOptions{
				Enable:                                 config.EnableContentsCheck,
				EnhancedRegionChecks:                   config.EnableEnhancedRegionCheck,
				RequireAttributesSection:               contents.Required,
				RequireImportSection:                   contents.Forbidden,
				RequireSchemaOrdering:                  config.RequireSchemaOrdering,
				IgnoreContentsCheck:                    ignoreContentsCheckEphemerals,
				IgnoreEnhancedRegionCheck:              ignoreEnhancedRegionCheckEphemerals,
				IgnoreEnhancedRegionCheckSubcategories: ignoreEnhancedRegionCheckSubcategories,
				ProviderName:                           config.ProviderName,
				TitleSectionPrefixes:                   []string{"Ephemeral"},
			},
			FileOptions: fileOpts,
			FrontMatter: &check.FrontMatterOptions{
				AllowedSubcategories: allowedResourceSubcategories,
				RequireSubcategory:   config.RequireResourceSubcategory,
			},
			ProviderName: config.ProviderName,
		},
		EphemeralFileMismatch: &check.FileMismatchOptions{
			IgnoreFileMismatch: ignoreFileMismatchEphemerals,
			IgnoreFileMissing:  ignoreFileMissingEphemerals,
			ProviderName:       config.ProviderName,
			ResourceType:       check.ResourceTypeEphemeral,
			ResourceNames:      ephemeralNames,
		},

		// function
		RegistryFunctionFile: &check.RegistryFunctionFileOptions{
			Contents: &check.ContentsOptions{
				Enable:                      config.EnableContentsCheck,
				IgnoreContentsCheck:         ignoreContentsCheckFunctions,
				ProviderName:                config.ProviderName,
				TitleSectionPrefixes:        []string{"Function"},
				RequireImportSection:        contents.Forbidden,
				ArgumentsHeadingTexts:       []string{"Arguments"},
				AllowArgumentsMissingByline: true,
				ArgumentsBylineTexts: []string{
					"This function supports the following arguments:",
					"This function does not support any arguments.",
				},
				RequireSignatureSection:    contents.Required,
				SignatureHeadingTexts:      []string{"Signature"},
				SignatureRequiresCodeBlock: true,
			},
			FileOptions: fileOpts,
			FrontMatter: &check.FrontMatterOptions{
				AllowedSubcategories: allowedResourceSubcategories,
				RequireSubcategory:   config.RequireResourceSubcategory,
			},
			ProviderName: config.ProviderName,
		},
		LegacyFunctionFile: &check.LegacyFunctionFileOptions{
			Contents: &check.ContentsOptions{
				Enable:                      config.EnableContentsCheck,
				IgnoreContentsCheck:         ignoreContentsCheckFunctions,
				ProviderName:                config.ProviderName,
				TitleSectionPrefixes:        []string{"Function"},
				RequireImportSection:        contents.Forbidden,
				ArgumentsHeadingTexts:       []string{"Arguments"},
				AllowArgumentsMissingByline: true,
				ArgumentsBylineTexts: []string{
					"This function supports the following arguments:",
					"This function does not support any arguments.",
				},
				RequireSignatureSection:    contents.Required,
				SignatureHeadingTexts:      []string{"Signature"},
				SignatureRequiresCodeBlock: true,
			},
			FileOptions: fileOpts,
			FrontMatter: &check.FrontMatterOptions{
				AllowedSubcategories: allowedResourceSubcategories,
				RequireSubcategory:   config.RequireResourceSubcategory,
			},
			ProviderName: config.ProviderName,
		},
		FunctionFileMismatch: &check.FileMismatchOptions{
			IgnoreFileMismatch: ignoreFileMismatchFunctions,
			IgnoreFileMissing:  ignoreFileMissingFunctions,
			ResourceType:       check.ResourceTypeFunction,
			ResourceNames:      functionNames,
		},

		// list resource
		RegistryListResourceFile: &check.RegistryListResourceFileOptions{
			Contents: &check.ContentsOptions{
				Enable:                                 config.EnableContentsCheck,
				EnhancedRegionChecks:                   config.EnableEnhancedRegionCheck,
				RequireAttributesSection:               contents.Forbidden,
				RequireTimeoutsSection:                 contents.Forbidden,
				RequireImportSection:                   contents.Forbidden,
				RequireSchemaOrdering:                  true,
				IgnoreContentsCheck:                    ignoreContentsCheckResources,
				IgnoreEnhancedRegionCheck:              ignoreEnhancedRegionCheckResources,
				IgnoreEnhancedRegionCheckSubcategories: ignoreEnhancedRegionCheckSubcategories,
				ProviderName:                           config.ProviderName,
				TitleSectionPrefixes:                   []string{"List Resource"},
			},
			FileOptions: fileOpts,
			FrontMatter: &check.FrontMatterOptions{
				AllowedSubcategories: allowedResourceSubcategories,
				RequireSubcategory:   config.RequireResourceSubcategory,
			},
			ProviderName: config.ProviderName,
		},
		LegacyListResourceFile: &check.LegacyListResourceFileOptions{
			Contents: &check.ContentsOptions{
				Enable:                                 config.EnableContentsCheck,
				EnhancedRegionChecks:                   config.EnableEnhancedRegionCheck,
				RequireAttributesSection:               contents.Forbidden,
				RequireTimeoutsSection:                 contents.Forbidden,
				RequireImportSection:                   contents.Forbidden,
				RequireSchemaOrdering:                  true,
				IgnoreContentsCheck:                    ignoreContentsCheckResources,
				IgnoreEnhancedRegionCheck:              ignoreEnhancedRegionCheckResources,
				IgnoreEnhancedRegionCheckSubcategories: ignoreEnhancedRegionCheckSubcategories,
				ProviderName:                           config.ProviderName,
				TitleSectionPrefixes:                   []string{"List Resource"},
			},
			FileOptions: fileOpts,
			FrontMatter: &check.FrontMatterOptions{
				AllowedSubcategories: allowedResourceSubcategories,
				RequireSubcategory:   config.RequireResourceSubcategory,
			},
			ProviderName: config.ProviderName,
		},
		ListResourceFileMismatch: &check.FileMismatchOptions{
			IgnoreFileMismatch: ignoreFileMismatchListResources,
			IgnoreFileMissing:  ignoreFileMissingListResources,
			ProviderName:       config.ProviderName,
			ResourceType:       check.ResourceTypeListResource,
			ResourceNames:      listResourceNames,
		},

		// resource
		RegistryResourceFile: &check.RegistryResourceFileOptions{
			Contents: &check.ContentsOptions{
				Enable:                                 config.EnableContentsCheck,
				EnhancedRegionChecks:                   config.EnableEnhancedRegionCheck,
				RequireAttributesSection:               contents.Required,
				RequireImportSection:                   contents.Optional,
				RequireSchemaOrdering:                  config.RequireSchemaOrdering,
				IgnoreContentsCheck:                    ignoreContentsCheckResources,
				IgnoreEnhancedRegionCheck:              ignoreEnhancedRegionCheckResources,
				IgnoreEnhancedRegionCheckSubcategories: ignoreEnhancedRegionCheckSubcategories,
				ProviderName:                           config.ProviderName,
				TitleSectionPrefixes:                   []string{"Resource"},
			},
			FileOptions: fileOpts,
			FrontMatter: &check.FrontMatterOptions{
				AllowedSubcategories: allowedResourceSubcategories,
				RequireSubcategory:   config.RequireResourceSubcategory,
			},
			ProviderName: config.ProviderName,
		},
		LegacyResourceFile: &check.LegacyResourceFileOptions{
			Contents: &check.ContentsOptions{
				Enable:                                 config.EnableContentsCheck,
				EnhancedRegionChecks:                   config.EnableEnhancedRegionCheck,
				RequireAttributesSection:               contents.Required,
				RequireImportSection:                   contents.Optional,
				RequireSchemaOrdering:                  config.RequireSchemaOrdering,
				IgnoreContentsCheck:                    ignoreContentsCheckResources,
				IgnoreEnhancedRegionCheck:              ignoreEnhancedRegionCheckResources,
				IgnoreEnhancedRegionCheckSubcategories: ignoreEnhancedRegionCheckSubcategories,
				ProviderName:                           config.ProviderName,
				TitleSectionPrefixes:                   []string{"Resource"},
			},
			FileOptions: fileOpts,
			FrontMatter: &check.FrontMatterOptions{
				AllowedSubcategories: allowedResourceSubcategories,
				RequireSubcategory:   config.RequireResourceSubcategory,
			},
			ProviderName: config.ProviderName,
		},
		ResourceFileMismatch: &check.FileMismatchOptions{
			IgnoreFileMismatch: ignoreFileMismatchResources,
			IgnoreFileMissing:  ignoreFileMissingResources,
			ProviderName:       config.ProviderName,
			ResourceType:       check.ResourceTypeResource,
			ResourceNames:      resourceNames,
		},

		// guide
		RegistryGuideFile: &check.RegistryGuideFileOptions{
			FileOptions: fileOpts,
			FrontMatter: &check.FrontMatterOptions{
				AllowedSubcategories: allowedGuideSubcategories,
				RequireSubcategory:   config.RequireGuideSubcategory,
			},
		},
		LegacyGuideFile: &check.LegacyGuideFileOptions{
			FileOptions: fileOpts,
			FrontMatter: &check.FrontMatterOptions{
				AllowedSubcategories: allowedGuideSubcategories,
				RequireSubcategory:   config.RequireGuideSubcategory,
			},
		},

		// index
		RegistryIndexFile: &check.RegistryIndexFileOptions{
			FileOptions: fileOpts,
		},
		LegacyIndexFile: &check.LegacyIndexFileOptions{
			FileOptions: fileOpts,
		},

		// general
		ProviderName:            config.ProviderName,
		ProviderSource:          config.ProviderSource,
		IgnoreCdktfMissingFiles: config.IgnoreCdktfMissingFiles,
	}

	if err := check.NewCheck(checkOpts).Run(directories); err != nil {
		c.Ui.Error(fmt.Sprintf("Error checking Terraform Provider documentation: %s", err))
		return 1
	}

	return 0
}

func (c *CheckCommand) Synopsis() string {
	return "Checks Terraform Provider documentation"
}

func allowedSubcategoriesFile(path string) ([]string, error) {
	log.Printf("[DEBUG] Loading allowed subcategories file: %s", path)

	file, err := os.Open(path)

	if err != nil {
		return nil, fmt.Errorf("error opening allowed subcategories file (%s): %w", path, err)
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)
	var allowedSubcategories []string

	for scanner.Scan() {
		allowedSubcategories = append(allowedSubcategories, scanner.Text())
	}

	if err != nil {
		return nil, fmt.Errorf("error reading allowed subcategories file (%s): %w", path, err)
	}

	return allowedSubcategories, nil
}

func providerNameFromCurrentDirectory() string {
	path, _ := os.Getwd()

	return providerNameFromPath(path)
}

func providerNameFromPath(path string) string {
	base := filepath.Base(path)

	if strings.ContainsAny(base, "./") {
		return ""
	}

	if !strings.HasPrefix(base, "terraform-provider-") {
		return ""
	}

	return strings.TrimPrefix(base, "terraform-provider-")
}

// providerSchemas reads, parses, and validates a provided terraform provider schema -json path.
func providerSchemas(path string) (*tfjson.ProviderSchemas, error) {
	log.Printf("[DEBUG] Loading providers schema JSON file: %s", path)

	content, err := os.ReadFile(path)

	if err != nil {
		return nil, fmt.Errorf("error reading providers schema JSON file (%s): %w", path, err)
	}

	var ps tfjson.ProviderSchemas

	if err := json.Unmarshal(content, &ps); err != nil {
		return nil, fmt.Errorf("error parsing providers schema JSON file (%s): %w", path, err)
	}

	if err := ps.Validate(); err != nil {
		return nil, fmt.Errorf("error validating providers schema JSON file (%s): %w", path, err)
	}

	return &ps, nil
}

// providerSchemasDataSources returns all data source names from a terraform providers schema -json provider.
func providerSchemasDataSources(ps *tfjson.ProviderSchemas, providerName string, providerSource string) []string {
	if ps == nil || ps.Schemas == nil {
		return nil
	}

	provider, ok := ps.Schemas[providerSource]

	if !ok {
		provider, ok = ps.Schemas[providerName]
	}

	if !ok {
		log.Printf("[WARN] Provider source (%s) and name (%s) not found in provider schema", providerSource, providerName)
		return nil
	}

	dataSources := make([]string, 0, len(provider.DataSourceSchemas))

	for name := range provider.DataSourceSchemas {
		dataSources = append(dataSources, name)
	}

	sort.Strings(dataSources)

	log.Printf("[DEBUG] Found provider schema data sources: %v", dataSources)

	return dataSources
}

// providerSchemasActions returns all action names from a terraform providers schema -json provider.
func providerSchemasActions(ps *tfjson.ProviderSchemas, providerName string, providerSource string) []string {
	if ps == nil || ps.Schemas == nil {
		return nil
	}

	provider, ok := ps.Schemas[providerSource]

	if !ok {
		provider, ok = ps.Schemas[providerName]
	}

	if !ok {
		log.Printf("[WARN] Provider source (%s) and name (%s) not found in provider schema", providerSource, providerName)
		return nil
	}

	actions := make([]string, 0, len(provider.ActionSchemas))

	for name := range provider.ActionSchemas {
		actions = append(actions, name)
	}

	sort.Strings(actions)

	log.Printf("[DEBUG] Found provider schema actions: %v", actions)

	return actions
}

// providerSchemasEphemerals returns all ephemeral names from a terraform providers schema -json provider.
func providerSchemasEphemerals(ps *tfjson.ProviderSchemas, providerName string, providerSource string) []string {
	if ps == nil || ps.Schemas == nil {
		return nil
	}

	provider, ok := ps.Schemas[providerSource]

	if !ok {
		provider, ok = ps.Schemas[providerName]
	}

	if !ok {
		log.Printf("[WARN] Provider source (%s) and name (%s) not found in provider schema", providerSource, providerName)
		return nil
	}

	ephemerals := make([]string, 0, len(provider.EphemeralResourceSchemas))

	for name := range provider.EphemeralResourceSchemas {
		ephemerals = append(ephemerals, name)
	}

	sort.Strings(ephemerals)

	log.Printf("[DEBUG] Found provider schema ephemerals: %v", ephemerals)

	return ephemerals
}

// providerSchemasFunctions returns all function names from a terraform providers schema -json provider.
func providerSchemasFunctions(ps *tfjson.ProviderSchemas, providerName string, providerSource string) []string {
	if ps == nil || ps.Schemas == nil {
		return nil
	}

	provider, ok := ps.Schemas[providerSource]

	if !ok {
		provider, ok = ps.Schemas[providerName]
	}

	if !ok {
		log.Printf("[WARN] Provider source (%s) and name (%s) not found in provider schema", providerSource, providerName)
		return nil
	}

	functions := make([]string, 0, len(provider.Functions))

	for name := range provider.Functions {
		functions = append(functions, name)
	}

	sort.Strings(functions)

	log.Printf("[DEBUG] Found provider schema functions: %v", functions)

	return functions
}

// providerSchemasListResources returns all list resource names from a terraform providers schema -json provider.
func providerSchemasListResources(ps *tfjson.ProviderSchemas, providerName string, providerSource string) []string {
	if ps == nil || ps.Schemas == nil {
		return nil
	}

	provider, ok := ps.Schemas[providerSource]

	if !ok {
		provider, ok = ps.Schemas[providerName]
	}

	if !ok {
		log.Printf("[WARN] Provider source (%s) and name (%s) not found in provider schema", providerSource, providerName)
		return nil
	}

	listResources := make([]string, 0, len(provider.ListResourceSchemas))

	for name := range provider.ListResourceSchemas {
		listResources = append(listResources, name)
	}

	sort.Strings(listResources)

	log.Printf("[DEBUG] Found provider schema list resources: %v", listResources)

	return listResources
}

// providerSchemasResources returns all resource names from a terraform providers schema -json provider.
func providerSchemasResources(ps *tfjson.ProviderSchemas, providerName string, providerSource string) []string {
	if ps == nil || ps.Schemas == nil {
		return nil
	}

	provider, ok := ps.Schemas[providerSource]

	if !ok {
		provider, ok = ps.Schemas[providerName]
	}

	if !ok {
		log.Printf("[WARN] Provider source (%s) and name (%s) not found in provider schema", providerSource, providerName)
		return nil
	}

	resources := make([]string, 0, len(provider.ResourceSchemas))

	for name := range provider.ResourceSchemas {
		resources = append(resources, name)
	}

	sort.Strings(resources)

	log.Printf("[DEBUG] Found provider schema resources: %v", resources)

	return resources
}
