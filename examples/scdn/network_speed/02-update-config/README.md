# Update Network Speed Config

This example demonstrates how to update network speed configuration for a template.

## Usage

1. Copy `terraform.tfvars.example` to `terraform.tfvars` and fill in your credentials:

```bash
cp terraform.tfvars.example terraform.tfvars
```

2. Edit `terraform.tfvars` with your actual values:

```hcl
access_key   = "your_access_key"
secret_key   = "your_secret_key"
business_id  = 1246
business_type = "tpl"

# Domain proxy configuration
domain_proxy_conf = {
  proxy_connect_timeout = 30
  fails_timeout         = 10
  keep_new_src_time     = 30
  max_fails             = 30
  proxy_keepalive       = 0
}

# HTTPS configuration
https_config = {
  status                   = "on"
  http2https               = "off"
  http2https_port          = 443
  http2                    = "off"
  hsts                     = "off"
  ocsp_stapling            = "off"
  min_version              = "TLSv1.2"
  ciphers_preset           = "default"
  custom_encrypt_algorithm = [
    "ECDHE-ECDSA",
    "ECDHE-RSA",
    "DHE-RSA",
    "EECDH+CHACHA20",
    "EECDH+AES128",
    "EECDH+AES256",
    "RSA+AES128",
    "RSA+AES256",
    "EECDH+3DES",
    "RSA+3DES"
  ]
}
```

3. Initialize and apply:

```bash
terraform init
terraform plan
terraform apply
```

## Notes

- This example demonstrates **dynamic configuration** using variables
- All configuration values can be customized through `terraform.tfvars`
- `business_id` is the template ID for `business_type = "tpl"`
- `business_id` is the user ID for `business_type = "global"`
- **Important**: You must have an existing template before configuring network speed settings
  - For `business_type = "tpl"`, you need to create a rule template first using `edgenext_scdn_rule_template` resource
  - You can get the template ID from the rule template list API or dashboard
  - The configuration will be applied to the existing template
- The resource supports updating all 15 configuration groups:
  - `domain_proxy_conf`, `upstream_redirect`, `customized_req_headers`, `resp_headers`
  - `upstream_uri_change`, `source_site_protect`, `slice`, `https`
  - `page_gzip`, `webp`, `upload_file`, `websocket`, `mobile_jump`, `custom_page`, `upstream_check`
- The resource will create or update the configuration based on the `business_id` and `business_type`

## Prerequisites

Before using this example, you need to:

1. **Create a rule template** first (if using `business_type = "tpl"`):
   ```hcl
   resource "edgenext_scdn_rule_template" "example" {
     name        = "my-template"
     description = "Template for network speed config"
     app_type    = "http"
   }
   ```

2. **Get the template ID** from the created template or list templates:
   ```hcl
   data "edgenext_scdn_rule_templates" "example" {
     # List all templates to find the ID
   }
   ```

3. **Use the template ID** as `business_id` in your `terraform.tfvars`

## Troubleshooting

If you encounter an error like "app not exist (code: 1010)":

1. **Verify the template exists**: Make sure you've created a rule template first
2. **Check the template ID**: Verify the `business_id` matches an existing template ID
3. **List templates**: Use the rule template data source to verify:
   ```hcl
   data "edgenext_scdn_rule_templates" "list" {}
   ```
4. **Verify permissions**: Ensure your API credentials have permission to access the template

## Configuration Groups

This example demonstrates updating all 15 configuration groups. Note:

- **mobile_jump**: When `status = "on"`, `jump_url` is **required** and must be a valid URL
- **domain_proxy_conf**: Domain proxy configuration (connection timeout, failure timeout, etc.)
- **https**: HTTPS configuration (TLS settings, HTTP2, HSTS, etc.)

## Outputs

- `config_id`: Network speed configuration resource ID
- `domain_proxy_conf`: Updated domain proxy configuration
- `https_config`: Updated HTTPS configuration

