## Import

Import Passings using the `name`. For example:

```terraform
import {
  to = test_passing_code_block_after.example
  id = "example"
}
```

```console
% terraform import test_passing_code_block_after.example example
```

Certain resource arguments, like `encryption_configuration` and `bucket`, do not have an API method for reading the information after creation. If the argument is set in the Terraform configuration on an imported resource, Terraform will always show a difference. To workaround this behavior, either omit the argument from the Terraform configuration or use [`ignore_changes`](https://www.terraform.io/docs/configuration/meta-arguments/lifecycle.html#ignore_changes) to hide the difference. For example:

```terraform
resource "aws_athena_database" "example" {
  name   = "database_name"
  bucket = test_passing_code_block_after.example.id

  # There is no API for reading bucket
  lifecycle {
    ignore_changes = [bucket]
  }
}
```