---
subcategory: "Test Action"
layout: "test"
page_title: "Test: test_action"
description: |-
  Manages a Test Action
---

# Action: test_action

Manages a Test Action.

## Example Usage

```terraform
action "test_action" "example" {
  config {
    example = "value"
  }
}
```

## Argument Reference

This action supports the following arguments:

* `example` - (Required) Example argument.
