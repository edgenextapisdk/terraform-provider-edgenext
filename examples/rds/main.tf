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

variable "rds_instance_id" {
  description = "RDS instance ID (data sources and edgenext_rds_database)"
  type        = string
}

variable "rds_database_name" {
  description = "Database name for edgenext_rds_database example resource"
  type        = string
}

data "edgenext_rds_instances" "mysql" {
  page_num         = 1
  page_size        = 1000
  datastore_type   = "mysql"
}

data "edgenext_rds_databases" "example" {
  instance_id = var.rds_instance_id
}

resource "edgenext_rds_database" "example" {
  instance_id   = var.rds_instance_id
  name          = var.rds_database_name
  character_set = "utf8"
  collate       = "utf8_general_ci"
}
