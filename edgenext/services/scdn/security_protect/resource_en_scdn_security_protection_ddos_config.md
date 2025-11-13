Provides a resource to manage SCDN security protection DDoS configuration.

Example Usage

Configure DDoS protection

```hcl
resource "edgenext_scdn_security_protection_ddos_config" "example" {
  business_id = 12345

  application_ddos_protection {
    status              = "on"
    ai_cc_status        = "on"
    type                = "normal"
    need_attack_detection = 1
    ai_status           = "on"
  }
}
```

Import

SCDN security protection DDoS configuration can be imported using the business ID:

```shell
terraform import edgenext_scdn_security_protection_ddos_config.example 12345
```

