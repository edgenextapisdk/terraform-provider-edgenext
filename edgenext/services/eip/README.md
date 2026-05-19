# EdgeNext EIP Services

This package provides Terraform resources and data sources for EdgeNext Elastic IP (floating IP) operations: listing EIPs and associating an existing EIP with an ECS instance, ELB load balancer, or network interface (port).

EIP lifecycle (allocate/release bandwidth, orders, and so on) is outside this module; use the platform or future resources for creating floating IPs. This module focuses on **association** and **discovery**.

The former data source `edgenext_ecs_floating_ips` was moved here as `edgenext_eip_floating_ips` (same ECS list API).

## Resources

### EIP Association
- **Resource**: `edgenext_eip_association` (`ResourceENEIPAssociation`)
- **File**: `resource_en_eip_association.go`
- **Description**: Bind an existing EIP to a private fixed IP on a port. Resolves `port_id` and `fixed_ip_address` from the target type, then calls the ECS floating IP relation API. No update; all arguments force replacement.

## EIP Association Behavior

### APIs

| Operation | Path | Notes |
|-----------|------|--------|
| Associate | `POST /ecs/openapi/v2/floatingips/ports/fixed_ip/relation` | `action: add` with `floating_ip.id`, `port_id`, `fixed_ip_address` |
| Disassociate | same path | `action: remove` with `floating_ip.id` only |
| Read / list | `POST /ecs/openapi/v2/floatingips/list` | Association exists when `port_id` and `fixed_ip_address` are set |

### `instance_type` and target resolution

| `instance_type` | `instance_id` meaning | Resolution API |
|-----------------|----------------------|----------------|
| `EcsInstance` (default) | ECS server ID | `POST /ecs/openapi/v2/instance/list_ports` → pick port and fixed IP |
| `ElbInstance` | Load balancer ID | `POST /elb/openapi/v2/load_balancers/get` → `vip_port_id`, `vip_address` |
| `NetworkInterface` | Network interface (port) ID | `POST /ecs/openapi/v2/ports/internal_ip_list` → fixed IP on that port |

When `fixed_ip_address` is omitted, the first suitable fixed IP is chosen (first port / first IP). When set, it must exist on the resolved target.

### Create validation

Before associate, the provider lists the EIP and fails if it is **already bound** (`port_id` and `fixed_ip_address` both non-empty).

### ForceNew / no update

There is no `Update` handler. Changing `allocation_id`, `instance_id`, `instance_type`, or `fixed_ip_address` replaces the association (destroy then create).

| Argument | ForceNew | Computed on read |
|----------|----------|------------------|
| `allocation_id` | yes | synced from API |
| `instance_id` | yes | — |
| `instance_type` | yes | — |
| `fixed_ip_address` | yes | yes (from API `fixed_ip_address`) |
| `floating_ip_address` | — | yes |

Terraform resource ID is the EIP `allocation_id`.

## Data Sources

### EIP Floating IPs
- **Data Source**: `edgenext_eip_floating_ips` (`DataSourceENEIPFloatingIps`)
- **File**: `data_source_en_eip_floating_ips.go`
- **Description**: List floating IPs via `floatingips/list`. Replaces `edgenext_ecs_floating_ips`. Nested `network_interface_id` maps from API `port_id`.

## File Structure

```
edgenext/services/eip/
├── README.md                                      # This documentation
├── resource_en_eip_association.go                 # EIP association resource
├── resource_en_eip_association.md                 # Per-resource doc (gendoc / registry)
├── data_source_en_eip_floating_ips.go             # Floating IP list data source
└── data_source_en_eip_floating_ips.md             # Per-data-source doc
```

## Usage Examples

See `examples/eip/main.tf`. Typical flow:

### List EIPs and associate to an ECS instance

```hcl
data "edgenext_eip_floating_ips" "all" {
  limit = 50
}

resource "edgenext_eip_association" "web" {
  allocation_id    = data.edgenext_eip_floating_ips.all.floating_ips[0].id
  instance_id      = var.ecs_instance_id
  fixed_ip_address = "172.31.0.32"
}
```

### Associate to an ELB load balancer

```hcl
resource "edgenext_eip_association" "lb" {
  allocation_id = var.eip_id
  instance_id   = var.loadbalancer_id
  instance_type = "ElbInstance"
}
```

### Associate to a network interface (port)

```hcl
resource "edgenext_eip_association" "eni" {
  allocation_id    = var.eip_id
  instance_id      = var.network_interface_id
  instance_type    = "NetworkInterface"
  fixed_ip_address = "172.31.0.32"
}
```

## Import

| Resource | Import ID | Notes |
|----------|-----------|--------|
| `edgenext_eip_association` | `allocation_id` | Set `instance_id`, `instance_type`, and `fixed_ip_address` in configuration to match the existing binding |

```shell
terraform import edgenext_eip_association.web f5390261-241e-43e0-ab1b-5bffaf334c81
```

