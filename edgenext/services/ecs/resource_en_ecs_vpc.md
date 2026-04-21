Use this resource to create and manage ECS VPC networks.

Example Usage

```hcl
resource "edgenext_ecs_vpc" "example" {
  region      = "tokyo-a"
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

Import format is `region/network_id`.

```shell
terraform import edgenext_ecs_vpc.example tokyo-a/0e07db22-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```

Argument Reference

* `region` - (Required) Region.
* `name` - (Required) VPC name.
* `description` - (Optional) VPC description.
* `subnet` - (Required, ForceNew) Initial subnet block with:
  * `name` - (Required) Subnet name.
  * `ip_version` - (Optional) IP version, default `4`.
  * `cidr` - (Required) Subnet CIDR.

Attributes Reference

* `id` - VPC network ID.
* `cidr`, `status`, `total_ips`, `used_ips`, `project_id`
* `created_at`, `updated_at`
