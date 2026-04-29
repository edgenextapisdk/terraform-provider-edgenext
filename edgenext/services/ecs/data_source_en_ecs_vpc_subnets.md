Use this data source to query subnets under an ECS VPC network.

Example Usage

```hcl
data "edgenext_ecs_vpc_subnets" "example" {
  vpc_id = data.edgenext_ecs_vpcs.all.vpcs[0].id
}

data "edgenext_ecs_vpcs" "all" {
  limit = 1
}
```
