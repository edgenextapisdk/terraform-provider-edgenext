---
subcategory: "Elastic Compute Service (ECS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_ecs_vpc_subnets"
sidebar_current: "docs-edgenext-datasource-ecs_vpc_subnets"
description: |-
  Use this data source to query subnets under an ECS VPC network.
---

# edgenext_ecs_vpc_subnets

Use this data source to query subnets under an ECS VPC network.

## Example Usage

```hcl
data "edgenext_ecs_vpc_subnets" "example" {
  region     = "tokyo-a"
  network_id = "0e07db22-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  router_id  = "f9883769-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}

output "vpc_subnet_total" {
  value = data.edgenext_ecs_vpc_subnets.example.total
}
```

## Argument Reference

The following arguments are supported:

* `network_id` - (Required, String) The VPC network ID.
* `region` - (Required, String) region description
* `router_id` - (Optional, String) The router ID to filter subnets.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `subnets` - A list of VPC subnets.
  * `allocation_pools` - Allocation pools.
    * `end` - End IP.
    * `start` - Start IP.
  * `cidr` - CIDR block.
  * `created_at` - Creation time.
  * `description` - Description.
  * `dns_nameservers` - DNS nameservers.
  * `enable_dhcp` - Whether DHCP is enabled.
  * `gateway_ip` - Gateway IP.
  * `host_routes` - Host routes.
    * `destination` - Route destination.
    * `nexthop` - Route next hop.
  * `id` - Subnet ID.
  * `ip_version` - IP version.
  * `ipv6_address_mode` - IPv6 address mode.
  * `ipv6_ra_mode` - IPv6 RA mode.
  * `name` - Subnet name.
  * `network_id` - Network ID.
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
* `total` - Total number of matched subnets.


