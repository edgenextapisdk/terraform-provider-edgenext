package oss

import "time"

// BucketInfo 存储桶信息
type BucketInfo struct {
	Name         string    `json:"name"`
	Region       string    `json:"region"`
	CreationDate time.Time `json:"creation_date,omitempty"`
	Owner        string    `json:"owner,omitempty"`
}

// ObjectInfo 对象信息
type ObjectInfo struct {
	Key          string    `json:"key"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"last_modified"`
	ETag         string    `json:"etag"`
	ContentType  string    `json:"content_type,omitempty"`
}
