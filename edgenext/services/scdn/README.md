# SCDN Service

This package provides Terraform resources and data sources for managing EdgeNext SCDN (Secure Content Delivery Network) domains and related configurations.

## Resources

### edgenext_scdn_domain

Manages SCDN domains with origin server configurations.

#### Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `domain` | `string` | Yes | The domain name to be added to SCDN |
| `group_id` | `number` | No | The ID of the domain group |
| `exclusive_resource_id` | `number` | No | The ID of the exclusive resource package |
| `remark` | `string` | No | The remark for the domain |
| `tpl_id` | `number` | No | The template ID to be applied to the domain |
| `protect_status` | `string` | No | The edge node type. Valid values: `back_source`, `scdn`, `exclusive` (default: `scdn`) |
| `tpl_recommend` | `string` | No | The template recommendation status |
| `app_type` | `string` | No | The application type |
| `origins` | `list(object)` | Yes | The origin server configuration |

#### Origins Block

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `protocol` | `number` | Yes | The origin protocol. Valid values: `0` (HTTP), `1` (HTTPS) |
| `listen_ports` | `list(number)` | Yes | The listening ports of the origin server |
| `origin_protocol` | `number` | Yes | The origin protocol. Valid values: `0` (HTTP), `1` (HTTPS), `2` (Follow) |
| `load_balance` | `number` | Yes | The load balancing method. Valid values: `0` (IP hash), `1` (Round robin), `2` (Cookie) |
| `origin_type` | `number` | Yes | The origin type. Valid values: `0` (IP), `1` (Domain) |
| `records` | `list(object)` | Yes | The origin records |

#### Records Block

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `view` | `string` | Yes | The view of the record |
| `value` | `string` | Yes | The value of the record (IP address or domain) |
| `port` | `number` | Yes | The port of the record |
| `priority` | `number` | Yes | The priority of the record |

#### Attributes

| Name | Type | Description |
|------|------|-------------|
| `id` | `number` | The ID of the domain |
| `access_progress` | `string` | The access progress status |
| `access_mode` | `string` | The access mode (ns or cname) |
| `ei_forward_status` | `string` | The explicit/implicit forwarding status |
| `cname` | `object` | The CNAME information |
| `use_my_cname` | `number` | The CNAME resolution status |
| `use_my_dns` | `number` | The DNS hosting status |
| `ca_status` | `string` | The certificate binding status |
| `access_progress_desc` | `string` | The description of the access progress status |
| `has_origin` | `bool` | Whether the domain has origin configuration |
| `ca_id` | `number` | The certificate ID |
| `created_at` | `string` | The creation timestamp |
| `updated_at` | `string` | The last update timestamp |
| `pri_domain` | `string` | The primary domain |

#### Example

```hcl
resource "edgenext_scdn_domain" "example" {
  domain        = "example.com"
  group_id      = 123
  remark        = "Example SCDN domain"
  protect_status = "scdn"
  
  origins {
    protocol        = 0
    listen_ports    = [80, 443]
    origin_protocol = 0
    load_balance    = 1
    origin_type     = 0
    
    records {
      view     = "default"
      value    = "1.1.1.1"
      port     = 80
      priority = 10
    }
  }
}
```

### edgenext_scdn_origin

Manages individual origin servers for SCDN domains.

#### Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `domain_id` | `number` | Yes | The ID of the domain to add origins to |
| `protocol` | `number` | Yes | The origin protocol. Valid values: `0` (HTTP), `1` (HTTPS) |
| `listen_ports` | `list(number)` | Yes | The listening ports of the origin server |
| `origin_protocol` | `number` | Yes | The origin protocol. Valid values: `0` (HTTP), `1` (HTTPS), `2` (Follow) |
| `load_balance` | `number` | Yes | The load balancing method. Valid values: `0` (IP hash), `1` (Round robin), `2` (Cookie) |
| `origin_type` | `number` | Yes | The origin type. Valid values: `0` (IP), `1` (Domain) |
| `records` | `list(object)` | Yes | The origin records |

#### Attributes

| Name | Type | Description |
|------|------|-------------|
| `id` | `number` | The ID of the origin |
| `listen_port` | `number` | The listening port of the origin server (single port from API) |

### edgenext_scdn_cert_binding

Manages certificate bindings for SCDN domains.

#### Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `domain_id` | `number` | Yes | The ID of the domain to bind the certificate to |
| `ca_id` | `number` | Yes | The ID of the certificate to bind |

#### Attributes

| Name | Type | Description |
|------|------|-------------|
| `id` | `string` | The unique identifier for this certificate binding |

## Data Sources

### edgenext_scdn_domain

Retrieves information about a specific SCDN domain.

#### Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `domain` | `string` | Yes | The domain name to query |
| `result_output_file` | `string` | No | Used to save results to a file |

#### Attributes

All attributes from the `edgenext_scdn_domain` resource are available as computed attributes.

### edgenext_scdn_domains

Retrieves a list of SCDN domains with optional filtering.

#### Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `page` | `number` | No | The page number for pagination (default: 1) |
| `page_size` | `number` | No | The page size for pagination (default: 100) |
| `access_progress` | `string` | No | Filter by access progress status |
| `group_id` | `number` | No | Filter by domain group ID |
| `domain` | `string` | No | Filter by domain name (fuzzy search) |
| `remark` | `string` | No | Filter by remark (fuzzy search) |
| `origin_ip` | `string` | No | Filter by origin IP |
| `ca_status` | `string` | No | Filter by certificate binding status |
| `access_mode` | `string` | No | Filter by access mode |
| `protect_status` | `string` | No | Filter by edge node type |
| `exclusive_resource_id` | `number` | No | Filter by exclusive resource package ID |
| `result_output_file` | `string` | No | Used to save results to a file |

#### Attributes

| Name | Type | Description |
|------|------|-------------|
| `total` | `number` | The total number of domains |
| `domains` | `list(object)` | The list of domains |

### edgenext_scdn_origin

Retrieves information about a specific SCDN origin.

#### Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `origin_id` | `number` | Yes | The ID of the origin to query |
| `domain_id` | `number` | Yes | The ID of the domain that owns the origin |
| `result_output_file` | `string` | No | Used to save results to a file |

#### Attributes

All attributes from the `edgenext_scdn_origin` resource are available as computed attributes.

### edgenext_scdn_origins

Retrieves a list of SCDN origins for a specific domain.

#### Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `domain_id` | `number` | Yes | The ID of the domain to query origins for |
| `result_output_file` | `string` | No | Used to save results to a file |

#### Attributes

| Name | Type | Description |
|------|------|-------------|
| `total` | `number` | The total number of origins |
| `origins` | `list(object)` | The list of origins |

## API Endpoints

This service implements the following EdgeNext SCDN API v5 endpoints:

- `GET /api/v5/domains` - List domains
- `GET /api/v5/domains/simple` - List simple domains
- `POST /api/v5/domains` - Create domain
- `PUT /api/v5/domains` - Update domain
- `POST /api/v5/domains/bind_cert` - Bind certificate to domain
- `POST /api/v5/domains/unbind_cert` - Unbind certificate from domain
- `DELETE /api/v5/domains` - Delete domain
- `POST /api/v5/domains_disable` - Disable domain
- `POST /api/v5/domains_enable` - Enable domain
- `POST /api/v5/domains/access_refresh` - Refresh domain access status
- `POST /api/v5/domains/domains_export` - Export domains
- `POST /api/v5/domains/origins` - Add origins
- `PUT /api/v5/domains/origins` - Update origins
- `DELETE /api/v5/domains/origins` - Delete origins
- `GET /api/v5/domains/origins` - List origins
- `POST /api/v5/domains/nodes_switch` - Switch domain nodes
- `POST /api/v5/domains/access_switch` - Switch domain access mode
- `GET /api/v5/domains/access_progress` - Get access progress status
- `PUT /api/v5/domains/base_settings` - Update domain base settings
- `GET /api/v5/domains/base_settings` - Get domain base settings
- `POST /api/v5/brief_domains` - List brief domains
- `GET /api/v5/domains/templates` - Get domain templates
- `POST /api/v5/domains/access_info_download` - Download access information

## Error Handling

The service includes comprehensive error handling for common scenarios:

- **Not Found Errors**: When resources don't exist
- **Rate Limit Errors**: When API rate limits are exceeded
- **Authentication Errors**: When API credentials are invalid
- **Validation Errors**: When required parameters are missing or invalid

## Testing

The package includes comprehensive unit tests covering:

- Service layer functionality
- Resource and data source schemas
- API request/response handling
- Error scenarios
- Data validation

Run tests with:

```bash
go test ./edgenext/services/scdn/...
```

## Examples

### Complete SCDN Setup

```hcl
# Create SCDN domain
resource "edgenext_scdn_domain" "main" {
  domain        = "cdn.example.com"
  group_id      = 123
  remark        = "Main CDN domain"
  protect_status = "scdn"
  
  origins {
    protocol        = 0
    listen_ports    = [80, 443]
    origin_protocol = 0
    load_balance    = 1
    origin_type     = 0
    
    records {
      view     = "default"
      value    = "1.1.1.1"
      port     = 80
      priority = 10
    }
    
    records {
      view     = "default"
      value    = "2.2.2.2"
      port     = 80
      priority = 20
    }
  }
}

# Add additional origin
resource "edgenext_scdn_origin" "backup" {
  domain_id      = edgenext_scdn_domain.main.id
  protocol       = 1
  listen_ports   = [443]
  origin_protocol = 1
  load_balance   = 1
  origin_type    = 1
  
  records {
    view     = "default"
    value    = "backup.example.com"
    port     = 443
    priority = 10
  }
}

# Bind SSL certificate
resource "edgenext_scdn_cert_binding" "ssl" {
  domain_id = edgenext_scdn_domain.main.id
  ca_id     = 456
}

# Query domain information
data "edgenext_scdn_domain" "info" {
  domain = "cdn.example.com"
}

# Query all domains
data "edgenext_scdn_domains" "all" {
  access_progress = "enabled"
  protect_status  = "scdn"
}
```

### Using Data Sources

```hcl
# Get specific domain
data "edgenext_scdn_domain" "example" {
  domain = "example.com"
}

# Get all domains with filters
data "edgenext_scdn_domains" "filtered" {
  page            = 1
  page_size       = 50
  access_progress = "enabled"
  protect_status  = "scdn"
}

# Get origins for a domain
data "edgenext_scdn_origins" "domain_origins" {
  domain_id = 123
}

# Get specific origin
data "edgenext_scdn_origin" "specific" {
  origin_id = 456
  domain_id = 123
}
```
