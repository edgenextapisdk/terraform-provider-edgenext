Use this resource to create and manage ECS routers.

Example Usage

```hcl
resource "edgenext_ecs_router" "example" {
  name                = "example-router"
  description         = "router for app network"
  external_network_id = data.edgenext_ecs_external_gateways.all.external_gateways[0].id
}
```

data "edgenext_ecs_external_gateways" "all" {
  limit = 1
}

Import

Import format is `router_id`.

```shell
terraform import edgenext_ecs_router.example f9883769-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```

Argument Reference

* `name` - (Required) Router name.
* `description` - (Optional) Router description.
* `external_network_id` - (Optional) External gateway network ID.

Attributes Reference

* `id` - Router ID.
* `tenant_id`, `admin_state_up`, `status`, `project_id`
* `created_at`, `updated_at`, `revision_number`
