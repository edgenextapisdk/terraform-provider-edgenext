# Create Origin Group Example

This example demonstrates how to create an origin group.

## Usage

1. Copy `terraform.tfvars.example` to `terraform.tfvars` and fill in your credentials:

```bash
cp terraform.tfvars.example terraform.tfvars
```

2. Edit `terraform.tfvars` with your actual values

3. Initialize Terraform:

```bash
terraform init
```

4. Review the execution plan:

```bash
terraform plan
```

5. Apply the configuration:

```bash
terraform apply
```

6. To destroy the configuration:

```bash
terraform destroy
```

## Configuration Options

- `name`: Origin group name (2-16 characters)
- `remark`: Remark (2-64 characters, optional)
- `origins`: Origin list (at least 1)
  - `origin_type`: Origin type (0-IP, 1-domain)
  - `records`: Origin record list
    - `value`: Origin address
    - `port`: Origin port (1-65535)
    - `priority`: Weight (1-100)
    - `view`: Origin type (primary-backup, backup-backup)
    - `host`: Origin Host (optional)
  - `protocol_ports`: Protocol port mapping
    - `protocol`: Protocol (0-http, 1-https)
    - `listen_ports`: Listen port list
  - `origin_protocol`: Origin protocol (0-http, 1-https, 2-follow)
  - `load_balance`: Load balance strategy (0-ip_hash, 1-round_robin, 2-cookie)

