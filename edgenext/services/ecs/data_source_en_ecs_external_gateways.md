Use this data source to query ECS external gateway networks.

Example Usage

```hcl
data "edgenext_ecs_external_gateways" "example" {
  limit  = 10
}

resource "edgenext_ecs_router" "example" {
  name                = "router-with-external-gateway"
  external_network_id = data.edgenext_ecs_external_gateways.example.external_gateways[0].id
}
```
