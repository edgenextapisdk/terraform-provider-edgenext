package connectivity

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

// ECSClient represents EdgeNext ECS API client.
type ECSClient struct {
	client    *resty.Client
	accessKey string
	secretKey string
	region    string
}

// NewECSClient creates a new ECS API client instance.
func NewECSClient(accessKey, secretKey, endpoint, region string) *ECSClient {
	normalizedRegion := strings.ToLower(strings.TrimSpace(region))
	client := resty.New().
		SetBaseURL(endpoint).
		SetTimeout(30*time.Second).
		SetRetryCount(3).
		SetRetryWaitTime(1*time.Second).
		SetRetryMaxWaitTime(5*time.Second).
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "terraform-provider-edgenext/1.0.0").
		SetHeader("X-Region", normalizedRegion)

	return &ECSClient{
		client:    client,
		accessKey: accessKey,
		secretKey: secretKey,
		region:    normalizedRegion,
	}
}

// Get executes GET request with ECS signature authentication.
func (c *ECSClient) Get(ctx context.Context, path string, query map[string]interface{}, result interface{}) error {
	values := setQueryParamsFromValues(query)
	payload := []byte(values.Encode())
	timestamp, signature := c.sign(payload)

	resp, err := c.client.R().
		SetContext(ctx).
		SetHeaders(c.authHeaders(timestamp, signature)).
		SetQueryParamsFromValues(values).
		SetResult(result).
		Get(path)
	if err != nil {
		return fmt.Errorf("ECS GET request failed: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("ECS GET request returned error status code: %d, response: %s", resp.StatusCode(), resp.String())
	}

	return nil
}

// Post executes POST request with ECS signature authentication.
func (c *ECSClient) Post(ctx context.Context, path string, body interface{}, result interface{}) error {
	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal ECS POST body: %w", err)
	}

	timestamp, signature := c.sign(payload)
	resp, err := c.client.R().
		SetContext(ctx).
		SetHeaders(c.authHeaders(timestamp, signature)).
		SetBody(body).
		SetResult(result).
		Post(path)
	if err != nil {
		return fmt.Errorf("ECS POST request failed: %w", err)
	}

	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated {
		return fmt.Errorf("ECS POST request returned error status code: %d, response: %s", resp.StatusCode(), resp.String())
	}

	return nil
}

func (c *ECSClient) authHeaders(timestamp, signature string) map[string]string {
	return map[string]string{
		"Authorization":      fmt.Sprintf("Bearer %s", c.accessKey),
		"Edgenext-Timestamp": timestamp,
		"Signature":          signature,
	}
}

func (c *ECSClient) Region() string {
	return c.region
}

func (c *ECSClient) sign(payload []byte) (string, string) {
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	key := fmt.Sprintf("%s-%s", c.secretKey, timestamp)

	expectedSignature := hmac.New(sha256.New, []byte(key))
	expectedSignature.Write(payload)
	signature := base64.StdEncoding.EncodeToString(expectedSignature.Sum(nil))

	return timestamp, signature
}

func setQueryParamsFromValues(data map[string]interface{}) url.Values {
	values := url.Values{}
	for key, value := range data {
		switch v := value.(type) {
		case string:
			values.Add(key, v)
		default:
			values.Add(key, fmt.Sprintf("%v", v))
		}
	}

	return values
}
