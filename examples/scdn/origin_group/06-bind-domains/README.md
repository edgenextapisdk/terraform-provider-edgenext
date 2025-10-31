# Bind Origin Group to Domains Example

This example demonstrates how to bind an origin group to domains.

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

- `origin_group_id`: Origin group ID to bind
- `domain_ids`: Domain ID array (optional)
- `domain_group_ids`: Domain group ID array (optional)
- `domains`: Domain array (optional)

At least one of `domain_ids`, `domain_group_ids`, or `domains` must be provided.

