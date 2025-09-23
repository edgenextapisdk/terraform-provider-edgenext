Provides a resource to create and manage CDN domain configuration.

Example Usage

Basic CDN domain configuration

```hcl
resource "edgenext_cdn_domain" "example" {
  domain = "example.com"
  area   = "global"
  type   = "page"
  
  config {
    origin {
      default_master = "origin.example.com"
      origin_mode    = "default"
    }
  }
}
```

Advanced CDN domain configuration with cache rules and HTTPS

```hcl
resource "edgenext_cdn_domain" "example" {
  domain = "example.com"
  area   = "global"
  type   = "page"
  
  config {
    origin {
      default_master = "origin.example.com"
      default_slave  = "backup.example.com"
      origin_mode    = "custom"
      port           = 443
      ori_https      = "yes"
    }
    
    cache_rule {
      type       = 1
      pattern    = "jpg,png,gif"
      time       = 86400
      timeunit   = "s"
      ignore_query = "on"
    }
    
    cache_rule {
      type       = 1  
      pattern    = "css,js"
      time       = 3600
      timeunit   = "s"
      ignore_query = "off"
    }
    
    https {
      cert_id      = 123
      http2        = "on"
      force_https  = "302"
    }
    
    referer {
      type        = 2
      list        = ["*.example.com", "example.org"]
      allow_empty = true
    }
  }
}
```

Import

CDN domain configuration can be imported using the domain name:

```shell
terraform import edgenext_cdn_domain.example example.com
```
