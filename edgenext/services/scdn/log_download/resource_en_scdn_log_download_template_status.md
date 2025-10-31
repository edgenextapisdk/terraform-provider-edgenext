Provides a resource to manage the status (enable/disable) of SCDN log download templates.

Example Usage

Enable log download template

```hcl
resource "edgenext_scdn_log_download_template_status" "example" {
  template_id = 12345
  status      = 1  # 1: enabled
}
```

Disable log download template

```hcl
resource "edgenext_scdn_log_download_template_status" "example" {
  template_id = 12345
  status      = 0  # 0: disabled
}
```

Import

SCDN log download template status can be imported using the template ID:

```shell
terraform import edgenext_scdn_log_download_template_status.example 12345
```

