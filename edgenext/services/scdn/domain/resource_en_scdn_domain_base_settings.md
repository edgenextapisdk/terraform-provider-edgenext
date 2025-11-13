Provides a resource to manage base settings of an SCDN domain.

Example Usage

Update domain base settings

```hcl
resource "edgenext_scdn_domain_base_settings" "example" {
  domain_id = 12345

  proxy_host {
    proxy_host     = "example.com"
    proxy_host_type = "custom"
  }

  proxy_sni {
    proxy_sni = "example.com"
    status    = "on"
  }

  domain_redirect {
    status   = "on"
    jump_to  = "https://www.example.com"
    jump_type = "301"
  }
}
```

Import

SCDN domain base settings can be imported using the domain ID:

```shell
terraform import edgenext_scdn_domain_base_settings.example 12345
```

