package oss

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOSSObject() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOSSObjectRead,

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the bucket",
			},
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The object key",
			},
			// Computed attributes
			"content_length": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Size of the body in bytes",
			},
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A standard MIME type describing the format of the object data",
			},
			"content_encoding": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "What content encodings have been applied to the object",
			},
			"content_disposition": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Presentational information for the object",
			},
			"content_language": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Language the content is in",
			},
			"cache_control": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies caching behavior",
			},
			"expires": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date and time at which the object is no longer cacheable",
			},
			"etag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ETag generated for the object",
			},
			"last_modified": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last modified date of the object",
			},
			"storage_class": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Storage class of the object",
			},
			"acl": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ACL of the object",
			},
			"metadata": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "A map of metadata stored with the object",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Presigned URL of the object",
			},
		},
	}
}

func dataSourceOSSObjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ossClient, err := client.OSSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	bucketName := d.Get("bucket").(string)
	key := d.Get("key").(string)

	// Prepare GetObject input
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	}

	if versionID, ok := d.GetOk("version_id"); ok {
		input.VersionId = aws.String(versionID.(string))
	}

	// Get object
	output, err := ossClient.GetObject(ctx, input)
	if err != nil {
		return diag.Errorf("failed to get OSS object %s/%s: %s", bucketName, key, err)
	}

	// Set resource ID
	d.SetId(fmt.Sprintf("%s/%s", bucketName, key))

	// Set attributes
	d.Set("bucket", bucketName)
	d.Set("key", key)
	d.Set("content_length", aws.ToInt64(output.ContentLength))

	if output.ContentType != nil {
		d.Set("content_type", aws.ToString(output.ContentType))
	}
	if output.ContentEncoding != nil {
		d.Set("content_encoding", aws.ToString(output.ContentEncoding))
	}
	if output.ContentDisposition != nil {
		d.Set("content_disposition", aws.ToString(output.ContentDisposition))
	}
	if output.ContentLanguage != nil {
		d.Set("content_language", aws.ToString(output.ContentLanguage))
	}
	if output.CacheControl != nil {
		d.Set("cache_control", aws.ToString(output.CacheControl))
	}
	if output.ExpiresString != nil {
		t, err := time.Parse(time.RFC1123, *output.ExpiresString)
		if err != nil {
			return diag.Errorf("failed to parse expires time %s: %s", *output.ExpiresString, err)
		}
		d.Set("expires", t.Format(time.RFC3339))
	}
	if output.ETag != nil {
		d.Set("etag", strings.Trim(aws.ToString(output.ETag), "\""))
	}
	if output.LastModified != nil {
		d.Set("last_modified", output.LastModified.Format(time.RFC3339))
	}
	if output.StorageClass != "" {
		d.Set("storage_class", string(output.StorageClass))
	}

	// Set metadata
	if len(output.Metadata) > 0 {
		d.Set("metadata", output.Metadata)
	}

	// Get object ACL
	aclOutput, err := ossClient.GetObjectAcl(ctx, &s3.GetObjectAclInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if err == nil {
		acl := determineObjectACLFromGrants(aclOutput.Grants)
		d.Set("acl", acl)
	}

	// Get presigned URL
	url, err := ossClient.PresignObject(ctx, input, 900)
	if err != nil {
		return diag.Errorf("failed to presign OSS object %s/%s: %s", bucketName, key, err)
	}
	d.Set("url", url)

	return nil
}
