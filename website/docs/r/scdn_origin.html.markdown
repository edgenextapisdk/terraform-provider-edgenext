---
subcategory: "Security CDN (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_origin"
sidebar_current: "docs-edgenext-resource-scdn_origin"
description: |-
  Provides a resource to create and manage SCDN origin servers for a domain.
---

# edgenext_scdn_origin

Provides a resource to create and manage SCDN origin servers for a domain.

## Example Usage

### Create origin with IP records

```hcl
resource "edgenext_scdn_origin" "example" {
  domain_id       = 12345
  protocol        = 0 # HTTP
  listen_ports    = [80, 443]
  origin_protocol = 0 # HTTP
  load_balance    = 1 # Round Robin
  origin_type     = 0 # IP

  records {
    view     = "primary"
    value    = "1.2.3.4"
    port     = 80
    priority = 10
  }
}
```

### Create origin with domain records

```hcl
resource "edgenext_scdn_origin" "example" {
  domain_id       = 12345
  protocol        = 0
  listen_ports    = [80]
  origin_protocol = 0
  load_balance    = 1
  origin_type     = 1 # Domain

  records {
    view     = "primary"
    value    = "origin.example.com"
    port     = 80
    priority = 10
  }
}
```

## Argument Reference

The following arguments are supported:

* `domain_id` - (Required, Int, ForceNew) The ID of the domain to add origins to
* `listen_ports` - (Required, List: [`Int`]) The listening ports of the origin server
* `load_balance` - (Required, Int) The load balancing method. Valid values: 0 (IP hash), 1 (Round robin), 2 (Cookie)
* `origin_protocol` - (Required, Int) The origin protocol. Valid values: 0 (HTTP), 1 (HTTPS), 2 (Follow)
* `origin_type` - (Required, Int) The origin type. Valid values: 0 (IP), 1 (Domain)
* `protocol` - (Required, Int) The origin protocol. Valid values: 0 (HTTP), 1 (HTTPS)
* `records` - (Required, List) The origin records

The `records` object supports the following:

* `port` - (Required, Int) The port of the record
* `priority` - (Required, Int) The priority of the record
* `value` - (Required, String) The value of the record (IP address or domain)
* `view` - (Required, String) The view of the record

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the origin
* `listen_port` - The listening port of the origin server (single port from API)


## Import

SCDN origins can be imported using the origin ID:

```shell
terraform import edgenext_scdn_origin.example 67890
```

