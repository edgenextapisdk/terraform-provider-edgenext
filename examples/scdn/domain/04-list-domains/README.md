# Example 4: List SCDN Domains

This example demonstrates how to query a list of SCDN domains with pagination support.

## Prerequisites

1. Valid EdgeNext API credentials (access_key and secret_key)

## Files

- `main.tf` - Main Terraform configuration
- `terraform.tfvars.example` - Example configuration file

## Usage

1. Copy `terraform.tfvars.example` to `terraform.tfvars`:
   ```bash
   cp terraform.tfvars.example terraform.tfvars
   ```

2. Update `terraform.tfvars` with your credentials

3. Optionally adjust pagination parameters:
   - `page` - Page number to retrieve (default: 1)
   - `page_size` - Number of items per page (default: 10)
   - `group_id` - Filter by group ID (optional)

4. Initialize Terraform:
   ```bash
   terraform init
   ```

5. Query the domain list:
   ```bash
   terraform apply
   ```

## Outputs

After successful query, you'll see:
- `domains_list` - List of domains with their basic information
- `total_count` - Number of domains in the current page

## Pagination

You can query different pages by changing the `page` variable:
- Page 1: `page = 1`
- Page 2: `page = 2`
- etc.
