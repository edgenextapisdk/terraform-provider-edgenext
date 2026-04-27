variable "domain_id" {
  description = "Domain ID for the domain template."
  type        = number
}

variable "template_source_id" {
  description = "Source template ID. If not specified, uses the global template."
  type        = number
  default     = null
}

variable "ddos_config" {
  description = "Optional DDoS protection configuration."
  type = object({
    application_ddos_protection = optional(object({
      status                = optional(string)
      ai_cc_status          = optional(string)
      type                  = optional(string)
      need_attack_detection = optional(number)
      ai_status             = optional(string)
    }))
    visitor_authentication = optional(object({
      status           = optional(string)
      auth_token       = optional(string)
      pass_still_check = optional(number)
    }))
  })
  default = null
}

variable "waf_rule_config" {
  description = "Optional WAF protection configuration."
  type = object({
    waf_rule_config = optional(object({
      status          = optional(string)
      ai_status       = optional(string)
      waf_level       = optional(string)
      waf_mode        = optional(string)
      waf_strategy_id = optional(number)
    }))
    waf_intercept_page = optional(object({
      status  = optional(string)
      type    = optional(string)
      content = optional(string)
    }))
  })
  default = null
}

variable "precise_access_control_config" {
  description = "Optional precise access control configuration."
  type        = any
  default     = null
}
