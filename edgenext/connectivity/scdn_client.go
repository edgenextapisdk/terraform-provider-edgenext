package connectivity

import (
	"context"
	"fmt"
	"time"

	edgenext "github.com/edgenextapisdk/edgenext-go"
)

// ScdnClient SCDN API client for domain_v5 interfaces
type ScdnClient struct {
	sdk       *edgenext.Sdk
	baseURL   string
	apiKey    string
	apiSecret string
	timeout   time.Duration
}

// ScdnRequest represents a request to SCDN API
type ScdnRequest struct {
	Data   map[string]interface{} `json:"data,omitempty"`
	Query  map[string]interface{} `json:"query,omitempty"`
	Header map[string]string      `json:"header,omitempty"`
}

// ScdnResponse represents a response from SCDN API
type ScdnResponse struct {
	Status ScdnStatus  `json:"status"`
	Data   interface{} `json:"data,omitempty"`
}

// ScdnStatus represents the status in SCDN API response
type ScdnStatus struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewScdnClient creates a new SCDN API client using official edgenext-go SDK
func NewScdnClient(baseURL, apiKey, apiSecret string, timeout time.Duration) *ScdnClient {
	// Initialize the official SDK
	sdk := &edgenext.Sdk{
		AppId:     apiKey,
		AppSecret: apiSecret,
		ApiPre:    baseURL,
		Timeout:   int(timeout.Seconds()),
	}

	return &ScdnClient{
		sdk:       sdk,
		baseURL:   baseURL,
		apiKey:    apiKey,
		apiSecret: apiSecret,
		timeout:   timeout,
	}
}

// NewScdnClientFromConfig creates a new SCDN API client from EdgeNextClient config
// Deprecated: Use EdgeNextClient.ScdnClient() instead
func NewScdnClientFromConfig(client *EdgeNextClient) *ScdnClient {
	// This function is deprecated and should not be used
	// Use EdgeNextClient.ScdnClient() instead for proper integration
	return NewScdnClient("https://api.edgenextscdn.com", "", "", 30*time.Second)
}

// Get performs a GET request to SCDN API
func (c *ScdnClient) Get(ctx context.Context, api string, req *ScdnRequest) (*ScdnResponse, error) {
	return c.doRequest(ctx, "GET", api, req)
}

// Post performs a POST request to SCDN API
func (c *ScdnClient) Post(ctx context.Context, api string, req *ScdnRequest) (*ScdnResponse, error) {
	return c.doRequest(ctx, "POST", api, req)
}

// Put performs a PUT request to SCDN API
func (c *ScdnClient) Put(ctx context.Context, api string, req *ScdnRequest) (*ScdnResponse, error) {
	return c.doRequest(ctx, "PUT", api, req)
}

// Delete performs a DELETE request to SCDN API
func (c *ScdnClient) Delete(ctx context.Context, api string, req *ScdnRequest) (*ScdnResponse, error) {
	return c.doRequest(ctx, "DELETE", api, req)
}

// doRequest performs the actual HTTP request using official SDK
func (c *ScdnClient) doRequest(ctx context.Context, method, api string, req *ScdnRequest) (*ScdnResponse, error) {
	// Convert ScdnRequest to SDK ReqParams
	reqParams := edgenext.ReqParams{
		Data:    req.Data,
		Query:   req.Query,
		Headers: req.Header,
	}

	// Add language header if not provided
	if reqParams.Headers == nil {
		reqParams.Headers = make(map[string]string)
	}
	if _, exists := reqParams.Headers["X-Lang"]; !exists {
		reqParams.Headers["X-Lang"] = "en"
	}

	var resp *edgenext.Response
	var err error

	// Call appropriate SDK method based on HTTP method
	switch method {
	case "GET":
		resp, err = c.sdk.Get(api, reqParams)
	case "POST":
		resp, err = c.sdk.Post(api, reqParams)
	case "PUT":
		resp, err = c.sdk.Put(api, reqParams)
	case "DELETE":
		resp, err = c.sdk.Delete(api, reqParams)
	default:
		return nil, fmt.Errorf("unsupported HTTP method: %s", method)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}

	// Check business status code
	if resp.BizCode != 1 {
		return nil, fmt.Errorf("API error: %s (code: %d)", resp.BizMsg, resp.BizCode)
	}

	// Parse response data
	var scdnResp ScdnResponse
	scdnResp.Status.Code = resp.BizCode
	scdnResp.Status.Message = resp.BizMsg

	// Parse business data if available
	if resp.BizData != nil {
		// Convert interface{} to map[string]interface{} for proper JSON handling
		if dataMap, ok := resp.BizData.(map[string]interface{}); ok {
			scdnResp.Data = dataMap
		} else {
			// If it's not a map, wrap it in a map
			scdnResp.Data = map[string]interface{}{
				"data": resp.BizData,
			}
		}
	}

	return &scdnResp, nil
}

// SetTimeout sets the client timeout
func (c *ScdnClient) SetTimeout(timeout time.Duration) {
	c.timeout = timeout
	c.sdk.Timeout = int(timeout.Seconds())
}

// GetTimeout returns the current client timeout
func (c *ScdnClient) GetTimeout() time.Duration {
	return c.timeout
}

// SetBaseURL sets the base URL for the client
func (c *ScdnClient) SetBaseURL(baseURL string) {
	c.baseURL = baseURL
	c.sdk.ApiPre = baseURL
}

// GetBaseURL returns the current base URL
func (c *ScdnClient) GetBaseURL() string {
	return c.baseURL
}

// SetCredentials sets the API credentials
func (c *ScdnClient) SetCredentials(apiKey, apiSecret string) {
	c.apiKey = apiKey
	c.apiSecret = apiSecret
	c.sdk.AppId = apiKey
	c.sdk.AppSecret = apiSecret
}

// GetCredentials returns the current API credentials
func (c *ScdnClient) GetCredentials() (string, string) {
	return c.apiKey, c.apiSecret
}

// IsHealthy checks if the SCDN API is healthy
func (c *ScdnClient) IsHealthy(ctx context.Context) error {
	req := &ScdnRequest{
		Query: map[string]interface{}{
			"page":      1,
			"page_size": 1,
		},
	}

	_, err := c.Get(ctx, "/api/v5/domains", req)
	return err
}

// GetAPIVersion returns the API version being used
func (c *ScdnClient) GetAPIVersion() string {
	return "v5"
}

// GetServiceName returns the service name
func (c *ScdnClient) GetServiceName() string {
	return "scdn"
}
