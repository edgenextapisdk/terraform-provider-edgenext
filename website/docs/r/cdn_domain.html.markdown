---
subcategory: "Content Delivery Network(CDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_cdn_domain"
sidebar_current: "docs-edgenext-resource-cdn_domain"
description: |-
  Provides a resource to create and manage CDN domain configuration.
---

# edgenext_cdn_domain

Provides a resource to create and manage CDN domain configuration.

## Example Usage

### Basic CDN domain configuration

```hcl
resource "edgenext_cdn_domain" "example" {
  domain = "example.com"
  area   = "global"
  type   = "page"

  config {
    origin {
      default_master = "origin.example.com"
      origin_mode    = "default"
    }
  }
}
```

### Advanced CDN domain configuration with cache rules and HTTPS

```hcl
resource "edgenext_cdn_domain" "example" {
  domain = "example.com"
  area   = "global"
  type   = "page"

  config {
    origin {
      default_master = "origin.example.com"
      default_slave  = "backup.example.com"
      origin_mode    = "custom"
      port           = 443
      ori_https      = "yes"
    }

    cache_rule {
      type         = 1
      pattern      = "jpg,png,gif"
      time         = 86400
      timeunit     = "s"
      ignore_query = "on"
    }

    cache_rule {
      type         = 1
      pattern      = "css,js"
      time         = 3600
      timeunit     = "s"
      ignore_query = "off"
    }

    https {
      cert_id     = 123
      http2       = "on"
      force_https = "302"
    }

    referer {
      type        = 2
      list        = ["*.example.com", "example.org"]
      allow_empty = true
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `config` - (Required, List) Domain configuration items
* `domain` - (Required, String, ForceNew) Accelerated domain name for setting functions
* `type` - (Required, String, ForceNew) Domain type.
	- page
	- download
	- video_demand
	- dynamic
* `area` - (Optional, String, ForceNew) Acceleration area.
	- mainland_china
	- outside_mainland_china
	- global

The `config` object supports the following:

* `origin` - (Required, List) Origin configuration
* `add_back_source_head` - (Optional, List) Add origin request headers
* `add_response_head` - (Optional, List) Add response headers
* `cache_rule_list` - (Optional, List) New cache rules
* `cache_rule` - (Optional, List) Cache rule list
* `cache_share` - (Optional, List) Shared cache
* `compress_response` - (Optional, List) Compress response
* `connect_timeout` - (Optional, List) Origin connection timeout
* `deny_url` - (Optional, List) Block illegal URLs
* `head_control` - (Optional, List) HTTP header control
* `https` - (Optional, List) HTTPS configuration
* `ip_black_list` - (Optional, List) IP blacklist
* `ip_white_list` - (Optional, List) IP whitelist
* `origin_host` - (Optional, List) Origin HOST
* `rate_limit` - (Optional, List) Rate limit
* `referer` - (Optional, List) Referer blacklist and whitelist
* `speed_limit` - (Optional, List) Speed limit
* `timeout` - (Optional, List) Origin read timeout

The `add_back_source_head` object of `config` supports the following:

* `head_name` - (Required, String) Origin request header name, maximum 20 can be added
* `head_value` - (Required, String) Origin request header value
* `write_when_exists` - (Optional, String) Whether to overwrite when the same request header exists, default value is yes. 
	- yes: Overwrite 
	- no: Do not overwrite

The `add_response_head` object of `config` supports the following:

* `list` - (Required, List) Response header list

The `cache_rule_list` object of `config` supports the following:

* `expire` - (Required, Int) Cache time, used with expire_unit, maximum time not exceeding 2 years. When time=0, no caching for the specified pattern (i.e., caching disabled)
* `match_method` - (Required, String) Cache type. 
	- ext: file extension 
	- dir: directory 
	- route: full path matching 
	- regex: regular expression
* `pattern` - (Required, String) Cache rules, multiple separated by commas. Default is cached with parameters, ignoring expiration time, not ignoring no-cache headers. 
	- when type=ext: jpg,png,gif 
	- when type=dir: /product/index,/test/index,/user/index 
	- when type=route: /index.html,/test/*.jpg,/user/get?index 
	- when type=regex: set the corresponding regex as described below. After setting regex, the request URL is matched against the regex, and if matched, this cache rule is used.
* `cache_or_not` - (Optional, String) Whether to cache, not set means use expire to determine whether to cache. 
	- yes
	- no
* `case_ignore` - (Optional, String) Whether to ignore case, Default is yes. 
	- yes: Ignore 
	- no: Do not ignore
* `expire_unit` - (Optional, String) Cache time unit, default value is s. 
	- Y year 
	- M month 
	- D day 
	- h hour 
	- i minute 
	- s second
* `follow_expired` - (Optional, String) Whether to follow origin server cache time, default value is no. 
	- yes
	- no
* `ignore_no_cache_headers` - (Optional, String) Whether to ignore no-cache information in origin server response headers, such as Cache-Control:no-cache, default value is no. 
	- yes
	- no
* `params` - (Optional, String) Only takes effect when query_params_op=customer, parameter list
* `priority` - (Optional, Int) Sort value, lower priority value means higher priority, duplicates are not allowed
* `query_params_op_way` - (Optional, String) Only takes effect when query_params_op=customer. 
	- keep 
	- remove
* `query_params_op_when` - (Optional, String) Only takes effect when query_params_op=customer. 
	- cache: only process during caching 
	- cache_back_source: process both caching and origin requests
* `query_params_op` - (Optional, String) Query parameter operation mode, default value is no. 
	- no
	- yes
	- customer

The `cache_rule` object of `config` supports the following:

* `pattern` - (Required, String) Cache rules, multiple separated by commas. Default is cached with parameters, ignoring expiration time, not ignoring no-cache headers. For example: 
	- When type=1: jpg,png,gif 
	- When type=2: /product/index,/test/index,/user/index 
	- When type=3: /index.html,/test/*.jpg,/user/get?index 
	- When type=4: set the corresponding regex. After setting regex, the request URL is matched against the regex, and if matched, this cache rule is used.
* `time` - (Required, Int) Cache time, used with timeunit, maximum time not exceeding 2 years, when time=0 no caching for specified pattern
* `type` - (Required, Int) Cache type. 
	- 1: file extension 
	- 2: directory 
	- 3: full path matching 
	- 4: regex
* `ignore_expired` - (Optional, String) Valid when cache time is greater than 0, ignore origin server expiration time, default value is on. 
	- on 
	- off
* `ignore_no_cache` - (Optional, String) Valid when cache time is greater than 0, ignore origin server no-cache headers, default value is off. 
	- on 
	- off
* `ignore_query` - (Optional, String) Valid when cache time is greater than 0, ignore parameters for caching and ignore parameters for origin requests, default value is off. 
	- on 
	- off
* `timeunit` - (Optional, String) Cache time unit, default value is s. 
	- Y: year 
	- M: month 
	- D: day 
	- h: hour 
	- i: minute 
	- s: second

The `cache_share` object of `config` supports the following:

* `domain` - (Required, String) domain to be shared, this item only takes effect when share_way is cross_single_share and cross_all_share
* `share_way` - (Required, String) sharing method. 
	- inner_share: HTTP and HTTPS share cache within this domain; 
	- cross_single_share: HTTP and HTTPS separately share cache between different domains 
	- cross_all_share: HTTP and HTTPS all share cache between different domains

The `compress_response` object of `config` supports the following:

* `content_type` - (Required, List) Corresponding headers, e.g., ['text/plain','application/x-javascript']
* `min_size_unit` - (Required, String) Minimum size unit. 
	- KB
	- MB
* `min_size` - (Required, Int) Minimum size, used with min_size_unit, indicates the minimum file size to start compression

The `connect_timeout` object of `config` supports the following:

* `origin_connect_timeout` - (Required, String) origin connection timeout duration, unit is s, value range: [5-60]

The `deny_url` object of `config` supports the following:

* `urls` - (Required, List) blocked URL list

The `head_control` object of `config` supports the following:

* `list` - (Required, List) HTTP header control rules list

The `https` object of `config` supports the following:

* `cert_id` - (Required, Int) Specify the bound certificate ID, which can be obtained through the certificate query interface. 
When cert_id=0, HTTPS service will be disabled for the domain.
* `force_https` - (Optional, String) Redirect HTTP requests to HTTPS protocol, Default value is 0. 
0: No redirect 
302: HTTP request 302 redirect to HTTPS request 
301: HTTP request 301 redirect to HTTPS request
* `http2` - (Optional, String) HTTP2 feature: 
	- on: Enable 
	- off: Disable
* `ocsp` - (Optional, String) No change when not specified. 
	- on: Enable 
	- off: Disable
* `ssl_protocol` - (Optional, List) ssl protocol. 
	- TLSv1 
	- TLSv1.1 
	- TLSv1.2 
	- TLSv1.3

The `ip_black_list` object of `config` supports the following:

* `list` - (Required, List) IP blacklist. IP format supports /8, /16, /24 subnet formats, IPs between subnets cannot overlap; maximum 500 IP formats can be set, multiple IP formats separated by commas; IP blacklist cannot coexist with IP whitelist, setting IP blacklist will clear IP whitelist functionality.

The `ip_white_list` object of `config` supports the following:

* `list` - (Required, List) IP whitelist. IP format supports /8, /16, /24 subnet formats, IPs between subnets cannot overlap; maximum 500 IP formats can be set, multiple IP formats separated by commas; IP whitelist cannot coexist with IP blacklist, setting IP whitelist will clear IP blacklist functionality.

The `origin_host` object of `config` supports the following:

* `host` - (Required, String) Origin HOST

The `origin` object of `config` supports the following:

* `default_master` - (Optional, String) Primary origin address, can fill multiple IPs or domains.
Multiple IPs or domains separated by comma(,); primary and backup origins cannot have same IP or domain.
* `default_slave` - (Optional, String) Backup origin address, can fill multiple IPs or domains.
Multiple IPs or domains separated by comma(,); primary and backup origins cannot have same IP or domain.
* `ori_https` - (Optional, String) Whether to enable HTTPS protocol origin, this value needs to be set when origin_mode=custom. 
	- yes 
	- no
* `origin_mode` - (Optional, String) Origin mode: 
	- default: Origin with user request protocol and port
	- http: Origin with http protocol on port 80
	- https: Origin with https protocol on port 443
	- custom: Origin with custom protocol(ori_https) and port(port)
* `port` - (Optional, Int) This value needs to be set when origin_mode=custom. 
Origin port, valid value range [1-65535].

The `rate_limit` object of `config` supports the following:

* `max_rate_count` - (Required, Int) rate limit value
* `leading_flow_count` - (Optional, Int) how many bytes at the beginning are not rate limited
* `leading_flow_unit` - (Optional, String) unit for how many bytes at the beginning are not rate limited, default value is MB. 
	- KB
	- MB
* `max_rate_unit` - (Optional, String) rate limit unit, default value is MB. 
	- KB
	- MB

The `referer` object of `config` supports the following:

* `list` - (Required, List) Referer list, maximum 200 entries, multiple separated by commas; regex not supported; for wildcard domains, start with *, e.g., *.example2.com, including any matching host headers and empty host headers.
* `type` - (Required, Int) Anti-hotlinking type: 
	- 1: referer blacklist 
	- 2: referer whitelist
* `allow_empty` - (Optional, Bool) Whether to allow empty referer, default value is true. 
	- true
	- false

The `speed_limit` object of `config` supports the following:

* `begin_time` - (Required, String) Speed limit effective start time, format is HH:ii, e.g., (08:30) 24-hour format
* `end_time` - (Required, String) Speed limit effective end time, format is HH:ii, e.g., (08:30) 24-hour format
* `pattern` - (Required, String) URL matching rules, multiple separated by commas. 
	- When type=all: only supports .* 
	- When type=ext: jpg,png,gif 
	- When type=dir: /product/index,/test/index,/user/index 
	- When type=route: /index.html,/test/*.jpg,/user/get?index 
	- When type=regex: set the corresponding regex, after setting regex, match the request URL against the regex, if matched then use this speed limit rule.
* `speed` - (Required, Int) Speed limit value, unit is Kbps, actual effect will be converted to KB
* `type` - (Optional, String) URL matching type. 
	- ext: file extension
	- dir: directory
	- route: full path matching
	- regex: regular expression
	- all: all

The `timeout` object of `config` supports the following:

* `time` - (Required, String) timeout duration, unit is s, value range: [5-300]

The `list` object of `add_response_head` supports the following:

* `name` - (Required, String) Response header name
* `value` - (Required, String) Response header value

The `list` object of `head_control` supports the following:

* `head_direction` - (Required, String) direction. 
	- CLI_REQ: client request header 
	- CLI_REP: client response header 
	- SER_REQ: origin request header 
	- SER_REP: origin response header
* `head_op` - (Required, String) operation content. 
	- ADD: add 
	- DEL: delete 
	- ALT: modify
* `head` - (Required, String) http header name
* `order` - (Required, Int) priority
* `regex` - (Required, String) matching URL
* `value` - (Required, String) http header value

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `cname` - CNAME
* `create_time` - Creation time
* `https` - HTTPS
* `icp_num` - ICP filing number
* `icp_status` - ICP filing status
* `status` - Domain status
* `update_time` - Update time


## Import

CDN domain configuration can be imported using the domain name:

```shell
terraform import edgenext_cdn_domain.example example.com
```

