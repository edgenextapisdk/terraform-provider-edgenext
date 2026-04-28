variable "access_key" {
  description = "EdgeNext Access Key."
  type        = string
  sensitive   = true
}

variable "secret_key" {
  description = "EdgeNext Secret Key."
  type        = string
  sensitive   = true
}

variable "endpoint" {
  description = "EdgeNext SCDN API endpoint."
  type        = string
  default     = "https://api.edgenextscdn.com"
}

variable "security_template_name" {
  description = "Security protection template name."
  type        = string
  default     = null
}

variable "security_template_remark" {
  description = "Security protection template remark."
  type        = string
  default     = ""
}

variable "security_template_source_id" {
  description = "Existing security protection template ID to clone from."
  type        = number
  default     = null
}

variable "security_domain_ids" {
  description = "Domain IDs to bind when creating the security protection template."
  type        = list(number)
  default     = []
}

variable "security_group_ids" {
  description = "Domain group IDs to bind when creating the security protection template."
  type        = list(number)
  default     = []
}

variable "security_domains" {
  description = "Domains to bind when creating the security protection template."
  type        = list(string)
  default     = []
}

variable "security_bind_all" {
  description = "Whether to bind all domains when creating the security protection template."
  type        = bool
  default     = false
}

variable "security_ddos_config" {
  description = "Optional DDoS protection config passed to the security protection module."
  type        = any
  default     = null
}

variable "security_waf_rule_config" {
  description = "Optional WAF protection config passed to the security protection module."
  type        = any
  default     = null
}

variable "security_precise_access_control_config" {
  description = "Optional precise access control config passed to the security protection module."
  type        = any
  default     = null
}

# Domain template variables (single domain)
variable "domain_template_id" {
  description = "Domain ID for single domain template creation."
  type        = number
  default     = null
}

variable "domain_template_source_id" {
  description = "Source template ID for single domain template. If not specified, uses the global template."
  type        = number
  default     = null
}

variable "domain_template_ddos_config" {
  description = "Optional DDoS protection config for single domain template."
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

variable "domain_template_waf_rule_config" {
  description = "Optional WAF protection config for single domain template."
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

variable "domain_template_precise_access_control_config" {
  description = "Optional precise access control config for single domain template."
  type        = any
  default     = null
}