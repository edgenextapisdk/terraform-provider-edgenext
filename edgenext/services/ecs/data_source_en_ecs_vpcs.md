Use this data source to query ECS VPC networks.

Example Usage

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
