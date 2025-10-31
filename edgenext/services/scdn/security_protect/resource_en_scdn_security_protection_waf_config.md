Provides a resource to manage SCDN security protection WAF configuration.

Example Usage

Configure WAF protection

```hcl
resource "edgenext_scdn_security_protection_waf_config" "example" {
  business_id = 12345

  waf_rule_config {
    status    = "on"
    ai_status = "on"
    waf_level = "strict"
    waf_mode  = "block"
  }
}
```

Import

SCDN security protection WAF configuration can be imported using the business ID:

```shell
terraform import edgenext_scdn_security_protection_waf_config.example 12345
```

