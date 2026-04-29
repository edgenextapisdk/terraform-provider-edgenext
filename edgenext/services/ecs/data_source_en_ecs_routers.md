Use this data source to query ECS routers.

Example Usage

```hcl
data "edgenext_ecs_routers" "example" {
  router_name = edgenext_ecs_router.example.name
  limit       = 10
}

resource "edgenext_ecs_router" "example" {
  name = "default-router"
}
```
