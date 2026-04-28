# Terraform Provider for EdgeNext

[![Go Report Card](https://goreportcard.com/badge/github.com/edgenextapisdk/terraform-provider-edgenext)](https://goreportcard.com/report/github.com/edgenextapisdk/terraform-provider-edgenext)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.23+-blue.svg)](https://golang.org)
[![Terraform](https://img.shields.io/badge/Terraform-1.0+-purple.svg)](https://terraform.io)

Terraform provider for managing EdgeNext services, including CDN, SSL, OSS, ECS, SDNS, and SCDN.

## Supported Services

- CDN: Domain configuration and cache operations
- SSL: Certificate lifecycle management
- OSS: Bucket and object management
- ECS: Network, security, tag, and instance operation resources/data sources
- SDNS: Domain group and record management
- SCDN: Domain/origin/template/cache/security/log modules

Service-level documentation:

- [CDN](edgenext/services/cdn/README.md)
- [SSL](edgenext/services/ssl/README.md)
- [OSS](edgenext/services/oss/README.md)
- [ECS](edgenext/services/ecs/README.md)
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

## ECS Registration Snapshot

Current ECS resources registered in `edgenext/provider.go`:

- `edgenext_ecs_key_pair`
- `edgenext_ecs_vpc`
- `edgenext_ecs_vpc_subnet`
- `edgenext_ecs_router`
- `edgenext_ecs_router_port`
- `edgenext_ecs_network_interface`
- `edgenext_ecs_network_interface_instance_binding`
- `edgenext_ecs_network_interface_floating_ip_binding`
- `edgenext_ecs_security_group`
- `edgenext_ecs_security_group_rule`
- `edgenext_ecs_tag`
- `edgenext_ecs_instance_tag`
- `edgenext_ecs_instance_power`
- `edgenext_ecs_instance_reboot`

Current ECS data sources registered in `edgenext/provider.go`:

- `edgenext_ecs_instances`
- `edgenext_ecs_images`
- `edgenext_ecs_key_pairs`
- `edgenext_ecs_vpcs`
- `edgenext_ecs_external_gateways`
- `edgenext_ecs_vpc_subnets`
- `edgenext_ecs_routers`
- `edgenext_ecs_router_ports`
- `edgenext_ecs_floating_ips`
- `edgenext_ecs_network_interfaces`
- `edgenext_ecs_security_groups`
- `edgenext_ecs_disks`
- `edgenext_ecs_tags`
- `edgenext_ecs_security_group_rules`
- `edgenext_ecs_instance_tags`

ECS resources currently present in code but not registered in provider:

- `edgenext_ecs_instance`
- `edgenext_ecs_image`
- `edgenext_ecs_floating_ip`
- `edgenext_ecs_disk`

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
