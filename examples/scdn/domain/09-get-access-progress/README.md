# Example 9: Query Access Progress Status

This example demonstrates how to query available access progress status options for SCDN domains.

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

3. Initialize Terraform:
   ```bash
   terraform init
   ```

4. Query the access progress options:
   ```bash
   terraform apply
   ```

## Outputs

After successful query, you'll see:
- `access_progress_options` - List of available access progress status options, each containing:
  - `key` - Progress status key/identifier
  - `name` - Progress status name/description
- `progress_count` - Total number of available progress status options

## Use Cases

This data source is useful for:
- Understanding available domain access progress statuses
- Filtering domains by access progress when using `data.edgenext_scdn_domains`
- Getting the valid values for the `access_progress` parameter

## Example Output

```
access_progress_options = [
  {
    key  = "pending"
    name = "Pending"
  },
  {
    key  = "processing"
    name = "Processing"
  },
  {
    key  = "completed"
    name = "Completed"
  }
]
```

