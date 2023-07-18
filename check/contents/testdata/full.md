---
subcategory: "Test Full"
layout: "test"
page_title: "Test: test_full"
description: |-
  Manages a Test Full
---

# Resource: test_full

Manages a Test Full.

## Example Usage

```terraform
resource "test_full" "example" {
  name = "example"
}
```

## Argument Reference

This resource supports the following arguments:

* `name` - (Required) Name of thing.
* `tags` - (Optional) Key-value map of resource tags.
* `type` - (Optional) Type of thing.

## Attribute Reference

This resource exports the following attributes in addition to the arguments above:

* `id` - Name of thing.

## Timeouts

`test_full` provides the following [Timeouts](/docs/configuration/resources.html#timeouts)
configuration options:

* `create` - (Default `10m`) How long to wait for the thing to be created.

## Import

Import Fulls using `name`. For example:

```
$ terraform import test_full.example example
```
