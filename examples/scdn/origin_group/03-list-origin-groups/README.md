# List Origin Groups Example

This example demonstrates how to list origin groups.

## Usage

1. Copy `terraform.tfvars.example` to `terraform.tfvars` and fill in your credentials

2. Edit `terraform.tfvars` with your actual values

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

- `page`: Page number (default: 1)
- `page_size`: Page size (default: 20)
- `name`: Origin group name filter (optional)

