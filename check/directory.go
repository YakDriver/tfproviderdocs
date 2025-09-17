package check

import (
	"fmt"
	"path/filepath"
	"slices"

	"github.com/bmatcuk/doublestar"
)

const (
	CdktfIndexDirectory = `cdktf`

	DocumentationGlobPattern = `{docs/index.md,docs/{,cdktf/}{actions,data-sources,ephemeral-resources,functions,guides,list-resources,resources},website/docs}/**/*`

	LegacyIndexDirectory         = `website/docs`
	LegacyActionsDirectory       = `actions`
	LegacyDataSourcesDirectory   = `d`
	LegacyEphemeralsDirectory    = `ephemeral-resources`
	LegacyFunctionsDirectory     = `functions`
	LegacyGuidesDirectory        = `guides`
	LegacyListResourcesDirectory = `list-resources`
	LegacyResourcesDirectory     = `r`

	RegistryIndexDirectory         = `docs`
	RegistryActionsDirectory       = `actions`
	RegistryDataSourcesDirectory   = `data-sources`
	RegistryEphemeralsDirectory    = `ephemeral-resources`
	RegistryFunctionsDirectory     = `functions`
	RegistryGuidesDirectory        = `guides`
	RegistryListResourcesDirectory = `list-resources`
	RegistryResourcesDirectory     = `resources`
)

var ValidLegacyDirectories = []string{
	LegacyIndexDirectory,
	LegacyIndexDirectory + "/" + LegacyActionsDirectory,
	LegacyIndexDirectory + "/" + LegacyDataSourcesDirectory,
	LegacyIndexDirectory + "/" + LegacyEphemeralsDirectory,
	LegacyIndexDirectory + "/" + LegacyFunctionsDirectory,
	LegacyIndexDirectory + "/" + LegacyGuidesDirectory,
	LegacyIndexDirectory + "/" + LegacyListResourcesDirectory,
	LegacyIndexDirectory + "/" + LegacyResourcesDirectory,
}

var ValidRegistryDirectories = []string{
	RegistryIndexDirectory,
	RegistryIndexDirectory + "/" + RegistryActionsDirectory,
	RegistryIndexDirectory + "/" + RegistryDataSourcesDirectory,
	RegistryIndexDirectory + "/" + RegistryEphemeralsDirectory,
	RegistryIndexDirectory + "/" + RegistryFunctionsDirectory,
	RegistryIndexDirectory + "/" + RegistryGuidesDirectory,
	RegistryIndexDirectory + "/" + RegistryListResourcesDirectory,
	RegistryIndexDirectory + "/" + RegistryResourcesDirectory,
}

var ValidCdktfLanguages = []string{
	"csharp",
	"go",
	"java",
	"python",
	"typescript",
}

var ValidLegacySubdirectories = []string{
	LegacyActionsDirectory,
	LegacyDataSourcesDirectory,
	LegacyEphemeralsDirectory,
	LegacyFunctionsDirectory,
	LegacyGuidesDirectory,
	LegacyListResourcesDirectory,
	LegacyResourcesDirectory,
}

var ValidRegistrySubdirectories = []string{
	RegistryActionsDirectory,
	RegistryDataSourcesDirectory,
	RegistryEphemeralsDirectory,
	RegistryFunctionsDirectory,
	RegistryGuidesDirectory,
	RegistryListResourcesDirectory,
	RegistryResourcesDirectory,
}

func InvalidDirectoriesCheck(directories map[string][]string) error {
	for directory := range directories {
		if IsValidRegistryDirectory(directory) {
			continue
		}

		if IsValidLegacyDirectory(directory) {
			continue
		}

		if IsValidCdktfDirectory(directory) {
			continue
		}

		return fmt.Errorf("invalid Terraform Provider documentation directory found: %s", directory)
	}

	return nil
}

func MixedDirectoriesCheck(directories map[string][]string) error {
	var legacyDirectoryFound bool
	var registryDirectoryFound bool
	err := fmt.Errorf("mixed Terraform Provider documentation directory layouts found, must use only legacy or registry layout")

	for directory := range directories {
		// Allow docs/ with other files
		if IsValidRegistryDirectory(directory) && directory != RegistryIndexDirectory {
			registryDirectoryFound = true

			if legacyDirectoryFound {
				return err
			}
		}

		if IsValidLegacyDirectory(directory) {
			legacyDirectoryFound = true

			if registryDirectoryFound {
				return err
			}
		}
	}

	return nil
}

func GetDirectories(basepath string) (map[string][]string, error) {
	globPattern := DocumentationGlobPattern

	if basepath != "" {
		globPattern = fmt.Sprintf("%s/%s", basepath, globPattern)
	}

	files, err := doublestar.Glob(globPattern)

	if err != nil {
		return nil, fmt.Errorf("error globbing Terraform Provider documentation directories: %w", err)
	}

	if basepath != "" {
		for index, file := range files {
			files[index], _ = filepath.Rel(basepath, file)
		}
	}

	directories := make(map[string][]string)

	for _, file := range files {
		// Simple skip of glob matches that are known directories
		if IsValidRegistryDirectory(file) || IsValidLegacyDirectory(file) || IsValidCdktfDirectory(file) {
			continue
		}

		if filepath.Base(file) == ".keep" {
			continue
		}

		directory := filepath.Dir(file)

		// Skip handling of docs/ files except index.md
		// if directory == RegistryIndexDirectory && filepath.Base(file) != "index.md" {
		// 	continue
		// }

		// Skip handling of docs/** outside valid Registry Directories

		directories[directory] = append(directories[directory], file)
	}

	return directories, nil
}

func IsValidLegacyDirectory(directory string) bool {
	return slices.Contains(ValidLegacyDirectories, directory)
}

func IsValidRegistryDirectory(directory string) bool {
	return slices.Contains(ValidRegistryDirectories, directory)
}

func IsValidCdktfDirectory(directory string) bool {
	if directory == fmt.Sprintf("%s/%s", LegacyIndexDirectory, CdktfIndexDirectory) {
		return true
	}

	if directory == fmt.Sprintf("%s/%s", RegistryIndexDirectory, CdktfIndexDirectory) {
		return true
	}

	for _, validCdktfLanguage := range ValidCdktfLanguages {

		if directory == fmt.Sprintf("%s/%s/%s", LegacyIndexDirectory, CdktfIndexDirectory, validCdktfLanguage) {
			return true
		}

		if directory == fmt.Sprintf("%s/%s/%s", RegistryIndexDirectory, CdktfIndexDirectory, validCdktfLanguage) {
			return true
		}

		for _, validLegacySubdirectory := range ValidLegacySubdirectories {
			if directory == fmt.Sprintf("%s/%s/%s/%s", LegacyIndexDirectory, CdktfIndexDirectory, validCdktfLanguage, validLegacySubdirectory) {
				return true
			}
		}

		for _, validRegistrySubdirectory := range ValidRegistrySubdirectories {
			if directory == fmt.Sprintf("%s/%s/%s/%s", RegistryIndexDirectory, CdktfIndexDirectory, validCdktfLanguage, validRegistrySubdirectory) {
				return true
			}
		}
	}

	return false
}
