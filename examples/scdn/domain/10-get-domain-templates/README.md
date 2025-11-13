# Example 10: Query Domain Templates

This example demonstrates how to query templates that are bound to an SCDN domain.

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

4. Query the domain templates:
   ```bash
   terraform apply
   ```

## Outputs

After successful query, you'll see:
- `domain_templates` - Templates bound to the domain including:
  - `domain_id` - Domain ID
  - `binded_templates` - List of bound templates, each containing:
    - `business_id` - Business ID
    - `business_type` - Business type
    - `app_type` - Application type
    - `name` - Template name
- `templates_count` - Total number of templates bound to the domain

## Use Cases

This data source is useful for:
- Checking which templates are applied to a domain
- Understanding domain configuration inheritance
- Auditing template usage across domains

