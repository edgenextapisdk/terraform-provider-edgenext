# EdgeNext SSL Certificate Management Example
# Terraform configuration for EdgeNext SSL Certificate API

# =============================================================================
# Provider Configuration
# =============================================================================

terraform {
  required_providers {
    edgenext = {
      source  = "edgenextapisdk/edgenext"
      version = "~> 1.0"
    }
  }
}

provider "edgenext" {
  access_key  = "your-edgenext-api-key-here"
  secret_key  = "your-edgenext-api-secret-here"
  endpoint    = "your-edgenext-api-endpoint-here"
}

# =============================================================================
# SSL Certificate Resource Examples
# =============================================================================

# 1. Import certificate
# resource "edgenext_ssl_certificate" "import_cert" {
#   # Empty resource block for import
# }

# 2. Main website SSL certificate
resource "edgenext_ssl_certificate" "web_cert" {
  name        = "example_cert"
  certificate = file("${path.module}/files/www.barqplay.com.crt")
  key         = file("${path.module}/files/www.barqplay.com.key")
}



# =============================================================================
# SSL Certificate Data Source Examples
# =============================================================================

# 1. Query single SSL certificate (by certificate ID)
# data "edgenext_ssl_certificate" "web_cert_info" {
#   cert_id = edgenext_ssl_certificate.web_cert.cert_id

#   depends_on = [ edgenext_ssl_certificate.web_cert ]
# }

# # 2. Query all SSL certificates (page 1, 10 items per page)
# data "edgenext_ssl_certificates" "all_certs_page1" {
#   page_number = 1
#   page_size   = 10

#   depends_on = [ edgenext_ssl_certificate.web_cert ]
# }
