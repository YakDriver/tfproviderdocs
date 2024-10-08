# tfproviderdocs

A documentation tool for [Terraform Provider](https://www.terraform.io/docs/providers/index.html) code.

## Install

### Local Install

Release binaries are available in the [Releases](https://github.com/YakDriver/tfproviderdocs/releases) section.

To instead use Go to install into your `$GOBIN` directory (e.g. `$GOPATH/bin`):

```shell
go install github.com/YakDriver/tfproviderdocs
```

### Docker Install

```shell
docker pull moab4x/tfproviderdocs
```

### Homebrew Install

```shell
brew install YakDriver/tap/tfproviderdocs
```

## Usage

Additional information about usage and configuration options can be found by passing the `help` argument:

```shell
tfproviderdocs help
```

### Local Usage

Change into the directory of the Terraform Provider code and run:

```shell
tfproviderdocs
```

### Docker Usage

Change into the directory of the Terraform Provider code and run:

```shell
docker run -v $(pwd):/src moab4x/tfproviderdocs
```

## Available Commands

### check Command

The `tfproviderdocs check` command verifies the Terraform Provider documentation against the [specifications from Terraform Registry documentation](https://www.terraform.io/docs/registry/providers/docs.html) and common practices across official Terraform Providers. This includes the following checks:

- Verifies that no invalid directories are found in the documentation directory structure.
- Ensures that there is not a mix (legacy and Terraform Registry) of directory structures, which is not supported during Terraform Registry documentation ingress.
- Verifies number of documentation files is below Terraform Registry storage limits.
- Verifies all known data sources and resources have an associated documentation file (if `-providers-schema-json` is provided)
- Verifies no extraneous or incorrectly named documentation files exist (if `-providers-schema-json` is provided)
- Verifies each file in the documentation directories is valid.

The validity of files is checked with the following rules:

- Proper file extensions are used (e.g. `.md` for Terraform Registry).
- Verifies size of file is below Terraform Registry storage limits.
- YAML frontmatter can be parsed and matches expectations.

The YAML frontmatter checks include some defaults (e.g. no `layout` field for Terraform Registry), but there are some useful flags that can be passed to the command to tune the behavior, especially for larger Terraform Providers.

The validity of files can also be experimentally checked (via the `-enable-contents-check` flag) with the following rules:

- Ensures all expected headings are present.
- Verifies heading levels and text.
- Verifies schema attribute lists are ordered (if `-require-schema-ordering` is provided). Only supports section level lists (not sub-section level lists) currently.
- Verifies resource type is present in code blocks (e.g. examples and import sections).

For additional information about check flags, you can run `tfproviderdocs check -help`.

## Development and Testing

This project uses [Go Modules](https://github.com/golang/go/wiki/Modules) for dependency management.

### Go Compatibility

This project follows the [Go support policy](https://golang.org/doc/devel/release.html#policy) for versions. The two latest major releases of Go are supported by the project.

Currently, that means Go **1.19** or later must be used when including this project as a dependency.

### Updating Dependencies

Dependency updates are managed via Dependabot.

### Unit Testing

```shell
go test ./...
```

### Local Install Testing

```shell
go install .
```
