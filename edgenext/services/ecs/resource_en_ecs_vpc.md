Use this resource to create and manage ECS VPC networks.

Example Usage

```hcl
resource "edgenext_ecs_vpc" "example" {
  name        = "example-vpc"
  description = "vpc for app"

  subnet {
    name       = "example-subnet"
    ip_version = 4
    cidr       = "192.168.0.0/24"
  }
}
```

Import

Import format is `vpc_id`.

```shell
terraform import edgenext_ecs_vpc.example 0e07db22-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```
