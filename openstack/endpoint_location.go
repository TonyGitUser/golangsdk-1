package openstack

import (
	"github.com/gophercloud/gophercloud"
	tokens2 "github.com/gophercloud/gophercloud/openstack/identity/v2/tokens"
	tokens3 "github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
	"errors"
	"strings"
)

// service have same endpoint address in different location, refer to https://developer.huaweicloud.com/endpoint
var endpointSchemaList = map[string]string{
	"COMPUTE": "https://ecs.%(region)s.%(domain)s/v2/%(projectId)s/",
	"IMAGE": "https://ims.%(region)s.%(domain)s/",
	"NETWORK": "https://vpc.%(region)s.%(domain)s/",
	"BLOCK_STORAGE": "https://evs.%(region)s.%(domain)s/",
	"ANTIDDOS": "https://antiddos.%(region)s.%(domain)s/",
}


/*
V2EndpointURL discovers the endpoint URL for a specific service from a
ServiceCatalog acquired during the v2 identity service.

The specified EndpointOpts are used to identify a unique, unambiguous endpoint
to return. It's an error both when multiple endpoints match the provided
criteria and when none do. The minimum that can be specified is a Type, but you
will also often need to specify a Name and/or a Region depending on what's
available on your OpenStack deployment.
*/
func V2EndpointURL(catalog *tokens2.ServiceCatalog, opts gophercloud.EndpointOpts) (string, error) {
	// Extract Endpoints from the catalog entries that match the requested Type, Name if provided, and Region if provided.
	var endpoints = make([]tokens2.Endpoint, 0, 1)
	for _, entry := range catalog.Entries {
		if (entry.Type == opts.Type) && (opts.Name == "" || entry.Name == opts.Name) {
			for _, endpoint := range entry.Endpoints {
				if opts.Region == "" || endpoint.Region == opts.Region {
					endpoints = append(endpoints, endpoint)
				}
			}
		}
	}

	// Report an error if the options were ambiguous.
	if len(endpoints) > 1 {
		err := &ErrMultipleMatchingEndpointsV2{}
		err.Endpoints = endpoints
		return "", err
	}

	// Extract the appropriate URL from the matching Endpoint.
	for _, endpoint := range endpoints {
		switch opts.Availability {
		case gophercloud.AvailabilityPublic:
			return gophercloud.NormalizeURL(endpoint.PublicURL), nil
		case gophercloud.AvailabilityInternal:
			return gophercloud.NormalizeURL(endpoint.InternalURL), nil
		case gophercloud.AvailabilityAdmin:
			return gophercloud.NormalizeURL(endpoint.AdminURL), nil
		default:
			err := &ErrInvalidAvailabilityProvided{}
			err.Argument = "Availability"
			err.Value = opts.Availability
			return "", err
		}
	}

	// Report an error if there were no matching endpoints.
	err := &gophercloud.ErrEndpointNotFound{}
	return "", err
}
/*
   GetEndpointURLForAKSKAuth discovers the endpoint  from V3EndpointURL function firstly,
   if the endpoint is null then concat the service type and domain as the endpoint
 */
func GetEndpointURLForAKSKAuth(catalog *tokens3.ServiceCatalog, opts gophercloud.EndpointOpts,akskOptions gophercloud.AKSKAuthOptions) (string, error) {
	if opts.Region == ""{
		opts.Region = akskOptions.Region
	}

	endpoint, err := V3EndpointURL(catalog,opts)

   if err != nil || endpoint == ""{
	   if akskOptions.Domain == ""{
		   return "", errors.New("ServiceDomainName can not be empty.")
	   }

	   if endpointSchema, ok := endpointSchemaList[strings.ToUpper(opts.Type)]; ok {
		   endpoint = strings.Replace(endpointSchema,"%(domain)s",akskOptions.Domain,1)
		   endpoint = strings.Replace(endpoint,"%(region)s",opts.Region,1)
		   endpoint = strings.Replace(endpoint,"%(projectId)s",akskOptions.ProjectId,1)

		   return endpoint,nil
	   }
   }

	return endpoint,err
}

/*
V3EndpointURL discovers the endpoint URL for a specific service from a Catalog
acquired during the v3 identity service.

The specified EndpointOpts are used to identify a unique, unambiguous endpoint
to return. It's an error both when multiple endpoints match the provided
criteria and when none do. The minimum that can be specified is a Type, but you
will also often need to specify a Name and/or a Region depending on what's
available on your OpenStack deployment.
*/
func V3EndpointURL(catalog *tokens3.ServiceCatalog, opts gophercloud.EndpointOpts) (string, error) {
	// Extract Endpoints from the catalog entries that match the requested Type, Interface,
	// Name if provided, and Region if provided.
	var endpoints = make([]tokens3.Endpoint, 0, 1)
	for _, entry := range catalog.Entries {
		if (entry.Type == opts.Type) && (opts.Name == "" || entry.Name == opts.Name) {
			for _, endpoint := range entry.Endpoints {
				if opts.Availability != gophercloud.AvailabilityAdmin &&
					opts.Availability != gophercloud.AvailabilityPublic &&
					opts.Availability != gophercloud.AvailabilityInternal {
					err := &ErrInvalidAvailabilityProvided{}
					err.Argument = "Availability"
					err.Value = opts.Availability
					return "", err
				}
				if (opts.Availability == gophercloud.Availability(endpoint.Interface)) &&
					(opts.Region == "" || endpoint.Region == opts.Region) {
					endpoints = append(endpoints, endpoint)
				}
			}
		}
	}

	// Report an error if the options were ambiguous.
	if len(endpoints) > 1 {
		return "", ErrMultipleMatchingEndpointsV3{Endpoints: endpoints}
	}

	// Extract the URL from the matching Endpoint.
	for _, endpoint := range endpoints {
		return gophercloud.NormalizeURL(endpoint.URL), nil
	}

	// Report an error if there were no matching endpoints.
	err := &gophercloud.ErrEndpointNotFound{}
	return "", err
}