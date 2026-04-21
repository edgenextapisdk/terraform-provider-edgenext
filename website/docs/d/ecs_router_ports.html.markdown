---
subcategory: "Elastic Compute Service (ECS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_ecs_router_ports"
sidebar_current: "docs-edgenext-datasource-ecs_router_ports"
description: |-
  Use this data source to query ports attached to a specific ECS router.
---

# edgenext_ecs_router_ports

Use this data source to query ports attached to a specific ECS router.

## Example Usage

```hcl
data "edgenext_ecs_router_ports" "example" {
  region = "tokyo-a"
  id     = "f9883769-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}

output "router_port_total" {
  value = data.edgenext_ecs_router_ports.example.total
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required, String) The router ID.
* `region` - (Required, String) region description

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `ports` - A list of router ports.
  * `created_at` - Creation time.
  * `id` - Port ID.
  * `ip_address` - Port IP address.
  * `mac_address` - Port MAC address.
  * `name` - Port name.
  * `network_name` - Network name.
  * `status` - Port status.
* `total` - Total number of router ports.


