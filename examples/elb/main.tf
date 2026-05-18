# EdgeNext ELB Terraform examples. Region is taken from the provider block only.

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

# =============================================================================
# Data sources
# =============================================================================

# data "edgenext_elb_load_balancers" "example" {
#   name  = ""
#   limit = 10
# }
# Computed: balancers, total, hasmore (balancers include target_groups, listeners, spec, …).

# data "edgenext_elb_certificates" "example" {
#   page_num   = 1
#   page_size  = 10
#   name       = ""
# }
# Computed: certificates, total.

# data "edgenext_elb_listeners" "example" {
#   loadbalancer_id = "65292e79-0e9b-44d4-8d38-8e8fc3b93b36"
#   limit           = 1000
# }
# Computed: listeners, total, hasmore (listeners include default_target_group_id, insert_headers).

# data "edgenext_elb_target_groups" "example" {
#   loadbalancer_id = "65292e79-0e9b-44d4-8d38-8e8fc3b93b36"
# }
# Computed: target_groups (target id list, health_monitor when present).

# data "edgenext_elb_target_group_attachments" "example" {
#   target_group_id = "d6d1b6a1-74fa-4f51-b490-4fd488700256"
# }
# Computed: targets.

# data "edgenext_elb_l7_policies" "example" {
#   listener_id = "f3318153-bc30-4b6f-bbd0-03aaa36fec50"
#   limit       = 100
#   sort_key    = "position"
#   sort_dir    = "asc"
# }
# Computed: l7_policies, total, hasmore.

# data "edgenext_elb_l7_rules" "example" {
#   # policy_id = edgenext_elb_l7_policy.redirect.id
#   policy_id                    = "c6f3f52d-e331-487f-bb9f-c33c3f7f51ba"
#   limit                        = 100
#   sort_key                     = "created_at"
#   sort_dir                     = "desc"
#   filter_id                    = "2920e86d-ae63-40cd-860a-22c625f0dea9"
#   filter_type                  = "PATH"
#   filter_provisioning_status   = "ACTIVE"
# }
# Computed: rules, total, hasmore, is_complete_task.

# =============================================================================
# Resources
# =============================================================================

# resource "edgenext_elb_listener" "example" {
#   loadbalancer_id             = "9d2da79f-431e-4d7a-a5af-92082c41de01"
#   name                        = "tyd-listener"
#   description                 = ""
#   protocol                    = "TERMINATED_HTTPS"
#   protocol_port               = 8080
#   connection_limit            = -1
#   insert_headers              = { "X-Forwarded-For" = "true" }
#   default_tls_certificate_ref = "https://barbican.openstack.svc.frankfurt-b.com/v1/containers/a037e728-a23f-4937-ae7d-a70e1fedf828"
# }

# resource "edgenext_elb_target_group" "example" {
#   listener_id  = "0c03c22d-b9b5-49c8-a5e1-aedd9bc6e6b9"
#   name         = "tyd-target-group"
#   description  = ""
#   protocol     = "HTTP"
#   lb_algorithm = "ROUND_ROBIN"
#   health_monitor {
#     name           = "tyd-check"
#     type           = "HTTP"
#     max_retries    = 3
#     delay          = 10
#     timeout        = 5
#     http_method    = "GET"
#     url_path       = "/check"
#     expected_codes = "200"
#   }
# }

# resource "edgenext_elb_target_group_attachment" "b1" {
#   target_group_id = edgenext_elb_target_group.example.id
#   address       = "192.168.0.244"
#   protocol_port = 80
#   name          = "target-1"
#   weight        = 10
# }

# resource "edgenext_elb_target_group_attachment" "extra_backend" {
#   target_group_id = "78525404-e6eb-43e3-8b41-c279c5aa8ae1"
#   address       = "192.168.0.28"
#   protocol_port = 8080
#   name          = "app-server-2"
#   weight        = 5
# }

# resource "edgenext_elb_l7_rule" "extra_path" {
#   l7policy_id    = edgenext_elb_l7_policy.redirect.id
#   type           = "PATH"
#   compare_type   = "STARTS_WITH"
#   value          = "/api"
#   invert         = false
#   admin_state_up = true
# }

# resource "edgenext_elb_l7_policy" "redirect" {
#   listener_id              = "0c03c22d-b9b5-49c8-a5e1-aedd9bc6e6b9"
#   name                     = "redirect-www"
#   description              = "ggg"
#   action                   = "REDIRECT_TO_TARGET_GROUP" # or REDIRECT_TO_URL
#   # redirect_url             = "https://www.barqplay.com/"
#   redirect_target_group_id = "b49ba19e-d737-4c6e-8067-62fe8f27448a"
# }

# resource "edgenext_elb_l7_rule" "redirect_path" {
#   l7policy_id  = edgenext_elb_l7_policy.redirect.id
#   type         = "PATH"
#   compare_type = "EQUAL_TO"
#   value        = "/old"
# }

# resource "edgenext_elb_certificate" "example" {
#   name = "my-cert"
#   certificate = file("~/Downloads/www.barqplay.com/www.barqplay.com.crt")
#   private_key = file("~/Downloads/www.barqplay.com/www.barqplay.com.key")
#   # certificate = <<-EOT
#   # -----BEGIN CERTIFICATE-----
#   # ... PEM certificate / chain ...
#   # -----END CERTIFICATE-----
#   # EOT
#   # private_key = <<-EOT
#   # -----BEGIN PRIVATE KEY-----
#   # ... PEM private key ...
#   # -----END PRIVATE KEY-----
#   # EOT
# }

# =============================================================================
# Import
# =============================================================================

# terraform import edgenext_elb_listener.imported <listener_id>
# resource "edgenext_elb_listener" "imported" {}

# terraform import edgenext_elb_target_group_attachment.imported <target_group_id>/<target_id>
# resource "edgenext_elb_target_group_attachment" "imported" {}

# terraform import edgenext_elb_l7_rule.imported <l7policy_id>/<l7rule_id>
# resource "edgenext_elb_l7_rule" "imported" {}

# terraform import edgenext_elb_l7_policy.imported <l7policy_id>
# resource "edgenext_elb_l7_policy" "imported" {}

# terraform import edgenext_elb_certificate.imported <certificate_id>
# resource "edgenext_elb_certificate" "imported" {}
