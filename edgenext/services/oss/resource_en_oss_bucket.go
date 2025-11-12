package oss

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceOSSBucket() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOSSBucketCreate,
		ReadContext:   resourceOSSBucketRead,
		UpdateContext: resourceOSSBucketUpdate,
		DeleteContext: resourceOSSBucketDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
				Description: "The name of the bucket (3-63 characters)",
			},
			"acl": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "private",
				Description: "The canned ACL to apply to the bucket (private, public-read, public-read-write, authenticated-read)",
			},
			"force_destroy": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "A boolean that indicates all objects should be deleted from the bucket so that the bucket can be destroyed without error",
			},
		},
	}
}

func resourceOSSBucketCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ossClient, err := client.OSSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	bucketName := d.Get("bucket").(string)
	acl := d.Get("acl").(string)

	// Convert ACL string to BucketCannedACL type
	var bucketACL types.BucketCannedACL
	switch acl {
	case "private":
		bucketACL = types.BucketCannedACLPrivate
	case "public-read":
		bucketACL = types.BucketCannedACLPublicRead
	case "public-read-write":
		bucketACL = types.BucketCannedACLPublicReadWrite
	case "authenticated-read":
		bucketACL = types.BucketCannedACLAuthenticatedRead
	default:
		bucketACL = types.BucketCannedACLPrivate
	}

	// Create bucket
	input := &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
		ACL:    bucketACL,
	}

	err = ossClient.CreateBucket(ctx, input)
	if err != nil {
		return diag.Errorf("failed to create OSS bucket %s: %s", bucketName, err)
	}

	// Put bucket cors
	err = ossClient.PutBucketCors(ctx, &s3.PutBucketCorsInput{
		Bucket: input.Bucket,
		CORSConfiguration: &types.CORSConfiguration{
			CORSRules: []types.CORSRule{
				{
					AllowedOrigins: []string{"*"},
					AllowedMethods: []string{"POST", "GET", "PUT", "DELETE", "HEAD"},
					AllowedHeaders: []string{"*"},
					ExposeHeaders:  []string{"*"},
					MaxAgeSeconds:  aws.Int32(3600),
				},
			},
		},
	})
	if err != nil {
		// Delete bucket if failed to put bucket cors
		err = ossClient.DeleteBucket(ctx, &s3.DeleteBucketInput{
			Bucket: input.Bucket,
		})
		if err != nil {
			return diag.Errorf("failed to delete bucket %s after failed to put bucket cors: %s, and failed to delete bucket: %s", bucketName, err, err)
		}
		return diag.Errorf("failed to put bucket cors %s: %s", bucketName, err)
	}

	d.SetId(bucketName)

	return resourceOSSBucketRead(ctx, d, m)
}

func resourceOSSBucketRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ossClient, err := client.OSSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	bucketName := d.Id()

	// Check if bucket exists
	exists, err := ossClient.BucketExists(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil || !exists {
		d.SetId("")
		return nil
	}

	// Set bucket name
	d.Set("bucket", bucketName)

	// Get bucket ACL
	aclOutput, err := ossClient.GetBucketAcl(ctx, &s3.GetBucketAclInput{
		Bucket: aws.String(bucketName),
	})
	if err == nil {
		// Determine ACL from grants
		acl := determineACLFromGrants(aclOutput.Grants)
		d.Set("acl", acl)
	}

	return nil
}

func resourceOSSBucketUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ossClient, err := client.OSSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	bucketName := d.Id()

	// Update ACL if changed
	if d.HasChange("acl") {
		acl := d.Get("acl").(string)

		var bucketACL types.BucketCannedACL
		switch acl {
		case "private":
			bucketACL = types.BucketCannedACLPrivate
		case "public-read":
			bucketACL = types.BucketCannedACLPublicRead
		case "public-read-write":
			bucketACL = types.BucketCannedACLPublicReadWrite
		case "authenticated-read":
			bucketACL = types.BucketCannedACLAuthenticatedRead
		default:
			bucketACL = types.BucketCannedACLPrivate
		}

		err = ossClient.PutBucketAcl(ctx, &s3.PutBucketAclInput{
			Bucket: aws.String(bucketName),
			ACL:    bucketACL,
		})
		if err != nil {
			return diag.Errorf("failed to update OSS bucket ACL %s: %s", bucketName, err)
		}
	}

	return resourceOSSBucketRead(ctx, d, m)
}

func resourceOSSBucketDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ossClient, err := client.OSSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	bucketName := d.Id()
	forceDestroy := d.Get("force_destroy").(bool)

	// If force_destroy is enabled, delete all objects first
	if forceDestroy {
		err := emptyBucket(ctx, ossClient, bucketName)
		if err != nil {
			return diag.Errorf("failed to empty OSS bucket %s: %s", bucketName, err)
		}
	}

	// Delete bucket
	err = ossClient.DeleteBucket(ctx, &s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return diag.Errorf("failed to delete OSS bucket %s: %s", bucketName, err)
	}

	return nil
}

// emptyBucket deletes all objects in a bucket
func emptyBucket(ctx context.Context, client *connectivity.OSSClient, bucketName string) error {
	// List all objects
	listOutput, err := client.ListObjects(ctx, &s3.ListObjectsInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return fmt.Errorf("failed to list objects: %w", err)
	}

	// Delete all objects
	for _, obj := range listOutput.Contents {
		err := client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(bucketName),
			Key:    obj.Key,
		})
		if err != nil {
			return fmt.Errorf("failed to delete object %s: %w", aws.ToString(obj.Key), err)
		}
	}

	// Handle pagination if there are more objects
	for listOutput.IsTruncated != nil && *listOutput.IsTruncated {
		listOutput, err = client.ListObjects(ctx, &s3.ListObjectsInput{
			Bucket: aws.String(bucketName),
			Marker: listOutput.NextMarker,
		})
		if err != nil {
			return fmt.Errorf("failed to list more objects: %w", err)
		}

		for _, obj := range listOutput.Contents {
			err := client.DeleteObject(ctx, &s3.DeleteObjectInput{
				Bucket: aws.String(bucketName),
				Key:    obj.Key,
			})
			if err != nil {
				return fmt.Errorf("failed to delete object %s: %w", aws.ToString(obj.Key), err)
			}
		}
	}

	return nil
}

// determineACLFromGrants converts grants to a simple ACL string
func determineACLFromGrants(grants []types.Grant) string {
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
