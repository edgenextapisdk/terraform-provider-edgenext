---
subcategory: "Elastic IP (EIP)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_eip_association"
sidebar_current: "docs-edgenext-resource-eip_association"
description: |-
  Associates an EdgeNext EIP with an ECS instance, ELB load balancer, or network interface.
---

# edgenext_eip_association

Associates an EdgeNext EIP with an ECS instance, ELB load balancer, or network interface.

## Example Usage

```hcl
resource "edgenext_eip_association" "example" {
  allocation_id    = "f5390261-241e-43e0-ab1b-5bffaf334c81"
  instance_id      = "518031d1-f66d-416d-a66a-96ab91f4def9"
  fixed_ip_address = "172.31.0.32"
}
```

## Argument Reference

The following arguments are supported:

* `allocation_id` - (Required, String, ForceNew) EIP (floating IP) ID.
* `instance_id` - (Required, String, ForceNew) Target instance ID. For EcsInstance and ElbInstance this is the server or load balancer ID; for NetworkInterface this is the network interface (port) ID.
* `fixed_ip_address` - (Optional, String, ForceNew) Fixed (private) IP used for the association. Set on create when the target has multiple network interfaces or fixed IPs; otherwise the first available fixed IP is chosen. On read, populated from the API.
* `instance_type` - (Optional, String, ForceNew) Instance type: EcsInstance, ElbInstance, or NetworkInterface.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `floating_ip_address` - Public floating IP address from the API.


## Import

Import format is `allocation_id`.

```shell
terraform import edgenext_eip_association.example f5390261-241e-43e0-ab1b-5bffaf334c81
```

