package network

import (
  "net/http"
  "io/ioutil"
  "gopkg.in/yaml.v2"
  "encoding/json"

  "tsai.eu/solar/controller/openstackController/common"
)

//------------------------------------------------------------------------------

// Sample request:
//
// Request: 5cf2625e-abcd-473a-9b8c-5838b014dab4
// Domain: demo
// Solution: app
// Version: V1.0.0
// Element: oam
// Cluster: V1.0.0
// Instance: 5cf2625e-6fb0-473a-9b8c-5838b0140ea3
// Component: os-network
// State: active
// Configuration: |
//   Name:         oam
//   IPv4CIDR:     10.0.1.0/24
//   IPv4Gateway:  10.0.1.1
//   IPv4Start:    10.0.1.100
//   IPv4Stop:     10.0.1.199
// Relationships:
//   - Relationship:  tenant
//     Dependency:    tenant
//     Configuration: ""
//     Endpoint: |-
//       Tenant:   "demo"
//       UUID: 15d3656801644ed38507de7ee8d43733
//       Token: 33fa80123e594cd881965d73ac694600
//       Identity: https://192.168.178.20/identity/v2.0
//       Compute: https://192.168.178.20/compute/v2/15d3656801644ed38507de7ee8d43733
//       Image: https://192.168.178.20/image
//       Volume: https://192.168.178.20/volumev2/v1/15d3656801644ed38507de7ee8d43733
//       Network: https://192.168.178.20/network
//       Username: "admin"
//       Password: "secret"
//       URL:      "https://192.168.178.20"
//       Proxy:    ""
//
// Sample response:
//
// Request: 5cf2625e-abcd-473a-9b8c-5838b014dab4
// Code: 200
// Status: ""
// Domain: demo
// Solution: app
// Version: V1.0.0
// Element: oam
// Cluster: V1.0.0
// Instance: 5cf2625e-6fb0-473a-9b8c-5838b0140ea3
// Component: os-network
// State: active
// Configuration: |
//   Name:         oam
//   IPv4CIDR:     10.0.1.0/24
//   IPv4Gateway:  10.0.1.1
//   IPv4Start:    10.0.1.100
//   IPv4Stop:     10.0.1.199
// Endpoint: |-
//   Name:         oam
//   Subnet:       oam_subnet
//   UUID1:        abc
//   UUID2:        def
//   IPv4CIDR:     10.0.1.0/24
//   IPv4Gateway:  10.0.1.1
//   IPv4Start:    10.0.1.100
//   IPv4Stop:     10.0.1.199

//------------------------------------------------------------------------------

// ResponseNetwork is part of the response
type ResponseNetwork struct {
  Name    string   `yaml:"name"`
  ID      string   `yaml:"id"`
  Status  string   `yaml:"status"`
  Subnets []string `yaml:"subnets"`
}

// ResponseNetworkData is part of the response
type ResponseNetworkData struct {
  Networks []ResponseNetwork `yaml:"networks"`
}

// ResponsePool is part of the response
type ResponsePool struct {
  Start    string       `json:"start"`
  End      string       `json:"end"`
}

// ResponseSubnet is part of the response
type ResponseSubnet struct {
  Name      string       `json:"name"`
  ID        string       `json:"id"`
  Version   int          `json:"ip_version"`
  CIDR      string       `json:"cidr"`
  Gateway   string       `json:"gateway_ip"`
  Pools   []ResponsePool `json:"allocation_pools" json:"allocation_pools"`
}

// ResponseSubnetData is part of the response
type ResponseSubnetData struct {
  Subnets []ResponseSubnet `yaml:"subnets"`
}

// GetSubnetID retrieves a subnet ID from a ResponseNetworkData
func (responseNetworkData * ResponseNetworkData) GetSubnetID() string {
  for _, networkData := range responseNetworkData.Networks {
    for _, subnetData := range networkData.Subnets {
      return subnetData
    }
  }

  // found nothing
  return ""
}

//------------------------------------------------------------------------------

// status determines the status of a network
func status(request *common.Request, response *common.Response) {
  // determine configuration
  config := &common.NetworkConfiguration{}
  err := yaml.Unmarshal([]byte(request.Configuration), &config )
  if err != nil {
    response.SetStatus(http.StatusBadRequest, "Unable to parse configuration:\n" + err.Error())
    return
  }

  // get network information
  resp, err := common.DoGet(request, "network", "/v2.0/networks?name=" + config.Name)
	if err != nil {
    response.SetStatus(http.StatusBadRequest, "Unable to get network information:\n" + err.Error())
    return
	}

  // read response body and check for status code
  respBody, _ := ioutil.ReadAll(resp.Body)
  if resp.StatusCode < 200 || 300 <= resp.StatusCode {
    response.State = common.InitialState
    response.SetStatus(resp.StatusCode, "Network information error:\n" + string(respBody))
    return
  }

  // parse result
  responseNetworkData := ResponseNetworkData{}
  err = json.Unmarshal(respBody, &responseNetworkData)
  if err != nil {
    response.SetStatus(http.StatusInternalServerError, "Unable to parse network information response:\n" + string(respBody) + "\n" + err.Error())
    return
  }

  // validate subnet ID
  subnetID := responseNetworkData.GetSubnetID()
  if subnetID == "" {
    response.State = common.InitialState
    response.SetStatus(http.StatusOK, "")
    return
  }

  // check status of subnet
  resp, err = common.DoGet(request, "network", "/v2.0/subnets?id=" + subnetID)
	if err != nil {
    response.SetStatus(http.StatusInternalServerError, "Unable to get subnet information:\n" + err.Error())
    return
	}

  // read response body and check for status code
  respBody, _ = ioutil.ReadAll(resp.Body)
  if resp.StatusCode < 200 || 300 <= resp.StatusCode {
    response.State = common.InitialState
    response.SetStatus(resp.StatusCode, "Subnet information error:\n" + string(respBody))
    return
  }

  // parse result
  responseSubnetData := ResponseSubnetData{}
  err = json.Unmarshal(respBody, &responseSubnetData)
  if err != nil {
    response.SetStatus(http.StatusInternalServerError, "Unable to parse subnet information response:\n" + string(respBody) + "\n" + err.Error())
    return
  }

  // construct endpoint
  networkEndpoint := common.NetworkEndpoint{
    Name:         responseNetworkData.Networks[0].Name,
    UUID1:        responseNetworkData.Networks[0].ID,
    UUID2:        responseSubnetData.Subnets[0].ID,
    IPv4CIDR:     responseSubnetData.Subnets[0].CIDR,
    IPv4Gateway:  responseSubnetData.Subnets[0].Gateway,
    IPv4Start:    responseSubnetData.Subnets[0].Pools[0].Start,
    IPv4Stop:     responseSubnetData.Subnets[0].Pools[0].End,
  }
  endpointYaml, _ := yaml.Marshal(networkEndpoint)

  // update state and endpoint
  response.State    = common.ActiveState
  response.Endpoint = string(endpointYaml)
  response.SetStatus(http.StatusOK, "")
}

//------------------------------------------------------------------------------
