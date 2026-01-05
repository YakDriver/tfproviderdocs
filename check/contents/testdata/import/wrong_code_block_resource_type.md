<!-- Copyright (c) IBM Corp. 2019-2026 -->
<!-- SPDX-License-Identifier: MPL-2.0 -->

<!-- Copyright (c) HashiCorp, Inc. 2019-2026 -->
## Import

Import Test Wrong Code Block Resource Types using the `name`. For example:

```terraform
import {
  to = test_wrong_code_block_resource_type.example
  id = "example"
}
```

```
$ terraform import test_wrong_code_block_resource_type.example example
```
