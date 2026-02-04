Provides a resource to switch domains to a different SCDN Rule Template.

Example Usage

Switch domains to a new template

```hcl
resource "edgenext_scdn_rule_template_switch" "example" {
  app_type     = "network_speed"
  new_tpl_id   = 12345
  new_tpl_type = "more_domain"
  domain_ids   = [1001, 1002]
}
```

Argument Reference

The following arguments are supported:

* `app_type` - (Optional) Application type the template applies to. Default: `network_speed`.
* `domain_ids` - (Required) List of domain IDs to switch templates.
* `new_tpl_id` - (Required) New template ID to switch to. When `new_tpl_type` is `global`, pass 0.
* `new_tpl_type` - (Required) New template type. Valid values: `more_domain`, `global`.

Import

Rule Template Switch instances can be imported using an ID with the format `template_id:domain_id1,domain_id2,...`:

```shell
terraform import edgenext_scdn_rule_template_switch.example 12345:1001,1002
```
