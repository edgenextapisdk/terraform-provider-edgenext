Use this resource to attach a subnet to an ECS router.

Example Usage

```hcl
data "edgenext_ecs_external_gateways" "all" {
  limit = 1
}

resource "edgenext_ecs_router" "example" {
  name                = "example-router"
  external_network_id = data.edgenext_ecs_external_gateways.all.external_gateways[0].id
}

data "edgenext_ecs_vpcs" "all" {
  limit = 1
}

data "edgenext_ecs_vpc_subnets" "all" {
  vpc_id = data.edgenext_ecs_vpcs.all.vpcs[0].id
}

resource "edgenext_ecs_router_port" "example" {
  router_id = edgenext_ecs_router.example.id
  vpc_id    = data.edgenext_ecs_vpcs.all.vpcs[0].id
  subnet_id = data.edgenext_ecs_vpc_subnets.all.subnets[0].id
}
```

Import

Import format is `router_id/port_id`.

```shell
terraform import edgenext_ecs_router_port.example f9883769-xxxx-xxxx-xxxx-xxxxxxxxxxxx/74f3a422-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```

Argument Reference

* `router_id` - (Required) Router ID. Cannot be changed after creation.
* `vpc_id` - (Required) VPC ID. Cannot be changed after creation.
* `subnet_id` - (Required) Subnet ID. Cannot be changed after creation.

Attributes Reference

* `id` - Router port ID.
* `port_id` - Same as router port ID.
* `name`, `ip_address`, `mac_address`, `vpc_name`, `status`, `created_at`
