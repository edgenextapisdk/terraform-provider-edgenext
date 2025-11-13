# Unbind Domain from Rule Template

This example demonstrates how to unbind domains from a rule template using Terraform.

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

- This example performs the actual unbind operation using the `edgenext_scdn_rule_template_domain_unbind` resource
- The domains specified in `domain_ids` will be unbound from the template
- The operation is idempotent - running `terraform apply` multiple times will have the same effect
- To re-bind domains, you would need to remove this resource and use the bind operation

## Unbinding Domains

To unbind domains, simply apply the configuration:

```bash
terraform apply
```

The `edgenext_scdn_rule_template_domain_unbind` resource will:
1. Unbind the specified domain IDs from the template
2. Create a resource in Terraform state to track the unbind operation

## Removing the Unbind Resource

To remove the unbind resource from state (without re-binding):

```bash
terraform destroy
```

Note: This only removes the resource from Terraform state. The domains remain unbound.

## Outputs

- `template_info`: Template information (ID, name, description)
- `domains_before_unbind`: List of domains that were bound before the unbind operation
- `unbind_operation_id`: ID of the unbind operation resource
- `domains_unbound`: List of domain IDs that were unbound

