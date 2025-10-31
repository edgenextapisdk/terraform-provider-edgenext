# Query All Origin Groups Example

This example demonstrates how to query all origin groups for domain configuration.

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

- `protect_status`: Protection status (required)
  - `scdn`: Shared nodes
  - `exclusive`: Dedicated nodes

