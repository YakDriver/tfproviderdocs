<!-- Copyright IBM Corp. 2019, 2026 -->
<!-- SPDX-License-Identifier: MPL-2.0 -->

## Import

Import Passings using the `name`. For example:

```terraform
import {
  to = test_wrong_code_block_order.example
  id = "example"
}
```

```console
% terraform import test_wrong_code_block_order.example example
```

```terraform
import {
  to = test_wrong_code_block_order.example
  id = "example"
}
```
