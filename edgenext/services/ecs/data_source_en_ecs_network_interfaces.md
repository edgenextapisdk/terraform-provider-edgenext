Use this data source to query ECS network interfaces (ports).

Example Usage

```hcl
data "edgenext_ecs_network_interfaces" "example" {
  network_interface_name = edgenext_ecs_network_interface.example.name
  limit                  = 10
}

resource "edgenext_ecs_network_interface" "example" {
  name      = "example-eni"
  vpc_id    = data.edgenext_ecs_vpcs.all.vpcs[0].id
  subnet_id = data.edgenext_ecs_vpc_subnets.all.subnets[0].id
}

data "edgenext_ecs_vpcs" "all" {
  limit = 1
}

data "edgenext_ecs_vpc_subnets" "all" {
  vpc_id = data.edgenext_ecs_vpcs.all.vpcs[0].id
}
```
