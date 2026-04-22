# EdgeNext ECS Services

This package provides Terraform resources and data sources for managing EdgeNext ECS networking, storage, and metadata.

## Resources

### ECS Key Pair
- **Resource**: `edgenext_ecs_key_pair` (`ResourceENECSKeyPair`)
- **File**: `resource_en_ecs_key_pair.go`
- **Description**: Manage ECS key pairs

### ECS VPC
- **Resource**: `edgenext_ecs_vpc` (`ResourceENECSVpc`)
- **File**: `resource_en_ecs_vpc.go`
- **Description**: Manage ECS VPC networks

### ECS VPC Subnet
- **Resource**: `edgenext_ecs_vpc_subnet` (`ResourceENECSVpcSubnet`)
- **File**: `resource_en_ecs_vpc_subnet.go`
- **Description**: Manage ECS VPC subnets

### ECS Router
- **Resource**: `edgenext_ecs_router` (`ResourceENECSRouter`)
- **File**: `resource_en_ecs_router.go`
- **Description**: Manage ECS routers

### ECS Router Port
- **Resource**: `edgenext_ecs_router_port` (`ResourceENECSRouterPort`)
- **File**: `resource_en_ecs_router_port.go`
- **Description**: Attach or detach subnets on routers

### ECS Network Interface
- **Resource**: `edgenext_ecs_network_interface` (`ResourceENECSNetworkInterface`)
- **File**: `resource_en_ecs_network_interface.go`
- **Description**: Manage ECS ENI lifecycle and relation settings

### ECS Security Group
- **Resource**: `edgenext_ecs_security_group` (`ResourceENECSSecurityGroup`)
- **File**: `resource_en_ecs_security_group.go`
- **Description**: Manage ECS security groups

### ECS Security Group Rule
- **Resource**: `edgenext_ecs_security_group_rule` (`ResourceENECSSecurityGroupRule`)
- **File**: `resource_en_ecs_security_group_rule.go`
- **Description**: Manage ECS security group rules

### ECS Tag
- **Resource**: `edgenext_ecs_tag` (`ResourceENECSTag`)
- **File**: `resource_en_ecs_tag.go`
- **Description**: Manage ECS tags

### ECS Resource Tag
- **Resource**: `edgenext_ecs_resource_tag` (`ResourceENECSResourceTag`)
- **File**: `resource_en_ecs_resource_tag.go`
- **Description**: Bind tags to ECS resources

## ECS Update Behavior

Several ECS resources now reject immutable argument changes directly during plan/apply instead of replacing resources automatically:

- `edgenext_ecs_key_pair`: `name`, `public_key`
- `edgenext_ecs_network_interface`: `network_id`, `subnet_id`
- `edgenext_ecs_tag`: `key`, `value`
- `edgenext_ecs_resource_tag`: `resource_uuid`, `resource_name`, `resource_type`
- `edgenext_ecs_router_port`: `router_id`, `network_id`, `subnet_id`
- `edgenext_ecs_security_group_rule`: all arguments except `region`
- `edgenext_ecs_vpc_subnet`: all arguments except `region`
- `edgenext_ecs_vpc`: `subnet` and all nested subnet fields

## Data Sources

### ECS Instances
- **Data Source**: `edgenext_ecs_instances` (`DataSourceENECSInstances`)
- **File**: `data_source_en_ecs_instances.go`
- **Description**: Query ECS instance list

### ECS Images
- **Data Source**: `edgenext_ecs_images` (`DataSourceENECSImages`)
- **File**: `data_source_en_ecs_images.go`
- **Description**: Query ECS image list

### ECS Key Pairs
- **Data Source**: `edgenext_ecs_key_pairs` (`DataSourceENECSKeyPairs`)
- **File**: `data_source_en_ecs_key_pairs.go`
- **Description**: Query ECS key pair list

### ECS VPCs
- **Data Source**: `edgenext_ecs_vpcs` (`DataSourceENECSVpcs`)
- **File**: `data_source_en_ecs_vpcs.go`
- **Description**: Query ECS VPC list

### ECS External Gateways
- **Data Source**: `edgenext_ecs_external_gateways` (`DataSourceENECSExternalGateways`)
- **File**: `data_source_en_ecs_external_gateways.go`
- **Description**: Query external gateway networks

### ECS VPC Subnets
- **Data Source**: `edgenext_ecs_vpc_subnets` (`DataSourceENECSVpcSubnets`)
- **File**: `data_source_en_ecs_vpc_subnets.go`
- **Description**: Query VPC subnet list

### ECS Routers
- **Data Source**: `edgenext_ecs_routers` (`DataSourceENECSRouters`)
- **File**: `data_source_en_ecs_routers.go`
- **Description**: Query ECS router list

### ECS Router Ports
- **Data Source**: `edgenext_ecs_router_ports` (`DataSourceENECSRouterPorts`)
- **File**: `data_source_en_ecs_router_ports.go`
- **Description**: Query router port list

### ECS Floating IPs
- **Data Source**: `edgenext_ecs_floating_ips` (`DataSourceENECSFloatingIps`)
- **File**: `data_source_en_ecs_floating_ips.go`
- **Description**: Query floating IP list

### ECS Network Interfaces
- **Data Source**: `edgenext_ecs_network_interfaces` (`DataSourceENECSNetworkInterfaces`)
- **File**: `data_source_en_ecs_network_interfaces.go`
- **Description**: Query ENI list

### ECS Security Groups
- **Data Source**: `edgenext_ecs_security_groups` (`DataSourceENECSSecurityGroups`)
- **File**: `data_source_en_ecs_security_groups.go`
- **Description**: Query security group list

### ECS Disks
- **Data Source**: `edgenext_ecs_disks` (`DataSourceENECSDisks`)
- **File**: `data_source_en_ecs_disks.go`
- **Description**: Query disk list

### ECS Tags
- **Data Source**: `edgenext_ecs_tags` (`DataSourceENECSTags`)
- **File**: `data_source_en_ecs_tags.go`
- **Description**: Query tag list

### ECS Security Group Rules
- **Data Source**: `edgenext_ecs_security_group_rules` (`DataSourceENECSSecurityGroupRules`)
- **File**: `data_source_en_ecs_security_group_rules.go`
- **Description**: Query security group rule list

### ECS Resource Tags
- **Data Source**: `edgenext_ecs_resource_tags` (`DataSourceENECSResourceTags`)
- **File**: `data_source_en_ecs_resource_tags.go`
- **Description**: Query resource tag relations

## File Structure

```
edgenext/services/ecs/
├── README.md                                      # This documentation
├── resource_en_ecs_*.go                           # ECS resource implementations
├── resource_en_ecs_*.md                           # ECS resource documentation
├── data_source_en_ecs_*.go                        # ECS data source implementations
└── data_source_en_ecs_*.md                        # ECS data source documentation
```

## Usage Examples

### Create VPC and Subnet

```hcl
resource "edgenext_ecs_vpc" "example" {
  region = "tokyo-a"
  name   = "example-vpc"
  subnet {
    name       = "example-subnet"
    ip_version = 4
    cidr       = "172.31.1.0/24"
  }
}

resource "edgenext_ecs_vpc_subnet" "example" {
  region     = "tokyo-a"
  network_id = edgenext_ecs_vpc.example.id
  name       = "example-subnet-2"
  cidr       = "172.31.2.0/24"
}
```

### Query Existing Network Interfaces

```hcl
data "edgenext_ecs_network_interfaces" "example" {
  region = "tokyo-a"
  name   = "example-eni"
  limit  = 20
}

output "eni_ids" {
  value = data.edgenext_ecs_network_interfaces.example.network_interfaces[*].id
}
```

### Create Security Group and Rule

```hcl
resource "edgenext_ecs_security_group" "example" {
  region = "tokyo-a"
  name   = "example-sg"
}

resource "edgenext_ecs_security_group_rule" "ssh" {
  region            = "tokyo-a"
  security_group_id = edgenext_ecs_security_group.example.id
  protocol          = "tcp"
  direction         = "ingress"
  ethertype         = "IPv4"
  port_range_min    = 22
  port_range_max    = 22
  remote_ip_prefix  = "0.0.0.0/0"
}
```
