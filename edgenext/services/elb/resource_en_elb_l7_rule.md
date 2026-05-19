Use this resource to manage a standalone L7 rule under an EdgeNext ELB L7 policy.

Updatable in place: `type`, `compare_type`, `value`. `l7policy_id` forces replacement. Read uses `l7rules/list` with an id filter. Delete waits until the rule is removed before create on replacement.

Example Usage

```hcl
resource "edgenext_elb_l7_policy" "api_redirect" {
  listener_id    = edgenext_elb_listener.https.id
  name           = "api-redirect"
  action         = "REDIRECT_TO_URL"
  redirect_url   = "https://api.example.com/"
}

resource "edgenext_elb_l7_rule" "api_prefix" {
  l7policy_id  = edgenext_elb_l7_policy.api_redirect.id
  type         = "PATH"
  compare_type = "STARTS_WITH"
  value        = "/api"
}
```

Import

Import format is `l7policy_id/l7rule_id`. You may also use `listener_id/l7rule_id`; the provider resolves the policy on that listener.

```shell
terraform import edgenext_elb_l7_rule.api_prefix c6f3f52d-e331-487f-bb9f-c33c3f7f51ba/afea81e4-0a08-45e6-9d68-2ffefc10ba54
```
