# Bind Domain to Rule Template

This example demonstrates how to bind domains to a rule template using Terraform.

## Usage

1. Copy `terraform.tfvars.example` to `terraform.tfvars` and fill in your credentials:

```bash
cp terraform.tfvars.example terraform.tfvars
```

2. Edit `terraform.tfvars` with your actual values:

```hcl
access_key  = "your_access_key"
secret_key  = "your_secret_key"
template_id = 123
domain_ids  = [101753]
```

3. Initialize and apply:

```bash
terraform init
terraform plan
terraform apply
```

## Notes

- This example performs the actual bind operation using the `edgenext_scdn_rule_template_domain_bind` resource
- The domains specified in `domain_ids` will be bound to the template
- The operation is idempotent - running `terraform apply` multiple times will have the same effect
- To unbind domains, you can use `terraform destroy` or use the unbind resource (see example 07)

## Binding Domains

To bind domains, simply apply the configuration:

```bash
terraform apply
```

The `edgenext_scdn_rule_template_domain_bind` resource will:
1. Bind the specified domain IDs to the template
2. Create a resource in Terraform state to track the bind operation

## Unbinding Domains

To unbind domains, you can:

**Method 1: Destroy the bind resource**
```bash
terraform destroy
```

This will unbind the domains from the template.

**Method 2: Use the unbind resource** (see example 07)

## Outputs

- `template_info`: Template information (ID, name, description)
- `domains_before_bind`: List of domains that were bound before the bind operation
- `bind_operation_id`: ID of the bind operation resource
- `domains_bound`: List of domain IDs that were bound

