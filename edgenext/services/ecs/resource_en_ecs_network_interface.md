Use this resource to create and manage ECS network interfaces (ENI / port).

Example Usage

```hcl
resource "edgenext_ecs_network_interface" "example" {
  region                = "tokyo-a"
  name                  = "example-eni"
  description           = "for application"
  network_id            = "0e07db22-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  subnet_id             = "50a0f20a-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  device_id             = "80e47fca-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  port_security_enabled = true
  security_groups       = ["aa2b7c0d-xxxx-xxxx-xxxx-xxxxxxxxxxxx"]
  floating_ip_address   = "148.222.161.86"
}
```

Import

Import format is `region/port_id`.

```shell
terraform import edgenext_ecs_network_interface.example tokyo-a/29faf396-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```

Argument Reference

* `region` - (Required) Region.
* `name` - (Required) Port name.
* `description` - (Optional) Port description.
* `network_id` - (Required) Network ID. Cannot be changed after creation.
* `subnet_id` - (Required) Subnet ID. Cannot be changed after creation.
* `device_id` - (Optional) Server ID to attach.
* `port_security_enabled` - (Optional) Port security switch.
* `security_groups` - (Optional) Security group IDs.
* `floating_ip_address` - (Optional) Floating IP address to bind.

Attributes Reference

* `id` - Importable ID in format `region/port_id`.
* `tenant_id`, `project_id`, `status`, `device_owner`
* `fixed_ips` - Fixed IP list with `subnet_id`, `ip_address`, `floating_ip`.
* `qos_policy_id`, `tags`, `created_at`, `updated_at`, `revision_number`
* `mac_address`, `binding_vnic_type`, `server_name`, `network_name`, `ipv4`, `ipv6`
