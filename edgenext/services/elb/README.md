# EdgeNext ELB Services

This package provides Terraform resources and data sources for managing EdgeNext Elastic Load Balancer (ELB) objects: listeners, target groups, backends, TLS certificates, and L7 routing policies.

API paths use OpenStack/Octavia-style names internally (for example `pools` for target groups). User-facing Terraform field names and documentation use **target group** and **target** instead of pool/member.

## Resources

### ELB Certificate
- **Resource**: `edgenext_elb_certificate` (`ResourceENELBCertificate`)
- **File**: `resource_en_elb_certificate.go`
- **Description**: Upload and delete Barbican TLS certificates used by `TERMINATED_HTTPS` listeners. No update API; changing PEM material forces replacement.

### ELB Listener
- **Resource**: `edgenext_elb_listener` (`ResourceENELBListener`)
- **File**: `resource_en_elb_listener.go`
- **Description**: Create and manage a listener on a load balancer. Updatable: `name`, `description`, `default_tls_certificate_ref`. Delete waits until the listener is gone before returning (async Octavia).

### ELB Target Group
- **Resource**: `edgenext_elb_target_group` (`ResourceENELBTargetGroup`)
- **File**: `resource_en_elb_target_group.go`
- **Description**: Target group plus health monitor on a listener. Backends are **not** inline; use `edgenext_elb_target_group_attachment`. Updatable: `name`, `description` only. Create/delete use Octavia immutable retry and provisioning waits.

### ELB Target Group Attachment
- **Resource**: `edgenext_elb_target_group_attachment` (`ResourceENELBTargetGroupAttachment`)
- **File**: `resource_en_elb_target_group_attachment.go`
- **Description**: One backend (target) on a target group. Updatable: `name`, `weight`, `protocol_port`. `target_group_id` and `address` are ForceNew. Delete waits until the target is removed from the API.

### ELB L7 Policy
- **Resource**: `edgenext_elb_l7_policy` (`ResourceENELBL7Policy`)
- **File**: `resource_en_elb_l7_policy.go`
- **Description**: L7 policy on a listener (redirect to URL or target group). Match rules are managed with `edgenext_elb_l7_rule`, not nested blocks. `position` is computed from the API.

### ELB L7 Rule
- **Resource**: `edgenext_elb_l7_rule` (`ResourceENELBL7Rule`)
- **File**: `resource_en_elb_l7_rule.go`
- **Description**: Standalone L7 rule under a policy. Updatable: `type`, `compare_type`, `value`. Delete waits until the rule disappears before create on replacement (async load balancer).

## ELB Update and ForceNew Behavior

| Resource | In-place update | ForceNew (replace) |
|----------|-----------------|-------------------|
| `edgenext_elb_certificate` | — (no update) | `name`, `certificate`, `private_key` |
| `edgenext_elb_listener` | `name`, `description`, `default_tls_certificate_ref` | `loadbalancer_id`, `protocol`, `protocol_port`, `connection_limit`, `insert_headers` |
| `edgenext_elb_target_group` | `name`, `description` | `listener_id`, `protocol`, `lb_algorithm`, `health_monitor` |
| `edgenext_elb_target_group_attachment` | `name`, `weight`, `protocol_port` | `target_group_id`, `address` |
| `edgenext_elb_l7_policy` | `name`, `description`, `action`, `redirect_url`, `redirect_target_group_id` | `listener_id` |
| `edgenext_elb_l7_rule` | `type`, `compare_type`, `value` | `l7policy_id` |

## Octavia Async Operations

Several operations hit transient `409 immutable` responses while Octavia finishes a prior change. The package uses `elbPOSTWithOctaviaImmutableRetry` and post-delete polling where needed:

- **Target group create** — waits for `provisioning_status ACTIVE` before creating the health monitor; create POSTs are retried.
- **Target group / listener / target attachment / L7 rule delete** — polls until the object is gone before returning (configure `timeouts` on the resource).
- **Target attachment / L7 rule create** — POST retries after delete on replacement.

If replace-after-delete fails with immutable errors, increase `timeouts.delete` (and sometimes `timeouts.create`) and re-apply. See `elb_octavia_retry.go`.

## Data Sources

### ELB Load Balancers
- **Data Source**: `edgenext_elb_load_balancers` (`DataSourceENELBLoadBalancers`)
- **File**: `data_source_en_elb_load_balancers.go`
- **Description**: List load balancers with listeners, target groups, VIP, spec, and status.

### ELB Certificates
- **Data Source**: `edgenext_elb_certificates` (`DataSourceENELBCertificates`)
- **File**: `data_source_en_elb_certificates.go`
- **Description**: Paginated certificate list with optional name filter.

### ELB Listeners
- **Data Source**: `edgenext_elb_listeners` (`DataSourceENELBListeners`)
- **File**: `data_source_en_elb_listeners.go`
- **Description**: List listeners for a load balancer.

### ELB Target Groups
- **Data Source**: `edgenext_elb_target_groups` (`DataSourceENELBTargetGroups`)
- **File**: `data_source_en_elb_target_groups.go`
- **Description**: List target groups for a load balancer, including health monitor summary when available.

### ELB Target Group Attachments
- **Data Source**: `edgenext_elb_target_group_attachments` (`DataSourceENELBTargetGroupAttachments`)
- **File**: `data_source_en_elb_target_group_attachments.go`
- **Description**: List backends for one target group.

### ELB L7 Policies
- **Data Source**: `edgenext_elb_l7_policies` (`DataSourceENELBL7Policies`)
- **File**: `data_source_en_elb_l7_policies.go`
- **Description**: List L7 policies on a listener.

### ELB L7 Rules
- **Data Source**: `edgenext_elb_l7_rules` (`DataSourceENELBL7Rules`)
- **File**: `data_source_en_elb_l7_rules.go`
- **Description**: List L7 rules for a policy with optional filters.

## File Structure

```
edgenext/services/elb/
├── README.md                                      # This documentation
├── elb_octavia_retry.go                           # Shared Octavia 409 immutable retry
├── resource_en_elb_*.go                           # ELB resource implementations
├── resource_en_elb_*.md                           # Per-resource docs (gendoc / registry)
├── data_source_en_elb_*.go                        # ELB data source implementations
└── data_source_en_elb_*.md                        # Per-data-source docs
```

## Usage Examples

See `examples/elb/main.tf` for full HCL. Minimal flow:

### Listener, target group, and backends

```hcl
resource "edgenext_elb_listener" "https" {
  loadbalancer_id             = var.loadbalancer_id
  name                        = "https"
  protocol                    = "TERMINATED_HTTPS"
  protocol_port               = 443
  default_tls_certificate_ref = edgenext_elb_certificate.site.certificate_ref
}

resource "edgenext_elb_target_group" "app" {
  listener_id  = edgenext_elb_listener.https.id
  name         = "app"
  protocol     = "HTTP"
  lb_algorithm = "ROUND_ROBIN"
  health_monitor {
    name           = "http-check"
    type           = "HTTP"
    max_retries    = 3
    delay          = 10
    timeout        = 5
    http_method    = "GET"
    url_path       = "/health"
    expected_codes = "200"
  }
}

resource "edgenext_elb_target_group_attachment" "app1" {
  target_group_id = edgenext_elb_target_group.app.id
  address         = "192.168.0.10"
  protocol_port   = 8080
  weight          = 10
}
```

### L7 redirect policy and rule

```hcl
resource "edgenext_elb_l7_policy" "redirect" {
  listener_id              = edgenext_elb_listener.https.id
  name                     = "redirect-old-path"
  action                   = "REDIRECT_TO_TARGET_GROUP"
  redirect_target_group_id = edgenext_elb_target_group.app.id
}

resource "edgenext_elb_l7_rule" "old_path" {
  l7policy_id  = edgenext_elb_l7_policy.redirect.id
  type         = "PATH"
  compare_type = "EQUAL_TO"
  value        = "/legacy"
}
```

### Query load balancers and listeners

```hcl
data "edgenext_elb_load_balancers" "all" {
  limit = 20
}

data "edgenext_elb_listeners" "lb" {
  loadbalancer_id = data.edgenext_elb_load_balancers.all.balancers[0].id
  limit           = 100
}
```

## Import

| Resource | Import ID format |
|----------|------------------|
| `edgenext_elb_certificate` | `certificate_id` |
| `edgenext_elb_listener` | `listener_id` |
| `edgenext_elb_target_group` | (no importer; adopt via data source + new resource) |
| `edgenext_elb_target_group_attachment` | `target_group_id/target_id` |
| `edgenext_elb_l7_policy` | `l7policy_id` |
| `edgenext_elb_l7_rule` | `l7policy_id/l7rule_id` or `listener_id/l7rule_id` (resolves policy) |
