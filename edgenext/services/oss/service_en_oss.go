package oss

import (
	"context"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Service provides OSS service operations
type Service struct {
	client *connectivity.EdgeNextClient
}

// NewService creates a new OSS service
func NewService(client *connectivity.EdgeNextClient) *Service {
	return &Service{
		client: client,
	}
}

// GetOSSClient returns the OSS client with error handling
func (s *Service) GetOSSClient() (*connectivity.OSSClient, diag.Diagnostics) {
	ossClient, err := s.client.OSSClient()
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return ossClient, nil
}

// ValidateOSSConfig validates OSS configuration
func ValidateOSSConfig(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	_, err := client.OSSClient()
	if err != nil {
		return diag.Errorf("OSS client initialization failed: %s", err)
	}
	return nil
}
