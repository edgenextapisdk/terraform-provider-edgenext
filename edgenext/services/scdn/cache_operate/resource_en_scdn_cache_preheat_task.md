Provides a resource to create SCDN cache preheat tasks.

Example Usage

Preheat cache for URLs

```hcl
resource "edgenext_scdn_cache_preheat_task" "example" {
  preheat_url = [
    "https://example.com/page1",
    "https://example.com/page2",
    "https://example.com/page3"
  ]
}
```

Preheat by group

```hcl
resource "edgenext_scdn_cache_preheat_task" "example" {
  group_id   = 1
  protocol   = "https"
  port       = "443"
  preheat_url = [
    "https://example.com/page1"
  ]
}
```

