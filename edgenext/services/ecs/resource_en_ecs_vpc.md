# en_ecs_vpc

Provides an EdgeNext ECS vpc resource. This allows you to manage vpcs within your ECS environment.

## Example Usage

```hcl
resource "edgenext_ecs_vpc" "example" {
  name = "example-vpc"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the vpc.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the vpc.
