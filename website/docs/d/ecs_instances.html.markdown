---
subcategory: "Elastic Compute Service (ECS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_ecs_instances"
sidebar_current: "docs-edgenext-datasource-ecs_instances"
description: |-
  Use this data source to query ECS instances.
---

# edgenext_ecs_instances

Use this data source to query ECS instances.

## Example Usage

```hcl
data "edgenext_ecs_instances" "example" {
  region = "tokyo-a"
  name   = "example-instance"
  limit  = 10
}

output "instance_count" {
  value = data.edgenext_ecs_instances.example.total
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Required, String) region description
* `limit` - (Optional, Int) Maximum number of instances to return.
* `name` - (Optional, String) The name to filter instances.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instances` - A list of ECS instances.
  * `created_at` - The creation time of the instance.
  * `fixed_addresses` - A list of fixed IP addresses.
  * `flavor_info` - Flavor detail information.
    * `ram` - The RAM size in MB.
    * `vcpus` - The number of vCPUs.
  * `flavor` - The flavor name of the instance.
  * `floating_addresses` - A list of floating IP addresses.
  * `id` - The ID of the instance.
  * `image_name` - The image name of the instance.
  * `instance_cost_info` - Instance billing and expiration information.
    * `billing_model` - The billing model code.
    * `instance_cost_type` - The instance billing type.
    * `instance_expiration_time` - The instance expiration time.
    * `network_cost_type` - The network billing type.
  * `name` - The name of the instance.
  * `status` - The status of the instance.
  * `tags` - A list of tag names.
* `total` - The total number of matched instances.


