package connectivity

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

// Client represents EdgeNext API client
type Client struct {
	apiKey   string
	secret   string
	endpoint string
	client   *resty.Client
}

// NewClient creates a new client instance
func NewClient(apiKey, secret, endpoint string) *Client {
	client := resty.New().
		SetBaseURL(endpoint).
		SetTimeout(30*time.Second).
		SetRetryCount(3).
		SetRetryWaitTime(1*time.Second).
		SetRetryMaxWaitTime(5*time.Second).
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "terraform-provider-edgenext/1.0.0")

	// Add authentication middleware
	client.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
		// Add API Key to request headers
		if apiKey != "" {
			req.SetHeader("X-API-Key", apiKey)
		}
		// Add Secret to request headers
		if secret != "" {
			req.SetHeader("X-API-Secret", secret)
		}
		// Copy Secret as token parameter and add to URL
		if secret != "" {
			req.SetQueryParam("token", secret)
		}
		return nil
	})

	return &Client{
		apiKey:   apiKey,
		secret:   secret,
		endpoint: endpoint,
		client:   client,
	}
}

// Get executes GET request
func (c *Client) Get(ctx context.Context, path string, result interface{}) error {
	resp, err := c.client.R().
		SetContext(ctx).
		SetResult(result).
		Get(path)

	if err != nil {
		return fmt.Errorf("GET request failed: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("GET request returned error status code: %d, response: %s", resp.StatusCode(), resp.String())
	}

	return nil
}

// Post executes POST request
func (c *Client) Post(ctx context.Context, path string, body interface{}, result interface{}) error {
	resp, err := c.client.R().
		SetContext(ctx).
		SetBody(body).
		SetResult(result).
		Post(path)

	if err != nil {
		return fmt.Errorf("POST request failed: %w", err)
	}

	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated {
		return fmt.Errorf("POST request returned error status code: %d, response: %s", resp.StatusCode(), resp.String())
	}

	return nil
}

// Put executes PUT request
func (c *Client) Put(ctx context.Context, path string, body interface{}, result interface{}) error {
	resp, err := c.client.R().
		SetContext(ctx).
		SetBody(body).
		SetResult(result).
		Put(path)

	if err != nil {
		return fmt.Errorf("PUT request failed: %w", err)
	}

	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusNoContent {
		return fmt.Errorf("PUT request returned error status code: %d, response: %s", resp.StatusCode(), resp.String())
	}

	return nil
}

// Delete executes DELETE request
func (c *Client) Delete(ctx context.Context, path string, result interface{}) error {
	resp, err := c.client.R().
		SetContext(ctx).
		SetResult(result).
		Delete(path)

	if err != nil {
		return fmt.Errorf("DELETE request failed: %w", err)
	}

	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusNoContent {
		return fmt.Errorf("DELETE request returned error status code: %d, response: %s", resp.StatusCode(), resp.String())
	}

	return nil
}

// DeleteWithBody executes DELETE request with request body
func (c *Client) DeleteWithBody(ctx context.Context, path string, body interface{}) error {
	resp, err := c.client.R().
		SetContext(ctx).
		SetBody(body).
		Delete(path)

	if err != nil {
		return fmt.Errorf("DELETE request failed: %w", err)
	}

	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusNoContent {
		return fmt.Errorf("DELETE request returned error status code: %d, response: %s", resp.StatusCode(), resp.String())
	}

	return nil
}

// DeleteWithBodyAndResult executes DELETE request with request body and response result
func (c *Client) DeleteWithBodyAndResult(ctx context.Context, path string, body interface{}, result interface{}) error {
	resp, err := c.client.R().
		SetContext(ctx).
		SetBody(body).
		SetResult(result).
		Delete(path)

	if err != nil {
		return fmt.Errorf("DELETE request failed: %w", err)
	}

	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusNoContent {
		return fmt.Errorf("DELETE request returned error status code: %d, response: %s", resp.StatusCode(), resp.String())
	}

	return nil
}

// Patch executes PATCH request
func (c *Client) Patch(ctx context.Context, path string, body interface{}, result interface{}) error {
	resp, err := c.client.R().
		SetContext(ctx).
		SetBody(body).
		SetResult(result).
		Patch(path)

	if err != nil {
		return fmt.Errorf("PATCH request failed: %w", err)
	}

	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusNoContent {
		return fmt.Errorf("PATCH request returned error status code: %d, response: %s", resp.StatusCode(), resp.String())
	}

	return nil
}

// GetWithQuery executes GET request with query parameters
func (c *Client) GetWithQuery(ctx context.Context, path string, query map[string]string, result interface{}) error {
	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParams(query).
		SetResult(result).
		Get(path)

	if err != nil {
		return fmt.Errorf("GET request failed: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("GET request returned error status code: %d, response: %s", resp.StatusCode(), resp.String())
	}

	return nil
}

// PostWithHeaders executes POST request with custom headers
func (c *Client) PostWithHeaders(ctx context.Context, path string, body interface{}, headers map[string]string, result interface{}) error {
	req := c.client.R().
		SetContext(ctx).
		SetBody(body).
		SetResult(result)

	for key, value := range headers {
		req.SetHeader(key, value)
	}

	resp, err := req.Post(path)

	if err != nil {
		return fmt.Errorf("POST request failed: %w", err)
	}

	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated {
		return fmt.Errorf("POST request returned error status code: %d, response: %s", resp.StatusCode(), resp.String())
	}

	return nil
}

// SetTimeout sets request timeout
func (c *Client) SetTimeout(timeout time.Duration) {
	c.client.SetTimeout(timeout)
}

// SetRetryCount sets retry count
func (c *Client) SetRetryCount(count int) {
	c.client.SetRetryCount(count)
}

// SetRetryWaitTime sets retry wait time
func (c *Client) SetRetryWaitTime(waitTime time.Duration) {
	c.client.SetRetryWaitTime(waitTime)
}

// SetRetryMaxWaitTime sets maximum retry wait time
func (c *Client) SetRetryMaxWaitTime(maxWaitTime time.Duration) {
	c.client.SetRetryMaxWaitTime(maxWaitTime)
}

// SetBaseURL sets base URL
func (c *Client) SetBaseURL(url string) {
	c.client.SetBaseURL(url)
}

// GetRestyClient gets the underlying resty client (for advanced usage)
func (c *Client) GetRestyClient() *resty.Client {
	return c.client
}

// Response represents the common structure of API response
type Response struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
	Code    int             `json:"code"`
}

// ErrorResponse represents the structure of error response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// IsSuccess checks if the response is successful
func (r *Response) IsSuccess() bool {
	return r.Success
}

// GetMessage gets response message
func (r *Response) GetMessage() string {
	return r.Message
}

// GetData gets response data
func (r *Response) GetData() json.RawMessage {
	return r.Data
}

// GetCode gets response code
func (r *Response) GetCode() int {
	return r.Code
}
