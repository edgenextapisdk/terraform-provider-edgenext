package scdn

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
)

// ScdnService SCDN domain service
type ScdnService struct {
	client *connectivity.EdgeNextClient
}

// NewScdnService creates a new SCDN domain service instance
func NewScdnService(client *connectivity.EdgeNextClient) *ScdnService {
	return &ScdnService{
		client: client,
	}
}

// callSCDNAPI is a helper function to make SCDN API calls
func (s *ScdnService) callSCDNAPI(ctx context.Context, method, endpoint string, reqData interface{}, responseData interface{}) error {
	// Get SCDN client from EdgeNextClient
	scdnClient, err := s.client.ScdnClient()
	if err != nil {
		return fmt.Errorf("failed to get SCDN client: %w", err)
	}

	// Convert request to ScdnRequest format
	scdnReq := &connectivity.ScdnRequest{}

	if reqData != nil {
		scdnReq.Data = make(map[string]interface{})
		reqBytes, err := json.Marshal(reqData)
		if err != nil {
			return fmt.Errorf("failed to marshal request: %w", err)
		}

		if err := json.Unmarshal(reqBytes, &scdnReq.Data); err != nil {
			return fmt.Errorf("failed to unmarshal request data: %w", err)
		}
	}

	// Call SCDN API
	var scdnResp *connectivity.ScdnResponse

	switch method {
	case MethodGET:
		scdnResp, err = scdnClient.Get(ctx, endpoint, scdnReq)
	case MethodPOST:
		scdnResp, err = scdnClient.Post(ctx, endpoint, scdnReq)
	case MethodPUT:
		scdnResp, err = scdnClient.Put(ctx, endpoint, scdnReq)
	case MethodDELETE:
		scdnResp, err = scdnClient.Delete(ctx, endpoint, scdnReq)
	default:
		return fmt.Errorf("unsupported HTTP method: %s", method)
	}

	if err != nil {
		return err
	}

	// Check business status code
	if scdnResp.Status.Code != 1 {
		return fmt.Errorf("API error: %s (code: %d)", scdnResp.Status.Message, scdnResp.Status.Code)
	}

	// Convert response
	if scdnResp != nil && responseData != nil {
		dataBytes, err := json.Marshal(scdnResp)
		if err != nil {
			return fmt.Errorf("failed to marshal response data: %w", err)
		}

		if err := json.Unmarshal(dataBytes, responseData); err != nil {
			return fmt.Errorf("failed to unmarshal response data: %w", err)
		}
	}

	return nil
}

// requestSCDNAPI is a helper function to make SCDN API calls
func (s *ScdnService) requestSCDNAPI(ctx context.Context, method, endpoint string, scdnReq *connectivity.ScdnRequest, responseData interface{}) error {
	// Get SCDN client from EdgeNextClient
	scdnClient, err := s.client.ScdnClient()
	if err != nil {
		return fmt.Errorf("failed to get SCDN client: %w", err)
	}

	// Call SCDN API
	var scdnResp *connectivity.ScdnResponse

	switch method {
	case MethodGET:
		scdnResp, err = scdnClient.Get(ctx, endpoint, scdnReq)
	case MethodPOST:
		scdnResp, err = scdnClient.Post(ctx, endpoint, scdnReq)
	case MethodPUT:
		scdnResp, err = scdnClient.Put(ctx, endpoint, scdnReq)
	case MethodDELETE:
		scdnResp, err = scdnClient.Delete(ctx, endpoint, scdnReq)
	default:
		return fmt.Errorf("unsupported HTTP method: %s", method)
	}

	if err != nil {
		return err
	}

	// Check business status code
	if scdnResp.Status.Code != 1 {
		return fmt.Errorf("API error: %s (code: %d)", scdnResp.Status.Message, scdnResp.Status.Code)
	}

	// Convert response
	if scdnResp != nil && responseData != nil {
		dataBytes, err := json.Marshal(scdnResp)
		if err != nil {
			return fmt.Errorf("failed to marshal response data: %w", err)
		}

		if err := json.Unmarshal(dataBytes, responseData); err != nil {
			return fmt.Errorf("failed to unmarshal response data: %w", err)
		}
	}

	return nil
}
