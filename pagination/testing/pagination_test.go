package testing

import (
	"github.com/huaweicloudsdk/golangsdk"
	"github.com/huaweicloudsdk/golangsdk/testhelper"
)

func createClient() *golangsdk.ServiceClient {
	return &golangsdk.ServiceClient{
		ProviderClient: &golangsdk.ProviderClient{TokenID: "abc123"},
		Endpoint:       testhelper.Endpoint(),
	}
}
