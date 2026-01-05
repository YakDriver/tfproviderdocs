<!-- Copyright IBM Corp. 2019, 2026 -->
<!-- SPDX-License-Identifier: MPL-2.0 -->

---
subcategory: "Example"
layout: "example"
page_title: "Example: example_thing"
description: |-
  Example description.
---

# Ephemeral: example_thing

Byline.

## Example Usage

```terraform
ephemeral "example_thing" "example" {
  name = "example"
}
```

## Argument Reference

* `name` - (Required) Name of thing.

## Attribute Reference

* `id` - Name of thing.
