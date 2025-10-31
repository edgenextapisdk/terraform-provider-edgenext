---
subcategory: "Secure Content Delivery Network (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_origin_group"
sidebar_current: "docs-edgenext-resource-scdn_origin_group"
description: |-
  Provides a resource to create and manage SCDN origin groups.
---

# edgenext_scdn_origin_group

Provides a resource to create and manage SCDN origin groups.

## Example Usage

### Create origin group with IP origins

```hcl
resource "edgenext_scdn_origin_group" "example" {
  name   = "my-origin-group"
  remark = "My origin group"

  origins {
    origin_type = 0 # IP

    records {
      value    = "1.2.3.4"
      port     = 80
      priority = 10
      view     = "primary"
    }

    protocol_ports {
      protocol = 0 # HTTP
      ports    = [80]
    }
  }
}
```

### Create origin group with domain origins

```hcl
resource "edgenext_scdn_origin_group" "example" {
  name   = "my-origin-group"
  remark = "My origin group"

  origins {
    origin_type = 1 # Domain

    records {
      value    = "origin.example.com"
      port     = 80
      priority = 10
      view     = "primary"
      host     = "example.com"
    }

    protocol_ports {
      protocol = 0 # HTTP
      ports    = [80]
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Origin group name (2-16 characters)
* `origins` - (Required, List) Origin list (at least 1)
* `origin_group_id` - (Optional, Int) Origin group ID. Required for update/delete, computed for create. If provided during create, will update existing origin group instead.
* `remark` - (Optional, String) Remark (2-64 characters)

The `origins` object supports the following:

* `load_balance` - (Required, Int) Load balance strategy: 0-ip_hash, 1-round_robin, 2-cookie
* `origin_protocol` - (Required, Int) Origin protocol: 0-http, 1-https, 2-follow
* `origin_type` - (Required, Int) Origin type: 0-IP, 1-domain
* `protocol_ports` - (Required, List) Protocol port mapping (at least 1)
* `records` - (Required, List) Origin record list (at least 1)
* `id` - (Optional, Int) Origin ID (0 for new, >0 for update)

The `protocol_ports` object of `origins` supports the following:

* `listen_ports` - (Required, List) Listen port list
* `protocol` - (Required, Int) Protocol: 0-http, 1-https

The `records` object of `origins` supports the following:

* `port` - (Required, Int) Origin port (1-65535)
* `priority` - (Required, Int) Weight (1-100)
* `value` - (Required, String) Origin address
* `view` - (Required, String) Origin type: primary-backup, backup-backup
* `host` - (Optional, String) Origin Host

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `created_at` - Creation time
* `member_id` - Member ID
* `updated_at` - Update time
* `username` - Username


## Import

SCDN origin groups can be imported using the origin group ID:

```shell
terraform import edgenext_scdn_origin_group.example 12345
```

