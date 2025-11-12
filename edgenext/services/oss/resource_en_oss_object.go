package oss

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceOSSObject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOSSObjectCreate,
		ReadContext:   resourceOSSObjectRead,
		UpdateContext: resourceOSSObjectUpdate,
		DeleteContext: resourceOSSObjectDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceOSSObjectImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the bucket to put the object in",
			},
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the object once it is in the bucket",
			},
			"source": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"content"},
				Description:   "The path to a file that will be read and uploaded as raw bytes for the object content, conflicts with content",
			},
			"content": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"source"},
				Description:   "Literal string value to use as the object content, conflicts with source",
			},
			"content_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "A standard MIME type describing the format of the object data",
			},
			"content_encoding": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies what content encodings have been applied to the object",
			},
			"content_disposition": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies presentational information for the object",
			},
			"cache_control": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies caching behavior along the request/reply chain",
			},
			"expires": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The date and time at which the object is no longer cacheable",
			},
			"metadata": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "A map of metadata to store with the object",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"acl": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "private",
				Description: "The canned ACL to apply to the object (private, public-read, public-read-write, authenticated-read)",
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
		},
	}
}

func resourceOSSObjectCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ossClient, err := client.OSSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	bucketName := d.Get("bucket").(string)
	key := d.Get("key").(string)

	// Get object content
	var body io.Reader
	var contentLength int64

	if source, ok := d.GetOk("source"); ok {
		// Read from file
		file, err := os.Open(source.(string))
		if err != nil {
			return diag.Errorf("failed to open source file %s: %s", source.(string), err)
		}
		defer file.Close()

		fileInfo, err := file.Stat()
		if err != nil {
			return diag.Errorf("failed to stat source file %s: %s", source.(string), err)
		}
		contentLength = fileInfo.Size()
		body = file
	} else if content, ok := d.GetOk("content"); ok {
		// Use content string
		contentStr := content.(string)
		body = strings.NewReader(contentStr)
		contentLength = int64(len(contentStr))
	} else {
		return diag.Errorf("either 'source' or 'content' must be specified")
	}

	// Prepare PutObject input
	input := &s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(key),
		Body:          body,
		ContentLength: aws.Int64(contentLength),
	}

	// Set optional parameters
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
	// Set expires
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

	// Set ACL
	if acl, ok := d.GetOk("acl"); ok {
		input.ACL = convertToObjectACL(acl.(string))
	}

	// Upload object
	err = ossClient.PutObject(ctx, input)
	if err != nil {
		return diag.Errorf("failed to upload OSS object %s/%s: %s", bucketName, key, err)
	}

	// Set resource ID
	d.SetId(fmt.Sprintf("%s/%s", bucketName, key))

	return resourceOSSObjectRead(ctx, d, m)
}

func resourceOSSObjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	// This prevents unnecessary diffs when API returns default values
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

func resourceOSSObjectUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	// If content or source changed, re-upload the object
	if d.HasChange("source") || d.HasChange("content") ||
		d.HasChange("content_type") || d.HasChange("content_encoding") ||
		d.HasChange("content_disposition") || d.HasChange("cache_control") ||
		d.HasChange("expires") || d.HasChange("metadata") {

		// Re-create the object
		return resourceOSSObjectCreate(ctx, d, m)
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

	return resourceOSSObjectRead(ctx, d, m)
}

func resourceOSSObjectDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

func resourceOSSObjectImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid import ID format, expected: bucket/key")
	}

	d.Set("bucket", parts[0])
	d.Set("key", parts[1])

	return []*schema.ResourceData{d}, nil
}

// convertToObjectACL converts ACL string to ObjectCannedACL type
func convertToObjectACL(acl string) types.ObjectCannedACL {
	switch acl {
	case "private":
		return types.ObjectCannedACLPrivate
	case "public-read":
		return types.ObjectCannedACLPublicRead
	case "public-read-write":
		return types.ObjectCannedACLPublicReadWrite
	case "authenticated-read":
		return types.ObjectCannedACLAuthenticatedRead
	default:
		return types.ObjectCannedACLPrivate
	}
}

// determineObjectACLFromGrants converts grants to a simple ACL string
func determineObjectACLFromGrants(grants []types.Grant) string {
	// ACL determination from grants is a best-effort approximation
	// because S3 ACLs can be more complex than the simple canned ACLs
	if len(grants) == 0 {
		return "private"
	}

	// S3 standard group URIs
	const (
		allUsersURI           = "http://acs.amazonaws.com/groups/global/AllUsers"
		authenticatedUsersURI = "http://acs.amazonaws.com/groups/global/AuthenticatedUsers"
	)

	hasPublicRead := false
	hasPublicWrite := false
	hasPublicFullControl := false
	hasAuthenticatedRead := false

	// Scan grants to determine permissions
	for _, grant := range grants {
		if grant.Grantee == nil {
			continue
		}

		// Check for group URIs (public access)
		if grant.Grantee.URI != nil {
			uri := aws.ToString(grant.Grantee.URI)

			switch uri {
			case allUsersURI:
				// Public (AllUsers) permissions
				switch grant.Permission {
				case types.PermissionRead:
					hasPublicRead = true
				case types.PermissionWrite:
					hasPublicWrite = true
				case types.PermissionFullControl:
					hasPublicFullControl = true
				case types.PermissionReadAcp, types.PermissionWriteAcp:
					// ACL permissions don't directly map to canned ACLs
					// but indicate some level of public access
				}

			case authenticatedUsersURI:
				// Authenticated users permissions
				switch grant.Permission {
				case types.PermissionRead:
					hasAuthenticatedRead = true
				case types.PermissionWrite:
					// Write permission for authenticated users
					hasAuthenticatedRead = true
				case types.PermissionFullControl:
					hasAuthenticatedRead = true
				}
			}
		}
	}

	// Determine the canned ACL based on the grants
	// Priority: most permissive first
	if hasPublicFullControl || (hasPublicRead && hasPublicWrite) {
		return "public-read-write"
	} else if hasPublicRead {
		return "public-read"
	} else if hasAuthenticatedRead {
		return "authenticated-read"
	}

	// Default to private if no public or authenticated access
	// Note: This might not be 100% accurate if there are custom grants
	return "private"
}
