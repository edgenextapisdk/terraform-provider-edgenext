Use this data source to query ports attached to a specific ECS router.

Example Usage

```hcl
data "edgenext_ecs_router_ports" "example" {
  router_id = data.edgenext_ecs_routers.example.routers[0].id
}

data "edgenext_ecs_routers" "example" {
  router_name = "default-router"
  limit       = 1
}
```
