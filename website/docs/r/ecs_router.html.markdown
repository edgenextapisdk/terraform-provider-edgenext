---
subcategory: "Elastic Compute Service (ECS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_ecs_router"
sidebar_current: "docs-edgenext-resource-ecs_router"
description: |-
  Use this resource to create and manage ECS routers.
---

# edgenext_ecs_router

Use this resource to create and manage ECS routers.

## Example Usage

```hcl
resource "edgenext_ecs_router" "example" {
  region              = "tokyo-a"
  name                = "example-router"
  description         = "router for app network"
  external_network_id = "5c83af33-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) name description
* `region` - (Required, String, ForceNew) region description
* `description` - (Optional, String) description description
* `external_network_id` - (Optional, String) external_network_id description

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `admin_state_up` - admin_state_up description
* `created_at` - created_at description
* `project_id` - project_id description
* `revision_number` - revision_number description
* `status` - status description
* `tenant_id` - tenant_id description
* `updated_at` - updated_at description


## Import

Import format is `region/router_id`.

```shell
terraform import edgenext_ecs_router.example tokyo-a/f9883769-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```

Argument Reference

* `region` - (Required) Region.
* `name` - (Required) Router name.
* `description` - (Optional) Router description.
* `external_network_id` - (Optional) External gateway network ID.

Attributes Reference

* `id` - Router ID.
* `tenant_id`, `admin_state_up`, `status`, `project_id`
* `created_at`, `updated_at`, `revision_number`

