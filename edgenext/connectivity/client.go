package connectivity

import (
	"fmt"
	"sync"
	"time"
)

// Config contains all the configuration for EdgeNext provider
type Config struct {
	AccessKey string
	SecretKey string
	Endpoint  string
	Region    string
}

// EdgeNextClient is the main client struct that holds all service clients
type EdgeNextClient struct {
	config     *Config
	apiClient  *APIClient  // For CDN/SSL
	ossClient  *OSSClient  // For OSS
	scdnClient *ScdnClient // For SCDN

	// Use sync.Once to ensure clients are initialized only once
	apiClientOnce  sync.Once
	ossClientOnce  sync.Once
	scdnClientOnce sync.Once

	// Store initialization errors
	apiClientErr  error
	ossClientErr  error
	scdnClientErr error
}

// Client returns the EdgeNext client
func (c *Config) Client() (*EdgeNextClient, error) {
	client := &EdgeNextClient{
		config: c,
	}

	return client, nil
}

// APIClient returns or initializes the API client
func (c *EdgeNextClient) APIClient() (*APIClient, error) {
	c.apiClientOnce.Do(func() {
		c.apiClient = NewAPIClient(c.config.AccessKey, c.config.SecretKey, c.config.Endpoint)
	})

	return c.apiClient, c.apiClientErr
}

// OSSClient returns or initializes the OSS S3 client
func (c *EdgeNextClient) OSSClient() (*OSSClient, error) {
	c.ossClientOnce.Do(func() {
		client, err := NewOSSClient(c.config.AccessKey, c.config.SecretKey, c.config.Endpoint, c.config.Region)
		if err != nil {
			c.ossClientErr = fmt.Errorf("failed to create OSS client: %w", err)
			return
		}

		c.ossClient = client
	})

	return c.ossClient, c.ossClientErr
}

// ScdnClient returns or initializes the SCDN API client
func (c *EdgeNextClient) ScdnClient() (*ScdnClient, error) {
	c.scdnClientOnce.Do(func() {
		// Use the same endpoint as API client but with SCDN-specific path
		scdnEndpoint := c.config.Endpoint
		if scdnEndpoint == "" {
			scdnEndpoint = "https://api.edgenextscdn.com"
		}

		c.scdnClient = NewScdnClient(scdnEndpoint, c.config.AccessKey, c.config.SecretKey, 30*time.Second)
	})

	return c.scdnClient, c.scdnClientErr
}
