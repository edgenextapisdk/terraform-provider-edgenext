package oss

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceOSSObjectCopy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOSSObjectCopyCreate,
		ReadContext:   resourceOSSObjectCopyRead,
		UpdateContext: resourceOSSObjectCopyUpdate,
		DeleteContext: resourceOSSObjectCopyDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"source_bucket": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the source bucket",
			},
			"source_key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The key of the source object",
			},
			"bucket": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the destination bucket",
			},
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The key of the destination object",
			},
			"acl": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "private",
				Description: "The canned ACL to apply to the object (private, public-read, public-read-write, authenticated-read)",
			},
			"metadata_directive": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "COPY",
				Description: "Specifies whether the metadata is copied from the source object or replaced with metadata provided in the request (COPY or REPLACE)",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if v != "COPY" && v != "REPLACE" {
						errs = append(errs, fmt.Errorf("%q must be either COPY or REPLACE, got: %s", key, v))
					}
					return
				},
			},
			"content_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "A standard MIME type describing the format of the object data, only used when metadata_directive is REPLACE",
			},
			"content_encoding": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies what content encodings have been applied to the object, only used when metadata_directive is REPLACE",
			},
			"content_disposition": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies presentational information for the object, only used when metadata_directive is REPLACE",
			},
			"cache_control": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies caching behavior along the request/reply chain, only used when metadata_directive is REPLACE",
			},
			"expires": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The date and time at which the object is no longer cacheable, only used when metadata_directive is REPLACE",
			},
			"metadata": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "A map of metadata to store with the object, only used when metadata_directive is REPLACE",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"etag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ETag generated for the object",
			},
			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The size of the object in bytes",
			},
			"last_modified": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last modified date of the object",
			},
			// "source_version_id": {
			// 	Type:        schema.TypeString,
			// 	Optional:    true,
			// 	ForceNew:    true,
			// 	Description: "Version ID of the source object (if versioning is enabled)",
			// },
		},
	}
}

func resourceOSSObjectCopyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ossClient, err := client.OSSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	sourceBucket := d.Get("source_bucket").(string)
	sourceKey := d.Get("source_key").(string)
	destBucket := d.Get("bucket").(string)
	destKey := d.Get("key").(string)

	// Build copy source string
	copySource := fmt.Sprintf("%s/%s", sourceBucket, sourceKey)
	// if sourceVersionId, ok := d.GetOk("source_version_id"); ok {
	// 	copySource = fmt.Sprintf("%s?versionId=%s", copySource, sourceVersionId.(string))
	// }

	// Prepare CopyObject input
	input := &s3.CopyObjectInput{
		Bucket:     aws.String(destBucket),
		Key:        aws.String(destKey),
		CopySource: aws.String(copySource),
	}

	// Set metadata directive
	if metadataDirective, ok := d.GetOk("metadata_directive"); ok {
		directive := metadataDirective.(string)
		if directive == "REPLACE" {
			input.MetadataDirective = types.MetadataDirectiveReplace
		} else {
			input.MetadataDirective = types.MetadataDirectiveCopy
		}
	}

	// Set optional parameters (only used when metadata_directive is REPLACE)
	metadataDirective := d.Get("metadata_directive").(string)
	if metadataDirective == "REPLACE" {
		if contentType, ok := d.GetOk("content_type"); ok {
			input.ContentType = aws.String(contentType.(string))
		}
		if contentEncoding, ok := d.GetOk("content_encoding"); ok {
			input.ContentEncoding = aws.String(contentEncoding.(string))
		}
		if contentDisposition, ok := d.GetOk("content_disposition"); ok {
			input.ContentDisposition = aws.String(contentDisposition.(string))
		}
		if cacheControl, ok := d.GetOk("cache_control"); ok {
			input.CacheControl = aws.String(cacheControl.(string))
		}
		if expires, ok := d.GetOk("expires"); ok {
			expiresTime, err := time.Parse(time.RFC3339, expires.(string))
			if err != nil {
				return diag.Errorf("failed to parse expires time %s: %s", expires.(string), err)
			}
			input.Expires = aws.Time(expiresTime)
		}

		// Set metadata
		if metadata, ok := d.GetOk("metadata"); ok {
			metadataMap := metadata.(map[string]interface{})
			input.Metadata = make(map[string]string)
			for k, v := range metadataMap {
				input.Metadata[k] = v.(string)
			}
		}
	}

	// Set ACL
	if acl, ok := d.GetOk("acl"); ok {
		input.ACL = convertToObjectACL(acl.(string))
	}

	// Copy object
	output, err := ossClient.CopyObject(ctx, input)
	if err != nil {
		return diag.Errorf("failed to copy OSS object from %s/%s to %s/%s: %s", sourceBucket, sourceKey, destBucket, destKey, err)
	}

	// Set resource ID
	d.SetId(fmt.Sprintf("%s/%s", destBucket, destKey))

	// Set computed attributes from copy result
	if output.CopyObjectResult != nil {
		if output.CopyObjectResult.ETag != nil {
			d.Set("etag", strings.Trim(aws.ToString(output.CopyObjectResult.ETag), "\""))
		}
		if output.CopyObjectResult.LastModified != nil {
			d.Set("last_modified", output.CopyObjectResult.LastModified.Format(time.RFC3339))
		}
	}

	return resourceOSSObjectCopyRead(ctx, d, m)
}

func resourceOSSObjectCopyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ossClient, err := client.OSSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	// Parse ID
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return diag.Errorf("invalid resource ID format: %s", d.Id())
	}
	bucketName := parts[0]
	key := parts[1]

	// Head object to get metadata
	headOutput, err := ossClient.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		// Object not found
		d.SetId("")
		return nil
	}

	// Set attributes
	d.Set("bucket", bucketName)
	d.Set("key", key)
	d.Set("etag", strings.Trim(aws.ToString(headOutput.ETag), "\""))
	d.Set("size", aws.ToInt64(headOutput.ContentLength))

	if headOutput.LastModified != nil {
		d.Set("last_modified", headOutput.LastModified.Format(time.RFC3339))
	}

	if headOutput.ContentType != nil {
		d.Set("content_type", aws.ToString(headOutput.ContentType))
	}

	// Only set optional headers if they were explicitly configured by user
	if _, ok := d.GetOk("content_encoding"); ok && headOutput.ContentEncoding != nil {
		d.Set("content_encoding", aws.ToString(headOutput.ContentEncoding))
	}
	if _, ok := d.GetOk("content_disposition"); ok && headOutput.ContentDisposition != nil {
		d.Set("content_disposition", aws.ToString(headOutput.ContentDisposition))
	}
	if _, ok := d.GetOk("cache_control"); ok && headOutput.CacheControl != nil {
		d.Set("cache_control", aws.ToString(headOutput.CacheControl))
	}
	if _, ok := d.GetOk("expires"); ok && headOutput.ExpiresString != nil {
		t, err := time.Parse(time.RFC1123, *headOutput.ExpiresString)
		if err != nil {
			return diag.Errorf("failed to parse expires time %s: %s", *headOutput.ExpiresString, err)
		}
		d.Set("expires", t.Format(time.RFC3339))
	}

	// Set metadata - only keep user-configured keys to avoid drift from server-added metadata
	if configMetadata, ok := d.GetOk("metadata"); ok && len(headOutput.Metadata) > 0 {
		configMap := configMetadata.(map[string]interface{})
		filteredMetadata := make(map[string]string)

		// Only include metadata keys that were in the original configuration
		for key := range configMap {
			if value, exists := headOutput.Metadata[key]; exists {
				filteredMetadata[key] = value
			}
		}

		// Only set if we have any matching keys
		if len(filteredMetadata) > 0 {
			d.Set("metadata", filteredMetadata)
		}
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

	return nil
}

func resourceOSSObjectCopyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ossClient, err := client.OSSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return diag.Errorf("invalid resource ID format: %s", d.Id())
	}
	bucketName := parts[0]
	key := parts[1]

	// If metadata-related fields changed, need to re-copy with REPLACE directive
	if d.HasChange("content_type") || d.HasChange("content_encoding") ||
		d.HasChange("content_disposition") || d.HasChange("cache_control") ||
		d.HasChange("expires") || d.HasChange("metadata") || d.HasChange("metadata_directive") {

		// Re-copy the object
		return resourceOSSObjectCopyCreate(ctx, d, m)
	}

	// If only ACL changed, update it
	if d.HasChange("acl") {
		acl := d.Get("acl").(string)
		err = ossClient.PutObjectAcl(ctx, &s3.PutObjectAclInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
			ACL:    convertToObjectACL(acl),
		})
		if err != nil {
			return diag.Errorf("failed to update OSS object ACL %s/%s: %s", bucketName, key, err)
		}
	}

	return resourceOSSObjectCopyRead(ctx, d, m)
}

func resourceOSSObjectCopyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ossClient, err := client.OSSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return diag.Errorf("invalid resource ID format: %s", d.Id())
	}
	bucketName := parts[0]
	key := parts[1]

	err = ossClient.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return diag.Errorf("failed to delete OSS object %s/%s: %s", bucketName, key, err)
	}

	return nil
}
