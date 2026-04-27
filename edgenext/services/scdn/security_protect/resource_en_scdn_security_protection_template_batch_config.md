Provides a resource to batch configure SCDN security protection templates.

Example Usage

Batch configure templates with precise access control

```hcl
resource "edgenext_scdn_security_protection_template_batch_config" "example" {
  template_ids = [12345, 67890]
  domains      = ["example.com"]

  ddos_config {
    application_ddos_protection {
      status                = "on"
      ai_cc_status          = "on"
      type                  = "strict"
      need_attack_detection = 1
      ai_status             = "on"
    }
  }

  waf_rule_config {
    waf_rule_config {
      status    = "on"
      ai_status = "on"
      waf_level = "strict"
      waf_mode  = "block"
    }
  }

  precise_access_control_config {
    action = "add"

    # Anti-CC protection for specific URL
    policies {
      action = "anticc"
      action_data = {
        level = "default"
      }
      type   = "plus"
      status = 1
      sort   = 1

      rules {
        rule_type = "url"
        logic     = "contains"
        data      = jsonencode(["/aaa"])
      }
    }

    # Deny access from specific referer
    policies {
      action = "deny"
      type   = "plus"
      from   = "aR"
      status = 1

      rules {
        rule_type = "referer_domain"
        logic     = "not_equals"
        data      = jsonencode(["home.example.com"])
      }

      rules {
        rule_type = "referer"
        logic     = "len_greater_than"
        data      = jsonencode({ len = 0 })
      }

      rules {
        rule_type = "postfix"
        logic     = "equals"
        data      = jsonencode(["css", "js", "png", "jpg"])
      }
    }

    # Region-based access control
    policies {
      action = "deny"
      type   = "plus"
      from   = "zL"
      status = 1

      rules {
        rule_type = "region"
        logic     = "not_belongs"
        data      = jsonencode({ province = [], country = ["CN", "JP"] })
      }
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

```hcl
terraform import edgenext_scdn_security_protection_template_batch_config.example 12345,67890