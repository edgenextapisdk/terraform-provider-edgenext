package ssl

import (
	"context"
	"fmt"
	"strconv"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
)

// SslCertificateService SSL certificate service
type SslCertificateService struct {
	client *connectivity.EdgeNextClient
}

// NewSslCertificateService creates a new SSL certificate service instance
func NewSslCertificateService(client *connectivity.EdgeNextClient) *SslCertificateService {
	return &SslCertificateService{client: client}
}

// SslCertificateRequest SSL certificate request
type SslCertificateRequest struct {
	Certificate string `json:"certificate"` // Certificate content
	Key         string `json:"key"`         // Private key content
	Name        string `json:"name"`        // Certificate name
	CertID      *int   `json:"cert_id"`     // Certificate ID (for modification)
}

type DeleteSslCertificateRequest struct {
	CertID int `json:"cert_id"`
}

// SslCertificateResponse SSL certificate response
type SslCertificateResponse struct {
	Code int                `json:"code"`
	Data SslCertificateData `json:"data"`
	Msg  string             `json:"msg,omitempty"`
}

// SslCertificateData SSL certificate data
type SslCertificateData struct {
	CertID         string   `json:"cert_id"`
	Name           string   `json:"name"`             // Certificate name
	Certificate    string   `json:"certificate"`      // Certificate content
	Key            string   `json:"key"`              // Private key content
	BindDomains    []string `json:"bind_domains"`     // Bound domains
	CertStartTime  string   `json:"cert_start_time"`  // Certificate start time
	CertExpireTime string   `json:"cert_expire_time"` // Certificate end time
}

type SslCertificateDataV2 struct {
	CertID            string   `json:"cert_id"`
	Name              string   `json:"name"`               // Certificate name
	Certificate       string   `json:"certificate"`        // Certificate content
	Key               string   `json:"key"`                // Private key content
	AssociatedDomains []string `json:"associated_domains"` // Associated domains
	IncludeDomains    []string `json:"include_domains"`    // Included domains
	CertStartTime     string   `json:"cert_start_time"`    // Certificate start time
	CertExpireTime    string   `json:"cert_expire_time"`   // Certificate end time
}

// SslCertificateListResponse SSL certificate list response
type SslCertificateListResponse struct {
	Code int                            `json:"code"`
	Data SslCertificateListResponseData `json:"data"`
	Msg  string                         `json:"msg,omitempty"`
}

type SslCertificateListResponseData struct {
	List        []SslCertificateDataV2 `json:"list"`
	TotalNumber int                    `json:"total_number"`
	PageNumber  int                    `json:"page_number"`
	PageSize    int                    `json:"page_size"`
}

// SslCertificateDeleteResponse SSL certificate delete response
type SslCertificateDeleteResponse struct {
	Code int    `json:"code"`
	Data string `json:"data"`
}

// CreateOrUpdateSslCertificate creates or updates SSL certificate
func (s *SslCertificateService) CreateOrUpdateSslCertificate(req SslCertificateRequest) (*SslCertificateResponse, error) {
	ctx := context.Background()

	var response SslCertificateResponse
	apiClient, err := s.client.APIClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create or update SSL certificate: %w", err)
	}
	err = apiClient.Post(ctx, "/v2/domain/certificate", req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to create or update SSL certificate: %w", err)
	}

	// Check API response status code, 0 indicates success
	if response.Code != 0 {
		return nil, fmt.Errorf("failed to create or update SSL certificate: %s (code: %d)", response.Msg, response.Code)
	}

	return &response, nil
}

// GetSslCertificate queries SSL certificate
func (s *SslCertificateService) GetSslCertificate(certID int) (*SslCertificateResponse, error) {
	ctx := context.Background()

	query := map[string]string{
		"cert_id": strconv.Itoa(certID),
	}

	var response SslCertificateResponse
	apiClient, err := s.client.APIClient()
	if err != nil {
		return nil, fmt.Errorf("failed to query SSL certificate: %w", err)
	}
	err = apiClient.GetWithQuery(ctx, "/v2/domain/certificate", query, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to query SSL certificate: %w", err)
	}

	// Check API response status code, 0 indicates success
	if response.Code != 0 {
		return nil, fmt.Errorf("failed to query SSL certificate: %s (code: %d)", response.Msg, response.Code)
	}

	return &response, nil
}

// ListSslCertificates queries SSL certificate list
func (s *SslCertificateService) ListSslCertificates(pageNumber int, pageSize int) (*SslCertificateListResponse, error) {
	ctx := context.Background()

	query := map[string]string{
		"page_number": strconv.Itoa(pageNumber),
		"page_size":   strconv.Itoa(pageSize),
	}

	var response SslCertificateListResponse
	apiClient, err := s.client.APIClient()
	if err != nil {
		return nil, fmt.Errorf("failed to query SSL certificate list: %w", err)
	}
	err = apiClient.GetWithQuery(ctx, "/v2/domain/certificate", query, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to query SSL certificate list: %w", err)
	}

	// Check API response status code, 0 indicates success
	if response.Code != 0 {
		return nil, fmt.Errorf("failed to query SSL certificate list: %s (code: %d)", response.Msg, response.Code)
	}

	return &response, nil
}

// DeleteSslCertificate deletes SSL certificate
func (s *SslCertificateService) DeleteSslCertificate(req DeleteSslCertificateRequest) error {
	ctx := context.Background()
	var response SslCertificateDeleteResponse
	apiClient, err := s.client.APIClient()
	if err != nil {
		return fmt.Errorf("failed to delete SSL certificate: %w", err)
	}
	err = apiClient.DeleteWithBodyAndResult(ctx, "/v2/domain/certificate", req, &response)
	if err != nil {
		return fmt.Errorf("failed to delete SSL certificate: %w", err)
	}

	if response.Code != 0 {
		return fmt.Errorf("failed to delete SSL certificate: %s (code: %d)", response.Data, response.Code)
	}

	return nil
}
