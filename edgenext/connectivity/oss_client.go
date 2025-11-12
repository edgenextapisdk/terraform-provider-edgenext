package connectivity

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// OSSClient represents OSS S3 client
type OSSClient struct {
	endpoint string
	region   string
	client   *s3.Client
}

// NewOSSClient creates a new OSS S3 client
func NewOSSClient(accessKey, secretKey, endpoint, region string) (*OSSClient, error) {

	if region == "" {
		region = "us-east-1"
	}

	// Create AWS credentials
	creds := credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")

	// Create S3 client configuration
	cfg := aws.Config{
		Region:      region,
		Credentials: creds,
	}

	// Create S3 client
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true // Use path-style addressing for compatibility
	})

	return &OSSClient{
		endpoint: endpoint,
		region:   region,
		client:   client,
	}, nil
}

// GetRegion returns the OSS region
func (c *OSSClient) GetRegion() string {
	return c.region
}

func (c *OSSClient) CreateBucket(ctx context.Context, input *s3.CreateBucketInput) error {
	_, err := c.client.CreateBucket(ctx, input)
	return err
}

func (c *OSSClient) PutBucketCors(ctx context.Context, input *s3.PutBucketCorsInput) error {
	_, err := c.client.PutBucketCors(ctx, input)
	return err
}

func (c *OSSClient) DeleteBucket(ctx context.Context, input *s3.DeleteBucketInput) error {
	_, err := c.client.DeleteBucket(ctx, input)
	return err
}

func (c *OSSClient) BucketExists(ctx context.Context, input *s3.HeadBucketInput) (bool, error) {
	_, err := c.client.HeadBucket(ctx, input)
	return err == nil, err
}

func (c *OSSClient) ListBuckets(ctx context.Context) (*s3.ListBucketsOutput, error) {
	return c.client.ListBuckets(ctx, &s3.ListBucketsInput{})
}

func (c *OSSClient) PutBucketAcl(ctx context.Context, input *s3.PutBucketAclInput) error {
	_, err := c.client.PutBucketAcl(ctx, input)
	return err
}

func (c *OSSClient) GetBucketAcl(ctx context.Context, input *s3.GetBucketAclInput) (*s3.GetBucketAclOutput, error) {
	return c.client.GetBucketAcl(ctx, input)
}

func (c *OSSClient) HeadBucket(ctx context.Context, input *s3.HeadBucketInput) (*s3.HeadBucketOutput, error) {
	output, err := c.client.HeadBucket(ctx, input)
	return output, err
}

func (c *OSSClient) PutObject(ctx context.Context, input *s3.PutObjectInput) error {
	_, err := c.client.PutObject(ctx, input)
	return err
}

func (c *OSSClient) GetObject(ctx context.Context, input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	return c.client.GetObject(ctx, input)
}

func (c *OSSClient) DeleteObject(ctx context.Context, input *s3.DeleteObjectInput) error {
	_, err := c.client.DeleteObject(ctx, input)
	return err
}

func (c *OSSClient) ListObjects(ctx context.Context, input *s3.ListObjectsInput) (*s3.ListObjectsOutput, error) {
	return c.client.ListObjects(ctx, input)
}

func (c *OSSClient) ListObjectsV2(ctx context.Context, input *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
	return c.client.ListObjectsV2(ctx, input)
}

func (c *OSSClient) HeadObject(ctx context.Context, input *s3.HeadObjectInput) (*s3.HeadObjectOutput, error) {
	return c.client.HeadObject(ctx, input)
}

func (c *OSSClient) PutObjectAcl(ctx context.Context, input *s3.PutObjectAclInput) error {
	_, err := c.client.PutObjectAcl(ctx, input)
	return err
}

func (c *OSSClient) GetObjectAcl(ctx context.Context, input *s3.GetObjectAclInput) (*s3.GetObjectAclOutput, error) {
	return c.client.GetObjectAcl(ctx, input)
}

func (c *OSSClient) PresignObject(ctx context.Context, input *s3.GetObjectInput, expiresIn int64) (string, error) {
	presignClient := s3.NewPresignClient(c.client)
	presignedRequest, err := presignClient.PresignGetObject(ctx, input, func(o *s3.PresignOptions) {
		o.Expires = time.Duration(expiresIn) * time.Second
	})
	if err != nil {
		return "", err
	}
	url := presignedRequest.URL
	return url, nil
}

func (c *OSSClient) CopyObject(ctx context.Context, input *s3.CopyObjectInput) (*s3.CopyObjectOutput, error) {
	return c.client.CopyObject(ctx, input)
}
