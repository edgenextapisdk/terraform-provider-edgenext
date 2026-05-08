# EdgeNext RDS Services

This package provides Terraform resources and data sources for managing EdgeNext RDS instances, users, database privileges, backups, and backup policies.

## Resources

### RDS Database
- **Resource**: `edgenext_rds_database` (`ResourceENRDSDatabase`)
- **File**: `resource_en_rds_database.go`
- **Description**: Manage one database on an RDS instance

### RDS Account
- **Resource**: `edgenext_rds_account` (`ResourceENRDSAccount`)
- **File**: `resource_en_rds_account.go`
- **Description**: Manage one database user on an RDS instance

### RDS Account Privilege
- **Resource**: `edgenext_rds_account_privilege` (`ResourceENRDSAccountPrivilege`)
- **File**: `resource_en_rds_account_privilege.go`
- **Description**: Manage granted database list for a user identified by `instance_id + user_name + host`

### RDS Account Root Password
- **Resource**: `edgenext_rds_account_root_password` (`ResourceENRDSAccountRootPassword`)
- **File**: `resource_en_rds_account_root_password.go`
- **Description**: Enable root access and set root password for an instance

### RDS Backup
- **Resource**: `edgenext_rds_backup` (`ResourceENRDSBackup`)
- **File**: `resource_en_rds_backup.go`
- **Description**: Create and delete manual backups

### RDS Backup Policy
- **Resource**: `edgenext_rds_backup_policy` (`ResourceENRDSBackupPolicy`)
- **File**: `resource_en_rds_backup_policy.go`
- **Description**: Manage backup policy lifecycle and schedule status

### RDS Backup Policy Associate Instance
- **Resource**: `edgenext_rds_backup_policy_associate_instance` (`ResourceENRDSBackupPolicyAssociateInstance`)
- **File**: `resource_en_rds_backup_policy_associate_instance.go`
- **Description**: Manage one backup-policy-to-instance association (`policy_id + instance_id`)

## Data Sources

### RDS Instances
- **Data Source**: `edgenext_rds_instances` (`DataSourceENRDSInstances`)
- **File**: `data_source_en_rds_instances.go`
- **Description**: Query RDS instance list (request always uses MySQL datastore type)

### RDS Databases
- **Data Source**: `edgenext_rds_databases` (`DataSourceENRDSDatabases`)
- **File**: `data_source_en_rds_databases.go`
- **Description**: Query database list for one instance

### RDS Accounts
- **Data Source**: `edgenext_rds_accounts` (`DataSourceENRDSAccounts`)
- **File**: `data_source_en_rds_accounts.go`
- **Description**: Query database user list for one instance

### RDS Backups
- **Data Source**: `edgenext_rds_backups` (`DataSourceENRDSBackups`)
- **File**: `data_source_en_rds_backups.go`
- **Description**: Query backup list with optional backup ID and name filters

### RDS Backup Policies
- **Data Source**: `edgenext_rds_backup_policies` (`DataSourceENRDSBackupPolicies`)
- **File**: `data_source_en_rds_backup_policies.go`
- **Description**: Query backup policy list with optional policy ID and name filters

### RDS Backup Policy Associate Instances
- **Data Source**: `edgenext_rds_backup_policy_associate_instances` (`DataSourceENRDSBackupPolicyAssociateInstances`)
- **File**: `data_source_en_rds_backup_policy_associate_instances.go`
- **Description**: Query associated instances of one backup policy

## File Structure

```text
edgenext/services/rds/
├── README.md                                                # This documentation
├── resource_en_rds_*.go                                     # RDS resource implementations
├── data_source_en_rds_*.go                                  # RDS data source implementations
├── data_source_en_rds_instances.md                          # Instances data source note
├── rds_mysql_charset_collations.json                        # Charset/collation catalog for database resource validation
└── rds_mysql_charsets.go                                    # Validation helpers for database charset/collation
```

## Usage Examples

### Query RDS Instances

```hcl
data "edgenext_rds_instances" "mysql" {
  page_num    = 1
  page_size   = 1000
  instance_id = ""
  name        = ""
}
```

### Create Database User and Grant Privileges

```hcl
resource "edgenext_rds_account" "example" {
  instance_id = "b4a406bc-2859-430d-8f95-9bc1e054d347"
  user_name   = "app_user"
  host        = "%"
  password    = "Strong@1234"
}

resource "edgenext_rds_account_privilege" "example" {
  instance_id = edgenext_rds_account.example.instance_id
  user_name   = edgenext_rds_account.example.user_name
  host        = edgenext_rds_account.example.host
  databases   = ["app_db"]
}
```

### Associate Instance to Backup Policy

```hcl
resource "edgenext_rds_backup_policy_associate_instance" "example" {
  policy_id   = "backup-policy-d0ce79aeedfd76697daf6569"
  instance_id = "b4a406bc-2859-430d-8f95-9bc1e054d347"
}
```
