# EdgeNext ECS Terraform examples.
# Schemas match edgenext/services/ecs. Region is taken from the provider block only.
#
# The following resources are implemented under services/ecs but commented out in
# edgenext/provider.go until registered: edgenext_ecs_instance, edgenext_ecs_image,
# edgenext_ecs_floating_ip, edgenext_ecs_disk.

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

#
# VPC and networking (registered in provider)
#

# resource "edgenext_ecs_vpc" "example" {
#   name        = "example-vpc"
#   description = "Example VPC for ECS instances"
#   subnet {
#     name       = "example-subnet"
#     ip_version = 4
#     cidr       = "192.168.0.0/24"
#   }
# }

# Subnet under an existing VPC. Arguments cannot be changed after creation (replace on change).
# resource "edgenext_ecs_vpc_subnet" "example" {
#   vpc_id     = data.edgenext_ecs_vpcs.all.vpcs[0].id
#   name       = "example-subnet-extra"
#   ip_version = 4
#   cidr       = "192.168.1.0/24"
# }

# resource "edgenext_ecs_router" "example" {
#   name                = "example-router"
#   description         = "Example router"
#   external_network_id = data.edgenext_ecs_external_gateways.all.external_gateways[0].id
# }

# Attach subnet to router. router_id, vpc_id, subnet_id are ForceNew.
# resource "edgenext_ecs_router_port" "example" {
#   router_id = edgenext_ecs_router.example.id
#   # vpc_id    = edgenext_ecs_vpc.example.id
#   vpc_id    = "e3ab4cfe-1f43-4698-87cb-1eb3c5c78806"
#   subnet_id = "17dd56f7-94f7-4617-92b1-029ea78c3cc8" # from edgenext_ecs_vpc_subnet or data.edgenext_ecs_vpc_subnets
# }

# ENI (Neutron port). vpc_id and subnet_id are ForceNew; name/description update in place.
# resource "edgenext_ecs_network_interface" "example" {
#   name        = "example-eni"
#   description = "Example ENI"
#   vpc_id      = "06b20160-6d47-4d80-84f1-6127c65f6f19"
#   subnet_id   = "e9ada4b5-bccd-4b88-86c1-b515e6730ef6"
#   # port_security_enabled = true
#   # security_groups       = ["2af2b1e5-344f-4184-9173-cf1b5d43bf7d"]
# }

# Bind an instance to an ENI (separate from ENI create).
# resource "edgenext_ecs_network_interface_instance_binding" "example" {
#   network_interface_id = edgenext_ecs_network_interface.example.id
#   instance_id          = "cb24704a-dda5-4287-978f-cec28ee1e816"
# }

# Bind floating IP to an ENI (separate resource; edgenext_ecs_floating_ip create is optional in provider).
# resource "edgenext_ecs_network_interface_floating_ip_binding" "example" {
#   network_interface_id  = "47b88552-2a31-4446-9306-6abd692051bd"
#   floating_ip_address   = "156.246.18.218"
# }


# resource "edgenext_ecs_security_group" "example" {
#   name        = "example-sg-managed"
#   description = "A standard example security group"
# }

# Managed rule; argument changes force replacement.
# resource "edgenext_ecs_security_group_rule" "example" {
#   security_group_id = data.edgenext_ecs_security_groups.all.security_groups[0].id
#   protocol          = "tcp"
#   direction         = "ingress"
#   ethertype         = "IPv4"
#   port_range_min    = 100
#   port_range_max    = 300
#   remote_ip_prefix  = "192.168.0.0/24"
#   # remote_group_id  = data.edgenext_ecs_security_groups.all.security_groups[0].id
# }

# resource "edgenext_ecs_key_pair" "example" {
#   name       = "example-key"
#   # public_key = file("~/.ssh/id_rsa.pub")
# }

#
# Instance lifecycle helpers (no full instance resource unless enabled in provider)
#

# resource "edgenext_ecs_instance_power" "example" {
#   instance_id   = "cb24704a-dda5-4287-978f-cec28ee1e816"
#   desired_state = "ACTIVE" # or "SHUTOFF"
# }

# resource "edgenext_ecs_instance_reboot" "example" {
#   instance_id = "0d4dd8b5-e581-4097-92df-7085b4c94953"
#   reboot_type = "reboot_soft"
#   trigger     = "2026-04-24T00:00:00Z"
# }

#
# Tags (edgenext_ecs_tag uses tag_key / tag_value; binding uses edgenext_ecs_instance_tag)
#

locals {
  ecs_tags = {
    "test-key" = "test-value"
    "team"     = "platform"
    "env"      = "dev"
  }
}

# resource "edgenext_ecs_tag" "example" {
#   for_each  = local.ecs_tags
#   tag_key   = each.key
#   tag_value = each.value
# }

# Bind tag IDs to an instance (instance_id and instance_name are ForceNew as a pair in practice).
# resource "edgenext_ecs_instance_tag" "example_binding" {
#   instance_id   = "0d4dd8b5-e581-4097-92df-7085b4c94953"
#   instance_name = "monitor-vnc"
#   tag_ids       = [3,12,6]
# }

#
# Resources present in services/ecs but not registered in provider.go (uncomment there to use)
#

# resource "edgenext_ecs_disk" "example" {
#   name        = "example-disk"
#   size        = 50
#   volume_type = "SSD"
# }

# resource "edgenext_ecs_floating_ip" "example" {
#   bandwidth = 10
# }

# resource "edgenext_ecs_instance" "example" {
#   name            = "example-instance"
#   flavor_ref      = "s1.small"
#   image_ref       = "centos-7.9"
#   admin_pass      = "SecurePass123!"
#   bandwidth       = 5
#   key_name        = edgenext_ecs_key_pair.example.name
#   networks        = [edgenext_ecs_vpc.example.id]
#   security_groups = [edgenext_ecs_security_group.example.id]
# }

# resource "edgenext_ecs_image" "example" {
#   name         = "my-custom-image"
#   instance_id  = "" # optional source instance
#   description  = "From instance snapshot"
# }

#
# Data sources
#

# data "edgenext_ecs_floating_ips" "all" {
#   limit             = 10
#   floating_ip_id    = "c48bf957-9dad-4aea-b572-911bed2ab5d3"
#   floating_ip_address = "156.246.18.218"
# }

# data "edgenext_ecs_vpcs" "all" {
#   limit   = 10
#   # vpc_id  = ""
#   # name    = ""
# }

# data "edgenext_ecs_vpc_subnets" "by_vpc" {
#   vpc_id    = "e3ab4cfe-1f43-4698-87cb-1eb3c5c78806"
# }

# data "edgenext_ecs_routers" "all" {
#   # router_id    = "ae6e88bb-b818-4ae4-8a15-763461fc08a1"
#   # router_name  = "test"
#   limit = 10
# }

# data "edgenext_ecs_router_ports" "ports" {
#   router_id = "ae6e88bb-b818-4ae4-8a15-763461fc08a1"
# }

# data "edgenext_ecs_external_gateways" "all" {
#   limit = 10
# }

# data "edgenext_ecs_network_interfaces" "all_ports" {
#   network_interface_name  = "test9"
#   limit = 10
# }

# data "edgenext_ecs_instances" "all" {
#   # instance_name = ""
#   # instance_id   = ""
#   limit         = 10
# }

# data "edgenext_ecs_key_pairs" "all" {
#   limit = 10
# }

# data "edgenext_ecs_images" "all" {
#   visibility = "public"
#   name       = "OpenClaw 2026.03.08 64-bit"
#   page_num   = 1
#   page_size  = 10
# }

# data "edgenext_ecs_tags" "all" {
#   tag_key    = "zyx_test"
#   tag_value  = "zyx_test1"
#   page_num   = 1
#   page_size  = 10
# }

# data "edgenext_ecs_instance_tags" "by_tag" {
#   # tag_id     = 6
#   # tag_key    = "zyx_test"
#   # tag_value  = "zyx_test2"
#   page_num   = 1
#   page_size  = 10
# }

# data "edgenext_ecs_security_groups" "all" {
#   name  = ""
#   limit = 10
# }

# data "edgenext_ecs_security_group_rules" "for_sg" {
#   security_group_id = "2af2b1e5-344f-4184-9173-cf1b5d43bf7d"
# }

# GET /ecs/openapi/v2/volume/list — query params name, page_num, page_size
# data "edgenext_ecs_disks" "all" {
#   name      = ""
#   page_num  = 1
#   page_size = 10
# }


#
# Import examples (all in one place)
# Registered resources (define stub first, then import):
#
# resource "edgenext_ecs_vpc" "imported_vpc" {}
# terraform import edgenext_ecs_vpc.imported_vpc '<vpc_id>'
#
# resource "edgenext_ecs_vpc_subnet" "imported_vpc_subnet" {}
# terraform import edgenext_ecs_vpc_subnet.imported_vpc_subnet '<vpc_id>/<subnet_id>'
#
# resource "edgenext_ecs_router" "imported_router" {}
# terraform import edgenext_ecs_router.imported_router '<router_id>'
#
# resource "edgenext_ecs_router_port" "imported_router_port" {}
# terraform import edgenext_ecs_router_port.imported_router_port '<router_id>/<router_port_id>'
#
# resource "edgenext_ecs_network_interface" "imported_network_interface" {}
# terraform import edgenext_ecs_network_interface.imported_network_interface '<network_interface_id>'
#
# resource "edgenext_ecs_network_interface_instance_binding" "imported_eni_instance_binding" {}
# terraform import edgenext_ecs_network_interface_instance_binding.imported_eni_instance_binding '<network_interface_id>/<instance_id>'
#
# resource "edgenext_ecs_network_interface_floating_ip_binding" "imported_eni_floating_ip_binding" {}
# terraform import edgenext_ecs_network_interface_floating_ip_binding.imported_eni_floating_ip_binding '<network_interface_id>/<floating_ip_address>'
#
# resource "edgenext_ecs_security_group" "imported_security_group" {}
# terraform import edgenext_ecs_security_group.imported_security_group '<security_group_id>'
#
# resource "edgenext_ecs_security_group_rule" "imported_security_group_rule" {}
# terraform import edgenext_ecs_security_group_rule.imported_security_group_rule '<security_group_id>/<rule_id>'
#
# resource "edgenext_ecs_key_pair" "imported_key_pair" {}
# terraform import edgenext_ecs_key_pair.imported_key_pair '<key_pair_name>'
#
# resource "edgenext_ecs_tag" "imported_tag" {}
# terraform import edgenext_ecs_tag.imported_tag '<tag_id>/<tag_key>/<tag_value>'
#
# resource "edgenext_ecs_instance_tag" "imported_instance_tag" {}
# terraform import edgenext_ecs_instance_tag.imported_instance_tag '<instance_id>'