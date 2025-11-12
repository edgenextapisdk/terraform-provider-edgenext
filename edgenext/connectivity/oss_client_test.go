package connectivity

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// testConfig holds the shared test configuration
var testConfig = &Config{
	AccessKey: "your-access-key-here",
	SecretKey: "your-secret-key-here",
	Endpoint:  "your-edgenext-api-endpoint-here",
	Region:    "your-region-here",
}

// setupOSSClient creates a test OSS client for integration tests
func setupOSSClient(t *testing.T) *OSSClient {
	client, err := testConfig.Client()
	if err != nil {
		t.Fatalf("Failed to create EdgeNext client: %v", err)
	}

	ossClient, err := client.OSSClient()
	if err != nil {
		t.Fatalf("Failed to get OSS client: %v", err)
	}

	return ossClient
}

// generateTestBucketName generates a unique bucket name for testing
func generateTestBucketName(prefix string) string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%s-test-%d", prefix, timestamp)
}

// cleanupBucket deletes all objects in a bucket and then the bucket itself
func cleanupBucket(ctx context.Context, client *OSSClient, bucketName string) error {
	// List and delete all objects
	listOutput, err := client.ListObjects(ctx, &s3.ListObjectsInput{
		Bucket: aws.String(bucketName),
	})
	if err == nil {
		for _, obj := range listOutput.Contents {
			_ = client.DeleteObject(ctx, &s3.DeleteObjectInput{
				Bucket: aws.String(bucketName),
				Key:    obj.Key,
			})
		}
	}

	// Delete bucket
	return client.DeleteBucket(ctx, &s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	})
}

// TestOSSClientListBuckets tests listing buckets
func TestOSSClientListBuckets(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := setupOSSClient(t)
	ctx := context.Background()

	output, err := client.ListBuckets(ctx)
	if err != nil {
		t.Fatalf("Failed to list buckets: %v", err)
	}

	t.Logf("Successfully listed %d buckets", len(output.Buckets))

	// Log first few bucket names
	for i, bucket := range output.Buckets {
		if i < 5 {
			t.Logf("Bucket %d: %s (Created: %v)", i+1, aws.ToString(bucket.Name), bucket.CreationDate)
		}
	}
}

// TestOSSClientBucketLifecycle tests complete bucket lifecycle
func TestOSSClientBucketLifecycle(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := setupOSSClient(t)
	ctx := context.Background()

	bucketName := generateTestBucketName("oss-lifecycle")
	t.Logf("Testing bucket lifecycle with: %s", bucketName)

	// Step 1: Create bucket
	t.Run("CreateBucket", func(t *testing.T) {
		err := client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			t.Fatalf("Failed to create bucket: %v", err)
		}
		t.Logf("✓ Successfully created bucket: %s", bucketName)
	})

	// Ensure cleanup
	defer func() {
		_ = cleanupBucket(ctx, client, bucketName)
		t.Logf("✓ Cleanup completed for bucket: %s", bucketName)
	}()

	// Step 3: BucketExists
	t.Run("BucketExists", func(t *testing.T) {
		exists, err := client.BucketExists(ctx, &s3.HeadBucketInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			t.Fatalf("Failed to check bucket existence: %v", err)
		}
		if !exists {
			t.Fatalf("Bucket should exist but doesn't: %s", bucketName)
		}
		t.Logf("✓ Verified bucket exists: %s", bucketName)
	})
}

// TestOSSClientBucketACL tests bucket ACL operations
func TestOSSClientBucketACL(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := setupOSSClient(t)
	ctx := context.Background()

	bucketName := generateTestBucketName("oss-acl")
	t.Logf("Testing bucket ACL with: %s", bucketName)

	// Create bucket with private ACL
	err := client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
		ACL:    types.BucketCannedACLPrivate,
	})
	if err != nil {
		t.Fatalf("Failed to create bucket: %v", err)
	}

	defer cleanupBucket(ctx, client, bucketName)

	// Test GetBucketAcl
	t.Run("GetBucketAcl", func(t *testing.T) {
		output, err := client.GetBucketAcl(ctx, &s3.GetBucketAclInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			t.Fatalf("Failed to get bucket ACL: %v", err)
		}

		t.Logf("✓ Bucket ACL Owner: %s", aws.ToString(output.Owner.DisplayName))
		t.Logf("✓ Number of grants: %d", len(output.Grants))

		for i, grant := range output.Grants {
			t.Logf("  Grant %d: Grantee=%+v Permission=%s", i+1, grant.Grantee, grant.Permission)
		}
	})

	// Test PutBucketAcl
	t.Run("PutBucketAcl", func(t *testing.T) {
		err := client.PutBucketAcl(ctx, &s3.PutBucketAclInput{
			Bucket: aws.String(bucketName),
			ACL:    types.BucketCannedACLPublicRead,
		})
		if err != nil {
			t.Fatalf("Failed to put bucket ACL: %v", err)
		}
		t.Logf("✓ Successfully set bucket ACL")
	})

	// Test GetBucketAcl
	t.Run("GetBucketAcl", func(t *testing.T) {
		output, err := client.GetBucketAcl(ctx, &s3.GetBucketAclInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			t.Fatalf("Failed to get bucket ACL: %v", err)
		}

		t.Logf("✓ Bucket ACL Owner: %s", aws.ToString(output.Owner.DisplayName))
		t.Logf("✓ Number of grants: %d", len(output.Grants))

		for i, grant := range output.Grants {
			t.Logf("  Grant %d: Grantee=%+v Permission=%s", i+1, grant.Grantee, grant.Permission)
		}
	})
}

// TestOSSClientObjectOperations tests complete object lifecycle
func TestOSSClientObjectOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := setupOSSClient(t)
	ctx := context.Background()

	bucketName := generateTestBucketName("oss-obj")
	objectKey := "config/test-object.txt"
	objectContent := "Hello, EdgeNext OSS! This is a test object."

	// Create bucket
	err := client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
		ACL:    types.BucketCannedACLPublicRead,
	})
	if err != nil {
		t.Fatalf("Failed to create bucket: %v", err)
	}
	t.Logf("Created test bucket: %s", bucketName)

	err = client.PutBucketCors(ctx, &s3.PutBucketCorsInput{
		Bucket: aws.String(bucketName),
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
		t.Fatalf("Failed to put bucket cors: %v", err)
	}
	t.Logf("Successfully put bucket cors")

	defer cleanupBucket(ctx, client, bucketName)

	// Test PutObject
	t.Run("PutObject", func(t *testing.T) {
		err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:      aws.String(bucketName),
			Key:         aws.String(objectKey),
			Body:        strings.NewReader(objectContent),
			ContentType: aws.String("text/plain"),
			ACL:         types.ObjectCannedACLPublicRead,
		})
		if err != nil {
			t.Fatalf("Failed to put object: %v", err)
		}
		t.Logf("✓ Successfully uploaded object: %s", objectKey)
	})

	// Test HeadObject
	t.Run("HeadObject", func(t *testing.T) {
		output, err := client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
		})
		if err != nil {
			t.Fatalf("Failed to head object: %v", err)
		}
		t.Logf("✓ Object metadata: ContentLength=%d, ContentType=%s, ETag=%s",
			aws.ToInt64(output.ContentLength),
			aws.ToString(output.ContentType),
			aws.ToString(output.ETag))
	})

	// Test GetObject
	t.Run("GetObject", func(t *testing.T) {
		output, err := client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
		})
		if err != nil {
			t.Fatalf("Failed to get object: %v", err)
		}
		defer output.Body.Close()

		t.Logf("✓ Successfully retrieved object: %s (ContentLength=%d)",
			objectKey, aws.ToInt64(output.ContentLength))
	})

	// Test PresignObject
	t.Run("PresignObject", func(t *testing.T) {
		url, err := client.PresignObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
		}, 300)
		if err != nil {
			t.Fatalf("Failed to presign object: %v", err)
		}
		t.Logf("✓ Presigned object URL: %s", url)
	})

	// Test ListObjects
	t.Run("ListObjects", func(t *testing.T) {
		output, err := client.ListObjects(ctx, &s3.ListObjectsInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			t.Fatalf("Failed to list objects: %v", err)
		}

		if len(output.Contents) == 0 {
			t.Fatal("Expected at least one object, got none")
		}

		t.Logf("✓ Found %d objects in bucket", len(output.Contents))
		for _, obj := range output.Contents {
			t.Logf("  - %s (Size: %d bytes, Modified: %v)",
				aws.ToString(obj.Key),
				aws.ToInt64(obj.Size),
				obj.LastModified)
		}
	})

	// Test ListObjectsV2
	t.Run("ListObjectsV2", func(t *testing.T) {
		output, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			t.Fatalf("Failed to list objects v2: %v", err)
		}

		if len(output.Contents) == 0 {
			t.Fatal("Expected at least one object, got none")
		}

		t.Logf("✓ Found %d objects in bucket (V2 API)", len(output.Contents))
		t.Logf("  KeyCount: %d", aws.ToInt32(output.KeyCount))
	})

	// Test GetObjectAcl
	t.Run("GetObjectAcl", func(t *testing.T) {
		output, err := client.GetObjectAcl(ctx, &s3.GetObjectAclInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
		})
		if err != nil {
			t.Logf("⚠ Warning: Failed to get object ACL: %v", err)
		} else {
			t.Logf("✓ Object ACL retrieved, Grants: %d", len(output.Grants))
		}
	})

	// Test PutObjectAcl
	t.Run("PutObjectAcl", func(t *testing.T) {
		err := client.PutObjectAcl(ctx, &s3.PutObjectAclInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
			ACL:    types.ObjectCannedACLAuthenticatedRead,
		})
		if err != nil {
			t.Logf("⚠ Warning: Failed to put object ACL: %v", err)
		} else {
			t.Logf("✓ Successfully set object ACL")
		}
	})

	// Test GetObjectAcl
	t.Run("GetObjectAcl", func(t *testing.T) {
		output, err := client.GetObjectAcl(ctx, &s3.GetObjectAclInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
		})
		if err != nil {
			t.Logf("⚠ Warning: Failed to get object ACL: %v", err)
		} else {
			t.Logf("✓ Object ACL retrieved, Grants: %d", len(output.Grants))
		}
	})

	// Test DeleteObject
	t.Run("DeleteObject", func(t *testing.T) {
		err := client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
		})
		if err != nil {
			t.Fatalf("Failed to delete object: %v", err)
		}
		t.Logf("✓ Successfully deleted object: %s", objectKey)
	})
}

// TestOSSClientConcurrentOperations tests concurrent operations
func TestOSSClientConcurrentOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := setupOSSClient(t)
	ctx := context.Background()

	bucketName := generateTestBucketName("oss-concurrent")

	// Create bucket
	err := client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		t.Fatalf("Failed to create bucket: %v", err)
	}
	t.Logf("Created test bucket: %s", bucketName)

	defer cleanupBucket(ctx, client, bucketName)

	// Upload multiple objects concurrently
	const numObjects = 20
	var wg sync.WaitGroup
	wg.Add(numObjects)

	errors := make([]error, numObjects)

	t.Logf("Starting concurrent upload of %d objects...", numObjects)
	startTime := time.Now()

	for i := 0; i < numObjects; i++ {
		go func(index int) {
			defer wg.Done()

			objectKey := fmt.Sprintf("concurrent-object-%d.txt", index)
			content := fmt.Sprintf("Content for concurrent object %d", index)

			err := client.PutObject(ctx, &s3.PutObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(objectKey),
				Body:   strings.NewReader(content),
			})

			errors[index] = err
		}(i)
	}

	wg.Wait()
	duration := time.Since(startTime)

	// Check for errors
	successCount := 0
	errorCount := 0
	for i, err := range errors {
		if err != nil {
			t.Errorf("Failed to upload object %d: %v", i, err)
			errorCount++
		} else {
			successCount++
		}
	}

	t.Logf("✓ Concurrent uploads completed in %v", duration)
	t.Logf("  Success: %d/%d, Failed: %d/%d", successCount, numObjects, errorCount, numObjects)

	// Verify objects exist by listing
	listOutput, err := client.ListObjects(ctx, &s3.ListObjectsInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		t.Fatalf("Failed to list objects: %v", err)
	}

	if len(listOutput.Contents) != successCount {
		t.Errorf("Expected %d objects, got %d", successCount, len(listOutput.Contents))
	} else {
		t.Logf("✓ Verified %d objects in bucket", len(listOutput.Contents))
	}
}

// TestOSSClientErrorHandling tests error scenarios
func TestOSSClientErrorHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := setupOSSClient(t)
	ctx := context.Background()

	t.Run("DeleteNonExistentBucket", func(t *testing.T) {
		nonExistentBucket := fmt.Sprintf("non-existent-bucket-%d", time.Now().Unix())
		err := client.DeleteBucket(ctx, &s3.DeleteBucketInput{
			Bucket: aws.String(nonExistentBucket),
		})
		if err == nil {
			t.Error("Expected error when deleting non-existent bucket, got nil")
		} else {
			t.Logf("✓ Got expected error: %v", err)
		}
	})

	t.Run("GetNonExistentObject", func(t *testing.T) {
		// Create a temporary bucket
		bucketName := generateTestBucketName("oss-error-test")
		_ = client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(bucketName),
		})
		defer cleanupBucket(ctx, client, bucketName)

		// Try to get non-existent object
		_, err := client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(fmt.Sprintf("non-existent-object-%d.txt", time.Now().Unix())),
		})
		if err == nil {
			t.Error("Expected error when getting non-existent object, got nil")
		} else {
			t.Logf("✓ Got expected error: %v", err)
		}
	})

	t.Run("BucketExistsNonExistent", func(t *testing.T) {
		nonExistentBucket := fmt.Sprintf("non-existent-bucket-%d", time.Now().Unix())
		_, err := client.BucketExists(ctx, &s3.HeadBucketInput{
			Bucket: aws.String(nonExistentBucket),
		})
		if err == nil {
			t.Error("Expected error when checking non-existent bucket")
		} else {
			t.Logf("✓ Got expected error: %v", err)
		}
	})

	t.Run("HeadNonExistentObject", func(t *testing.T) {
		bucketName := generateTestBucketName("oss-error-head")
		_ = client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(bucketName),
		})
		defer cleanupBucket(ctx, client, bucketName)

		_, err := client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(fmt.Sprintf("non-existent-%d.txt", time.Now().Unix())),
		})
		if err == nil {
			t.Error("Expected error when heading non-existent object")
		} else {
			t.Logf("✓ Got expected error: %v", err)
		}
	})
}

// BenchmarkOSSClientPutObject benchmarks object upload
func BenchmarkOSSClientPutObject(b *testing.B) {
	client, _ := testConfig.Client()
	ossClient, _ := client.OSSClient()
	ctx := context.Background()

	bucketName := generateTestBucketName("oss-bench")
	content := "Benchmark test content for object upload performance testing"

	// Create bucket
	_ = ossClient.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})

	defer cleanupBucket(ctx, ossClient, bucketName)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		objectKey := fmt.Sprintf("bench-object-%d.txt", i)
		_ = ossClient.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
			Body:   strings.NewReader(content),
		})
	}
}

// BenchmarkOSSClientGetObject benchmarks object download
func BenchmarkOSSClientGetObject(b *testing.B) {
	client, _ := testConfig.Client()
	ossClient, _ := client.OSSClient()
	ctx := context.Background()

	bucketName := generateTestBucketName("oss-bench-get")
	objectKey := "benchmark-object.txt"
	content := "Benchmark test content for get operations"

	// Setup: Create bucket and upload object
	_ = ossClient.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})
	_ = ossClient.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   strings.NewReader(content),
	})

	defer cleanupBucket(ctx, ossClient, bucketName)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		output, _ := ossClient.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
		})
		if output != nil && output.Body != nil {
			output.Body.Close()
		}
	}
}

// BenchmarkOSSClientListObjects benchmarks object listing
func BenchmarkOSSClientListObjects(b *testing.B) {
	client, _ := testConfig.Client()
	ossClient, _ := client.OSSClient()
	ctx := context.Background()

	bucketName := generateTestBucketName("oss-bench-list")

	// Setup: Create bucket and upload objects
	_ = ossClient.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})

	for i := 0; i < 50; i++ {
		_ = ossClient.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(fmt.Sprintf("object-%d.txt", i)),
			Body:   strings.NewReader("test content"),
		})
	}

	defer cleanupBucket(ctx, ossClient, bucketName)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = ossClient.ListObjects(ctx, &s3.ListObjectsInput{
			Bucket: aws.String(bucketName),
		})
	}
}

// BenchmarkOSSClientConcurrentPutObject benchmarks concurrent uploads
func BenchmarkOSSClientConcurrentPutObject(b *testing.B) {
	client, _ := testConfig.Client()
	ossClient, _ := client.OSSClient()
	ctx := context.Background()

	bucketName := generateTestBucketName("oss-bench-concurrent")
	content := "Concurrent benchmark test content"

	// Setup
	_ = ossClient.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})

	defer cleanupBucket(ctx, ossClient, bucketName)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			objectKey := fmt.Sprintf("concurrent-bench-%d.txt", i)
			_ = ossClient.PutObject(ctx, &s3.PutObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(objectKey),
				Body:   strings.NewReader(content),
			})
			i++
		}
	})
}
