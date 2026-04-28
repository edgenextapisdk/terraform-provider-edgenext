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

Argument Reference

* `name` - (Required) Key pair name. Cannot be changed after creation.
* `public_key` - (Optional) Public key content. Cannot be changed after creation.

Attributes Reference

* `id` - Key pair name (provider ID).
* `private_key` - Private key returned by API (if generated server side).
* `key_id` - Numeric key pair ID returned by API field `id`.
* `fingerprint` - Public key fingerprint.
* `created_at`, `updated_at`, `deleted`, `deleted_at`, `user_id`
