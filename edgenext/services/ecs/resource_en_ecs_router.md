Use this resource to create and manage ECS routers.

Example Usage

```hcl
resource "edgenext_ecs_router" "example" {
  region              = "tokyo-a"
  name                = "example-router"
  description         = "router for app network"
  external_network_id = "5c83af33-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

Import

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
