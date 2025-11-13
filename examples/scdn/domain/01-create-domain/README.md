# Example 1: Create SCDN Domain

This example demonstrates how to create a basic SCDN domain with a single origin server.

## Prerequisites

1. Valid EdgeNext API credentials (access_key and secret_key)
2. A domain name that you want to configure

## Files

- `main.tf` - Main Terraform configuration
- `terraform.tfvars.example` - Example configuration file

## Usage

1. Copy `terraform.tfvars.example` to `terraform.tfvars`:
   ```bash
   cp terraform.tfvars.example terraform.tfvars
   ```

2. Update `terraform.tfvars` with your credentials and domain information:
   - `domain_name` - Your domain name (e.g., "example.com")
   - `group_id` - Group ID (default: 1)
   - `origin_ip` - Your origin server IP address
   - `origin_port` - Origin server port (default: 80)

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

## Outputs

After successful creation, you'll see:
- `domain_id` - The ID of the created domain
- `domain_name` - The domain name
- `cname` - CNAME record that needs to be configured in your DNS
- `access_progress` - Domain access configuration progress

## Next Steps

After creating the domain, you can:
- Use the `domain_id` in other examples to add origins
- Query the domain information using Example 3
- Configure the CNAME record in your DNS provider
