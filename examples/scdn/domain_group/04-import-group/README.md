# Import Domain Group Example

This example demonstrates how to import an existing SCDN Domain Group into Terraform management.

## Use Case

Use this when you have domain groups created outside of Terraform (e.g., via console or API) and want to manage them with Terraform.

## Features

- Import existing domain group
- Query group information
- List domains in the group
- Modify group after import

## Usage

1. Copy `terraform.tfvars.example` to `terraform.tfvars`:

```bash
cp terraform.tfvars.example terraform.tfvars
```

2. Edit `terraform.tfvars` with your credentials and the group ID to import.

3. Initialize Terraform:

```bash
terraform init
```

4. Import the existing group:

```bash
terraform import edgenext_scdn_domain_group.imported <group_id>
```

Example:
```bash
terraform import edgenext_scdn_domain_group.imported 123
```

5. After import, update the resource in `main.tf` with actual values from the import.

6. Verify and apply:

```bash
terraform plan
terraform apply
```

## Post-Import

After importing, you can:
- Modify the group name or remark
- Add or remove domain bindings
- Manage the group lifecycle with Terraform

## Outputs

- `group_info` - Complete information about the imported group
- `group_domains` - List of domains in the group
