# en_ecs_key_pair

Provides an EdgeNext ECS key_pair resource. This allows you to manage key_pairs within your ECS environment.

## Example Usage

```hcl
resource "edgenext_ecs_key_pair" "example" {
  name = "example-key_pair"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the key_pair.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the key_pair.
