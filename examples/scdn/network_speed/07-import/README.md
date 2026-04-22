# Import Network Speed Config

This example demonstrates how to import an existing SCDN network speed configuration into Terraform with **automatic ID parsing**.

## Usage

1. Copy `terraform.tfvars.example` to `terraform.tfvars` and fill in your credentials:

```bash
cp terraform.tfvars.example terraform.tfvars
```

2. Initialize Terraform:

```bash
terraform init
```

3. Import the existing configuration directly:

You only need to provide the resource ID in the format `<business_id>-<business_type>`. The provider will automatically parse these fields.

```bash
# Example command:
terraform import edgenext_scdn_network_speed_config.example 2676-tpl
```

4. Verify the result:

```bash
terraform state show edgenext_scdn_network_speed_config.example
```

## How it works

The provider now includes a custom importer and validation logic:

1.  **Automated Import**: The importer splits the ID (e.g., `2676-tpl`) and automatically populates the `business_id` and `business_type`.
2.  **Robust Validation**: Although the fields are `Optional` in the schema (to allow the `import` command to skip configuration check), we've added a **`CustomizeDiff`** validator.
    *   **In Create Mode**: If you are creating a NEW resource, the provider will enforce that both `business_id` and `business_type` are provided in your `main.tf`.
    *   **In Import Mode**: The validator detects the existing ID and skips the mandatory check, allowing for a seamless import experience without pre-defining the variables.


