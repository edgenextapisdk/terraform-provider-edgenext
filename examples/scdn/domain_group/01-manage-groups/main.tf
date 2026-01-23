terraform {
  required_providers {
    edgenext = {
      source  = "edgenextapisdk/edgenext"
      version = "0.2.1"
    }
  }
}

provider "edgenext" {
  access_key = var.access_key
  secret_key = var.secret_key
  endpoint   = var.endpoint
}

variable "access_key" {
  type        = string
  description = "EdgeNext access key"
  sensitive   = true
}

variable "secret_key" {
  type        = string
  description = "EdgeNext secret key"
  sensitive   = true
}

variable "endpoint" {
  type        = string
  description = "EdgeNext API endpoint"
  default     = "https://api.edgenextscdn.com"
}

variable "group_name" {
  type        = string
  description = "Name of the domain group"
  default     = "terraform-example-group"
}

variable "remark" {
  type        = string
  description = "Remark for the domain group"
  default     = "Created by Terraform"
}

variable "domains" {
  type        = list(string)
  description = "List of domains to bind to the group"
  default     = []
}

# Create a Domain Group
resource "edgenext_scdn_domain_group" "example" {
  group_name = var.group_name
  remark     = var.remark
  domains    = var.domains
}

# Query the created group using data source
data "edgenext_scdn_domain_groups" "example_query" {
  group_name = var.group_name
  depends_on = [edgenext_scdn_domain_group.example]
}

# Query domains in the group
data "edgenext_scdn_domain_group_domains" "example_domains" {
  group_id = tonumber(edgenext_scdn_domain_group.example.id)
}

output "group_id" {
  value = edgenext_scdn_domain_group.example.id
}

output "group_query_result" {
  value = data.edgenext_scdn_domain_groups.example_query.list
}

output "group_domains" {
  value = data.edgenext_scdn_domain_group_domains.example_domains.list
}
