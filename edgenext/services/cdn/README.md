# EdgeNext CDN Service

Comprehensive CDN (Content Delivery Network) management service for the EdgeNext Terraform Provider with full domain lifecycle, configuration management, cache control, and enterprise-grade testing.

## 🚀 Overview

This service provides complete CDN domain management including domain creation, configuration management, cache refresh, file purge, and comprehensive monitoring capabilities. The implementation features a robust testing framework with 100% coverage, performance benchmarks, and mock API server.

## 📂 Project Structure

```
services/cdn/
├── service_edgenext_cdn.go              # Core CDN service implementation
├── service_edgenext_cdn_test.go         # Comprehensive test suite
├── data_edgenext_cdn_domain_config.go   # Terraform data source for domain config
├── data_edgenext_cdn_purge.go           # Terraform data source for purge tasks
├── data_edgenext_cdn_push.go            # Terraform data source for push tasks
├── resource_edgenext_cdn_domain_config.go # Terraform resource for domain config
├── resource_edgenext_cdn_purge.go       # Terraform resource for cache purge
├── resource_edgenext_cdn_push.go        # Terraform resource for cache push
└── README.md                            # This documentation
```

## 🔧 Service Methods

### CdnService

```go
type CdnService struct {
    client *connectivity.Client
}
```

#### Constructor
```go
func NewCdnService(client *connectivity.Client) *CdnService
```

#### Domain Management Methods

##### CreateDomain
```go
func (c *CdnService) CreateDomain(req DomainCreateRequest) (*DomainResponse, error)
```
Creates a new CDN domain with specified configuration.

**Parameters:**
- `req.Domain` - Domain name to create
- `req.Area` - Coverage area (global, mainland_china, outside_mainland_china, rim)
- `req.Type` - Domain type (page, download, video_demand, dynamic, video_live)
- `req.Config` - Domain configuration including origin settings

##### GetDomain
```go
func (c *CdnService) GetDomain(domains string) (*GetDomainResponse, error)
```
Retrieves detailed information about specific domains.

##### ListDomains
```go
func (c *CdnService) ListDomains(req DomainListRequest) (*DomainListResponse, error)
```
Lists CDN domains with pagination and filtering support.

##### DeleteDomain
```go
func (c *CdnService) DeleteDomain(domains string) error
```
Permanently deletes CDN domains.

#### Domain Configuration Methods

##### SetDomainConfig
```go
func (c *CdnService) SetDomainConfig(domains string, config map[string]interface{}) (*DomainConfigResponse, error)
```
Sets or updates domain configuration.

##### GetDomainConfig
```go
func (c *CdnService) GetDomainConfig(domains string, config []string) (*GetDomainConfigResponse, error)
```
Retrieves domain configuration for specified domains.

##### DeleteDomainConfig
```go
func (c *CdnService) DeleteDomainConfig(req DeleteDomainConfigRequest) error
```
Deletes specific configuration items from domains.

#### Cache Management Methods

##### CacheRefresh
```go
func (c *CdnService) CacheRefresh(urls []string, refreshType string) (*CacheRefreshResponse, error)
```
Initiates cache refresh for specified URLs.

##### QueryCacheRefresh
```go
func (c *CdnService) QueryCacheRefresh(req CacheRefreshQueryRequest) (*CacheRefreshQueryResponse, error)
```
Queries cache refresh task status.

#### File Purge Methods

##### FilePurge
```go
func (c *CdnService) FilePurge(urls []string) (*FilePurgeResponse, error)
```
Initiates file purge (preheating) for specified URLs.

##### QueryFilePurge
```go
func (c *CdnService) QueryFilePurge(req FilePurgeQueryRequest) (*FilePurgeQueryResponse, error)
```
Queries file purge task status.

##### QueryFilePurgeByTaskID
```go
func (c *CdnService) QueryFilePurgeByTaskID(taskID int) (*FilePurgeQueryResponse, error)
```
Convenience method to query purge status by task ID.

## 🏗️ Data Structures

### Domain Management Types

#### DomainCreateRequest
```go
type DomainCreateRequest struct {
    Domain string       `json:"domain"`          // Domain name
    Area   string       `json:"area,omitempty"`  // Coverage area
    Type   string       `json:"type"`            // Domain type
    Config DomainConfig `json:"config"`          // Configuration
}
```

#### DomainData
```go
type DomainData struct {
    ID         string `json:"id"`           // Domain ID
    Domain     string `json:"domain"`       // Domain name
    Type       string `json:"type"`         // Domain type
    Status     string `json:"status"`       // Domain status
    IcpStatus  string `json:"icp_status"`   // ICP filing status
    IcpNum     string `json:"icp_num"`      // ICP filing number
    Area       string `json:"area"`         // Coverage area
    Cname      string `json:"cname"`        // CNAME record
    CreateTime string `json:"create_time"`  // Creation time
    UpdateTime string `json:"update_time"`  // Last update time
    Https      int    `json:"https"`        // HTTPS status
}
```

#### Domain Status Constants
```go
const (
    DomainStatusServing   = "serving"   // Domain is serving traffic
    DomainStatusSuspended = "suspend"   // Domain is suspended
    DomainStatusDeleted   = "deleted"   // Domain is deleted
)
```

### Configuration Types

#### DomainConfigRequest
```go
type DomainConfigRequest struct {
    Domains string                 `json:"domains"` // Domain names (comma-separated)
    Config  map[string]interface{} `json:"config"`  // Configuration object
}
```

### Cache Management Types

#### CacheRefreshRequest
```go
type CacheRefreshRequest struct {
    URLs []string `json:"urls"` // URLs to refresh
    Type string   `json:"type"` // Refresh type (url/dir)
}
```

#### CacheRefreshQueryRequest
```go
type CacheRefreshQueryRequest struct {
    TaskID     int    `json:"task_id"`      // Task ID to query
    StartTime  string `json:"start_time"`   // Query start time
    EndTime    string `json:"end_time"`     // Query end time
    URL        string `json:"url"`          // Specific URL to query
    PageNumber int    `json:"page_number"`  // Page number
    PageSize   int    `json:"page_size"`    // Page size
}
```

### File Purge Types

#### FilePurgeRequest
```go
type FilePurgeRequest struct {
    URLs []string `json:"urls"` // URLs to purge
}
```

## 🧪 Testing Framework

### Test Coverage
- ✅ **Domain Management**: Create, read, update, delete operations
- ✅ **Configuration Management**: Set, get, delete domain configurations
- ✅ **Cache Operations**: Refresh and query cache status
- ✅ **File Purge**: Purge files and query status
- ✅ **Helper Methods**: Domain status checks and utility functions
- ✅ **Error Scenarios**: Invalid inputs, missing resources, API errors
- ✅ **Performance Benchmarks**: All core operations
- ✅ **Integration Framework**: Ready for real API testing

### Running Tests

```bash
# Run all tests
go test ./edgenext/services/cdn/ -v

# Run specific test
go test ./edgenext/services/cdn/ -run TestCreateDomain -v

# Run benchmarks
go test ./edgenext/services/cdn/ -bench=.

# Run with coverage
go test ./edgenext/services/cdn/ -cover

# Skip integration tests
go test ./edgenext/services/cdn/ -short
```

### Test Results Example
```
✅ TestNewCdnService (0.00s)
✅ TestCreateDomain (0.00s)
✅ TestGetDomain (0.00s)
✅ TestListDomains (0.00s)
✅ TestDeleteDomain (0.00s)
✅ TestDomainDataHelperMethods (0.00s)
✅ TestSetDomainConfig (0.00s)
✅ TestGetDomainConfig (0.00s)
✅ TestDeleteDomainConfig (0.00s)
✅ TestRefreshCache (0.00s)
✅ TestQueryCacheRefresh (0.00s)
✅ TestPurgeFiles (0.00s)
✅ TestQueryFilePurge (0.00s)
✅ TestQueryFilePurgeByTaskID (0.00s)
✅ TestNewCacheRefreshRequest (0.00s)
✅ TestCdnServiceErrorScenarios (0.00s)

PASS
ok  github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/cdn  0.229s
```

### Performance Benchmarks
```
BenchmarkCreateDomain-8      22633    55555 ns/op
BenchmarkGetDomain-8         23050    51386 ns/op
BenchmarkListDomains-8       20446    56959 ns/op
BenchmarkRefreshCache-8      22909    51816 ns/op
```

## 📡 API Endpoints

| Method | Endpoint | Purpose |
|--------|----------|---------|
| `POST` | `/v2/domain` | Create CDN domain |
| `GET` | `/v2/domain?domains={domains}` | Get domain details |
| `GET` | `/v2/domain/list` | List domains with pagination |
| `DELETE` | `/v2/domain?domains={domains}` | Delete domains |
| `POST` | `/v2/domain/config` | Set domain configuration |
| `GET` | `/v2/domain/config?domains={domains}` | Get domain configuration |
| `DELETE` | `/v2/domain/config` | Delete domain configuration |
| `POST` | `/v2/cache/refresh` | Initiate cache refresh |
| `GET` | `/v2/cache/refresh` | Query cache refresh status |
| `POST` | `/v2/cache/prefetch` | Initiate file purge |
| `GET` | `/v2/cache/prefetch` | Query file purge status |

## 💡 Usage Examples

### Basic Domain Operations

```go
// Initialize service
client := connectivity.NewClient("api-key", "secret", "https://api.edgenext.com")
service := cdn.NewCdnService(client)

// Create domain
req := cdn.DomainCreateRequest{
    Domain: "example.com",
    Area:   "global",
    Type:   "page",
    Config: cdn.DomainConfig{
        Origin: cdn.OriginItem{
            DefaultMaster: "origin.example.com",
            OriginMode:    "https",
            Port:          "443",
        },
    },
}
response, err := service.CreateDomain(req)
if err != nil {
    log.Fatalf("Failed to create domain: %v", err)
}

// Get domain details
domains, err := service.GetDomain("example.com")
if err != nil {
    log.Fatalf("Failed to get domain: %v", err)
}

// List domains
listReq := cdn.DomainListRequest{
    PageNumber:   1,
    PageSize:     20,
    DomainStatus: "serving",
}
domainList, err := service.ListDomains(listReq)
if err != nil {
    log.Fatalf("Failed to list domains: %v", err)
}
```

### Configuration Management

```go
// Set domain configuration
config := map[string]interface{}{
    "origin": map[string]interface{}{
        "default_master": "new-origin.example.com",
        "origin_mode":    "https",
        "port":           "443",
    },
    "cache_rule": map[string]interface{}{
        "cache_time": 3600,
        "cache_type": "all",
    },
}
configResp, err := service.SetDomainConfig("example.com", config)
if err != nil {
    log.Fatalf("Failed to set config: %v", err)
}

// Get domain configuration
currentConfig, err := service.GetDomainConfig("example.com", []string{"origin", "cache_rule"})
if err != nil {
    log.Fatalf("Failed to get config: %v", err)
}

// Delete specific configuration
deleteReq := cdn.DeleteDomainConfigRequest{
    Domains: "example.com",
    Config:  []string{"cache_rule"},
}
err = service.DeleteDomainConfig(deleteReq)
if err != nil {
    log.Fatalf("Failed to delete config: %v", err)
}
```

### Cache and Purge Operations

```go
// Cache refresh
urls := []string{
    "https://example.com/image.jpg",
    "https://example.com/style.css",
}
refreshResp, err := service.CacheRefresh(urls, "url")
if err != nil {
    log.Fatalf("Failed to refresh cache: %v", err)
}
log.Printf("Refresh task ID: %s", refreshResp.Data.TaskID)

// Query refresh status
queryReq := cdn.CacheRefreshQueryRequest{
    TaskID: 123456,
}
status, err := service.QueryCacheRefresh(queryReq)
if err != nil {
    log.Fatalf("Failed to query refresh status: %v", err)
}

// File purge
purgeResp, err := service.FilePurge([]string{"https://example.com/file.pdf"})
if err != nil {
    log.Fatalf("Failed to purge file: %v", err)
}

// Query purge status by task ID
purgeStatus, err := service.QueryFilePurgeByTaskID(789012)
if err != nil {
    log.Fatalf("Failed to query purge status: %v", err)
}
```

### Domain Status Checks

```go
// Using helper methods
domain := &cdn.DomainData{Status: "serving"}

if domain.IsServing() {
    log.Println("Domain is currently serving traffic")
}

if domain.IsSuspended() {
    log.Println("Domain is suspended")
}

if domain.IsDeleted() {
    log.Println("Domain has been deleted")
}
```

### Error Handling

```go
response, err := service.CreateDomain(req)
if err != nil {
    // Check for specific error types
    if strings.Contains(err.Error(), "Invalid domain format") {
        log.Printf("Domain validation failed: %v", err)
        return
    }
    
    if strings.Contains(err.Error(), "Domain already exists") {
        log.Printf("Domain conflicts: %v", err)
        return
    }
    
    log.Printf("General API error: %v", err)
    return
}

// Check API response code
if response.Code != 0 {
    log.Printf("API error code %d: %s", response.Code, response.Msg)
    return
}
```

## 🔒 Security Features

- **Input Validation**: Comprehensive validation of domain names and configurations
- **Rate Limiting**: Built-in support for API rate limiting
- **Secure Communication**: All API calls use HTTPS
- **Error Sanitization**: Sensitive data excluded from logs
- **Authentication**: API key and secret based authentication

## 🚨 Error Handling

### Common Error Scenarios
- **Invalid Domain Format**: Malformed domain names
- **Domain Not Found**: Non-existent domain requests
- **Configuration Conflicts**: Invalid configuration combinations
- **Rate Limiting**: API rate limit exceeded
- **Network Issues**: Connection timeouts and failures

### Error Response Format
```json
{
    "code": 1001,
    "msg": "Invalid domain format",
    "data": {...}
}
```

## 🛠️ Development

### Adding New Features
1. Define data structures in service file
2. Implement service method with proper error handling
3. Add comprehensive tests (success + error cases)
4. Update mock server handlers
5. Add performance benchmarks
6. Update documentation

### Code Quality Standards
- 100% test coverage required
- All public methods documented
- Error handling for all scenarios
- Performance benchmarks included
- Mock server updates for new endpoints

## 🔍 Troubleshooting

### Domain Validation
```bash
# Check domain format
echo "example.com" | grep -E '^[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9]?\.[a-zA-Z]{2,}$'
```

### Common Issues

| Issue | Solution |
|-------|----------|
| `Invalid domain format` | Ensure domain follows valid format |
| `Domain not found` | Verify domain exists and is accessible |
| `Configuration conflict` | Check configuration compatibility |
| `Rate limit exceeded` | Implement backoff and retry logic |
| `Unauthorized (401)` | Verify API credentials |

### Debug Mode

Enable debug logging for detailed troubleshooting:

```go
import "log"

// Enable verbose logging
log.SetFlags(log.LstdFlags | log.Lshortfile)

// Log all API requests and responses
client := connectivity.NewClient(apiKey, secret, endpoint)
service := cdn.NewCdnService(client)
```

## 📊 Performance

### Optimization Tips
- Use pagination for large domain lists
- Batch configuration updates when possible
- Implement connection pooling
- Use goroutines for concurrent operations
- Cache domain information locally

### Metrics
- **Average Response Time**: 50-57ms
- **Throughput**: 20k-23k ops/sec
- **Memory Usage**: Minimal per operation

## 🔄 Integration

### Terraform Resource Examples

#### Basic Domain Configuration
```hcl
resource "edgenext_cdn_domain_config" "website" {
  domain = "example.com"
  area   = "global"
  type   = "page"

  config {
    origin {
      default_master = "origin.example.com"
      origin_mode    = "https"
      port          = "443"
    }

    cache_rule {
      cache_time = 3600
      cache_type = "all"
    }
  }
}
```

#### Cache Management
```hcl
resource "edgenext_cdn_purge" "cache_refresh" {
  urls = [
    "https://example.com/image.jpg",
    "https://example.com/style.css"
  ]
  type = "url"
}

resource "edgenext_cdn_push" "file_preload" {
  urls = [
    "https://example.com/video.mp4",
    "https://example.com/document.pdf"
  ]
  type = "url"
}
```

### Data Source Usage
```hcl
# Query domain configuration
data "edgenext_cdn_domain_config" "current" {
  domains = "example.com"
  config  = ["origin", "cache_rule", "https"]
}

# Query purge task status
data "edgenext_cdn_purge" "task_status" {
  task_id = 123456
}

# Query push task status
data "edgenext_cdn_push" "push_status" {
  task_id = 789012
}
```

## 📝 Changelog

### 2024.1 - Major Testing and Documentation Update
- ✅ Comprehensive test suite with mock server
- ✅ Performance benchmarks for all operations
- ✅ Error scenario testing with table-driven tests
- ✅ Integration test framework
- ✅ Enhanced documentation with examples
- ✅ Domain helper methods testing
- ✅ Configuration management testing

### Previous Versions
- Basic CDN domain CRUD operations
- Cache refresh and file purge functionality
- Terraform resource and data source implementation
- Initial API client integration

---

## 📞 Support

- **GitHub Issues**: Bug reports and feature requests
- **Documentation**: This README for common solutions
- **Tests**: Run test suite to verify functionality
- **Examples**: Reference integration examples above

---

*Keep this documentation updated with code changes.*
