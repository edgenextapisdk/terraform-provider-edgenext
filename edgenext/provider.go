package edgenext

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/cdn"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/ssl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider returns the EdgeNext CDN Terraform Provider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "EdgeNext API key for authentication",
				Sensitive:   true,
			},
			"secret": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "EdgeNext secret for authentication",
				Sensitive:   true,
			},
			"endpoint": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "EdgeNext API endpoint address",
			},
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     300,
				Description: "API request timeout in seconds",
			},
			"retry_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     3,
				Description: "API request retry count",
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
		},
		ConfigureContextFunc: providerConfigure,
	}
}

// ProviderConfigure configures the provider and returns a client instance
func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Get configuration parameters
	apiKey := d.Get("api_key").(string)
	secret := d.Get("secret").(string)
	endpoint := d.Get("endpoint").(string)
	timeout := d.Get("timeout").(int)
	retryCount := d.Get("retry_count").(int)

	// Validate configuration parameters
	if err := validateProviderConfig(apiKey, secret, endpoint); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Provider configuration validation failed",
			Detail:   err.Error(),
		})
		return nil, diags
	}

	// Create client instance
	client, err := createClient(apiKey, secret, endpoint, timeout, retryCount)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to create client",
			Detail:   err.Error(),
		})
		return nil, diags
	}

	// Test connection
	// if err := testConnection(ctx, client); err != nil {
	// 	diags = append(diags, diag.Diagnostic{
	// 		Severity: diag.Warning,
	// 		Summary:  "Connection test failed",
	// 		Detail:   fmt.Sprintf("Unable to connect to EdgeNext API: %v", err),
	// 	})
	// 	// Connection test failure doesn't block provider from continuing, only shows warning
	// }

	return client, diags
}

// validateProviderConfig validates provider configuration parameters
func validateProviderConfig(apiKey, secret, endpoint string) error {
	// Validate API key
	if strings.TrimSpace(apiKey) == "" {
		return fmt.Errorf("API key cannot be empty")
	}
	if len(apiKey) < 8 {
		return fmt.Errorf("API key length cannot be less than 8 characters")
	}

	// Validate secret
	if strings.TrimSpace(secret) == "" {
		return fmt.Errorf("secret cannot be empty")
	}
	if len(secret) < 8 {
		return fmt.Errorf("secret length cannot be less than 8 characters")
	}

	// Validate endpoint
	if strings.TrimSpace(endpoint) == "" {
		return fmt.Errorf("API endpoint cannot be empty")
	}
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		return fmt.Errorf("API endpoint must start with http:// or https://")
	}

	return nil
}

// createClient creates an EdgeNext client instance
func createClient(apiKey, secret, endpoint string, timeout, retryCount int) (*connectivity.Client, error) {
	// Create base client
	client := connectivity.NewClient(apiKey, secret, endpoint)

	// Configure timeout settings
	if timeout > 0 {
		client.SetTimeout(time.Duration(timeout) * time.Second)
	}

	// Configure retry settings
	if retryCount > 0 {
		client.SetRetryCount(retryCount)
	}

	// Configure retry wait time
	client.SetRetryWaitTime(1 * time.Second)
	client.SetRetryMaxWaitTime(10 * time.Second)

	return client, nil
}

// testConnection tests the connection to EdgeNext API
func testConnection(ctx context.Context, client *connectivity.Client) error {
	// Try to send a simple health check request
	// This can be adjusted based on the actual API endpoint
	var response interface{}
	err := client.Get(ctx, "/health", &response)
	if err != nil {
		return fmt.Errorf("connection test failed: %w", err)
	}
	return nil
}

// GetClient gets the client instance from provider configuration
func GetClient(meta interface{}) (*connectivity.Client, error) {
	client, ok := meta.(*connectivity.Client)
	if !ok {
		return nil, fmt.Errorf("invalid client type: %T", meta)
	}
	return client, nil
}

// GetClientWithContext gets the client instance from provider configuration (with context)
func GetClientWithContext(ctx context.Context, meta interface{}) (*connectivity.Client, error) {
	client, err := GetClient(meta)
	if err != nil {
		return nil, err
	}

	// Check if context is cancelled
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return client, nil
	}
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
