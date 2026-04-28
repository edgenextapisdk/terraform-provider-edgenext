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
  vpc_id = data.edgenext_ecs_vpcs.all.vpcs[0].id
}

data "edgenext_ecs_vpcs" "all" {
  limit = 1
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required, String) The VPC ID to filter subnets.

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
  * `vpc_id` - VPC ID.
* `total` - Total number of matched subnets.


