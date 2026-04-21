---
subcategory: "Elastic Compute Service (ECS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_ecs_key_pair"
sidebar_current: "docs-edgenext-resource-ecs_key_pair"
description: |-
  Use this resource to create and manage ECS key pairs.
---

# edgenext_ecs_key_pair

Use this resource to create and manage ECS key pairs.

## Example Usage

```hcl
resource "edgenext_ecs_key_pair" "example" {
  region     = "tokyo-a"
  name       = "example-key"
  public_key = file("~/.ssh/id_rsa.pub")
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, ForceNew) name description
* `region` - (Required, String, ForceNew) region description
* `public_key` - (Optional, String, ForceNew) public_key description

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `private_key` - private_key description


## Import

Import format is `region/name`.

```shell
terraform import edgenext_ecs_key_pair.example tokyo-a/example-key
```

Argument Reference

* `region` - (Required) Region.
* `name` - (Required) Key pair name.
* `public_key` - (Optional) Public key content.

Attributes Reference

* `id` - Key pair name (provider ID).
* `private_key` - Private key returned by API (if generated server side).

