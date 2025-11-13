Provides a resource to create and manage SCDN domains.

Example Usage

Basic domain creation

```hcl
resource "edgenext_scdn_domain" "example" {
  domain         = "example.com"
  group_id       = 1
  remark         = "My SCDN domain"
  protect_status = "scdn"
  app_type       = "web"

  origins {
    protocol        = 0  # HTTP
    listen_ports    = [80, 443]
    origin_protocol = 0  # HTTP
    load_balance    = 1  # Round Robin
    origin_type     = 0  # IP

    records {
      view     = "primary"
      value    = "1.2.3.4"
      port     = 80
      priority = 10
    }
  }
}
```

Domain with multiple origin records

```hcl
resource "edgenext_scdn_domain" "example" {
  domain         = "example.com"
  protect_status = "scdn"

  origins {
    protocol        = 0
    listen_ports    = [80, 443]
    origin_protocol = 0
    load_balance    = 1
    origin_type     = 0

    records {
      view     = "primary"
      value    = "1.2.3.4"
      port     = 80
      priority = 10
    }

    records {
      view     = "backup"
      value    = "5.6.7.8"
      port     = 80
      priority = 20
    }
  }
}
```

Domain with domain-type origin

```hcl
resource "edgenext_scdn_domain" "example" {
  domain         = "example.com"
  protect_status = "scdn"

  origins {
    protocol        = 0
    listen_ports    = [80]
    origin_protocol = 0
    load_balance    = 1
    origin_type     = 1  # Domain

    records {
      view     = "primary"
      value    = "origin.example.com"
      port     = 80
      priority = 10
    }
  }
}
```

Import

SCDN domains can be imported using the domain ID:

```shell
terraform import edgenext_scdn_domain.example 12345
```

