<!-- Copyright (c) IBM Corp. 2019-2026 -->
<!-- SPDX-License-Identifier: MPL-2.0 -->

<!-- Copyright (c) HashiCorp, Inc. 2019-2026 -->
## Import

Import Passings using the `name`. For example:

```terraform
import {
  to = test_passing.example
  id = "example"
}
```

```console
% terraform import test_passing.example example
```