---
subcategory: "Secure Content Delivery Network (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_domain"
sidebar_current: "docs-edgenext-resource-scdn_domain"
description: |-
  Provides a resource to create and manage SCDN domains.
---

# edgenext_scdn_domain

Provides a resource to create and manage SCDN domains.

## Example Usage

### Basic domain creation

```hcl
resource "edgenext_scdn_domain" "example" {
  domain         = "example.com"
  group_id       = 1
  remark         = "My SCDN domain"
  protect_status = "scdn"
  app_type       = "web"

  origins {
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
}
```

### Domain with multiple origin records

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

### Domain with domain-type origin

```hcl
resource "edgenext_scdn_domain" "example" {
  domain         = "example.com"
  protect_status = "scdn"

  origins {
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
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String, ForceNew) The domain name to be added to SCDN
* `origins` - (Required, List) The origin server configuration
* `app_type` - (Optional, String) The application type
* `exclusive_resource_id` - (Optional, Int) The ID of the exclusive resource package
* `group_id` - (Optional, Int) The ID of the domain group
* `protect_status` - (Optional, String) The edge node type. Valid values: back_source, scdn, exclusive
* `remark` - (Optional, String) The remark for the domain
* `tpl_id` - (Optional, Int) The template ID to be applied to the domain
* `tpl_recommend` - (Optional, String) The template recommendation status

The `origins` object supports the following:

* `listen_ports` - (Required, List) The listening ports of the origin server
* `load_balance` - (Required, Int) The load balancing method. Valid values: 0 (IP hash), 1 (Round robin), 2 (Cookie)
* `origin_protocol` - (Required, Int) The origin protocol. Valid values: 0 (HTTP), 1 (HTTPS), 2 (Follow)
* `origin_type` - (Required, Int) The origin type. Valid values: 0 (IP), 1 (Domain)
* `protocol` - (Required, Int) The origin protocol. Valid values: 0 (HTTP), 1 (HTTPS)
* `records` - (Required, List) The origin records

The `records` object of `origins` supports the following:

* `port` - (Required, Int) The port of the record
* `priority` - (Required, Int) The priority of the record
* `value` - (Required, String) The value of the record (IP address or domain)
* `view` - (Required, String) The view of the record

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `access_mode` - The access mode (ns or cname)
* `access_progress_desc` - The description of the access progress status
* `access_progress` - The access progress status
* `ca_id` - The certificate ID
* `ca_status` - The certificate binding status
* `cname` - The CNAME information
  * `master` - The master CNAME record
  * `slaves` - The slave CNAME records
* `created_at` - The creation timestamp
* `ei_forward_status` - The explicit/implicit forwarding status
* `has_origin` - Whether the domain has origin configuration
* `id` - The ID of the domain
* `pri_domain` - The primary domain
* `updated_at` - The last update timestamp
* `use_my_cname` - The CNAME resolution status
* `use_my_dns` - The DNS hosting status


## Import

SCDN domains can be imported using the domain ID:

```shell
terraform import edgenext_scdn_domain.example 12345
```

