# Sort Network Speed Rules

This example demonstrates how to sort network speed rules for a specific configuration group.

## Usage

1. Copy `terraform.tfvars.example` to `terraform.tfvars` and fill in your credentials:

```bash
cp terraform.tfvars.example terraform.tfvars
```

2. Edit `terraform.tfvars` with your actual values:

```hcl
access_key   = "your_access_key"
secret_key   = "your_secret_key"
business_id  = 1246
business_type = "tpl"
config_group  = "customized_req_headers_rule"

# Rule IDs in the desired order
rule_ids = [609, 608, 607]
```

3. Initialize and apply:

```bash
terraform init
terraform plan
terraform apply
```

## Prerequisites

Before using this example, you need to:

1. **Have existing rules** in the specified `config_group`
   - You can query existing rules using the `edgenext_scdn_network_speed_rules` data source
   - Example:
     ```hcl
     data "edgenext_scdn_network_speed_rules" "list" {
       business_id   = 1246
       business_type = "tpl"
       config_group  = "customized_req_headers_rule"
     }
     ```

2. **Get the rule IDs** from the data source output or dashboard

3. **Specify the desired order** in the `rule_ids` variable

## Notes

- The `ids` array defines the order in which rules will be sorted
- Rules will be reordered to match the exact sequence specified in `ids`
- All rule IDs in the `ids` array must exist in the specified `config_group`
- The resource is idempotent - running `terraform apply` multiple times with the same `ids` order will not change anything

## Config Groups

Available `config_group` values:

- **`custom_page`**: Error page customization
- **`upstream_uri_change_rule`**: Upstream URI rewrite
- **`resp_headers_rule`**: Custom response headers
- **`customized_req_headers_rule`**: Custom upstream request headers

## Example Workflow

1. **Query existing rules**:
   ```hcl
   data "edgenext_scdn_network_speed_rules" "example" {
     business_id   = 1246
     business_type = "tpl"
     config_group  = "customized_req_headers_rule"
   }
   
   output "current_rule_ids" {
     value = [for rule in data.edgenext_scdn_network_speed_rules.example.list : rule.id]
   }
   ```

2. **Sort rules** (reverse order in this example):
   ```hcl
   resource "edgenext_scdn_network_speed_rules_sort" "example" {
     business_id   = 1246
     business_type = "tpl"
     config_group  = "customized_req_headers_rule"
     ids           = [609, 608, 607]  # Reverse order
   }
   ```

## Outputs

- `sorted_ids`: The rule IDs after sorting (should match the `ids` input)

