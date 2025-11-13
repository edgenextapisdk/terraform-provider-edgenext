Provides a resource to create and manage SCDN origin servers for a domain.

Example Usage

Create origin with IP records

```hcl
resource "edgenext_scdn_origin" "example" {
  domain_id      = 12345
  protocol       = 0  # HTTP
  listen_ports   = [80, 443]
  origin_protocol = 0  # HTTP
  load_balance   = 1  # Round Robin
  origin_type    = 0  # IP

  records {
    view     = "primary"
    value    = "1.2.3.4"
    port     = 80
    priority = 10
  }
}
```

Create origin with domain records

```hcl
resource "edgenext_scdn_origin" "example" {
  domain_id      = 12345
  protocol       = 0
  listen_ports   = [80]
  origin_protocol = 0
  load_balance   = 1
  origin_type    = 1  # Domain

  records {
    view     = "primary"
    value    = "origin.example.com"
    port     = 80
    priority = 10
  }
}
```

Import

SCDN origins can be imported using the origin ID:

```shell
terraform import edgenext_scdn_origin.example 67890
```

