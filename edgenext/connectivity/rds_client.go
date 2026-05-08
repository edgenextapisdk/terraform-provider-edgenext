package connectivity

// RDSClient is an alias for ECSClient: same EdgeNext signing, headers, Get/Post, and region handling.
// Use this when calling RDS OpenAPI paths (for example /rds/openapi/...) with the same auth model as ECS.
type RDSClient = ECSClient

// NewRDSClient creates a new RDS API client. It is equivalent to NewECSClient with the same arguments.
func NewRDSClient(accessKey, secretKey, endpoint, region string) *RDSClient {
	return newServiceClient("RDS", accessKey, secretKey, endpoint, region)
}
