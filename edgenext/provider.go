package edgenext

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/cdn"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/ssl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider 返回EdgeNext CDN Terraform Provider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "EdgeNext API密钥，用于身份验证",
				Sensitive:   true,
			},
			"secret": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "EdgeNext密钥，用于身份验证",
				Sensitive:   true,
			},
			"endpoint": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "EdgeNext API端点地址",
			},
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     300,
				Description: "API请求超时时间（秒）",
			},
			"retry_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     3,
				Description: "API请求重试次数",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			// CDN域名及配置管理资源
			"edgenext_cdn_domain_config": cdn.ResourceEdgenextCdnDomainConfig(),

			// CDN缓存刷新和文件预热资源
			"edgenext_cdn_push":  cdn.ResourceEdgenextCdnPush(),
			"edgenext_cdn_purge": cdn.ResourceEdgenextCdnPurge(),

			// SSL证书管理资源
			"edgenext_ssl_certificate": ssl.ResourceEdgenextSslCertificate(),
		},
		DataSourcesMap: map[string]*schema.Resource{

			// CDN域名及配置数据源
			"edgenext_cdn_domain_config": cdn.DataSourceEdgenextCdnDomainConfig(),

			// CDN缓存刷新数据源
			"edgenext_cdn_push":   cdn.DataSourceEdgenextCdnPush(),
			"edgenext_cdn_pushes": cdn.DataSourceEdgenextCdnPushes(),

			// CDN文件预热数据源
			"edgenext_cdn_purge":  cdn.DataSourceEdgenextCdnPurge(),
			"edgenext_cdn_purges": cdn.DataSourceEdgenextCdnPurges(),

			// SSL证书数据源
			"edgenext_ssl_certificate":  ssl.DataSourceEdgenextSslCertificate(),
			"edgenext_ssl_certificates": ssl.DataSourceEdgenextSslCertificates(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

// ProviderConfigure 配置Provider并返回客户端实例
func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	// 获取配置参数
	apiKey := d.Get("api_key").(string)
	secret := d.Get("secret").(string)
	endpoint := d.Get("endpoint").(string)
	timeout := d.Get("timeout").(int)
	retryCount := d.Get("retry_count").(int)

	// 验证配置参数
	if err := validateProviderConfig(apiKey, secret, endpoint); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Provider配置验证失败",
			Detail:   err.Error(),
		})
		return nil, diags
	}

	// 创建客户端实例
	client, err := createClient(apiKey, secret, endpoint, timeout, retryCount)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "创建客户端失败",
			Detail:   err.Error(),
		})
		return nil, diags
	}

	// 测试连接
	// if err := testConnection(ctx, client); err != nil {
	// 	diags = append(diags, diag.Diagnostic{
	// 		Severity: diag.Warning,
	// 		Summary:  "连接测试失败",
	// 		Detail:   fmt.Sprintf("无法连接到EdgeNext API: %v", err),
	// 	})
	// 	// 连接测试失败不阻止Provider继续工作，只显示警告
	// }

	return client, diags
}

// validateProviderConfig 验证Provider配置参数
func validateProviderConfig(apiKey, secret, endpoint string) error {
	// 验证API密钥
	if strings.TrimSpace(apiKey) == "" {
		return fmt.Errorf("API密钥不能为空")
	}
	if len(apiKey) < 8 {
		return fmt.Errorf("API密钥长度不能少于8个字符")
	}

	// 验证密钥
	if strings.TrimSpace(secret) == "" {
		return fmt.Errorf("密钥不能为空")
	}
	if len(secret) < 8 {
		return fmt.Errorf("密钥长度不能少于8个字符")
	}

	// 验证端点
	if strings.TrimSpace(endpoint) == "" {
		return fmt.Errorf("API端点不能为空")
	}
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		return fmt.Errorf("API端点必须以http://或https://开头")
	}

	return nil
}

// createClient 创建EdgeNext客户端实例
func createClient(apiKey, secret, endpoint string, timeout, retryCount int) (*connectivity.Client, error) {
	// 创建基础客户端
	client := connectivity.NewClient(apiKey, secret, endpoint)

	// 配置超时设置
	if timeout > 0 {
		client.SetTimeout(time.Duration(timeout) * time.Second)
	}

	// 配置重试设置
	if retryCount > 0 {
		client.SetRetryCount(retryCount)
	}

	// 配置重试等待时间
	client.SetRetryWaitTime(1 * time.Second)
	client.SetRetryMaxWaitTime(10 * time.Second)

	return client, nil
}

// testConnection 测试与EdgeNext API的连接
func testConnection(ctx context.Context, client *connectivity.Client) error {
	// 尝试发送一个简单的健康检查请求
	// 这里可以根据实际的API端点进行调整
	var response interface{}
	err := client.Get(ctx, "/health", &response)
	if err != nil {
		return fmt.Errorf("连接测试失败: %w", err)
	}
	return nil
}

// GetClient 从Provider配置中获取客户端实例
func GetClient(meta interface{}) (*connectivity.Client, error) {
	client, ok := meta.(*connectivity.Client)
	if !ok {
		return nil, fmt.Errorf("无效的客户端类型: %T", meta)
	}
	return client, nil
}

// GetClientWithContext 从Provider配置中获取客户端实例（带上下文）
func GetClientWithContext(ctx context.Context, meta interface{}) (*connectivity.Client, error) {
	client, err := GetClient(meta)
	if err != nil {
		return nil, err
	}

	// 检查上下文是否已取消
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return client, nil
	}
}

// IsNotFoundError 检查是否为"未找到"错误
func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	errMsg := strings.ToLower(err.Error())
	notFoundKeywords := []string{
		"not found", "notfound", "404", "不存在", "未找到",
		"domain not found", "certificate not found",
	}

	for _, keyword := range notFoundKeywords {
		if strings.Contains(errMsg, keyword) {
			return true
		}
	}

	return false
}

// IsRateLimitError 检查是否为限流错误
func IsRateLimitError(err error) bool {
	if err == nil {
		return false
	}

	errMsg := strings.ToLower(err.Error())
	rateLimitKeywords := []string{
		"rate limit", "ratelimit", "too many requests", "429",
		"请求过于频繁", "限流", "频率限制",
	}

	for _, keyword := range rateLimitKeywords {
		if strings.Contains(errMsg, keyword) {
			return true
		}
	}

	return false
}

// IsAuthenticationError 检查是否为认证错误
func IsAuthenticationError(err error) bool {
	if err == nil {
		return false
	}

	errMsg := strings.ToLower(err.Error())
	authKeywords := []string{
		"unauthorized", "401", "forbidden", "403",
		"invalid api key", "invalid secret", "authentication failed",
		"未授权", "认证失败", "无效的api密钥", "无效的密钥",
	}

	for _, keyword := range authKeywords {
		if strings.Contains(errMsg, keyword) {
			return true
		}
	}

	return false
}

// FormatError 格式化错误信息
func FormatError(operation string, err error) string {
	if err == nil {
		return fmt.Sprintf("%s成功", operation)
	}

	// 根据错误类型返回不同的错误信息
	if IsNotFoundError(err) {
		return fmt.Sprintf("%s失败: 资源不存在", operation)
	}

	if IsRateLimitError(err) {
		return fmt.Sprintf("%s失败: 请求过于频繁，请稍后重试", operation)
	}

	if IsAuthenticationError(err) {
		return fmt.Sprintf("%s失败: 认证失败，请检查API密钥和密钥", operation)
	}

	return fmt.Sprintf("%s失败: %v", operation, err)
}
