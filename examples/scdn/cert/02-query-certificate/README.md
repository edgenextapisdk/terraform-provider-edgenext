# Query SCDN Certificate

This example demonstrates how to query a SCDN certificate by ID using Terraform.

## Prerequisites

- Terraform >= 1.0
- EdgeNext account with API credentials
- A valid certificate ID

## Files

- `main.tf` - Main Terraform configuration
- `terraform.tfvars.example` - Example variables file

## Usage

1. Copy `terraform.tfvars.example` to `terraform.tfvars`
2. Fill in your EdgeNext API credentials and certificate ID
3. Run `terraform init` to initialize the provider
4. Run `terraform plan` to preview changes
5. Run `terraform apply` to query the certificate

## Variables

- `access_key` - EdgeNext Access Key (required)
- `secret_key` - EdgeNext Secret Key (required)
- `endpoint` - EdgeNext SCDN API Endpoint (optional, defaults to https://api.edgenextscdn.com)
- `certificate_id` - Certificate ID to query (required)

## Outputs

- `certificate_id` - The ID of the certificate
- `certificate_name` - The name of the certificate
- `issuer` - The certificate issuer
- `issuer_start_time` - The certificate start time
- `issuer_expiry_time` - The certificate expiry time
- `issuer_expiry_time_desc` - The certificate expiry time description
- `binded` - Whether the certificate is bound to any domain
- `ca_domain` - The domains in the certificate
- `apply_status` - The application status
- `ca_type` - The certificate type

