---
subcategory: "Elastic Compute Service (ECS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_ecs_vpc_subnet"
sidebar_current: "docs-edgenext-resource-ecs_vpc_subnet"
description: |-
  Use this resource to create and delete an ECS VPC subnet.
---

# edgenext_ecs_vpc_subnet

Use this resource to create and delete an ECS VPC subnet.

## Example Usage

```hcl
resource "edgenext_ecs_vpc_subnet" "example" {
  region     = "tokyo-a"
  network_id = "68451a78-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  name       = "example-subnet"
  ip_version = 4
  cidr       = "172.31.10.0/24"
}
```

## Argument Reference

The following arguments are supported:

* `cidr` - (Required, String, ForceNew) Subnet CIDR.
* `name` - (Required, String, ForceNew) Subnet name.
* `network_id` - (Required, String, ForceNew) The VPC network ID.
* `region` - (Required, String, ForceNew) The region of the subnet.
* `ip_version` - (Optional, Int, ForceNew) IP version.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `allocation_pools` - Allocation pools.
  * `end` - End IP.
  * `start` - Start IP.
* `created_at` - Creation time.
* `description` - Description.
* `dns_nameservers` - DNS nameservers.
* `enable_dhcp` - Whether DHCP is enabled.
* `gateway_ip` - Gateway IP.
* `host_routes` - Host routes.
  * `destination` - Route destination.
  * `nexthop` - Route next hop.
* `ipv6_address_mode` - IPv6 address mode.
* `ipv6_ra_mode` - IPv6 RA mode.
* `not_bind_reason` - Reason if subnet is not bindable.
* `port_num` - Port count.
* `project_id` - Project ID.
* `revision_number` - Revision number.
* `router_id` - Bound router ID.
* `service_types` - Service types.
* `subnetpool_id` - Subnet pool ID.
* `tags` - Tags.
* `tenant_id` - Tenant ID.
* `total_ips` - Total IP count.
* `updated_at` - Last update time.
* `used_ips` - Used IP count.


## Import

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

