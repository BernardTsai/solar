package common

import (
  "gopkg.in/yaml.v2"
)
//------------------------------------------------------------------------------

// TenantConfiguration holds the information required to address an openstack API endpoint
type TenantConfiguration struct {
  Tenant   string `yaml:"Tenant"`    // name of the tenant
  Username string `yaml:"Username"`  // username credential
  Password string `yaml:"Password"`  // password credential
  URL      string `yaml:"URL"`       // URL of the API endpoint
  Proxy    string `yaml:"Proxy"`     // proxy address is needed
}

//------------------------------------------------------------------------------

// TenantEndpoint holds the endpoint information of a tenant
type TenantEndpoint struct {
  Tenant   string `yaml:"Tenant"`    // name of the tenant
  UUID     string `yaml:"UUID"`      // unique universal identifier
  Token    string `yaml:"Token"`     // authentication token
  Identity string `yaml:"Identity"`  // URL of identity API
  Compute  string `yaml:"Compute"`   // URL of compute API
  Image    string `yaml:"Image"`     // URL of image API
  Volume   string `yaml:"Volume"`    // URL of volume API
  Network  string `yaml:"Network"`   // URL of network API
  Username string `yaml:"Username"`  // username credential
  Password string `yaml:"Password"`  // password credential
  URL      string `yaml:"URL"`       // URL of the API endpoint
  Proxy    string `yaml:"Proxy"`     // proxy address is needed
}

//------------------------------------------------------------------------------

// GetEndpointTokenAndProxy determines a specific endpoint from a tenant endpoint
func GetEndpointTokenAndProxy(tenantEndpointString string, service string) (string, string, string) {
  // convert yaml string into TenantEndpoint object
  tenantEndpoint := TenantEndpoint{}
  err := yaml.Unmarshal([]byte(tenantEndpointString), &tenantEndpoint)
  if err != nil {
    return "", "", ""
  }

  // determine endpoint
  switch service {
  case "identity":
      return tenantEndpoint.Identity, tenantEndpoint.Token, tenantEndpoint.Proxy
  case "compute":
      return tenantEndpoint.Compute, tenantEndpoint.Token, tenantEndpoint.Proxy
  case "image":
      return tenantEndpoint.Image, tenantEndpoint.Token, tenantEndpoint.Proxy
  case "volume":
      return tenantEndpoint.Volume, tenantEndpoint.Token, tenantEndpoint.Proxy
  case "network":
      return tenantEndpoint.Network, tenantEndpoint.Token, tenantEndpoint.Proxy
  }

  // nothing found
  return "", "", ""
}

//------------------------------------------------------------------------------
