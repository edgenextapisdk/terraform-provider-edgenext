# Switch Domain Template Example

This example demonstrates how to use the `edgenext_scdn_rule_template_switch` resource to switch domains to a different rule template.

## Usage

1. Initialize Terraform:
    ```bash
    terraform init
    ```

2. Create a `terraform.tfvars` file with your credentials and configuration:
    ```hcl
    access_key = "your_access_key"
    secret_key = "your_secret_key"
    api_host   = "https://api.edgenext.com"
    domain_ids = [12345, 67890]
    new_tpl_id = 54321
    ```

3. Plan the changes:
    ```bash
    terraform plan
    ```

4. Apply the changes:
    ```bash
    terraform apply
    ```

## Notes

- The `edgenext_scdn_rule_template_switch` resource performs an action (switching template) rather than maintaining a long-term state.
- Destroying the resource will attempt to unbind the domains from the target template if `new_tpl_id` is greater than 0.
