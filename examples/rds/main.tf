# EdgeNext RDS Terraform examples. Region is taken from the provider block only.

terraform {
  required_providers {
    edgenext = {
      source  = "edgenextapisdk/edgenext"
      version = ">= 1.0.0"
    }
  }
}

provider "edgenext" {
  access_key = var.access_key
  secret_key = var.secret_key
  endpoint   = var.endpoint
  region     = var.region
}

variable "access_key" {
  description = "EdgeNext Access Key"
  type        = string
}

variable "secret_key" {
  description = "EdgeNext Secret Key"
  type        = string
}

variable "endpoint" {
  description = "EdgeNext API endpoint"
  type        = string
}

variable "region" {
  description = "EdgeNext region"
  type        = string
}

# data "edgenext_rds_instances" "mysql" {
#   page_num    = 1
#   page_size   = 1000
#   instance_id = ""
#   name        = ""
# }

# Replace instance_id with a value from your account (for example from data.edgenext_rds_instances.mysql).
# data "edgenext_rds_databases" "example" {
#   instance_id = "da3a5c51-bf3b-4bc3-9c6d-a9095092f942"
# }

# data "edgenext_rds_accounts" "example" {
#   instance_id = "da3a5c51-bf3b-4bc3-9c6d-a9095092f942"
# }

# # Optional filters backup_id and name: use empty strings to list without those filters.
# data "edgenext_rds_backups" "example" {
#   page_num  = 1
#   page_size = 1000
#   backup_id = ""
#   name      = "8.0"
# }

# # Optional filters policy_id and name: use empty strings to list without those filters.
# data "edgenext_rds_backup_policies" "example" {
#   page_num  = 1
#   page_size = 10
#   policy_id = ""
#   name      = ""
# }

# data "edgenext_rds_backup_policy_associate_instances" "example" {
#   policy_id = "backup-policy-d0ce79aeedfd76697daf6569"
# }

# resource "edgenext_rds_database" "example" {
#   instance_id   = "da3a5c51-bf3b-4bc3-9c6d-a9095092f942"
#   name          = "tyd_db"
#   character_set = "utf8"
#   collate       = "utf8_general_ci"
# }

# resource "edgenext_rds_account" "example" {
#   instance_id = "b4a406bc-2859-430d-8f95-9bc1e054d347"
#   user_name   = "user"
#   host        = "%"
#   password    = "Pw@123456"
# }

# resource "edgenext_rds_account_root_password" "example" {
#   instance_id = "b4a406bc-2859-430d-8f95-9bc1e054d347"
#   password    = "Root@123456"
# }

# resource "edgenext_rds_account_privilege" "example" {
#   instance_id = "b4a406bc-2859-430d-8f95-9bc1e054d347"
#   user_name   = "user"
#   host        = "%"
#   databases   = ["db"]
# }

# # Replace instance_id with your RDS instance ID. Creates a full backup (API type is always full).
# resource "edgenext_rds_backup" "example" {
#   instance_id = "1429bdc3-96fa-4b81-b13c-b306aa8b41ff"
#   name        = "tyd-backup"
# }

# resource "edgenext_rds_backup_policy" "example" {
#   name             = "tyd-backup-policy"
#   cycle_type       = "weekly"
#   weekdays         = [1, 7]
#   schedule_times   = ["00:00:00", "01:00:00"]
#   retention_type   = "days"
#   retention_value  = 30
#   stop_type        = "end_time"
#   stop_value       = "2026-05-31 14:32:20"
#   schedule_enabled = true
# }

# resource "edgenext_rds_backup_policy_associate_instance" "example" {
#   # policy_id = edgenext_rds_backup_policy.example.id
#   policy_id   = "backup-policy-d0ce79aeedfd76697daf6569"
#   instance_id = "2505e6e9-bd5a-4847-b4a8-43ddaf948d34"
# }

# resource "edgenext_rds_backup_policy_associate_instance" "example2" {
#   # policy_id = edgenext_rds_backup_policy.example.id
#   policy_id   = "backup-policy-451fd8b1ecd3e54422bfdff3"
#   instance_id = "2505e6e9-bd5a-4847-b4a8-43ddaf948d34"
# }

# Import examples (all in one place)
# Registered resources (define stub first, then import):
#
# resource "edgenext_rds_database" "imported_database" {}
# terraform import edgenext_rds_database.imported_database '<instance_id>/<database_name>'
#
# resource "edgenext_rds_account" "imported_account" {}
# terraform import edgenext_rds_account.imported_account '<instance_id>/<user_name>/<host>'
#
# resource "edgenext_rds_account_privilege" "imported_account_privilege" {}
# terraform import edgenext_rds_account_privilege.imported_account_privilege '<instance_id>/<user_name>/<host>'
#
# resource "edgenext_rds_backup" "imported_backup" {}
# terraform import edgenext_rds_backup.imported_backup '<backup_id>'
#
# resource "edgenext_rds_backup_policy" "imported_backup_policy" {}
# terraform import edgenext_rds_backup_policy.imported_backup_policy '<policy_id>'
#
# resource "edgenext_rds_backup_policy_associate_instance" "imported_backup_policy_association" {}
# terraform import edgenext_rds_backup_policy_associate_instance.imported_backup_policy_association '<policy_id>/<instance_id>'
