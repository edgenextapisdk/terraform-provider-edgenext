Provides a resource to batch configure SCDN security protection templates.

Example Usage

Batch configure templates

```hcl
resource "edgenext_scdn_security_protection_template_batch_config" "example" {
  template_ids = [12345, 67890]

  ddos_config {
    application_ddos_protection {
      status    = "on"
      ai_cc_status = "on"
      type      = "normal"
    }
  }

  waf_config {
    waf_rule_config {
      status    = "on"
      waf_level = "strict"
      waf_mode  = "block"
    }
  }
}
```

Batch configure for all domains

```hcl
resource "edgenext_scdn_security_protection_template_batch_config" "example" {
  template_ids = [12345]
  all          = 1

  ddos_config {
    application_ddos_protection {
      status = "on"
    }
  }
}
```

Import

SCDN security protection template batch configuration can be imported using the template IDs:

```shell
terraform import edgenext_scdn_security_protection_template_batch_config.example 12345,67890
```

