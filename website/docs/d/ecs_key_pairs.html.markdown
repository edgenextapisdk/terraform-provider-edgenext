---
subcategory: "Elastic Compute Service (ECS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_ecs_key_pairs"
sidebar_current: "docs-edgenext-datasource-ecs_key_pairs"
description: |-
  Use this data source to query ECS key pairs.
---

# edgenext_ecs_key_pairs

Use this data source to query ECS key pairs.

## Example Usage

```hcl
data "edgenext_ecs_key_pairs" "example" {
  region = "tokyo-a"
  limit  = 10
}

output "first_key_name" {
  value = try(data.edgenext_ecs_key_pairs.example.key_pairs[0].name, null)
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Required, String) region description
* `limit` - (Optional, Int) Maximum number of key_pairs to return.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `key_pairs` - A list of ECS key_pairs.
  * `fingerprint` - The fingerprint of the key_pair.
  * `name` - The name of the key_pair.
  * `public_key` - The public key material.
  * `type` - The key type (e.g. ssh).


