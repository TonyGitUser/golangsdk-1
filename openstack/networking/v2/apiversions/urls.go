package apiversions

import (
	"strings"

	"github.com/huaweicloudsdk/golangsdk"
)

func apiVersionsURL(c *golangsdk.ServiceClient) string {
	return c.Endpoint
}

func apiInfoURL(c *golangsdk.ServiceClient, version string) string {
	return c.Endpoint + strings.TrimRight(version, "/") + "/"
}
