package gophercloud

// AKSKAuthOptions presents the required information for AK/SK auth
type AKSKAuthOptions struct {
	// IdentityEndpoint specifies the HTTP endpoint that is required to work with
	// the Identity API of the appropriate version. While it's ultimately needed by
	// all of the identity services, it will often be populated by a provider-level
	// function.
	//
	// The IdentityEndpoint is typically referred to as the "auth_url" or
	// "OS_AUTH_URL" in the information provided by the cloud operator.
	IdentityEndpoint string `json:"-"`

	// user project id
	ProjectId string

	// region
	Region string

	// cloud service domain, example: myhwclouds.com
	Domain string

	AccessKey           string //Access Key
	SecretKey           string //Secret key
}

// Implements the method of AuthOptionsProvider
func (opts AKSKAuthOptions) GetIdentityEndpoint() string {
	return opts.IdentityEndpoint
}