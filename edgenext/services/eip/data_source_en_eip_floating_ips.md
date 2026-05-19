Use this data source to list EdgeNext floating IPs (EIPs). Filter by `floating_ip_id` or `floating_ip_address`. Each list item includes association fields when bound (`network_interface_id`, `fixed_ip_address`).

Example Usage

```hcl
data "edgenext_eip_floating_ips" "all" {
  limit = 50
}

data "edgenext_eip_floating_ips" "one" {
  floating_ip_id      = "f5390261-241e-43e0-ab1b-5bffaf334c81"
  floating_ip_address = ""
  limit               = 10
}

output "unbound_eips" {
  value = [
    for f in data.edgenext_eip_floating_ips.all.floating_ips : f.id
    if f.network_interface_id == ""
  ]
}

resource "edgenext_eip_association" "web" {
  allocation_id = data.edgenext_eip_floating_ips.one.floating_ips[0].id
  instance_id   = var.ecs_instance_id
}
```

