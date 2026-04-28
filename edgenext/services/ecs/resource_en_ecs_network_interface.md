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

Import format is `port_id`.

```shell
terraform import edgenext_ecs_network_interface.example 29faf396-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```

Argument Reference

* `name` - (Required) Port name.
* `description` - (Optional) Port description.
* `vpc_id` - (Required) VPC ID. Cannot be changed after creation.
* `subnet_id` - (Required) Subnet ID. Cannot be changed after creation.
* `port_security_enabled` - (Optional) Port security switch.
* `security_groups` - (Optional) Security group IDs.

Attributes Reference

* `id` - Port ID.
* `instance_id`, `floating_ip_address`, `tenant_id`, `project_id`, `status`, `instance_owner`
* `fixed_ips` - Fixed IP list with `subnet_id`, `ip_address`, `floating_ip`.
* `qos_policy_id`, `tags`, `created_at`, `updated_at`, `revision_number`
* `mac_address`, `binding_vnic_type`, `instance_name`, `vpc_name`, `ipv4`, `ipv6`
