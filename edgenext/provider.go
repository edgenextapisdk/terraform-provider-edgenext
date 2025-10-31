package edgenext

import (
	"context"
	"fmt"
	"strings"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/cdn"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/oss"
	scdncache "github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn/cache"
	scdncacheoperate "github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn/cache_operate"
	scdncert "github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn/cert"
	scdndomain "github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn/domain"
	scdnlogdownload "github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn/log_download"
	scdnnetworkspeed "github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn/network_speed"
	scdnorigingroup "github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn/origin_group"
	scdnsecurityprotect "github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn/security_protect"
	scdntemplate "github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn/template"
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

	// Build resources map
	ResourcesMap := map[string]*schema.Resource{
		// CDN domain and configuration management resources
		"edgenext_cdn_domain": cdn.ResourceEdgenextCdnDomainConfig(),

		// CDN cache purge and file push resources
		"edgenext_cdn_push":  cdn.ResourceEdgenextCdnPush(),
		"edgenext_cdn_purge": cdn.ResourceEdgenextCdnPurge(),

		// SSL certificate management resources
		"edgenext_ssl_certificate": ssl.ResourceEdgenextSslCertificate(),

		// OSS object storage resources
		"edgenext_oss_bucket": oss.ResourceOSSBucket(),
		// OSS object management resources
		"edgenext_oss_object":      oss.ResourceOSSObject(),
		"edgenext_oss_object_copy": oss.ResourceOSSObjectCopy(),

		// SCDN domain management resources (from domain module)
		// Note: These resources are organized under scdn/domain/ for better module management
	}

	// Add domain module resources dynamically
	for k, v := range domainResources {
		ResourcesMap[k] = v
	}

	// Add cert module resources dynamically
	for k, v := range certResources {
		ResourcesMap[k] = v
	}

	// Add template module resources dynamically
	for k, v := range templateResources {
		ResourcesMap[k] = v
	}

	// Add network speed module resources dynamically
	for k, v := range networkSpeedResources {
		ResourcesMap[k] = v
	}

	// Add cache module resources dynamically
	for k, v := range cacheResources {
		ResourcesMap[k] = v
	}

	// Add security protection module resources dynamically
	for k, v := range securityProtectResources {
		ResourcesMap[k] = v
	}

	// Add origin group module resources dynamically
	for k, v := range originGroupResources {
		ResourcesMap[k] = v
	}

	// Add cache operate module resources dynamically
	for k, v := range cacheOperateResources {
		ResourcesMap[k] = v
	}

	// Add log download module resources dynamically
	for k, v := range logDownloadResources {
		ResourcesMap[k] = v
	}

	// Build data sources map
	DataSourcesMap := map[string]*schema.Resource{
		// CDN domain and configuration data sources
		"edgenext_cdn_domain":  cdn.DataSourceEdgenextCdnDomainConfig(),
		"edgenext_cdn_domains": cdn.DataSourceEdgenextCdnDomains(),

		// CDN cache push data sources
		"edgenext_cdn_push":   cdn.DataSourceEdgenextCdnPush(),
		"edgenext_cdn_pushes": cdn.DataSourceEdgenextCdnPushes(),

		// CDN file purge data sources
		"edgenext_cdn_purge":  cdn.DataSourceEdgenextCdnPurge(),
		"edgenext_cdn_purges": cdn.DataSourceEdgenextCdnPurges(),

		// SSL certificate data sources
		"edgenext_ssl_certificate":  ssl.DataSourceEdgenextSslCertificate(),
		"edgenext_ssl_certificates": ssl.DataSourceEdgenextSslCertificates(),

		// OSS bucket management data sources
		"edgenext_oss_buckets": oss.DataSourceOSSBuckets(),
		// OSS object management data sources
		"edgenext_oss_objects": oss.DataSourceOSSObjects(),
		// OSS object management data sources
		"edgenext_oss_object": oss.DataSourceOSSObject(),

		// SCDN domain data sources (from domain module)
		// Note: These data sources are organized under scdn/domain/ for better module management
	}

	// Add domain module data sources dynamically
	for k, v := range domainDataSources {
		DataSourcesMap[k] = v
	}

	// Add cert module data sources dynamically
	for k, v := range certDataSources {
		DataSourcesMap[k] = v
	}

	// Add template module data sources dynamically
	for k, v := range templateDataSources {
		DataSourcesMap[k] = v
	}

	// Add network speed module data sources dynamically
	for k, v := range networkSpeedDataSources {
		DataSourcesMap[k] = v
	}

	// Add cache module data sources dynamically
	for k, v := range cacheDataSources {
		DataSourcesMap[k] = v
	}

	// Add security protection module data sources dynamically
	for k, v := range securityProtectDataSources {
		DataSourcesMap[k] = v
	}

	// Add origin group module data sources dynamically
	for k, v := range originGroupDataSources {
		DataSourcesMap[k] = v
	}

	// Add cache operate module data sources dynamically
	for k, v := range cacheOperateDataSources {
		DataSourcesMap[k] = v
	}

	// Add log download module data sources dynamically
	for k, v := range logDownloadDataSources {
		DataSourcesMap[k] = v
	}

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
		ResourcesMap:         ResourcesMap,
		DataSourcesMap:       DataSourcesMap,
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
