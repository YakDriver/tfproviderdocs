# Repository Guidelines
<!-- Copyright (c) HashiCorp, Inc. 2019-2026 -->
<!-- SPDX-License-Identifier: MPL-2.0 -->

## Project Purpose & Usage
`tfproviderdocs` operates as a CI/CD safeguard for Terraform providers, validating user-facing docs for required structure, headings, and metadata. When supplied a provider schema in JSON (via `terraform providers schema -json`), it cross-references expected data sources and resources so engineers do not hand-maintain checklists. Providers such as AWS invoke the tool within release pipelines, for example:

```sh
echo 'data "aws_partition" "example" {}' > example.tf
terraform init -plugin-dir terraform-plugin-dir
mkdir -p terraform-providers-schema
terraform providers schema -json > terraform-providers-schema/schema.json
tfproviderdocs check \
  -allowed-resource-subcategories-file website/allowed-subcategories.txt \
  -enable-contents-check \
  -ignore-contents-check-data-sources aws_kms_secrets,aws_kms_secret \
  -ignore-file-missing-data-sources aws_alb,aws_alb_listener,aws_alb_target_group,aws_alb_trust_store,aws_alb_trust_store_revocation,aws_albs \
  -ignore-file-missing-resources aws_alb,aws_alb_listener,aws_alb_listener_certificate,aws_alb_listener_rule,aws_alb_target_group,aws_alb_target_group_attachment,aws_alb_trust_store,aws_alb_trust_store_revocation \
  -provider-source registry.terraform.io/hashicorp/aws \
  -providers-schema-json terraform-providers-schema/schema.json \
  -require-resource-subcategory \
  -ignore-cdktf-missing-files \
  -ignore-enhanced-region-check-subcategories-file website/ignore-enhanced-region-check-subcategories.txt \
  -ignore-enhanced-region-check-data-sources-file website/ignore-enhanced-region-check-data-sources.txt \
  -ignore-enhanced-region-check-resources-file website/ignore-enhanced-region-check-resources.txt \
  -enable-enhanced-region-check
```

## Project Structure & Module Organization
Core CLI wiring sits in `main.go`; command orchestration and fixtures live in `command/` and `command/testdata/`. The validation engine is in `check/` with matching `*_test.go` suites and fixture files under `check/testdata/`. Markdown helpers live in `markdown/`, version metadata in `version/`, and automation files (`GNUmakefile`, `.github/workflows/`) support builds and releases.

## Build, Test, and Development Commands
Run `go build ./...` for a fast compilation check and `make` to install into `$GOBIN`. Execute `go test ./...` before opening a PR, or narrow with `go test ./check/...` while iterating. Keep modules clean with `go mod tidy`; CI fails if `go.mod` or `go.sum` drift.

## Coding Style & Naming Conventions
Always commit `gofmt`-formatted Go (tabs, same-line braces). Favor domain-specific identifiers such as `registryResourceFile`. Exported symbols need doc comments, and multi-case assertions should use table-driven tests. Discuss any third-party dependency additions in review.

## Testing Guidelines
Pair every behavior change with package-level unit tests, storing fixtures in existing `testdata/` folders. Name tests `TestFunction_Scenario` and label table rows clearly. When extending documentation checks, mirror success and failure cases already covered in `check/` suites. Use `go test -run Name ./package` to focus on a single scenario.

## Commit & Pull Request Guidelines
Write concise, imperative commit subjects (e.g., `Add registry list resource guard`) and add body context when the diff warrants it. Squash noisy fixups before opening the PR. Describe the change, list executed tests (such as `go test ./...`), link issues, and attach CLI output or screenshots when UX or generated docs shift.
