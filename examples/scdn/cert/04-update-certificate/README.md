# Update SCDN Certificate

This example demonstrates how to update a SCDN certificate using Terraform.

## Prerequisites

- Terraform >= 1.0
- EdgeNext account with API credentials
- An existing certificate ID

## Files

- `main.tf` - Main Terraform configuration
- `terraform.tfvars.example` - Example variables file

## Usage

1. Copy `terraform.tfvars.example` to `terraform.tfvars`
2. Fill in your EdgeNext API credentials, certificate ID, and new certificate name
3. Run `terraform init` to initialize the provider
4. Run `terraform plan` to preview changes
5. Run `terraform apply` to update the certificate

**Note**: You can directly run `terraform apply` without importing. The `certificate_id` field in the resource will identify the existing certificate to update.

## Variables

- `access_key` - EdgeNext Access Key (required)
- `secret_key` - EdgeNext Secret Key (required)
- `endpoint` - EdgeNext SCDN API Endpoint (optional, defaults to https://api.edgenextscdn.com)
- `certificate_id` - Certificate ID to update (required)
- `new_ca_name` - New certificate name (required)

## Outputs

- `certificate_id` - The ID of the certificate
- `old_certificate_name` - The old certificate name
- `new_certificate_name` - The new certificate name

## Notes

- This example shows how to update the certificate name only. The `ca_cert` and `ca_key` fields are optional for updates.
- To update certificate content (ca_cert and ca_key), you need to:
  1. Import the existing certificate resource first: `terraform import edgenext_scdn_certificate.example <certificate_id>`
  2. Provide the new certificate values in the resource definition
- When updating certificate content, both `ca_cert` and `ca_key` must be provided together and must match.

