Use this resource to create and delete an ECS VPC subnet.

Example Usage

```hcl
resource "edgenext_ecs_vpc" "example" {
  name = "example-vpc"
  subnet {
    name = "seed-subnet"
    cidr = "172.31.1.0/24"
  }
}

resource "edgenext_ecs_vpc_subnet" "example" {
  vpc_id     = edgenext_ecs_vpc.example.id
  name       = "example-subnet"
  ip_version = 4
  cidr       = "172.31.10.0/24"
}
```

Import

Import format is `vpc_id/subnet_id`.

```shell
terraform import edgenext_ecs_vpc_subnet.example 68451a78-xxxx-xxxx-xxxx-xxxxxxxxxxxx/b34fe463-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```

Argument Reference

* `vpc_id` - (Required) VPC ID. Cannot be changed after creation.
* `name` - (Required) Subnet name. Cannot be changed after creation.
* `ip_version` - (Optional) IP version, default `4`. Cannot be changed after creation.
* `cidr` - (Required) Subnet CIDR. Cannot be changed after creation.

Attributes Reference

* `id` - Subnet ID.
* `tenant_id`, `subnetpool_id`, `enable_dhcp`, `ipv6_ra_mode`, `ipv6_address_mode`
* `gateway_ip`, `allocation_pools`, `host_routes`, `dns_nameservers`
* `description`, `service_types`, `tags`
* `created_at`, `updated_at`, `revision_number`, `project_id`
* `used_ips`, `total_ips`, `port_num`, `not_bind_reason`, `router_id`
