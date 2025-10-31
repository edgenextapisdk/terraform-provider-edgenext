# Create SCDN Certificate

This example demonstrates how to create a SCDN certificate using Terraform.

## Prerequisites

- Terraform >= 1.0
- EdgeNext account with API credentials
- Valid SSL certificate (PEM format) for upload

## Files

- `main.tf` - Main Terraform configuration
- `terraform.tfvars.example` - Example variables file

## Usage

1. Copy `terraform.tfvars.example` to `terraform.tfvars`
2. Fill in your EdgeNext API credentials and certificate details
3. Run `terraform init` to initialize the provider
4. Run `terraform plan` to preview changes
5. Run `terraform apply` to create the certificate

## Variables

- `access_key` - EdgeNext Access Key (required)
- `secret_key` - EdgeNext Secret Key (required)
- `endpoint` - EdgeNext SCDN API Endpoint (optional, defaults to https://api.edgenextscdn.com)
- `ca_name` - Certificate name (required)
- `ca_cert` - Certificate public key in PEM format (required)
- `ca_key` - Certificate private key in PEM format (required)

## Outputs

- `certificate_id` - The ID of the created certificate
- `certificate_name` - The name of the certificate
- `certificate_sn` - The serial number of the certificate
- `issuer` - The certificate issuer
- `issuer_expiry_time` - The certificate expiry time
- `binded` - Whether the certificate is bound to any domain

