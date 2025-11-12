package oss

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOSSObjects() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOSSObjectsRead,

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the bucket",
			},
			"prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Limits the response to keys that begin with the specified prefix",
			},
			"delimiter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A delimiter is a character you use to group keys",
			},
			"max_keys": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1000,
				Description: "Sets the maximum number of keys returned in the response",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`)",
			},
			"objects": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of objects",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The object key",
						},
						"size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The size of the object in bytes",
						},
						"etag": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ETag of the object",
						},
						"last_modified": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The last modified date of the object",
						},
						"storage_class": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The storage class of the object",
						},
					},
				},
			},
			"keys": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of object keys",
			},
			"common_prefixes": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of common prefixes (if delimiter is specified)",
			},
		},
	}
}

func dataSourceOSSObjectsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ossClient, err := client.OSSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	bucketName := d.Get("bucket").(string)

	// Prepare list input
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	}

	if prefix, ok := d.GetOk("prefix"); ok {
		input.Prefix = aws.String(prefix.(string))
	}

	if delimiter, ok := d.GetOk("delimiter"); ok {
		input.Delimiter = aws.String(delimiter.(string))
	}

	if maxKeys, ok := d.GetOk("max_keys"); ok {
		input.MaxKeys = aws.Int32(int32(maxKeys.(int)))
	}

	// List objects
	var objects []map[string]interface{}
	var keys []string
	var commonPrefixes []string

	// Use ListObjectsV2 for better pagination
	output, err := ossClient.ListObjectsV2(ctx, input)
	if err != nil {
		return diag.Errorf("failed to list OSS objects in bucket %s: %s", bucketName, err)
	}

	// Process objects
	for _, obj := range output.Contents {
		objectKey := aws.ToString(obj.Key)

		objectMap := map[string]interface{}{
			"key":  objectKey,
			"size": aws.ToInt64(obj.Size),
			"etag": strings.Trim(aws.ToString(obj.ETag), "\""),
		}

		if obj.LastModified != nil {
			objectMap["last_modified"] = obj.LastModified.Format(time.RFC3339)
		}

		if obj.StorageClass != "" {
			objectMap["storage_class"] = string(obj.StorageClass)
		}

		objects = append(objects, objectMap)
		keys = append(keys, objectKey)
	}

	// Process common prefixes (if delimiter was used)
	for _, prefix := range output.CommonPrefixes {
		commonPrefixes = append(commonPrefixes, aws.ToString(prefix.Prefix))
	}

	// Handle pagination if there are more results
	for output.IsTruncated != nil && *output.IsTruncated {
		input.ContinuationToken = output.NextContinuationToken
		output, err = ossClient.ListObjectsV2(ctx, input)
		if err != nil {
			return diag.Errorf("failed to list more OSS objects in bucket %s: %s", bucketName, err)
		}

		for _, obj := range output.Contents {
			objectKey := aws.ToString(obj.Key)

			objectMap := map[string]interface{}{
				"key":  objectKey,
				"size": aws.ToInt64(obj.Size),
				"etag": strings.Trim(aws.ToString(obj.ETag), "\""),
			}

			if obj.LastModified != nil {
				objectMap["last_modified"] = obj.LastModified.Format(time.RFC3339)
			}

			if obj.StorageClass != "" {
				objectMap["storage_class"] = string(obj.StorageClass)
			}

			objects = append(objects, objectMap)
			keys = append(keys, objectKey)
		}

		for _, prefix := range output.CommonPrefixes {
			commonPrefixes = append(commonPrefixes, aws.ToString(prefix.Prefix))
		}
	}

	d.SetId(fmt.Sprintf("%s-%s", bucketName, time.Now().UTC().Format(time.RFC3339)))
	d.Set("objects", objects)
	d.Set("keys", keys)
	d.Set("common_prefixes", commonPrefixes)

	// Write to output file if specified
	if outputFile := d.Get("output_file").(string); outputFile != "" {
		outputData := map[string]interface{}{
			"bucket":          bucketName,
			"prefix":          aws.ToString(input.Prefix),
			"delimiter":       aws.ToString(input.Delimiter),
			"max_keys":        int(aws.ToInt32(input.MaxKeys)),
			"total_count":     len(objects),
			"objects":         objects,
			"keys":            keys,
			"common_prefixes": commonPrefixes,
		}
		if err := helper.WriteToFile(d, outputData); err != nil {
			return diag.Errorf("failed to write output file: %s", err)
		}
	}
	return nil
}
