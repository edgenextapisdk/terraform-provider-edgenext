Use this resource to create and delete an ECS VPC subnet.

Example Usage

```hcl
resource "edgenext_ecs_vpc_subnet" "example" {
  region     = "tokyo-a"
  network_id = "68451a78-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  name       = "example-subnet"
  ip_version = 4
  cidr       = "172.31.10.0/24"
}
```

Import

Import format is `region/network_id/subnet_id`.

```shell
terraform import edgenext_ecs_vpc_subnet.example tokyo-a/68451a78-xxxx-xxxx-xxxx-xxxxxxxxxxxx/b34fe463-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```

Argument Reference

* `region` - (Required) Region.
* `network_id` - (Required, ForceNew) VPC network ID.
* `name` - (Required, ForceNew) Subnet name.
* `ip_version` - (Optional, ForceNew) IP version, default `4`.
* `cidr` - (Required, ForceNew) Subnet CIDR.

Attributes Reference

* `id` - Subnet ID.
* `tenant_id`, `subnetpool_id`, `enable_dhcp`, `ipv6_ra_mode`, `ipv6_address_mode`
* `gateway_ip`, `allocation_pools`, `host_routes`, `dns_nameservers`
* `description`, `service_types`, `tags`
* `created_at`, `updated_at`, `revision_number`, `project_id`
* `used_ips`, `total_ips`, `port_num`, `not_bind_reason`, `router_id`
