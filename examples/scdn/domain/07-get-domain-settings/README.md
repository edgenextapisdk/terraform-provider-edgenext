# Example 7: Query Domain Base Settings

This example demonstrates how to query domain base settings including proxy host, proxy SNI, and domain redirect configurations.

## Prerequisites

1. Valid EdgeNext API credentials (access_key and secret_key)
2. An existing SCDN domain ID

## Files

- `main.tf` - Main Terraform configuration
- `terraform.tfvars.example` - Example configuration file

## Usage

1. Copy `terraform.tfvars.example` to `terraform.tfvars`:
   ```bash
   cp terraform.tfvars.example terraform.tfvars
   ```

2. Update `terraform.tfvars` with your credentials and domain ID

3. Initialize Terraform:
   ```bash
   terraform init
   ```

4. Query the domain base settings:
   ```bash
   terraform apply
   ```

## Outputs

After successful query, you'll see:
- `domain_settings` - Domain base settings including:
  - `domain_id` - Domain ID
  - `proxy_host` - Proxy host configuration
  - `proxy_sni` - Proxy SNI configuration
  - `domain_redirect` - Domain redirect configuration

## Settings Explanation

### Proxy Host
- `proxy_host` - The proxy host header value
- `proxy_host_type` - The type of proxy host configuration

### Proxy SNI
- `proxy_sni` - The proxy SNI value
- `status` - The proxy SNI status (on/off)

### Domain Redirect
- `status` - Redirect status (on/off)
- `jump_to` - Target URL for redirect
- `jump_type` - Type of redirect (301/302)

