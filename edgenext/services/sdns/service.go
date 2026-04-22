package sdns

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
)

// SdnsService DNS service
type SdnsService struct {
	client *connectivity.EdgeNextClient
}

// NewSdnsService creates a new DNS service instance
func NewSdnsService(client *connectivity.EdgeNextClient) *SdnsService {
	return &SdnsService{
		client: client,
	}
}

// callAPI is a helper function to make DNS API calls
func (s *SdnsService) callAPI(ctx context.Context, method, endpoint string, reqData interface{}, responseData interface{}) error {
	// Get SCDN client from EdgeNextClient (DNS uses the same client/auth mechanism)
	dnsClient, err := s.client.ScdnClient()
	if err != nil {
		return fmt.Errorf("failed to get DNS client: %w", err)
	}

	// Convert request to ScdnRequest format
	dnsReq := &connectivity.ScdnRequest{
		Data:  make(map[string]interface{}),
		Query: make(map[string]interface{}),
	}

	if reqData != nil {
		reqBytes, err := json.Marshal(reqData)
		if err != nil {
			return fmt.Errorf("failed to marshal request: %w", err)
		}

		var wrapper map[string]interface{}
		if err := json.Unmarshal(reqBytes, &wrapper); err != nil {
			return fmt.Errorf("failed to unmarshal request data: %w", err)
		}

		if method == "GET" {
			dnsReq.Query = wrapper
		} else {
			dnsReq.Data = wrapper
		}
	}

	// Call API
	var dnsResp *connectivity.ScdnResponse

	switch method {
	case "GET":
		dnsResp, err = dnsClient.Get(ctx, endpoint, dnsReq)
	case "POST":
		dnsResp, err = dnsClient.Post(ctx, endpoint, dnsReq)
	case "PUT":
		dnsResp, err = dnsClient.Put(ctx, endpoint, dnsReq)
	case "DELETE":
		dnsResp, err = dnsClient.Delete(ctx, endpoint, dnsReq)
	default:
		return fmt.Errorf("unsupported HTTP method: %s", method)
	}

	if err != nil {
		return err
	}

	// Check business status code
	if dnsResp.Status.Code != 1 {
		return fmt.Errorf("API error: %s (code: %d)", dnsResp.Status.Message, dnsResp.Status.Code)
	}

	// Convert response
	if dnsResp != nil && responseData != nil {
		// Use JSON roundtrip to populate responseData including Status and Data
		fullResp := map[string]interface{}{
			"status": map[string]interface{}{
				"code":    dnsResp.Status.Code,
				"message": dnsResp.Status.Message,
			},
			"data": dnsResp.Data,
		}

		fullRespBytes, err := json.Marshal(fullResp)
		if err != nil {
			return fmt.Errorf("failed to marshal full response: %w", err)
		}

		if err := json.Unmarshal(fullRespBytes, responseData); err != nil {
			return fmt.Errorf("failed to unmarshal into response struct: %w", err)
		}
	}

	return nil
}
