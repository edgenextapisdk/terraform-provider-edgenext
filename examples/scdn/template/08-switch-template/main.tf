terraform {
  required_providers {
    edgenext = {
      source  = "edgenextapisdk/edgenext"
      version = "~> 1.0"
    }
  }
}

provider "edgenext" {
  access_key = var.access_key
  secret_key = var.secret_key
  endpoint   = var.endpoint
}

variable "access_key" {
  type = string
}

variable "secret_key" {
  type = string
}

variable "endpoint" {
  type = string
}

variable "domain_ids" {
  type    = list(number)
  default = [101753]
}

variable "new_tpl_id" {
  type    = number
  default = 1246
}

variable "new_tpl_type" {
  type    = string
  default = "only_domain"
}

resource "edgenext_scdn_rule_template_switch" "switch_test" {
  app_type     = "network_speed"
  new_tpl_id   = var.new_tpl_id
  new_tpl_type = var.new_tpl_type
  domain_ids   = var.domain_ids
}

output "switch_result_id" {
  value = edgenext_scdn_rule_template_switch.switch_test.id
}
