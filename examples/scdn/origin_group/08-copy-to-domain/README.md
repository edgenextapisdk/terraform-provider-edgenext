# Copy Origin Group to Domain Example

This example demonstrates how to copy an origin group to a domain.

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

- `origin_group_id`: Origin group ID to copy
- `domain_id`: Domain ID to copy to

