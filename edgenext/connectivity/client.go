package connectivity

import (
	"fmt"
	"strings"
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
	ecsClient  *ECSClient  // For ECS
	rdsClient  *RDSClient  // For RDS (same HTTP client stack as ECS)
	elbClient  *ELBClient  // For ELB (same HTTP client stack as ECS)

	// Use sync.Once to ensure clients are initialized only once
	apiClientOnce  sync.Once
	ossClientOnce  sync.Once
	scdnClientOnce sync.Once
	ecsClientOnce  sync.Once
	rdsClientOnce  sync.Once
	elbClientOnce  sync.Once

	// Store initialization errors
	apiClientErr  error
	ossClientErr  error
	scdnClientErr error
	ecsClientErr  error
	rdsClientErr  error
	elbClientErr  error
}

// Client returns the EdgeNext client
func (c *Config) Client() (*EdgeNextClient, error) {
	client := &EdgeNextClient{
		config: c,
	}

	return client, nil
}

// validateServiceClientConfig checks credentials and endpoint required for signed HTTP service clients.
func validateServiceClientConfig(accessKey, secretKey, endpoint string) error {
	if strings.TrimSpace(accessKey) == "" {
		return fmt.Errorf("access_key is required")
	}
	if strings.TrimSpace(secretKey) == "" {
		return fmt.Errorf("secret_key is required")
	}
	if strings.TrimSpace(endpoint) == "" {
		return fmt.Errorf("endpoint is required")
	}
	return nil
}

// APIClient returns or initializes the API client
func (c *EdgeNextClient) APIClient() (*APIClient, error) {
	c.apiClientOnce.Do(func() {
		if err := validateServiceClientConfig(c.config.AccessKey, c.config.SecretKey, c.config.Endpoint); err != nil {
			c.apiClientErr = fmt.Errorf("failed to create API client: %w", err)
			return
		}

		client := NewAPIClient(c.config.AccessKey, c.config.SecretKey, c.config.Endpoint)
		if client == nil {
			c.apiClientErr = fmt.Errorf("failed to create API client: client constructor returned nil")
			return
		}
		c.apiClient = client
	})

	if c.apiClientErr != nil {
		return nil, c.apiClientErr
	}
	return c.apiClient, nil
}

// OSSClient returns or initializes the OSS S3 client
func (c *EdgeNextClient) OSSClient() (*OSSClient, error) {
	c.ossClientOnce.Do(func() {
		if err := validateServiceClientConfig(c.config.AccessKey, c.config.SecretKey, c.config.Endpoint); err != nil {
			c.ossClientErr = fmt.Errorf("failed to create OSS client: %w", err)
			return
		}

		client, err := NewOSSClient(c.config.AccessKey, c.config.SecretKey, c.config.Endpoint, c.config.Region)
		if err != nil {
			c.ossClientErr = fmt.Errorf("failed to create OSS client: %w", err)
			return
		}
		if client == nil {
			c.ossClientErr = fmt.Errorf("failed to create OSS client: client constructor returned nil")
			return
		}

		c.ossClient = client
	})

	if c.ossClientErr != nil {
		return nil, c.ossClientErr
	}
	return c.ossClient, nil
}

// ScdnClient returns or initializes the SCDN API client
func (c *EdgeNextClient) ScdnClient() (*ScdnClient, error) {
	c.scdnClientOnce.Do(func() {
		if err := validateServiceClientConfig(c.config.AccessKey, c.config.SecretKey, c.config.Endpoint); err != nil {
			c.scdnClientErr = fmt.Errorf("failed to create SCDN client: %w", err)
			return
		}

		scdnEndpoint := c.config.Endpoint
		client := NewScdnClient(scdnEndpoint, c.config.AccessKey, c.config.SecretKey, 30*time.Second)
		if client == nil {
			c.scdnClientErr = fmt.Errorf("failed to create SCDN client: client constructor returned nil")
			return
		}
		c.scdnClient = client
	})

	if c.scdnClientErr != nil {
		return nil, c.scdnClientErr
	}
	return c.scdnClient, nil
}

// ECSClient returns or initializes the ECS API client
func (c *EdgeNextClient) ECSClient() (*ECSClient, error) {
	c.ecsClientOnce.Do(func() {
		if err := validateServiceClientConfig(c.config.AccessKey, c.config.SecretKey, c.config.Endpoint); err != nil {
			c.ecsClientErr = fmt.Errorf("failed to create ECS client: %w", err)
			return
		}

		client := NewECSClient(c.config.AccessKey, c.config.SecretKey, c.config.Endpoint, c.config.Region)
		if client == nil {
			c.ecsClientErr = fmt.Errorf("failed to create ECS client: client constructor returned nil")
			return
		}
		c.ecsClient = client
	})

	if c.ecsClientErr != nil {
		return nil, c.ecsClientErr
	}
	return c.ecsClient, nil
}

// RDSClient returns or initializes the RDS API client (same signing and transport as ECSClient).
func (c *EdgeNextClient) RDSClient() (*RDSClient, error) {
	c.rdsClientOnce.Do(func() {
		if err := validateServiceClientConfig(c.config.AccessKey, c.config.SecretKey, c.config.Endpoint); err != nil {
			c.rdsClientErr = fmt.Errorf("failed to create RDS client: %w", err)
			return
		}

		client := NewRDSClient(c.config.AccessKey, c.config.SecretKey, c.config.Endpoint, c.config.Region)
		if client == nil {
			c.rdsClientErr = fmt.Errorf("failed to create RDS client: client constructor returned nil")
			return
		}
		c.rdsClient = client
	})

	if c.rdsClientErr != nil {
		return nil, c.rdsClientErr
	}
	return c.rdsClient, nil
}

// ELBClient returns or initializes the ELB API client (same signing and transport as ECSClient).
func (c *EdgeNextClient) ELBClient() (*ELBClient, error) {
	c.elbClientOnce.Do(func() {
		if err := validateServiceClientConfig(c.config.AccessKey, c.config.SecretKey, c.config.Endpoint); err != nil {
			c.elbClientErr = fmt.Errorf("failed to create ELB client: %w", err)
			return
		}

		client := NewELBClient(c.config.AccessKey, c.config.SecretKey, c.config.Endpoint, c.config.Region)
		if client == nil {
			c.elbClientErr = fmt.Errorf("failed to create ELB client: client constructor returned nil")
			return
		}
		c.elbClient = client
	})

	if c.elbClientErr != nil {
		return nil, c.elbClientErr
	}
	return c.elbClient, nil
}
