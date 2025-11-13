# List SCDN Certificates

This example demonstrates how to list SCDN certificates with various filters using Terraform.

## Prerequisites

- Terraform >= 1.0
- EdgeNext account with API credentials

## Files

- `main.tf` - Main Terraform configuration
- `terraform.tfvars.example` - Example variables file

## Usage

1. Copy `terraform.tfvars.example` to `terraform.tfvars`
2. Fill in your EdgeNext API credentials and optional filters
3. Run `terraform init` to initialize the provider
4. Run `terraform plan` to preview changes
5. Run `terraform apply` to list certificates

## Variables

- `access_key` - EdgeNext Access Key (required)
- `secret_key` - EdgeNext Secret Key (required)
- `endpoint` - EdgeNext SCDN API Endpoint (optional, defaults to https://api.edgenextscdn.com)
- `domain` - Filter by domain name (optional)
- `ca_name` - Filter by certificate name (optional)
- `binded` - Filter by binding status: true-bound, false-unbound (optional)
- `apply_status` - Filter by application status: 1-applying, 2-issued, 3-review failed, 4-uploaded (optional)
- `expiry_time` - Filter by expiry status: true-expired, false-not expired, inno-about to expire (optional)

## Outputs

- `total` - Total number of certificates
- `issuer_list` - List of available issuers
- `certificates` - List of certificates matching the filters
- `certificate_count` - Number of certificates returned

