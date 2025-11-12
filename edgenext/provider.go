package edgenext

import (
	"context"
	"fmt"
	"strings"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/cdn"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/oss"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/ssl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider returns the EdgeNext CDN Terraform Provider
func Provider() *schema.Provider {
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
		ResourcesMap: map[string]*schema.Resource{
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
		},
		DataSourcesMap: map[string]*schema.Resource{

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
		},
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
