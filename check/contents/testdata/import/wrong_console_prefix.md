<!-- Copyright (c) HashiCorp, Inc. 2019-2026 -->
<!-- SPDX-License-Identifier: MPL-2.0 -->

## Import

Import Passings using the `name`. For example:

```terraform
import {
  to = test_wrong_console_prefix.example
  id = "example"
}
```

```console
$ terraform import test_wrong_console_prefix.example example
```