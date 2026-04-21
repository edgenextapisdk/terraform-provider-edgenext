---
subcategory: "Elastic Compute Service (ECS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_ecs_vpc"
sidebar_current: "docs-edgenext-resource-ecs_vpc"
description: |-
  Use this resource to create and manage ECS VPC networks.
---

# edgenext_ecs_vpc

Use this resource to create and manage ECS VPC networks.

## Example Usage

```hcl
resource "edgenext_ecs_vpc" "example" {
  region      = "tokyo-a"
  name        = "example-vpc"
  description = "vpc for app"

  subnet {
    name       = "example-subnet"
    ip_version = 4
    cidr       = "192.168.0.0/24"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) name description
* `region` - (Required, String, ForceNew) region description
* `subnet` - (Required, List, ForceNew) Subnet configuration used when creating the VPC.
* `description` - (Optional, String) description description

The `subnet` object supports the following:

* `cidr` - (Required, String, ForceNew) Subnet CIDR.
* `name` - (Required, String, ForceNew) Subnet name.
* `ip_version` - (Optional, Int, ForceNew) Subnet IP version.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `cidr` - The primary IPv4 CIDR returned by the API.
* `created_at` - Creation time.
* `project_id` - The project ID.
* `status` - The VPC status.
* `total_ips` - Total number of IPs in the VPC.
* `updated_at` - Last update time.
* `used_ips` - Used number of IPs in the VPC.


## Import

Import format is `region/network_id`.

```shell
terraform import edgenext_ecs_vpc.example tokyo-a/0e07db22-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```

Argument Reference

* `region` - (Required) Region.
* `name` - (Required) VPC name.
* `description` - (Optional) VPC description.
* `subnet` - (Required, ForceNew) Initial subnet block with:
  * `name` - (Required) Subnet name.
  * `ip_version` - (Optional) IP version, default `4`.
  * `cidr` - (Required) Subnet CIDR.

Attributes Reference

* `id` - VPC network ID.
* `cidr`, `status`, `total_ips`, `used_ips`, `project_id`
* `created_at`, `updated_at`

