---
subcategory: "Elastic Compute Service (ECS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_ecs_disks"
sidebar_current: "docs-edgenext-datasource-ecs_disks"
description: |-
  Use this data source to query ECS disks via **GET** `/ecs/openapi/v2/volume/list`.
---

# edgenext_ecs_disks

Use this data source to query ECS disks via **GET** `/ecs/openapi/v2/volume/list`.

## Example Usage

```hcl
data "edgenext_ecs_disks" "example" {
  name      = edgenext_ecs_disk.example.name
  page_num  = 1
  page_size = 10
}

resource "edgenext_ecs_disk" "example" {
  name        = "example-disk"
  volume_type = "SSD"
  size        = 50
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional, String) Disk name filter (empty string lists all names).
* `page_num` - (Optional, Int) Page number for listing.
* `page_size` - (Optional, Int) Page size for listing.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `disks` - Disks returned for the current page.
  * `attachment` - Attachment records when the disk is mounted on an instance.
    * `device` - Device path on the instance (e.g. /dev/vda).
    * `instance_id` - ID of the instance this disk is attached to.
    * `instance_name` - Name of the instance this disk is attached to.
  * `billing_model_name` - Billing model display name.
  * `billing_model` - Billing model code.
  * `created_at` - Creation timestamp.
  * `description` - Disk description.
  * `disk_label` - Disk label, e.g. System Disk.
  * `disk_status` - Disk status code.
  * `disk_type` - Disk product type (API field type), e.g. Quick Disk.
  * `expiration_time` - Expiration time if applicable.
  * `id` - Disk ID.
  * `image_name` - Image name associated with the disk.
  * `name` - Disk name.
  * `policy_names` - Backup or policy names attached to the disk.
  * `server_name` - Associated server name from API (may be empty).
  * `size` - Disk size in GB.
  * `status_name` - Human-readable status, e.g. in-use, available.
  * `status` - API status field.
  * `volume_type` - Volume type code (API field volumeType).
* `total` - Total number of disks reported by the API.


