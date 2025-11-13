# Example 5: Query SCDN Origins

This example demonstrates how to query all origins for a specific SCDN domain.

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
   - `domain_id` - The ID of an existing domain (you can get this from Example 1 or Example 3)

3. Initialize Terraform:
   ```bash
   terraform init
   ```

4. Query the origins:
   ```bash
   terraform apply
   ```

## Outputs

After successful query, you'll see:
- `origins_list` - List of all origins for the domain including:
  - `id` - Origin ID
  - `protocol` - Protocol type
  - `listen_port` - Listening port
  - `origin_protocol` - Origin protocol
  - `load_balance` - Load balancing method
  - `origin_type` - Origin type
  - `records` - Origin records (IP addresses/domains)
- `origins_count` - Total number of origins for the domain
