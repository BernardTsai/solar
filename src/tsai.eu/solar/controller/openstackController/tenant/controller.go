package tenant

import (
  "fmt"
  "time"
  "bytes"
  "net/http"
	"net/url"
  "crypto/tls"
  "gopkg.in/yaml.v2"
  "io/ioutil"
  "encoding/json"

  "tsai.eu/solar/controller/openstackController/common"
)

//------------------------------------------------------------------------------

// Sample request:
//
// Request:   5cf2625e-abcd-473a-9b8c-5838b014dab4
// Domain:    demo
// Solution:  app
// Version:   V1.0.0
// Element:   tenant
// Cluster:   V1.0.0
// Instance:  5cf2625e-6fb0-473a-9b8c-5838b0140ea3
// Component: os-tenant
// State:     active
// Configuration: |
//   Tenant:   "demo"
//   Username: "admin"
//   Password: "secret"
//   URL:      "https://192.168.178.20"
//   Proxy:    ""

// Sample response:
//
// Request:   5cf2625e-abcd-473a-9b8c-5838b014dab4
// Action:    status
// Code:      200
// Status:    ""
// Domain:    demo
// Solution:  app
// Version:   V1.0.0
// Element:   tenant
// Cluster:   V1.0.0
// Instance:  5cf2625e-6fb0-473a-9b8c-5838b0140ea3
// State:     active
// Configuration: \-
//   Tenant:   "demo"
//   Username: "admin"
//   Password: "secret"
//   URL:      "https://192.168.178.20"
//   Proxy:    ""
// Endpoint: |-
//   Tenant:   "demo"
//   UUID: 15d3656801644ed38507de7ee8d43733
//   Token: 33fa80123e594cd881965d73ac694600
//   Identity: https://192.168.178.20/identity/v2.0
//   Compute: https://192.168.178.20/compute/v2/15d3656801644ed38507de7ee8d43733
//   Image: https://192.168.178.20/image
//   Volume: https://192.168.178.20/volumev2/v1/15d3656801644ed38507de7ee8d43733
//   Network: https://192.168.178.20/network
//   Username: "admin"
//   Password: "secret"
//   URL:      "https://192.168.178.20"
//   Proxy:    ""
//------------------------------------------------------------------------------

// RequestBody holds the json template for a request
var RequestBody = `{
"auth": {
  "identity": {
    "methods": ["password"],
    "password": {
      "user": {
        "name": "%s",
        "domain": { "id": "default" },
        "password": "%s"
      }
    }
  },
  "scope": {
    "project": {
      "name": "%s",
      "domain": { "id": "default" }
    }
  }
}
}`

//------------------------------------------------------------------------------

// ResponseEndpoint is part of the response
type ResponseEndpoint struct {
  URL       string `yaml:"url"`
  Interface string `yaml:"interface"`
  Region    string `yaml:"region"`
  RegionID  string `yaml:"region_id"`
  ID        string `yaml:"id"`
}

// ResponseCatalog is part of the response
type ResponseCatalog struct {
  Endpoints []ResponseEndpoint `yaml:"endpoints"`
  Type      string             `yaml:"type"`
  ID        string             `yaml:"id"`
  Name      string             `yaml:"name"`
}

// ResponseProject is part of the response
type ResponseProject struct {
  ID string `yaml:"id"`
}

// ResponseToken is part of the response
type ResponseToken struct {
  Project ResponseProject   `yaml:"project"`
  Catalog []ResponseCatalog `yaml:"catalog"`
}

// ResponseData is part of the response
type ResponseData struct {
  Token ResponseToken `yaml:"token"`
}

func (responseData *ResponseData) getEndpoint(name string) string {
  for _, catalogEntry := range responseData.Token.Catalog {
    if catalogEntry.Type == name {
      for _, endpoint := range catalogEntry.Endpoints {
        if endpoint.Interface == "public" {
          return endpoint.URL
        }
      }
    }

  }

  // not found
  return ""
}

//------------------------------------------------------------------------------

// Process handles all incoming kubernetes cluster related requests
func Process(request *common.Request, response *common.Response) {
  // delegate to requested action
  switch response.Action {
  case "create":
    status(request, response)
  case "destroy":
    status(request, response)
  case "start":
    status(request, response)
  case "stop":
    status(request, response)
  case "configure":
    status(request, response)
  case "reconfigure":
    status(request, response)
  case "reset":
    status(request, response)
  case "status":
    status(request, response)
  }
}

//------------------------------------------------------------------------------

// status determines the status of an openstack tenant
func status(request *common.Request, response *common.Response) {
  var transport *http.Transport

  // determine configuration
  config := &common.TenantConfiguration{}
  err := yaml.Unmarshal([]byte(request.Configuration), &config )
  if err != nil {
    response.Configuration = request.Configuration
    response.Endpoint      = ""
    response.State         = common.UndefinedState
    response.Code          = http.StatusBadRequest
    response.Status        = "Unable to parse configuration:\n" + err.Error()
    return
  }

  // define transport parameters
  if config.Proxy != "" {
  	proxyURL, err := url.Parse(config.Proxy)
  	if err != nil {
      response.Configuration = request.Configuration
      response.Endpoint      = ""
      response.State         = common.UndefinedState
      response.Code          = http.StatusBadRequest
      response.Status        = "Invalid proxy setting"
      return
  	}

    transport = &http.Transport{
      TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
      Proxy:           http.ProxyURL(proxyURL),

    }
  } else {
    transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true},}
  }

  // request for a token
  timeout   := time.Duration(10 * time.Second)
  client    := &http.Client{Transport: transport, Timeout: timeout,}
  jsonBody  := fmt.Sprintf(RequestBody, config.Username, config.Password, config.Tenant )
  body      := []byte(jsonBody)

  req, _    := http.NewRequest("POST", config.URL + "/auth/tokens", bytes.NewBuffer(body))
  req.Header.Set("Content-Type", "application/json")

  resp, err := client.Do(req)
	if err != nil {
    response.Configuration = request.Configuration
    response.Endpoint      = ""
    response.State         = common.UndefinedState
    response.Code          = http.StatusBadRequest
    response.Status        = "Unable to access tenant:\n" + err.Error()
    return
	}

  // read response body
  respBody, err := ioutil.ReadAll(resp.Body)

  if err != nil {
    response.Configuration = request.Configuration
    response.Endpoint      = ""
    response.State         = common.UndefinedState
    response.Code          = http.StatusBadRequest
    response.Status        = "Strange response:\n" + err.Error()
    return
  }

  // check for status code
  if resp.StatusCode < 200 || 300 <= resp.StatusCode {
    response.Configuration = request.Configuration
    response.Endpoint      = ""
    response.State         = common.InitialState
    response.Code          = resp.StatusCode
    response.Status        = "Error:\n" + string(respBody)
    return
  }

  // construct endpoint
  responseData := ResponseData{}
  json.Unmarshal(respBody, &responseData)

  endpoint := common.TenantEndpoint{
    Tenant:   config.Tenant,
    UUID:     responseData.Token.Project.ID,
    Token:    resp.Header.Get("X-Subject-Token"),
    Identity: responseData.getEndpoint("identity"),
    Compute:  responseData.getEndpoint("compute"),
    Image:    responseData.getEndpoint("image"),
    Volume:   responseData.getEndpoint("volume"),
    Network:  responseData.getEndpoint("network"),
    Username: config.Username,
    Password: config.Password,
    URL:      config.URL,
    Proxy:    config.Proxy,
  }
  endpointYaml, _ := yaml.Marshal(endpoint)

  // update state and endpoint
  response.Configuration = request.Configuration
  response.State         = common.ActiveState
  response.Endpoint      = string(endpointYaml)
  response.Code          = http.StatusOK
  response.Status        = ""
}

//------------------------------------------------------------------------------
