Use this resource to associate an existing EdgeNext EIP (floating IP) with an ECS instance, ELB load balancer, or network interface.

There is no in-place update; changing any argument forces replacement. On create, the EIP must not already be bound to another port. `fixed_ip_address` is optional on create (first available private IP is used) and is refreshed from the API on read. Resource ID is `allocation_id`.

Example Usage

```hcl
data "edgenext_eip_floating_ips" "all" {
  limit = 50
}

data "edgenext_ecs_instances" "all" {
  limit = 10
}

data "edgenext_elb_load_balancers" "all" {
  limit = 10
}

data "edgenext_ecs_network_interfaces" "all" {
  limit = 10
}

# ECS instance (default instance_type)
resource "edgenext_eip_association" "ecs" {
  allocation_id    = data.edgenext_eip_floating_ips.all.floating_ips[0].id
  instance_id      = data.edgenext_ecs_instances.all.instances[0].id
  fixed_ip_address = data.edgenext_ecs_instances.all.instances[0].fixed_ip_addresses[0]
}

# ELB load balancer — instance_id is the load balancer ID
resource "edgenext_eip_association" "elb" {
  allocation_id = data.edgenext_eip_floating_ips.all.floating_ips[1].id
  instance_id   = data.edgenext_elb_load_balancers.all.balancers[0].id
  instance_type = "ElbInstance"
}

# Network interface — instance_id is the port (ENI) ID
resource "edgenext_eip_association" "eni" {
  allocation_id    = data.edgenext_eip_floating_ips.all.floating_ips[2].id
  instance_id      = data.edgenext_ecs_network_interfaces.all.network_interfaces[0].id
  instance_type    = "NetworkInterface"
  fixed_ip_address = data.edgenext_ecs_network_interfaces.all.network_interfaces[0].fixed_ips[0].ip_address
}
```

Import

Import format is `allocation_id`. Provide `instance_id`, `instance_type`, and `fixed_ip_address` in Terraform configuration to match the existing association.

```shell
terraform import edgenext_eip_association.ecs f5390261-241e-43e0-ab1b-5bffaf334c81
```
