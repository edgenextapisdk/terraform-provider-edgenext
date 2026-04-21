# en_ecs_images

This data source provides a list of EdgeNext ECS images in your current credentials.

## Example Usage

```hcl
data "edgenext_ecs_images" "example" {
  name_regex = "example"
}

output "image_id" {
  value = data.edgenext_ecs_images.example.images.0.id
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to filter results by name.
* `ids` - (Optional) A list of image IDs.

## Attributes Reference

The following attributes are exported:

* `id` - The provider-assigned unique ID for this managed resource.
* `images` - A list of images. Each element contains the following attributes:
  * `id` - The ID of the image.
  * `name` - The name of the image.
