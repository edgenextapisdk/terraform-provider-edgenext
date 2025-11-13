# Example 3: Query SCDN Domain

This example demonstrates how to query a single SCDN domain by its domain name.

## Prerequisites

1. Valid EdgeNext API credentials (access_key and secret_key)
2. An existing SCDN domain name

## Files

- `main.tf` - Main Terraform configuration
- `terraform.tfvars.example` - Example configuration file

## Usage

1. Copy `terraform.tfvars.example` to `terraform.tfvars`:
   ```bash
   cp terraform.tfvars.example terraform.tfvars
   ```

2. Update `terraform.tfvars` with your credentials and domain name

3. Initialize Terraform:
   ```bash
   terraform init
   ```

4. Query the domain:
   ```bash
   terraform apply
   ```

## Outputs

After successful query, you'll see:
- `domain_info` - Complete domain information including:
  - `id` - Domain ID
  - `domain` - Domain name
  - `remark` - Domain remark
  - `access_progress` - Access configuration progress
  - `protect_status` - Protection status
  - `cname` - CNAME record information
