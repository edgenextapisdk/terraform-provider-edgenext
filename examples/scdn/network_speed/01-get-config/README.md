# Get Network Speed Config

This example demonstrates how to query network speed configuration for a template.

## Usage

1. Copy `terraform.tfvars.example` to `terraform.tfvars` and fill in your credentials:

```bash
cp terraform.tfvars.example terraform.tfvars
```

2. Edit `terraform.tfvars` with your actual values:

```hcl
access_key   = "your_access_key"
secret_key   = "your_secret_key"
business_id  = 777
business_type = "tpl"
```

3. Initialize and apply:

```bash
terraform init
terraform plan
terraform apply
```

## Notes

- This example queries the network speed configuration for a template
- `business_id` is the template ID for `business_type = "tpl"`
- `business_id` is the user ID for `business_type = "global"`
- You can optionally specify `config_groups` to retrieve only specific configuration groups

## Outputs

This example outputs all network speed configuration groups:

- `domain_proxy_conf`: Domain proxy configuration (connection timeout, failure timeout, etc.)
- `upstream_redirect`: Upstream redirect configuration (301/302 follow)
- `customized_req_headers`: Customized request headers configuration
- `resp_headers`: Response headers configuration
- `upstream_uri_change`: Upstream URI change configuration
- `source_site_protect`: Source site protection configuration
- `slice`: Range request configuration
- `https_config`: HTTPS configuration (TLS settings, HTTP2, HSTS, etc.)
- `page_gzip`: Page Gzip configuration
- `webp`: WebP format configuration
- `upload_file`: Upload file configuration (file size limits)
- `websocket`: WebSocket configuration
- `mobile_jump`: Mobile jump configuration
- `custom_page`: Custom page configuration
- `upstream_check`: Upstream check configuration (health check settings)
- `all_configs`: All configurations combined in a single object


