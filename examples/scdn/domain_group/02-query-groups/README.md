# Query Domain Groups Example

This example demonstrates how to query SCDN Domain Groups using data sources.

## Features

- Query all domain groups
- Filter groups by name
- Filter groups by domain
- Pagination support

## Usage

1. Copy `terraform.tfvars.example` to `terraform.tfvars`:

```bash
cp terraform.tfvars.example terraform.tfvars
```

2. Edit `terraform.tfvars` with your credentials and optional filters.

3. Initialize and run:

```bash
terraform init
terraform plan
terraform apply
```

## Outputs

- `all_groups` - List of all domain groups
- `total_groups` - Total number of groups
- `filtered_by_name` - Groups matching the name filter
- `filtered_by_domain` - Groups containing the specified domain
