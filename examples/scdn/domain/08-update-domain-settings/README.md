# Example 8: Update Domain Base Settings

This example demonstrates how to update domain base settings including proxy host, proxy SNI, and domain redirect configurations.

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

2. Update `terraform.tfvars` with:
   - Your credentials
   - `domain_id` - The ID of an existing domain
   - Configure any settings you want to update (proxy_host, proxy_sni, or domain_redirect)

3. Initialize Terraform:
   ```bash
   terraform init
   ```

4. Review the execution plan:
   ```bash
   terraform plan
   ```

5. Apply the configuration:
   ```bash
   terraform apply
   ```

## Configuration Options

### Proxy Host
Configure custom proxy host header:
```hcl
enable_proxy_host = true
proxy_host       = "example.com"
proxy_host_type  = "custom"
```

### Proxy SNI
Configure proxy SNI settings:
```hcl
enable_proxy_sni = true
proxy_sni        = "example.com"
proxy_sni_status = "on"
```

### Domain Redirect
Configure domain redirect:
```hcl
enable_redirect    = true
redirect_status    = "on"
redirect_jump_to   = "https://www.example.com"
redirect_jump_type = "301"
```

## Outputs

After successful update, you'll see:
- `domain_id` - The domain ID

## Notes

- All configuration blocks (proxy_host, proxy_sni, domain_redirect) are optional
- You can configure one or more settings at a time
- Settings can be updated by modifying the configuration and reapplying
- To remove a setting, remove it from the configuration or set the enable flag to false

