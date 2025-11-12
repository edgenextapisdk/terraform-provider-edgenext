package connectivity

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

// APIClient represents EdgeNext API client
type APIClient struct {
	client *resty.Client
}

// NewAPIClient creates a new API client instance
func NewAPIClient(accessKey, secretKey, endpoint string) *APIClient {
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
		if accessKey != "" {
			req.SetHeader("X-API-Key", accessKey)
		}
		// Add Secret to request headers
		if secretKey != "" {
			req.SetHeader("X-API-Secret", secretKey)
		}
		// Copy Secret as token parameter and add to URL
		if secretKey != "" {
			req.SetQueryParam("token", secretKey)
		}
		return nil
	})

	return &APIClient{
		client: client,
	}
}

// Get executes GET request
func (c *APIClient) Get(ctx context.Context, path string, result interface{}) error {
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
func (c *APIClient) Post(ctx context.Context, path string, body interface{}, result interface{}) error {
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
func (c *APIClient) Put(ctx context.Context, path string, body interface{}, result interface{}) error {
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
func (c *APIClient) Delete(ctx context.Context, path string, result interface{}) error {
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
func (c *APIClient) DeleteWithBody(ctx context.Context, path string, body interface{}) error {
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
func (c *APIClient) DeleteWithBodyAndResult(ctx context.Context, path string, body interface{}, result interface{}) error {
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
func (c *APIClient) Patch(ctx context.Context, path string, body interface{}, result interface{}) error {
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
func (c *APIClient) GetWithQuery(ctx context.Context, path string, query map[string]string, result interface{}) error {
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
func (c *APIClient) PostWithHeaders(ctx context.Context, path string, body interface{}, headers map[string]string, result interface{}) error {
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
func (c *APIClient) SetTimeout(timeout time.Duration) {
	c.client.SetTimeout(timeout)
}

// SetRetryCount sets retry count
func (c *APIClient) SetRetryCount(count int) {
	c.client.SetRetryCount(count)
}

// SetRetryWaitTime sets retry wait time
func (c *APIClient) SetRetryWaitTime(waitTime time.Duration) {
	c.client.SetRetryWaitTime(waitTime)
}

// SetRetryMaxWaitTime sets maximum retry wait time
func (c *APIClient) SetRetryMaxWaitTime(maxWaitTime time.Duration) {
	c.client.SetRetryMaxWaitTime(maxWaitTime)
}

// SetBaseURL sets base URL
func (c *APIClient) SetBaseURL(url string) {
	c.client.SetBaseURL(url)
}

// GetRestyClient gets the underlying resty client (for advanced usage)
func (c *APIClient) GetRestyClient() *resty.Client {
	return c.client
}
