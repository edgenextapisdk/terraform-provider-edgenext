package scdn

// HTTP Methods
const (
	MethodGET    = "GET"
	MethodPOST   = "POST"
	MethodPUT    = "PUT"
	MethodDELETE = "DELETE"
)

// API Endpoints - Domain Management
const (

	// ============================================================================
	// Domain Management API Endpoints
	// ============================================================================
	// Domain CRUD operations
	EndpointDomains        = "/api/v5/domains"         // GET, POST, PUT, DELETE
	EndpointDomainsSimple  = "/api/v5/domains/simple"  // GET
	EndpointDomainsDisable = "/api/v5/domains_disable" // POST
	EndpointDomainsEnable  = "/api/v5/domains_enable"  // POST

	// Domain certificate operations
	EndpointDomainsBindCert   = "/api/v5/domains/bind_cert"   // POST
	EndpointDomainsUnbindCert = "/api/v5/domains/unbind_cert" // POST

	// Domain access operations
	EndpointDomainsAccessRefresh      = "/api/v5/domains/access_refresh"       // POST
	EndpointDomainsAccessSwitch       = "/api/v5/domains/access_switch"        // POST
	EndpointDomainsAccessProgress     = "/api/v5/domains/access_progress"      // GET
	EndpointDomainsAccessInfoDownload = "/api/v5/domains/access_info_download" // POST

	// Domain export
	EndpointDomainsExport = "/api/v5/domains/domains_export" // POST

	// Domain node operations
	EndpointDomainsNodesSwitch = "/api/v5/domains/nodes_switch" // POST

	// Domain base settings
	EndpointDomainsBaseSettings = "/api/v5/domains/base_settings" // GET, PUT

	// Domain templates
	EndpointDomainsTemplates = "/api/v5/domains/templates" // GET

	// Domain origins
	EndpointDomainsOrigins = "/api/v5/domains/origins" // GET, POST, PUT, DELETE

	// Brief domains
	EndpointBriefDomains = "/api/v5/brief_domains" // POST

	// ============================================================================
	// Certificate Management API Endpoints
	// ============================================================================
	// Certificate CRUD operations
	EndpointCASelfAdd     = "/api/v5/Web.ca.self.add"          // POST - Upload certificate
	EndpointCATextSave    = "/api/v5/Web.ca.text.save"         // POST - Add/Update certificate (text)
	EndpointCAInfoEdit    = "/api/v5/Web.ca.info.edit"         // POST - Edit certificate
	EndpointCASelfList    = "/api/v5/Web.ca.self.list"         // GET - List certificates
	EndpointCASelf        = "/api/v5/Web.ca.self"              // GET - Get certificate detail
	EndpointCASelfExport  = "/api/v5/Web.ca.self.export"       // GET - Export certificate
	EndpointCASelfDel     = "/api/v5/Web.ca.self.del"          // DELETE - Delete certificate
	EndpointCAEditName    = "/api/v5/Web.ca.self.editcaname"   // POST - Edit certificate name
	EndpointCABatchList   = "/api/v5/Web.Domain.batch.ca.list" // POST - List certificates by domains
	EndpointCAApplyAdd    = "/api/v5/Web.ca.apply.add"         // POST - Apply for certificate
	EndpointCABatchOperat = "/api/v5/Web.ca.batch.operat"      // POST - Batch certificate operations

	// ============================================================================
	// Rule Template Management API Endpoints
	// ============================================================================
	// Rule template CRUD operations
	EndpointRuleTemplates        = "/api/v5/ruletpls"               // GET, POST, PUT, DELETE
	EndpointRuleTemplatesBind    = "/api/v5/ruletpls/bind_domain"   // PUT - Bind domain to template
	EndpointRuleTemplatesUnbind  = "/api/v5/ruletpls/unbind_domain" // PUT - Unbind domain from template
	EndpointRuleTemplatesDomains = "/api/v5/ruletpls/domains"       // GET - List domains bound to template

	// ============================================================================
	// Network Speed Management API Endpoints
	// ============================================================================
	EndpointNetworkSpeedGetConfig    = "/api/v5/ruletpl/network_speed/get_conf"  // POST - Get template config
	EndpointNetworkSpeedUpdateConfig = "/api/v5/ruletpl/network_speed/conf"      // PUT - Update template config
	EndpointNetworkSpeedRules        = "/api/v5/ruletpl/network_speed/rules"     // GET - Get rules list, POST - Create rule
	EndpointNetworkSpeedRule         = "/api/v5/ruletpl/network_speed/rule"      // DELETE - Delete rule, POST - Update rule
	EndpointNetworkSpeedRuleSort     = "/api/v5/ruletpl/network_speed/rule_sort" // PUT - Sort rules

	// ============================================================================
	// Cache Rule Management API Endpoints
	// ============================================================================
	EndpointCacheRules      = "/api/v5/ruletpl/cache/rules"       // GET - Get rules list, POST - Create rule
	EndpointCacheRule       = "/api/v5/ruletpl/cache/rule"        // PUT - Update rule name/remark, DELETE - Delete rule
	EndpointCacheRuleConf   = "/api/v5/ruletpl/cache/rule/conf"   // PUT - Update rule configuration
	EndpointCacheRuleStatus = "/api/v5/ruletpl/cache/rule_status" // PUT - Update rule status (enable/disable)
	EndpointCacheRuleSort   = "/api/v5/ruletpl/cache/rule_sort"   // PUT - Sort rules
	EndpointCacheGlobalConf = "/api/v5/ruletpl/cache/global/conf" // GET - Get global cache config

	// ============================================================================
	// Security Protection API Endpoints
	// ============================================================================
	// DDoS Protection
	EndpointSecurityProtectionDdosConfigs = "/api/v5/security_protection/ddos/configs" // GET - Get DDoS config, PUT - Update DDoS config

	// WAF Rule Config
	EndpointSecurityProtectionWafRules = "/api/v5/security_protection/waf/rules" // GET - Get WAF config, PUT - Update WAF config

	// Security Protection Template
	EndpointSecurityProtectionTemplateMemberGlobal        = "/api/v5/security_protection/template/member/global"         // GET - Get member global template
	EndpointSecurityProtectionTemplate                    = "/api/v5/security_protection/template"                       // POST - Create template, PUT - Edit template, DELETE - Delete template
	EndpointSecurityProtectionTemplateDomain              = "/api/v5/security_protection/template/domain"                // POST - Create domain template
	EndpointSecurityProtectionTemplateSearch              = "/api/v5/security_protection/template/search"                // POST - Get template list
	EndpointSecurityProtectionTemplateDomainBindSearch    = "/api/v5/security_protection/template/domain/bind/search"    // POST - Get template bind domain list
	EndpointSecurityProtectionTemplateDomainBind          = "/api/v5/security_protection/template/domain/bind"           // POST - Bind template domain
	EndpointSecurityProtectionTemplateBatchConfig         = "/api/v5/security_protection/template/batch/config"          // POST - Batch config template
	EndpointSecurityProtectionTemplateDomainUnboundSearch = "/api/v5/security_protection/template/domain/unbound/search" // POST - Get unbound template domain list

	// Security Protection Iota
	EndpointSecurityProtectionIota = "/api/v5/security_protection/template/iota" // GET - Get iota enum values

	// ============================================================================
	// Origin Group Management API Endpoints
	// ============================================================================
	EndpointOriginGroups            = "/api/v5/origin_groups"                     // GET - List origin groups, POST - Create origin group, PUT - Update origin group, DELETE - Delete origin groups
	EndpointOriginGroupsDetail      = "/api/v5/origin_groups/detail"              // GET - Get origin group detail
	EndpointOriginGroupsBindDomains = "/api/v5/origin_groups/domains_bind"        // POST - Bind origin group to domains
	EndpointOriginGroupsAll         = "/api/v5/origin_groups/all"                 // GET - Get all origin groups
	EndpointOriginGroupsCopy        = "/api/v5/origin_groups/copy"                // POST - Copy origin group to domain
	EndpointOriginGroupsBindHistory = "/api/v5/origin_groups/bind_history/latest" // GET - Get latest bind history

	// ============================================================================
	// Cache Clean and Preheat API Endpoints
	// ============================================================================
	// Cache clean operations
	EndpointCacheCleanGetConfig  = "/api/v5/Web.Domain.DashBoard.getCache"           // GET - Get cache clean config list
	EndpointCacheCleanSave       = "/api/v5/Web.Domain.DashBoard.saveCache"          // PUT - Submit cache clean task
	EndpointCacheCleanTaskList   = "/api/v5/Web.Domain.DashBoard.cache.clean.list"   // GET - Get cache clean task list
	EndpointCacheCleanTaskDetail = "/api/v5/Web.Domain.DashBoard.cache.clean.detail" // GET - Get cache clean task detail

	// Cache preheat operations
	EndpointCachePreheatTaskList = "/api/v5/Web.Domain.DashBoard.get.preheat.cache.new.list" // POST - Get preheat task list (new API)
	EndpointCachePreheatSave     = "/api/v5/Web.Domain.DashBoard.save.preheat.cache"         // POST - Submit preheat task

	// ============================================================================
	// Log Download API Endpoints
	// ============================================================================
	// Log download task operations
	EndpointLogDownloadTaskList        = "/api/v5/soc.log.download.task.list"         // GET/POST - List log download tasks
	EndpointLogDownloadTaskAdd         = "/api/v5/soc.log.download.task.add"          // POST - Add log download task
	EndpointLogDownloadTaskCancel      = "/api/v5/soc.log.download.task.cancel"       // POST - Cancel log download task
	EndpointLogDownloadTaskBatchCancel = "/api/v5/soc.log.download.task.batch.cancel" // DELETE - Batch cancel log download tasks
	EndpointLogDownloadTaskDelete      = "/api/v5/soc.log.download.task.del"          // POST - Delete log download task
	EndpointLogDownloadTaskBatchDelete = "/api/v5/soc.log.download.task.batch.del"    // DELETE - Batch delete log download tasks
	EndpointLogDownloadTaskRegenerate  = "/api/v5/soc.log.download.task.regenerate"   // POST - Regenerate log download task

	// Log download fields
	EndpointLogDownloadFields = "/api/v5/soc.log.download.fields" // GET - Get log download fields

	// Log download template operations
	EndpointLogDownloadTemplateList              = "/api/v5/soc.log.download.template.list"                // GET/POST - List log download templates
	EndpointLogDownloadTemplateDomainList        = "/api/v5/soc.log.download.template.domain.list"         // GET - List template domains
	EndpointLogDownloadTemplateAdd               = "/api/v5/soc.log.download.template.add"                 // POST - Add log download template
	EndpointLogDownloadTemplateSave              = "/api/v5/soc.log.download.template.save"                // POST - Save (update) log download template
	EndpointLogDownloadTemplateDelete            = "/api/v5/soc.log.download.template.del"                 // DELETE - Delete log download template
	EndpointLogDownloadTemplateBatchDelete       = "/api/v5/soc.log.download.template.batch.del"           // DELETE - Batch delete log download templates
	EndpointLogDownloadTemplateChangeStatus      = "/api/v5/soc.log.download.template.change.status"       // POST - Change template status
	EndpointLogDownloadTemplateBatchChangeStatus = "/api/v5/soc.log.download.template.batch.change.status" // POST - Batch change template status
	EndpointLogDownloadTemplateAll               = "/api/v5/soc.log.download.template.all"                 // GET/POST - Get all templates (for adding tasks)
	EndpointLogDownloadTemplateGroupAll          = "/api/v5/soc.log.download.template.group.all"           // GET/POST - Get all template groups
)
