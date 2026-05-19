# EdgeNext EIP association examples.
# Associates an existing EIP (floating IP) with an ECS instance, ELB, or network interface.
# Region is taken from the provider block only.

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
# Data sources (optional)
# =============================================================================

# List floating IPs to find allocation_id (EIP id) and floating_ip_address.
data "edgenext_eip_floating_ips" "all" {
  limit = 50
}

# =============================================================================
# EIP association — ECS instance (default instance_type)
# =============================================================================

# Bind EIP to an ECS server. When the instance has multiple ENIs or fixed IPs,
# set fixed_ip_address to the target private IP.
# resource "edgenext_eip_association" "ecs" {
#   allocation_id = "f5390261-241e-43e0-ab1b-5bffaf334c81"
#   instance_id   = "518031d1-f66d-416d-a66a-96ab91f4def9"
#   # instance_type = "EcsInstance"  # default
#   fixed_ip_address = "172.31.0.32"
# }
# Computed: floating_ip_address, fixed_ip_address

# =============================================================================
# EIP association — ELB load balancer
# =============================================================================

# instance_id is the load balancer id.
# resource "edgenext_eip_association" "elb" {
#   allocation_id = "f5390261-241e-43e0-ab1b-5bffaf334c81"
#   instance_id   = "9d2da79f-431e-4d7a-a5af-92082c41de01"
#   instance_type = "ElbInstance"
#   # fixed_ip_address = "172.31.0.170"
# }

# =============================================================================
# EIP association — network interface (port)
# =============================================================================

# instance_id is the network interface (port) id, not the ECS server id.
# resource "edgenext_eip_association" "eni" {
#   allocation_id = "f5390261-241e-43e0-ab1b-5bffaf334c81"
#   instance_id   = "6292aaaa-388f-4a59-a379-83bd4f38246b"
#   instance_type = "NetworkInterface"
#   fixed_ip_address = "172.31.0.32"
# }

# =============================================================================
# Import
# =============================================================================

# terraform import edgenext_eip_association.imported <allocation_id>
# resource "edgenext_eip_association" "imported" {}
