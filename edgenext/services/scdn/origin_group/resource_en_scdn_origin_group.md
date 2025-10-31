Provides a resource to create and manage SCDN origin groups.

Example Usage

Create origin group with IP origins

```hcl
resource "edgenext_scdn_origin_group" "example" {
  name   = "my-origin-group"
  remark = "My origin group"

  origins {
    origin_type = 0  # IP

    records {
      value    = "1.2.3.4"
      port     = 80
      priority = 10
      view     = "primary"
    }

    protocol_ports {
      protocol = 0  # HTTP
      ports    = [80]
    }
  }
}
```

Create origin group with domain origins

```hcl
resource "edgenext_scdn_origin_group" "example" {
  name   = "my-origin-group"
  remark = "My origin group"

  origins {
    origin_type = 1  # Domain

    records {
      value    = "origin.example.com"
      port     = 80
      priority = 10
      view     = "primary"
      host     = "example.com"
    }

    protocol_ports {
      protocol = 0  # HTTP
      ports    = [80]
    }
  }
}
```

Import

SCDN origin groups can be imported using the origin group ID:

```shell
terraform import edgenext_scdn_origin_group.example 12345
```

