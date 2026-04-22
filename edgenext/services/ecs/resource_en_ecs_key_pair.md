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
* `name` - (Required) Key pair name. Cannot be changed after creation.
* `public_key` - (Optional) Public key content. Cannot be changed after creation.

Attributes Reference

* `id` - Key pair name (provider ID).
* `private_key` - Private key returned by API (if generated server side).
* `key_id` - Numeric key pair ID returned by API field `id`.
* `fingerprint` - Public key fingerprint.
* `created_at`, `updated_at`, `deleted`, `deleted_at`, `user_id`
