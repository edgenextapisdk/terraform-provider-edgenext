---
subcategory: "Elastic Compute Service (ECS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_ecs_images"
sidebar_current: "docs-edgenext-datasource-ecs_images"
description: |-
  Use this data source to query ECS images.
---

# edgenext_ecs_images

Use this data source to query ECS images.

## Example Usage

```hcl
data "edgenext_ecs_images" "example" {
  region     = "tokyo-a"
  visibility = "public"
  name       = "Debian"
  page_num   = 1
  page_size  = 10
}

output "first_image_id" {
  value = try(data.edgenext_ecs_images.example.images[0].id, null)
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Required, String) region description
* `name` - (Optional, String) The name to filter images.
* `page_num` - (Optional, Int) Page number for image listing.
* `page_size` - (Optional, Int) Page size for image listing.
* `status` - (Optional, String) Image status to filter by.
* `visibility` - (Optional, String) Image visibility to filter by.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `images` - A list of ECS images.
  * `created_at` - Creation time of the image.
  * `description` - The description of the image.
  * `id` - The ID of the image.
  * `image_type` - The image type.
  * `min_disk` - Minimum disk required.
  * `min_ram` - Minimum RAM required.
  * `name` - The name of the image.
  * `os_distro` - OS distribution of the image.
  * `os_version` - OS version of the image.
  * `size` - The size of the image in bytes.
  * `status` - The status of the image.
  * `updated_at` - Last update time of the image.
  * `visibility` - The visibility of the image.
* `total` - Total number of images.


