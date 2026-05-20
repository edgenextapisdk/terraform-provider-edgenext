package edgenext

import (
	"context"
	"fmt"
	"strings"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/cdn"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/ecs"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/eip"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/elb"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/oss"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/rds"
	scdncache "github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn/cache"
	scdncacheoperate "github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn/cache_operate"
	scdncert "github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn/cert"
	scdndomain "github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn/domain"
	scdndomaingroupdata "github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn/domain_group/data"
	scdndomaingroupresource "github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn/domain_group/resource"
	scdnipdata "github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn/ip/data"
	scdnipresource "github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn/ip/resource"
	scdnlogdownload "github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn/log_download"
	scdnnetworkspeed "github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn/network_speed"
	scdnorigingroup "github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn/origin_group"
	scdnsecurityprotect "github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn/security_protect"
	scdntemplate "github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn/template"
	sdnsdomain "github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/sdns/domain"
	sdnsgroup "github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/sdns/domain_group"
	sdnsrecord "github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/sdns/record"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/ssl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider returns the EdgeNext CDN Terraform Provider
func Provider() *schema.Provider {
	// Initialize domain module resources and data sources
	domainResources := scdndomain.Resources()
	domainDataSources := scdndomain.DataSources()

	// Initialize cert module resources and data sources
	certResources := scdncert.Resources()
	certDataSources := scdncert.DataSources()

	// Initialize template module resources and data sources
	templateResources := scdntemplate.Resources()
	templateDataSources := scdntemplate.DataSources()

	// Initialize network speed module resources and data sources
	networkSpeedResources := scdnnetworkspeed.Resources()
	networkSpeedDataSources := scdnnetworkspeed.DataSources()

	// Initialize cache module resources and data sources
	cacheResources := scdncache.Resources()
	cacheDataSources := scdncache.DataSources()

	// Initialize security protection module resources and data sources
	securityProtectResources := scdnsecurityprotect.Resources()
	securityProtectDataSources := scdnsecurityprotect.DataSources()

	// Initialize origin group module resources and data sources
	originGroupResources := scdnorigingroup.Resources()
	originGroupDataSources := scdnorigingroup.DataSources()

	// Initialize cache operate module resources and data sources
	cacheOperateResources := scdncacheoperate.Resources()
	cacheOperateDataSources := scdncacheoperate.DataSources()

	// Initialize log download module resources and data sources
	logDownloadResources := scdnlogdownload.Resources()
	logDownloadDataSources := scdnlogdownload.DataSources()

	// Initialize SDNS module resources and data sources
	sdnsDomainResources := sdnsdomain.Resources()
	sdnsDomainDataSources := sdnsdomain.DataSources()
	sdnsGroupResources := sdnsgroup.Resources()
	sdnsGroupDataSources := sdnsgroup.DataSources()
	sdnsRecordResources := sdnsrecord.Resources()
	sdnsRecordDataSources := sdnsrecord.DataSources()

	// Build resources map (merge with duplicate-key detection per source)
	resourcesMap := map[string]*schema.Resource{
		// CDN domain and configuration management resources
		"edgenext_cdn_domain": cdn.ResourceEdgenextCdnDomainConfig(),

		// CDN cache prefetch and file purge resources
		"edgenext_cdn_purge":    cdn.ResourceEdgenextCdnPurge(),
		"edgenext_cdn_prefetch": cdn.ResourceEdgenextCdnPrefetch(),

		// SSL certificate management resources
		"edgenext_ssl_certificate": ssl.ResourceEdgenextSslCertificate(),

		// OSS object storage resources
		"edgenext_oss_bucket": oss.ResourceOSSBucket(),
		// OSS object management resources
		"edgenext_oss_object":      oss.ResourceOSSObject(),
		"edgenext_oss_object_copy": oss.ResourceOSSObjectCopy(),

		// ECS resources
		"edgenext_ecs_key_pair":                           ecs.ResourceENECSKeyPair(),
		"edgenext_ecs_vpc":                                ecs.ResourceENECSVpc(),
		"edgenext_ecs_vpc_subnet":                         ecs.ResourceENECSVpcSubnet(),
		"edgenext_ecs_router":                             ecs.ResourceENECSRouter(),
		"edgenext_ecs_router_port":                        ecs.ResourceENECSRouterPort(),
		"edgenext_ecs_network_interface":                  ecs.ResourceENECSNetworkInterface(),
		"edgenext_ecs_network_interface_instance_binding": ecs.ResourceENECSNetworkInterfaceInstanceBinding(),
		"edgenext_ecs_security_group":                     ecs.ResourceENECSSecurityGroup(),
		"edgenext_ecs_security_group_rule":                ecs.ResourceENECSSecurityGroupRule(),
		"edgenext_ecs_tag":                                ecs.ResourceENECSTag(),
		"edgenext_ecs_instance_tag":                       ecs.ResourceENECSInstanceTag(),
		"edgenext_ecs_instance_power":                     ecs.ResourceENECSInstancePower(),
		"edgenext_ecs_instance_reboot":                    ecs.ResourceENECSInstanceReboot(),

		// RDS resources
		"edgenext_rds_backup":                           rds.ResourceENRDSBackup(),
		"edgenext_rds_backup_policy":                    rds.ResourceENRDSBackupPolicy(),
		"edgenext_rds_backup_policy_associate_instance": rds.ResourceENRDSBackupPolicyAssociateInstance(),
		"edgenext_rds_database":                         rds.ResourceENRDSDatabase(),
		"edgenext_rds_account":                          rds.ResourceENRDSAccount(),
		"edgenext_rds_account_privilege":                rds.ResourceENRDSAccountPrivilege(),
		"edgenext_rds_account_root_password":            rds.ResourceENRDSAccountRootPassword(),

		// ELB resources
		"edgenext_elb_certificate":             elb.ResourceENELBCertificate(),
		"edgenext_elb_listener":                elb.ResourceENELBListener(),
		"edgenext_elb_target_group":            elb.ResourceENELBTargetGroup(),
		"edgenext_elb_target_group_attachment": elb.ResourceENELBTargetGroupAttachment(),
		"edgenext_elb_l7_policy":               elb.ResourceENELBL7Policy(),
		"edgenext_elb_l7_rule":                 elb.ResourceENELBL7Rule(),

		// EIP resources
		"edgenext_eip_association": eip.ResourceENEIPAssociation(),

		// SCDN domain management resources (from domain module)
		// Note: These resources are organized under scdn/domain/ for better module management

		// Domain Group resources
		"edgenext_scdn_domain_group": scdndomaingroupresource.ResourceEdgenextScdnDomainGroup(),

		// User IP Intelligence resources
		"edgenext_scdn_user_ip":      scdnipresource.ResourceEdgenextScdnUserIp(),
		"edgenext_scdn_user_ip_item": scdnipresource.ResourceEdgenextScdnUserIpItem(),
	}
	mergeProviderResources(resourcesMap, "scdn/domain", domainResources)
	mergeProviderResources(resourcesMap, "scdn/cert", certResources)
	mergeProviderResources(resourcesMap, "scdn/template", templateResources)
	mergeProviderResources(resourcesMap, "scdn/network_speed", networkSpeedResources)
	mergeProviderResources(resourcesMap, "scdn/cache", cacheResources)
	mergeProviderResources(resourcesMap, "scdn/security_protect", securityProtectResources)
	mergeProviderResources(resourcesMap, "scdn/origin_group", originGroupResources)
	mergeProviderResources(resourcesMap, "scdn/cache_operate", cacheOperateResources)
	mergeProviderResources(resourcesMap, "scdn/log_download", logDownloadResources)
	mergeProviderResources(resourcesMap, "sdns/domain", sdnsDomainResources)
	mergeProviderResources(resourcesMap, "sdns/domain_group", sdnsGroupResources)
	mergeProviderResources(resourcesMap, "sdns/record", sdnsRecordResources)

	// Build data sources map (merge with duplicate-key detection per source)
	dataSourcesMap := map[string]*schema.Resource{
		// CDN domain and configuration data sources
		"edgenext_cdn_domain":  cdn.DataSourceEdgenextCdnDomainConfig(),
		"edgenext_cdn_domains": cdn.DataSourceEdgenextCdnDomains(),

		// CDN cache purge data sources
		"edgenext_cdn_purge":  cdn.DataSourceEdgenextCdnPurge(),
		"edgenext_cdn_purges": cdn.DataSourceEdgenextCdnPurges(),

		// CDN file prefetch data sources
		"edgenext_cdn_prefetch":   cdn.DataSourceEdgenextCdnPrefetch(),
		"edgenext_cdn_prefetches": cdn.DataSourceEdgenextCdnPrefetches(),

		// SSL certificate data sources
		"edgenext_ssl_certificate":  ssl.DataSourceEdgenextSslCertificate(),
		"edgenext_ssl_certificates": ssl.DataSourceEdgenextSslCertificates(),

		// OSS bucket management data sources
		"edgenext_oss_buckets": oss.DataSourceOSSBuckets(),
		// OSS object management data sources
		"edgenext_oss_objects": oss.DataSourceOSSObjects(),
		// OSS object management data sources
		"edgenext_oss_object": oss.DataSourceOSSObject(),

		// ECS data sources
		"edgenext_ecs_instances":            ecs.DataSourceENECSInstances(),
		"edgenext_ecs_images":               ecs.DataSourceENECSImages(),
		"edgenext_ecs_key_pairs":            ecs.DataSourceENECSKeyPairs(),
		"edgenext_ecs_vpcs":                 ecs.DataSourceENECSVpcs(),
		"edgenext_ecs_external_gateways":    ecs.DataSourceENECSExternalGateways(),
		"edgenext_ecs_vpc_subnets":          ecs.DataSourceENECSVpcSubnets(),
		"edgenext_ecs_routers":              ecs.DataSourceENECSRouters(),
		"edgenext_ecs_router_ports":         ecs.DataSourceENECSRouterPorts(),
		"edgenext_ecs_network_interfaces":   ecs.DataSourceENECSNetworkInterfaces(),
		"edgenext_ecs_security_groups":      ecs.DataSourceENECSSecurityGroups(),
		"edgenext_ecs_disks":                ecs.DataSourceENECSDisks(),
		"edgenext_ecs_tags":                 ecs.DataSourceENECSTags(),
		"edgenext_ecs_security_group_rules": ecs.DataSourceENECSSecurityGroupRules(),
		"edgenext_ecs_instance_tags":        ecs.DataSourceENECSInstanceTags(),

		// RDS data sources
		"edgenext_rds_instances":                         rds.DataSourceENRDSInstances(),
		"edgenext_rds_databases":                         rds.DataSourceENRDSDatabases(),
		"edgenext_rds_accounts":                          rds.DataSourceENRDSAccounts(),
		"edgenext_rds_backups":                           rds.DataSourceENRDSBackups(),
		"edgenext_rds_backup_policies":                   rds.DataSourceENRDSBackupPolicies(),
		"edgenext_rds_backup_policy_associate_instances": rds.DataSourceENRDSBackupPolicyAssociateInstances(),

		// ELB data sources
		"edgenext_elb_load_balancers":           elb.DataSourceENELBLoadBalancers(),
		"edgenext_elb_certificates":             elb.DataSourceENELBCertificates(),
		"edgenext_elb_listeners":                elb.DataSourceENELBListeners(),
		"edgenext_elb_target_groups":            elb.DataSourceENELBTargetGroups(),
		"edgenext_elb_target_group_attachments": elb.DataSourceENELBTargetGroupAttachments(),
		"edgenext_elb_l7_policies":              elb.DataSourceENELBL7Policies(),
		"edgenext_elb_l7_rules":                 elb.DataSourceENELBL7Rules(),

		// EIP data sources
		"edgenext_eip_floating_ips": eip.DataSourceENEIPFloatingIps(),

		// SCDN domain data sources (from domain module)
		// Note: These data sources are organized under scdn/domain/ for better module management

		// Domain Group data sources
		"edgenext_scdn_domain_groups":        scdndomaingroupdata.DataSourceEdgenextScdnDomainGroups(),
		"edgenext_scdn_domain_group_domains": scdndomaingroupdata.DataSourceEdgenextScdnDomainGroupDomains(),

		// User IP Intelligence data sources
		"edgenext_scdn_user_ips":      scdnipdata.DataSourceEdgenextScdnUserIps(),
		"edgenext_scdn_user_ip_items": scdnipdata.DataSourceEdgenextScdnUserIpItems(),
	}
	mergeProviderDataSources(dataSourcesMap, "scdn/domain", domainDataSources)
	mergeProviderDataSources(dataSourcesMap, "scdn/cert", certDataSources)
	mergeProviderDataSources(dataSourcesMap, "scdn/template", templateDataSources)
	mergeProviderDataSources(dataSourcesMap, "scdn/network_speed", networkSpeedDataSources)
	mergeProviderDataSources(dataSourcesMap, "scdn/cache", cacheDataSources)
	mergeProviderDataSources(dataSourcesMap, "scdn/security_protect", securityProtectDataSources)
	mergeProviderDataSources(dataSourcesMap, "scdn/origin_group", originGroupDataSources)
	mergeProviderDataSources(dataSourcesMap, "scdn/cache_operate", cacheOperateDataSources)
	mergeProviderDataSources(dataSourcesMap, "scdn/log_download", logDownloadDataSources)
	mergeProviderDataSources(dataSourcesMap, "sdns/domain", sdnsDomainDataSources)
	mergeProviderDataSources(dataSourcesMap, "sdns/domain_group", sdnsGroupDataSources)
	mergeProviderDataSources(dataSourcesMap, "sdns/record", sdnsRecordDataSources)

	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			// Unified authentication fields
			"access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("EDGENEXT_ACCESS_KEY", nil),
				Description: "EdgeNext access key for authentication",
				Sensitive:   true,
			},
			"secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("EDGENEXT_SECRET_KEY", nil),
				Description: "EdgeNext secret key for authentication",
				Sensitive:   true,
			},
			"endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("EDGENEXT_ENDPOINT", nil),
				Description: "EdgeNext API endpoint address",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("EDGENEXT_REGION", nil),
				Description: "EdgeNext region",
			},
		},
		ResourcesMap:         resourcesMap,
		DataSourcesMap:       dataSourcesMap,
		ConfigureContextFunc: providerConfigure,
	}
}

// ProviderConfigure configures the provider and returns a client instance
func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Get configuration parameters
	accessKey := d.Get("access_key").(string)
	secretKey := d.Get("secret_key").(string)
	endpoint := d.Get("endpoint").(string)
	region := d.Get("region").(string)

	// Validate that at least access_key and secret_key are provided
	if accessKey == "" || secretKey == "" || endpoint == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Provider configuration validation failed",
			Detail:   "access_key, secret_key and endpoint are required",
		})
		return nil, diags
	}

	// Create config
	config := &connectivity.Config{
		AccessKey: accessKey,
		SecretKey: secretKey,
		Endpoint:  endpoint,
		Region:    region,
	}

	// Create client
	client, err := config.Client()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to create client",
			Detail:   err.Error(),
		})
		return nil, diags
	}

	return client, diags
}

// GetClient gets the client instance from provider configuration
func GetClient(meta interface{}) (*connectivity.EdgeNextClient, error) {
	client, ok := meta.(*connectivity.EdgeNextClient)
	if !ok {
		return nil, fmt.Errorf("invalid client type: %T", meta)
	}
	return client, nil
}

// IsNotFoundError checks if it's a "not found" error
func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	errMsg := strings.ToLower(err.Error())
	notFoundKeywords := []string{
		"not found", "notfound", "404", "does not exist", "not found",
		"domain not found", "certificate not found",
	}

	for _, keyword := range notFoundKeywords {
		if strings.Contains(errMsg, keyword) {
			return true
		}
	}

	return false
}

// IsRateLimitError checks if it's a rate limit error
func IsRateLimitError(err error) bool {
	if err == nil {
		return false
	}

	errMsg := strings.ToLower(err.Error())
	rateLimitKeywords := []string{
		"rate limit", "ratelimit", "too many requests", "429",
		"requests too frequent", "rate limited", "frequency limit",
	}

	for _, keyword := range rateLimitKeywords {
		if strings.Contains(errMsg, keyword) {
			return true
		}
	}

	return false
}

// IsAuthenticationError checks if it's an authentication error
func IsAuthenticationError(err error) bool {
	if err == nil {
		return false
	}

	errMsg := strings.ToLower(err.Error())
	authKeywords := []string{
		"unauthorized", "401", "forbidden", "403",
		"invalid api key", "invalid secret", "authentication failed",
		"unauthorized access", "authentication error", "invalid credentials", "invalid key",
	}

	for _, keyword := range authKeywords {
		if strings.Contains(errMsg, keyword) {
			return true
		}
	}

	return false
}

// FormatError formats error messages
func FormatError(operation string, err error) string {
	if err == nil {
		return fmt.Sprintf("%s succeeded", operation)
	}

	// Return different error messages based on error type
	if IsNotFoundError(err) {
		return fmt.Sprintf("%s failed: resource not found", operation)
	}

	if IsRateLimitError(err) {
		return fmt.Sprintf("%s failed: requests too frequent, please retry later", operation)
	}

	if IsAuthenticationError(err) {
		return fmt.Sprintf("%s failed: authentication failed, please check API key and secret", operation)
	}

	return fmt.Sprintf("%s failed: %v", operation, err)
}
