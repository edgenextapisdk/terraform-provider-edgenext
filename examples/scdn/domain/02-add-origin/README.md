# Example 2: Add Origin to Existing Domain

This example demonstrates how to add an additional origin server to an existing SCDN domain.

## Prerequisites

1. Valid EdgeNext API credentials (access_key and secret_key)
2. An existing SCDN domain ID (you can get this from Example 1)

## Files

- `main.tf` - Main Terraform configuration
- `terraform.tfvars.example` - Example configuration file

## Usage

1. Copy `terraform.tfvars.example` to `terraform.tfvars`:
   ```bash
   cp terraform.tfvars.example terraform.tfvars
   ```

2. Update `terraform.tfvars` with:
   - `domain_id` - The ID of an existing domain (from Example 1 output)
   - `origin_ip` - IP address of the new origin server
   - `origin_port` - Port of the origin server (default: 80)
   - `priority` - Priority for load balancing (default: 15)

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
- `origin_id` - The ID of the created origin
- `domain_id` - The domain ID this origin belongs to

## Notes

- You can add multiple origins to the same domain by running this example multiple times with different IPs
- Higher priority values indicate higher priority in load balancing
