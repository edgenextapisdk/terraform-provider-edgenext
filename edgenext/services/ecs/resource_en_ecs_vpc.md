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

Argument Reference

* `name` - (Required) VPC name.
* `description` - (Optional) VPC description.
* `subnet` - (Required) Initial subnet block. Cannot be changed after creation:
  * `name` - (Required) Subnet name. Cannot be changed after creation.
  * `ip_version` - (Optional) IP version, default `4`. Cannot be changed after creation.
  * `cidr` - (Required) Subnet CIDR. Cannot be changed after creation.

Attributes Reference

* `id` - VPC network ID.
* `cidr`, `status`, `total_ips`, `used_ips`, `project_id`
* `created_at`, `updated_at`
