Use this resource to create and manage ECS network interfaces (ENI / port).

Example Usage

```hcl
data "edgenext_ecs_vpcs" "all" {
  limit = 1
}

data "edgenext_ecs_vpc_subnets" "all" {
  vpc_id = data.edgenext_ecs_vpcs.all.vpcs[0].id
}

data "edgenext_ecs_security_groups" "all" {
  limit = 1
}

resource "edgenext_ecs_network_interface" "example" {
  name                  = "example-eni"
  description           = "for application"
  vpc_id                = data.edgenext_ecs_vpcs.all.vpcs[0].id
  subnet_id             = data.edgenext_ecs_vpc_subnets.all.subnets[0].id
  port_security_enabled = true
  security_groups       = [data.edgenext_ecs_security_groups.all.security_groups[0].id]
}
```

Import

Import format is `network_interface_id`.

```shell
terraform import edgenext_ecs_network_interface.example 29faf396-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```
