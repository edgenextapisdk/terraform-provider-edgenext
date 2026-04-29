Use this resource to create and delete an ECS VPC subnet.

Example Usage

```hcl
resource "edgenext_ecs_vpc" "example" {
  name = "example-vpc"
  subnet {
    name = "seed-subnet"
    cidr = "172.31.1.0/24"
  }
}

resource "edgenext_ecs_vpc_subnet" "example" {
  vpc_id     = edgenext_ecs_vpc.example.id
  name       = "example-subnet"
  ip_version = 4
  cidr       = "172.31.10.0/24"
}
```

Import

Import format is `vpc_id/subnet_id`.

```shell
terraform import edgenext_ecs_vpc_subnet.example 68451a78-xxxx-xxxx-xxxx-xxxxxxxxxxxx/b34fe463-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```
