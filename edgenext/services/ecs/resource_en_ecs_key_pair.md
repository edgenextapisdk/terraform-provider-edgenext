Use this resource to create and manage ECS key pairs.

Example Usage

```hcl
resource "edgenext_ecs_key_pair" "example" {
  region     = "tokyo-a"
  name       = "example-key"
  public_key = file("~/.ssh/id_rsa.pub")
}
```

Import

Import format is `region/name`.

```shell
terraform import edgenext_ecs_key_pair.example tokyo-a/example-key
```

Argument Reference

* `region` - (Required) Region.
* `name` - (Required) Key pair name.
* `public_key` - (Optional) Public key content.

Attributes Reference

* `id` - Key pair name (provider ID).
* `private_key` - Private key returned by API (if generated server side).
