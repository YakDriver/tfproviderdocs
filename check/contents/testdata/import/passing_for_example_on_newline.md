<!-- Copyright IBM Corp. 2019, 2026 -->
<!-- SPDX-License-Identifier: MPL-2.0 -->

## Import

Import Passings using the `name`.
For example:

```terraform
import {
  to = test_passing_for_example_on_newline.example
  id = "example"
}
```

```console
% terraform import test_passing_for_example_on_newline.example example
```
