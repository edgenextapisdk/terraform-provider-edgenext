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


#
# VPC and Networking
#
# resource "edgenext_ecs_vpc" "example" {
#   region      = "tokyo-a"
#   name        = "example-vpc"
#   description = "Example VPC for ECS instances"
#   subnet {
#     name       = "example-subnet"
#     ip_version = 4
#     cidr       = "192.168.0.0/24"
#   }
# }

# Create an additional subnet under an existing VPC.
# Any argument change forces replacement because update is not supported.
# resource "edgenext_ecs_vpc_subnet" "example" {
#   region     = "tokyo-a"
#   network_id = data.edgenext_ecs_vpcs.all_vpcs.vpcs[0].id
#   name       = "example-subnet-extra"
#   ip_version = 4
#   cidr       = "192.168.1.0/24"
# }

# # Import format: region/network_id/subnet_id
# # terraform import edgenext_ecs_vpc_subnet.example 'tokyo-a/0e07db22-e210-4138-b40b-6cbdb3580b12/50a0f20a-16a8-46e1-8271-0d1660c739d5'
# resource "edgenext_ecs_vpc_subnet" "example" {}

# resource "edgenext_ecs_router" "example" {
#   region              = "tokyo-a"
#   name                = "example-router"
#   description         = "Example router"
#   external_network_id = ""
# }
#
# # Import format: region/router_id
# # terraform import edgenext_ecs_router.example 'tokyo-a/32253e60-4cdc-4bcb-9944-30fd8e2dcefb'
# resource "edgenext_ecs_router" "example" {}
#
# # Attach a subnet to router.
# resource "edgenext_ecs_router_port" "example" {
#   region     = "tokyo-a"
#   router_id  = "f9883769-85e0-4cdb-ae9a-c49da3a89edb"
#   # router_id  = edgenext_ecs_router.example.id
#   # network_id = data.edgenext_ecs_vpcs.all_vpcs.vpcs[0].id
#   # subnet_id  = data.edgenext_ecs_vpc_subnets.vpc_subnets.subnets[0].id
#   network_id = "68451a78-fcd6-439b-81f6-e01ed7525c16"
#   subnet_id  = "b34fe463-3930-4eb6-a177-c8039fd376b4"
# }
#
# # Import format: region/router_id/port_id
# # terraform import edgenext_ecs_router_port.example 'tokyo-a/1bac4223-a709-4639-92c7-e2f7eba9e33c/74f3a422-e0cd-4031-9f4e-cb6a2e85ca2b'
# resource "edgenext_ecs_router_port" "example" {}

# Standalone ENI (Neutron port): /ecs/openapi/v2/ports/add, detail, delete.
# Only name and description are updatable in place; other changes replace the resource.
resource "edgenext_ecs_network_interface" "example" {
  region      = "tokyo-a"
  name        = "tyd-eni"
  description = "for test"
  network_id  = "0e07db22-e210-4138-b40b-6cbdb3580b12"
  subnet_id   = "50a0f20a-16a8-46e1-8271-0d1660c739d5"
  # Optional: bind ENI to an existing ECS instance.
  device_id = "80e47fca-7822-4a71-9e51-6fe8bd232a18"
  # Optional: manage security relation.
  port_security_enabled = true
  security_groups       = ["aa2b7c0d-0e1e-4ab2-97af-14a5d5fe48cd"]
  # Optional: bind/unbind floating IP by address.
  floating_ip_address   = "148.222.161.86"
}
#
# # Import format: region/port_id
# # terraform import edgenext_ecs_network_interface.example 'tokyo-a/fa67c471-722b-4a0d-944b-9e2d741e5c5c'
# resource "edgenext_ecs_network_interface" "example" {}

# resource "edgenext_ecs_security_group" "example" {
#   region      = "tokyo-a"
#   name        = "example-sg-managed"
#   description = "A standard example security group"
# }

# Managed rule (create via /ecs/openapi/v2/security_group_rule/add). Any argument change forces replacement.
# resource "edgenext_ecs_security_group_rule" "example" {
#   region            = "tokyo-a"
#   security_group_id = data.edgenext_ecs_security_groups.all_security_groups.security_groups[0].id
#   protocol          = "tcp"
#   direction         = "ingress"
#   ethertype         = "IPv4"
#   port_range_min    = 100
#   port_range_max    = 300
#   remote_ip_prefix  = "192.168.0.0/24"
#   # remote_group_id  = ""
# }

# Import an existing rule into Terraform state. Import ID format: region/security_group_id/rule_id
# (region is normalized to lower-case by the provider). Uncomment the resource block above and run:
#
#   terraform import edgenext_ecs_security_group_rule.example 'tokyo-a/aa2b7c0d-0e1e-4ab2-97af-14a5d5fe48cd/2d259a40-4657-4324-9f3a-3c9aae529f88'
# resource "edgenext_ecs_security_group_rule" "example" {}

# #
# # SSH Key
# #
# resource "edgenext_ecs_key_pair" "example" {
#   region = "tokyo-a"
#   name = "example-key"
#   # public_key = file("~/.ssh/id_rsa.pub")
# }

# #
# # Disks
# #
# resource "edgenext_ecs_disk" "example" {
#   name        = "example-disk"
#   size        = 50
#   volume_type = "SSD"
# }

# #
# # Compute Instance
# #
# resource "edgenext_ecs_instance" "example" {
#   name            = "example-instance"
#   region          = var.region
#   flavor_ref      = "s1.small"
#   image_ref       = "centos-7.9"
#   admin_pass      = "SecurePass123!"
#   bandwidth       = 5
#   key_name        = edgenext_ecs_key_pair.example.name
#   networks        = [edgenext_ecs_vpc.example.id]
#   security_groups = [edgenext_ecs_security_group.example.id]
# }

# #
# # Tagging
# #
locals {
  ecs_tags = {
    "test-key" = "test-value"
    "team"     = "platform"
    "env"      = "dev"
  }
}

# resource "edgenext_ecs_tag" "example" {
#   for_each = local.ecs_tags

#   key   = each.key
#   value = each.value
# }

# #
# # Bind Tag IDs to a Target Resource
# #
# # Replace the values below with a real ECS resource.
# resource "edgenext_ecs_resource_tag" "example_binding" {
#   resource_uuid = "55d747cd-497b-460a-8c3d-2cb6ab2673cd"
#   resource_name = "ch11"
#   region        = "tokyo-a"
#   resource_type = 1
#   # tag_ids = [
#   #   for t in values(edgenext_ecs_tag.example) : tonumber(t.id)
#   # ]
#   tag_ids = [ 52,56,57 ]
# }

# #
# # Floating IP
# #
# resource "edgenext_ecs_floating_ip" "example" {
#   bandwidth = 10
# }
#
# data "edgenext_ecs_floating_ips" "all_floating_ips" {
#   region = "tokyo-a"
#   limit  = 10
#   eid    = "c1eae862-7725-4595-96de-6b97db7d48d9"
#   # floating_ip_address = "148.222.161.86"
# }
#
# output "floating_ip_total" {
#   value = data.edgenext_ecs_floating_ips.all_floating_ips.total
# }
#
# output "first_floating_ip_id" {
#   value = try(data.edgenext_ecs_floating_ips.all_floating_ips.floating_ips[0].id, null)
# }
#
# data "edgenext_ecs_vpcs" "all_vpcs" {
#   region = "tokyo-a"
#   limit  = 10
#   # network_id = "0e07db22-e210-4138-b40b-6cbdb3580b12"
#   # name       = "default-vpc"
# }
#
# output "vpc_total" {
#   value = data.edgenext_ecs_vpcs.all_vpcs.total
# }
#
# output "first_vpc_id" {
#   value = try(data.edgenext_ecs_vpcs.all_vpcs.vpcs[0].id, null)
# }
#
# data "edgenext_ecs_vpc_subnets" "vpc_subnets" {
#   region = "tokyo-a"
#   # network_id = data.edgenext_ecs_vpcs.all_vpcs.vpcs[0].id
#   network_id = "68451a78-fcd6-439b-81f6-e01ed7525c16"
#   router_id  = "f9883769-85e0-4cdb-ae9a-c49da3a89edb"
# }
#
# output "vpc_subnet_total" {
#   value = data.edgenext_ecs_vpc_subnets.vpc_subnets.total
# }
#
# output "first_vpc_subnet_id" {
#   value = try(data.edgenext_ecs_vpc_subnets.vpc_subnets.subnets[0].id, null)
# }
#
# data "edgenext_ecs_routers" "all_routers" {
#   region = "tokyo-a"
#   name   = "default-router"
#   limit  = 10
# }
#
# output "router_total" {
#   value = data.edgenext_ecs_routers.all_routers.total
# }
#
# output "first_router_id" {
#   value = try(data.edgenext_ecs_routers.all_routers.routers[0].id, null)
# }
#
# data "edgenext_ecs_router_ports" "router_ports" {
#   region = "tokyo-a"
#   # id     = data.edgenext_ecs_routers.all_routers.routers[0].id
#   id = "f9883769-85e0-4cdb-ae9a-c49da3a89edb"
# }

# output "router_port_total" {
#   value = data.edgenext_ecs_router_ports.router_ports.total
# }
#
# output "first_router_port_id" {
#   value = try(data.edgenext_ecs_router_ports.router_ports.ports[0].id, null)
# }
#
# data "edgenext_ecs_external_gateways" "all_external_gateways" {
#   region = "tokyo-a"
#   limit  = 10
# }
#
# output "external_gateway_total" {
#   value = data.edgenext_ecs_external_gateways.all_external_gateways.total
# }
#
# output "first_external_gateway_id" {
#   value = try(data.edgenext_ecs_external_gateways.all_external_gateways.external_gateways[0].id, null)
# }

# Network interfaces (Neutron ports) via /ecs/openapi/v2/ports/extension/list
# data "edgenext_ecs_network_interfaces" "all_ports" {
#   region = "tokyo-a"
#   name   = ""
#   limit  = 10
# }
#
# output "network_interface_total" {
#   value = data.edgenext_ecs_network_interfaces.all_ports.total
# }
#
# output "first_network_interface_id" {
#   value = try(data.edgenext_ecs_network_interfaces.all_ports.network_interfaces[0].id, null)
# }
#
# output "first_network_interface_server_name" {
#   value = try(data.edgenext_ecs_network_interfaces.all_ports.network_interfaces[0].server_name, null)
# }

#
# Data Source Usage
#
# data "edgenext_ecs_instances" "all_instances" {
#   region = "tokyo-a"
#   name = "ch11"
#   limit = 10
# }

# output "instance_id" {
#   value = try(data.edgenext_ecs_instances.all_instances.instances[0].id, null)
# }

# output "instance_id" {
#   value = edgenext_ecs_instance.example.id
# }

# data "edgenext_ecs_key_pairs" "example" {
#   region = "tokyo-a"
#   limit = 10
#   # depends_on = [ edgenext_ecs_key_pair.example ]
# }

# data "edgenext_ecs_images" "all_images" {
#   region = "tokyo-a"
#   visibility = "public"
#   name = "Debian 11.11 64-bit-v1.1"
#   page_num = 1
#   page_size = 10
# }

# output "image_id" {
#   value = data.edgenext_ecs_images.all_images.images.0.id
# }


# data "edgenext_ecs_tags" "all_tags" {
#   tag_key = "test-key"
#   tag_value = "test-value"
#   page_num = 1
#   page_size = 10
# }

# output "tag_id" {
#   value = data.edgenext_ecs_tags.all_tags.tags.0.id
# }

# data "edgenext_ecs_resource_tags" "all_resource_tags" {
#   region = "tokyo-a"
#   page_num = 1
#   page_size = 10
#   tag_key = "test-key"
#   tag_value = "test-value"
# }

# data "edgenext_ecs_security_groups" "all_security_groups" {
#   region = "tokyo-a"
#   name   = ""
#   limit  = 10
# }

# Rules for a specific security group (detail API). Requires a valid security group id.
# Example A: use the first group from the list query above (plan fails if the list is empty).
# data "edgenext_ecs_security_group_rules" "first_security_group_rules" {
#   region = "tokyo-a"
#   id     = data.edgenext_ecs_security_groups.all_security_groups.security_groups[0].id
# }

# Example B: fixed id (uncomment and replace when you already know the security group id).
# data "edgenext_ecs_security_group_rules" "example_sg_rules" {
#   region = "tokyo-a"
#   id     = "9758e912-c7f6-4609-8f33-a1e0dccf0240"
# }

# output "managed_security_group_rule_id" {
#   value = edgenext_ecs_security_group_rule.example.id
# }

# output "created_tag_ids" {
#   value = [
#     for t in values(edgenext_ecs_tag.example) : t.id
#   ]
# }

# output "resource_tag_total" {
#   value = data.edgenext_ecs_resource_tags.all_resource_tags.total
# }

# output "first_tagged_resource_id" {
#   value = try(data.edgenext_ecs_resource_tags.all_resource_tags.resource_tags[0].resource_id, null)
# }

# output "security_group_count" {
#   value = data.edgenext_ecs_security_groups.all_security_groups.count
# }

# output "first_security_group_id" {
#   value = try(data.edgenext_ecs_security_groups.all_security_groups.security_groups[0].id, null)
# }

# output "security_group_rule_count" {
#   value = length(data.edgenext_ecs_security_group_rules.first_security_group_rules.security_group_rules)
# }

# output "first_security_group_rule_id" {
#   value = try(data.edgenext_ecs_security_group_rules.first_security_group_rules.security_group_rules[0].id, null)
# }

# output "managed_security_group_id" {
#   value = edgenext_ecs_security_group.example.id
# }