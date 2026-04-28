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
  name  = edgenext_ecs_vpc.example.name
  limit = 10
}

resource "edgenext_ecs_vpc" "example" {
  name = "default-vpc"
  subnet {
    name = "default-subnet"
    cidr = "10.10.0.0/24"
  }
}
```

## Argument Reference

The following arguments are supported:

* `limit` - (Optional, Int) Maximum number of vpcs to return.
* `name` - (Optional, String) The name to filter vpcs.
* `vpc_id` - (Optional, String) The VPC ID to filter vpcs.

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


