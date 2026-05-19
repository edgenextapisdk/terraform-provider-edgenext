package connectivity

// ELBClient is an alias for ECSClient: same EdgeNext signing, headers, Get/Post, and region handling.
// Use this when calling ELB OpenAPI paths (for example /elb/openapi/...) with the same auth model as ECS.
type ELBClient = ECSClient

// NewELBClient creates a new ELB API client. It is equivalent to NewECSClient with the same arguments.
func NewELBClient(accessKey, secretKey, endpoint, region string) *ELBClient {
	return newServiceClient("ELB", accessKey, secretKey, endpoint, region)
}
