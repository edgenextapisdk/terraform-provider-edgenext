Associates an EdgeNext EIP with an ECS instance, ELB load balancer, or network interface.

Example Usage

```hcl
resource "edgenext_eip_association" "example" {
  allocation_id      = "f5390261-241e-43e0-ab1b-5bffaf334c81"
  instance_id        = "518031d1-f66d-416d-a66a-96ab91f4def9"
  fixed_ip_address   = "172.31.0.32"
}
```

Import

Import format is `allocation_id`.

```shell
terraform import edgenext_eip_association.example f5390261-241e-43e0-ab1b-5bffaf334c81
```
