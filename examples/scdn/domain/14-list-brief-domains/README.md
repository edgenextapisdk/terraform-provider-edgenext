# Example 14: List Brief Domains

This example demonstrates how to query a brief list of SCDN domains with minimal information (ID, domain name, member ID).

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

3. Optionally specify domain IDs to query specific domains:
   ```hcl
   domain_ids = [12345, 67890]
   ```
   Or leave empty to query all domains

4. Initialize Terraform:
   ```bash
   terraform init
   ```

5. Query the brief domains:
   ```bash
   terraform apply
   ```

## Outputs

After successful query, you'll see:
- `brief_domains_list` - List of brief domain information including:
  - `id` - Domain ID
  - `domain` - Domain name
  - `member_id` - Member ID
- `total_count` - Total number of domains

## Use Cases

This data source is useful when you only need basic domain information:
- Quick domain lookups by ID
- Getting domain names for a list of IDs
- Lightweight queries when full domain details are not needed
- Faster than full domain queries when you only need ID and name

## Comparison with Full Domain Query

- **Brief Domains**: Faster, returns only ID, domain name, and member ID
- **Full Domains**: Slower, returns complete domain information including settings, origins, etc.

Use brief domains when you only need basic information for faster queries.

