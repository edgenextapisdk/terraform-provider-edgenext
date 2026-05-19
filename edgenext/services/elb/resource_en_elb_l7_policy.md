Use this resource to manage an EdgeNext ELB L7 policy on a listener (for example redirect to a URL or another target group).

Match rules are **not** nested in this resource; create them with `edgenext_elb_l7_rule`. `position` is computed from the API. Actions use `REDIRECT_TO_TARGET_GROUP` in Terraform (mapped to the Octavia API as needed).

Example Usage

```hcl
resource "edgenext_elb_l7_policy" "redirect" {
  listener_id              = edgenext_elb_listener.https.id
  name                     = "redirect-legacy"
  description              = "Send legacy path to another pool"
  action                   = "REDIRECT_TO_TARGET_GROUP"
  redirect_target_group_id = edgenext_elb_target_group.legacy.id
}

resource "edgenext_elb_l7_rule" "legacy_path" {
  l7policy_id  = edgenext_elb_l7_policy.redirect.id
  type         = "PATH"
  compare_type = "EQUAL_TO"
  value        = "/old"
}
```

Import

Import format is `l7policy_id`.

```shell
terraform import edgenext_elb_l7_policy.redirect 6232fc9e-a76e-4e53-bb17-d641ac21e91a
```
