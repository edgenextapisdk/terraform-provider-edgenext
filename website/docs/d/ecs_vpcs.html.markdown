---
subcategory: "Elastic Compute Service (ECS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_ecs_vpcs"
sidebar_current: "docs-edgenext-datasource-ecs_vpcs"
description: |-
  Use this data source to query ECS VPC networks.
---

# edgenext_ecs_vpcs

Use this data source to query ECS VPC networks.

## Example Usage

```hcl
data "edgenext_ecs_vpcs" "example" {
  region     = "tokyo-a"
  network_id = ""
  name       = "default-vpc"
  limit      = 10
}

output "first_vpc_id" {
  value = try(data.edgenext_ecs_vpcs.example.vpcs[0].id, null)
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Required, String) region description
* `limit` - (Optional, Int) Maximum number of vpcs to return.
* `name` - (Optional, String) The name to filter vpcs.
* `network_id` - (Optional, String) The network ID to filter vpcs.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `total` - Total number of matched vpcs.
* `vpcs` - A list of ECS vpcs.
  * `created_at` - Creation time.
  * `description` - The vpc description.
  * `id` - The ID of the vpc.
  * `ipv4_cidrs` - A list of IPv4 CIDRs.
  * `name` - The name of the vpc.
  * `project_id` - The project ID.
  * `status` - The status of the vpc.
  * `updated_at` - Last update time.


