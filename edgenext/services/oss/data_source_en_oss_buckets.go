package oss

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOSSBuckets() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOSSBucketsRead,

		Schema: map[string]*schema.Schema{
			"bucket_prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A prefix string to filter results by bucket name",
			},
			"max_buckets": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1000,
				Description: "The maximum number of buckets to return.",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`)",
			},
			"buckets": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of buckets",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the bucket",
						},
						"creation_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The creation date of the bucket",
						},
						"acl": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The access control list (ACL) of the bucket",
						},
					},
				},
			},
		},
	}
}

func dataSourceOSSBucketsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ossClient, err := client.OSSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	// Get filter parameters
	bucketPrefix := d.Get("bucket_prefix").(string)
	maxBuckets := d.Get("max_buckets").(int)

	// List all buckets
	output, err := ossClient.ListBuckets(ctx)
	if err != nil {
		return diag.Errorf("failed to list OSS buckets: %s", err)
	}

	ids := make([]string, 0)
	var buckets []map[string]interface{}

	// Filter and limit buckets
	for _, bucket := range output.Buckets {
		bucketName := aws.ToString(bucket.Name)

		// Apply bucket_prefix filter
		if bucketPrefix != "" && !strings.HasPrefix(bucketName, bucketPrefix) {
			continue
		}

		// Apply max_buckets limit
		if len(buckets) >= maxBuckets {
			break
		}

		bucketMap := map[string]interface{}{
			"name": bucketName,
		}

		if bucket.CreationDate != nil {
			bucketMap["creation_date"] = bucket.CreationDate.Format(time.RFC3339)
		}

		// Get bucket ACL
		aclOutput, err := ossClient.GetBucketAcl(ctx, &s3.GetBucketAclInput{
			Bucket: aws.String(bucketName),
		})
		if err == nil {
			// Determine the ACL based on grants
			acl := determineACLFromGrants(aclOutput.Grants)
			bucketMap["acl"] = acl
		}

		buckets = append(buckets, bucketMap)
		ids = append(ids, bucketName)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	d.Set("buckets", buckets)

	// Write to output file if specified
	if outputFile := d.Get("output_file").(string); outputFile != "" {
		outputData := map[string]interface{}{
			"bucket_prefix": bucketPrefix,
			"max_buckets":   maxBuckets,
			"total_count":   len(buckets),
			"buckets":       buckets,
		}
		if err := helper.WriteToFile(d, outputData); err != nil {
			return diag.Errorf("failed to write output file: %s", err)
		}
	}

	return nil
}
