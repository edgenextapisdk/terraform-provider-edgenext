Provides a resource to switch the node type of an SCDN domain.

Example Usage

Switch to SCDN node

```hcl
resource "edgenext_scdn_domain_node_switch" "example" {
  domain_id      = 12345
  protect_status = "scdn"
}
```

Switch to exclusive node

```hcl
resource "edgenext_scdn_domain_node_switch" "example" {
  domain_id           = 12345
  protect_status      = "exclusive"
  exclusive_resource_id = 999
}
```

Import

SCDN domain node switch can be imported using the domain ID:

```shell
terraform import edgenext_scdn_domain_node_switch.example 12345
```

