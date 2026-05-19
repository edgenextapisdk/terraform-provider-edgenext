# Terraform Provider for EdgeNext

[![Go Report Card](https://goreportcard.com/badge/github.com/edgenextapisdk/terraform-provider-edgenext)](https://goreportcard.com/report/github.com/edgenextapisdk/terraform-provider-edgenext)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.23+-blue.svg)](https://golang.org)
[![Terraform](https://img.shields.io/badge/Terraform-1.0+-purple.svg)](https://terraform.io)

Terraform provider for managing EdgeNext services, including CDN, SSL, OSS, ECS, ELB, EIP, RDS, SDNS, and SCDN.

## Supported Services

- **CDN**: Domain configuration and cache operations
- **SSL**: Certificate lifecycle management
- **OSS**: Bucket and object management
- **ECS**: VPC, router, ENI, security groups, tags, and instance power/reboot resources and data sources
- **ELB**: Load balancers, listeners, target groups, backends, TLS certificates, and L7 policies/rules
- **EIP**: Floating IP list and association with ECS, ELB, or network interfaces
- **RDS**: Relational database resources and data sources
- **SDNS**: Domain group and record management
- **SCDN**: Domain/origin/template/cache/security/log modules

Service-level documentation:

- [CDN](edgenext/services/cdn/README.md)
- [SSL](edgenext/services/ssl/README.md)
- [OSS](edgenext/services/oss/README.md)
- [ECS](edgenext/services/ecs/README.md)
- [ELB](edgenext/services/elb/README.md)
- [EIP](edgenext/services/eip/README.md)
- [RDS](edgenext/services/rds/README.md)
- [SCDN](edgenext/services/scdn/README.md)

## Installation

### Terraform

```hcl
terraform {
  required_providers {
    edgenext = {
      source  = "edgenextapisdk/edgenext"
      version = "~> 1.0"
    }
  }
}
```

### Build from Source

```bash
git clone https://github.com/edgenextapisdk/terraform-provider-edgenext.git
cd terraform-provider-edgenext
go build -o terraform-provider-edgenext
```

## Provider Configuration

```hcl
provider "edgenext" {
  access_key = var.access_key
  secret_key = var.secret_key
  endpoint   = var.endpoint
  region     = var.region # optional
}
```

Environment variables are supported:

```bash
export EDGENEXT_ACCESS_KEY="your-access-key"
export EDGENEXT_SECRET_KEY="your-secret-key"
export EDGENEXT_ENDPOINT="https://cdn.api.edgenext.com"
export EDGENEXT_REGION="us-east-1"
```

Provider arguments:

- `access_key` (Required): EdgeNext access key.
- `secret_key` (Required): EdgeNext secret key.
- `endpoint` (Required): EdgeNext API endpoint.
- `region` (Optional): Default region.

## Quick Examples

### CDN

```hcl
resource "edgenext_cdn_domain" "example" {
  domain = "example.com"
  area   = "global"
  type   = "page"
}
```

### OSS

```hcl
resource "edgenext_oss_bucket" "assets" {
  bucket = "my-assets-bucket"
  acl    = "private"
}

resource "edgenext_oss_object" "logo" {
  bucket = edgenext_oss_bucket.assets.bucket
  key    = "images/logo.png"
  source = "${path.module}/assets/logo.png"
}
```

### ECS

```hcl
resource "edgenext_ecs_vpc" "example" {
  name = "example-vpc"
  subnet {
    name       = "example-subnet"
    ip_version = 4
    cidr       = "172.31.1.0/24"
  }
}

resource "edgenext_ecs_vpc_subnet" "extra" {
  vpc_id = edgenext_ecs_vpc.example.id
  name   = "example-subnet-2"
  cidr   = "172.31.2.0/24"
}
```

### ELB

```hcl
resource "edgenext_elb_listener" "https" {
  loadbalancer_id             = var.loadbalancer_id
  name                        = "https"
  protocol                    = "TERMINATED_HTTPS"
  protocol_port               = 443
  default_tls_certificate_ref = edgenext_elb_certificate.site.certificate_ref
}

resource "edgenext_elb_target_group" "app" {
  listener_id  = edgenext_elb_listener.https.id
  name         = "app"
  protocol     = "HTTP"
  lb_algorithm = "ROUND_ROBIN"
  health_monitor {
    name           = "check"
    type           = "HTTP"
    max_retries    = 3
    delay          = 10
    timeout        = 5
    http_method    = "GET"
    url_path       = "/health"
    expected_codes = "200"
  }
}

resource "edgenext_elb_target_group_attachment" "app1" {
  target_group_id = edgenext_elb_target_group.app.id
  address         = "192.168.0.10"
  protocol_port   = 8080
}
```

See [examples/elb](examples/elb/).

### EIP

```hcl
data "edgenext_eip_floating_ips" "all" {
  limit = 50
}

resource "edgenext_eip_association" "web" {
  allocation_id = var.eip_id
  instance_id   = var.ecs_instance_id
  # instance_type = "EcsInstance"  # default; or ElbInstance, NetworkInterface
  fixed_ip_address = "172.31.0.32"
}
```

See [examples/eip](examples/eip/).

## Provider Registration Snapshot

Registration matches `edgenext/provider.go` (see also [provider.md](edgenext/provider.md) for the registry website list).

### ECS resources

- `edgenext_ecs_key_pair`
- `edgenext_ecs_vpc`
- `edgenext_ecs_vpc_subnet`
- `edgenext_ecs_router`
- `edgenext_ecs_router_port`
- `edgenext_ecs_network_interface`
- `edgenext_ecs_network_interface_instance_binding`
- `edgenext_ecs_security_group`
- `edgenext_ecs_security_group_rule`
- `edgenext_ecs_tag`
- `edgenext_ecs_instance_tag`
- `edgenext_ecs_instance_power`
- `edgenext_ecs_instance_reboot`

### ECS data sources

- `edgenext_ecs_instances`
- `edgenext_ecs_images`
- `edgenext_ecs_key_pairs`
- `edgenext_ecs_vpcs`
- `edgenext_ecs_external_gateways`
- `edgenext_ecs_vpc_subnets`
- `edgenext_ecs_routers`
- `edgenext_ecs_router_ports`
- `edgenext_ecs_network_interfaces`
- `edgenext_ecs_security_groups`
- `edgenext_ecs_disks`
- `edgenext_ecs_tags`
- `edgenext_ecs_security_group_rules`
- `edgenext_ecs_instance_tags`

### ELB resources

- `edgenext_elb_certificate`
- `edgenext_elb_listener`
- `edgenext_elb_target_group`
- `edgenext_elb_target_group_attachment`
- `edgenext_elb_l7_policy`
- `edgenext_elb_l7_rule`

### ELB data sources

- `edgenext_elb_load_balancers`
- `edgenext_elb_certificates`
- `edgenext_elb_listeners`
- `edgenext_elb_target_groups`
- `edgenext_elb_target_group_attachments`
- `edgenext_elb_l7_policies`
- `edgenext_elb_l7_rules`

### EIP

- Data source: `edgenext_eip_floating_ips` (replaces former `edgenext_ecs_floating_ips`)
- Resource: `edgenext_eip_association`

### RDS resources

- `edgenext_rds_backup`
- `edgenext_rds_backup_policy`
- `edgenext_rds_backup_policy_associate_instance`
- `edgenext_rds_database`
- `edgenext_rds_account`
- `edgenext_rds_account_privilege`
- `edgenext_rds_account_root_password`

### RDS data sources

- `edgenext_rds_instances`
- `edgenext_rds_databases`
- `edgenext_rds_accounts`
- `edgenext_rds_backups`
- `edgenext_rds_backup_policies`
- `edgenext_rds_backup_policy_associate_instances`

### ECS (implemented, not registered)

- `edgenext_ecs_instance`
- `edgenext_ecs_image`
- `edgenext_ecs_floating_ip`
- `edgenext_ecs_disk`

CDN, SSL, OSS, SDNS, and SCDN modules include additional resources registered from subpackages; see `edgenext/provider.go` and [provider.md](edgenext/provider.md).

## Documentation

- [Provider doc source](edgenext/provider.md)
- [Generated website docs](website/docs/index.html.markdown)
- [Examples](examples/)
- [Changelog](CHANGELOG.md)

## Development

```bash
make build
make test
make testacc
make lint
make fmt
make doc
```

## Contributing

1. Open an issue to discuss bug fixes or feature requests.
2. Create a branch and implement changes with tests.
3. Update documentation when behavior changes.
4. Submit a pull request with a clear summary.

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE).
