Use this resource to create and manage ECS key pairs.

Example Usage

```hcl
resource "edgenext_ecs_key_pair" "example" {
  name       = "example-key"
  public_key = file("~/.ssh/id_rsa.pub")
}
```

Import

Import format is `name`.

```shell
terraform import edgenext_ecs_key_pair.example example-key
```
