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
  description = "EdgeNext Access Key"
  type        = string
  sensitive   = true
}

variable "secret_key" {
  description = "EdgeNext Secret Key"
  type        = string
  sensitive   = true
}

variable "endpoint" {
  description = "EdgeNext SCDN API Endpoint"
  type        = string
  default     = "https://api.edgenextscdn.com"
}

resource "edgenext_scdn_user_ip" "file_upload_demo" {
  name      = "terraform-file-upload-demo"
  remark    = "IP list created via file upload from Terraform"
  file_path = "${path.module}/ip_list.txt"
}

output "user_ip_id" {
  value = edgenext_scdn_user_ip.file_upload_demo.id
}
