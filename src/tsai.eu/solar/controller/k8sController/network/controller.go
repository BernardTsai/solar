package network

import (
  "errors"
  "strings"
  "io/ioutil"

  "net/http"
  "crypto/tls"
  "gopkg.in/yaml.v2"
  "github.com/cbroglie/mustache"

  "tsai.eu/solar/controller/k8sController/common"
)

//------------------------------------------------------------------------------

// template resembles an k8s bridge configuration
const template = `
apiVersion: "k8s.cni.cncf.io/v1"
kind: NetworkAttachmentDefinition
metadata:
  name: {{Name}}
spec:
  config: '{
    "cniVersion": "0.3.0",
    "name": "{{Name}}",
    "type": "bridge",
    "bridge": "{{Name}}",
    "isDefaultGateway": true,
    "forceAddress": false,
    "ipMasq": true,
    "hairpinMode": true,
    "ipam": {
      "type": "host-local",
      "subnet": "{{Subnet}}",
      "routes": [{"dst": "0.0.0.0/0"}]
    }
  }'`


//   "name": "{{Name}}",
//  "type": "macvlan",
//  "master": "eth0",
//  "mode": "bridge",
//  "ipam": {
//    "type": "host-local",
//    "subnet": "192.168.1.0/24",
//    "rangeStart": "192.168.1.200",
//    "rangeEnd": "192.168.1.216",
//    "routes": [
//      { "dst": "0.0.0.0/0" }
//    ],
//    "gateway": "192.168.1.1"
//  }

//   "name": "{{Name}}",
//   "type": "bridge",
//   "bridge": "{{Name}}",
//   "isDefaultGateway": true,
//   "forceAddress": false,
//   "ipMasq": true,
//   "hairpinMode": true,
//   "ipam": {
//     "type": "host-local",
//     "subnet": "{{Subnet}}",
//     "routes": [{"dst": "0.0.0.0/0"}]
//   }


//------------------------------------------------------------------------------

// Configuration holds the information required to address a k8s API endpoint
type Configuration struct {
  URL    string `yaml:"URL"`     // URL of API endpoint
  Token  string `yaml:"Token"`   // authorization bearer token
  Name   string `yaml:"Name"`    // URL of API endpoint
  Subnet string `yaml:"Subnet"`  // authorization bearer token
}

// Sample request:
//
// Request: 5cf2625e-abcd-473a-9b8c-5838b014dab4
// Domain: demo
// Solution: app
// Version: V1.0.0
// Element: oam
// Cluster: V1.0.0
// Instance: 5cf2625e-6fb0-473a-9b8c-5838b0140ea3
// Component: k8s-network
// State: active
// Configuration: |
//   URL: https://192.168.178.20:6443
//   Token: eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImRlbW8tYWRtaW4tdG9rZW4tbHAyY3EiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoiZGVtby1hZG1pbiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6IjA3ODJkZmYyLTY5ZjEtMTFlOS1hNzNkLWYwZGVmMWZjN2FkNSIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZWZhdWx0OmRlbW8tYWRtaW4ifQ.Z6WIRuaqNKhTG2u9cx4eXd8Ado6j1bNwN8JiZFTF9hPLC90i820auALiL9RtHBv-lOo_auSZKJtF9eiLFRxyoon7RsUSie-pZtKHupTKb3YbggEC9ZZxt-NGyKhMhNLRtJxri-FSkSzJqQkkXjZ_qX0CQSgi22-cIbH-ki99C2RvvySkXSJSXeQOhBVDyqSL1Skj042KITlMne20G4AxMB5HSjUsrWdYK8dGCdHW0ohxIRKW1nLvySMTk3G4Em-7Sy-Ug_bVo3KK_ObycmICHaYcbWchHj_6NHYh3ISx0yMpoSsLS3Q3Z6qDABGNVFidMXcvANHNaWtzeiLXNKeVAA
//   Name: "oam"
//   Subnet: "10.10.0.0/16"

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
// Component: k8s-network
// State: active
// Configuration: |
//   URL: https://192.168.178.20:6443
//   Token: eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImRlbW8tYWRtaW4tdG9rZW4tbHAyY3EiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoiZGVtby1hZG1pbiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6IjA3ODJkZmYyLTY5ZjEtMTFlOS1hNzNkLWYwZGVmMWZjN2FkNSIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZWZhdWx0OmRlbW8tYWRtaW4ifQ.Z6WIRuaqNKhTG2u9cx4eXd8Ado6j1bNwN8JiZFTF9hPLC90i820auALiL9RtHBv-lOo_auSZKJtF9eiLFRxyoon7RsUSie-pZtKHupTKb3YbggEC9ZZxt-NGyKhMhNLRtJxri-FSkSzJqQkkXjZ_qX0CQSgi22-cIbH-ki99C2RvvySkXSJSXeQOhBVDyqSL1Skj042KITlMne20G4AxMB5HSjUsrWdYK8dGCdHW0ohxIRKW1nLvySMTk3G4Em-7Sy-Ug_bVo3KK_ObycmICHaYcbWchHj_6NHYh3ISx0yMpoSsLS3Q3Z6qDABGNVFidMXcvANHNaWtzeiLXNKeVAA
//   Name: "oam"
//   Subnet: "10.10.0.0/16"
// Endpoint: |-
//   URL: https://192.168.178.20:6443
//   Token: eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImRlbW8tYWRtaW4tdG9rZW4tbHAyY3EiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoiZGVtby1hZG1pbiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6IjA3ODJkZmYyLTY5ZjEtMTFlOS1hNzNkLWYwZGVmMWZjN2FkNSIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZWZhdWx0OmRlbW8tYWRtaW4ifQ.Z6WIRuaqNKhTG2u9cx4eXd8Ado6j1bNwN8JiZFTF9hPLC90i820auALiL9RtHBv-lOo_auSZKJtF9eiLFRxyoon7RsUSie-pZtKHupTKb3YbggEC9ZZxt-NGyKhMhNLRtJxri-FSkSzJqQkkXjZ_qX0CQSgi22-cIbH-ki99C2RvvySkXSJSXeQOhBVDyqSL1Skj042KITlMne20G4AxMB5HSjUsrWdYK8dGCdHW0ohxIRKW1nLvySMTk3G4Em-7Sy-Ug_bVo3KK_ObycmICHaYcbWchHj_6NHYh3ISx0yMpoSsLS3Q3Z6qDABGNVFidMXcvANHNaWtzeiLXNKeVAA
//   Name: "oam"
//   Subnet: "10.10.0.0/16"
//------------------------------------------------------------------------------

// Process handles all incoming namespace related requests
func Process(request *common.Request, response *common.Response) {
  // delegate to requested action
  switch response.Action {
  case "create":
    create(request, response)
  case "destroy":
    destroy(request, response)
  case "start":
    create(request, response)
  case "stop":
    destroy(request, response)
  case "configure":
    destroy(request, response)
  case "reconfigure":
    reconfigure(request, response)
  case "reset":
    destroy(request, response)
  case "status":
    status(request, response)
  }
}

//------------------------------------------------------------------------------

// create instantiates a network
func create(request *common.Request, response *common.Response) {
  // check if the network has been created already
  status(request, response)
  if response.State == common.ActiveState {
    return
  }

  // determine configuration
  config := &Configuration{}
  err := yaml.Unmarshal([]byte(request.Configuration), &config )
  if err != nil {
    response.Code   = http.StatusBadRequest
    response.Status = "Unable to parse configuration parameters:\n" + err.Error()
    return
  }

  // construct request information
  body, err := mustache.Render(template, config)
  if err != nil {
    response.Code   = http.StatusInternalServerError
    response.Status = "Unable to render configuration:\n" + err.Error()
		return
  }

  // create network
  transport := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true},}
  client    := &http.Client{Transport: transport}
  req, _    := http.NewRequest("POST", config.URL + "/apis/k8s.cni.cncf.io/v1/namespaces/default/network-attachment-definitions/" + config.Name, strings.NewReader(body) )
  req.Header.Set("Authorization", "Bearer " + config.Token)
  req.Header.Set("Content-Type", "application/yaml")
  rsp, err := client.Do(req)
	if err != nil {
    response.Code   = http.StatusBadRequest
    response.Status = "Unable to access cluster:\n" + err.Error()
    return
	}

  // check if the network has been created
  if rsp.StatusCode != http.StatusCreated {
    response.Code   = rsp.StatusCode
    response.Status = rsp.Status
		return
  }

  // determine status
  status(request, response)
}

//------------------------------------------------------------------------------

// destroy removes a network
func destroy(request *common.Request, response *common.Response) {
  // check if the network needs to be deleted
  status(request, response)
  if response.State == common.InactiveState {
    return
  }

  // determine configuration
  config := &Configuration{}
  err := yaml.Unmarshal([]byte(request.Configuration), &config )
  if err != nil {
    response.Code   = http.StatusBadRequest
    response.Status = "Unable to parse configuration parameters:\n" + err.Error()
    return
  }

  // destroy network
  transport := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true},}
  client    := &http.Client{Transport: transport}
  req, _    := http.NewRequest("DELETE", config.URL + "/apis/k8s.cni.cncf.io/v1/namespaces/default/network-attachment-definitions/" + config.Name, nil )
  req.Header.Set("Authorization", "Bearer " + config.Token)
  rsp, err := client.Do(req)
	if err != nil {
    response.Code   = http.StatusBadRequest
    response.Status = "Unable to access cluster:\n" + err.Error()
    return
	}

  // verify execution
  if rsp.StatusCode != http.StatusOK {
    response.Code   = rsp.StatusCode
    response.Status = rsp.Status
		return
  }

  // determine status
	status(request, response)
}

//------------------------------------------------------------------------------

// reconfigure fails since it is not implemented
func reconfigure(request *common.Request, response *common.Response) {
  response.Code   = http.StatusNotImplemented
  response.Status = "Method is not implemented"
}

//------------------------------------------------------------------------------

// status determines the status of a network
func status(request *common.Request, response *common.Response) {
  // determine configuration
  config := &Configuration{}
  err := yaml.Unmarshal([]byte(request.Configuration), &config )
  if err != nil {
    response.Code   = http.StatusBadRequest
    response.Status = "Unable to parse configuration:\n" + err.Error()
    return
  }

  // check status of cluster
  transport := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true},}
  client    := &http.Client{Transport: transport}
  req, _    := http.NewRequest("GET", config.URL + "/apis/k8s.cni.cncf.io/v1/namespaces/default/network-attachment-definitions/" + config.Name, nil)
  req.Header.Set("Authorization", "Bearer " + config.Token)
  rsp, err := client.Do(req)
	if err != nil {
    response.State  = common.FailureState
    response.Code   = http.StatusBadRequest
    response.Status = "Unable to access cluster:\n" + err.Error()
    return
	}

  // check if the network has been found
  if rsp.StatusCode != http.StatusOK {
    response.State    = common.InactiveState
    response.Endpoint = ""
    response.Code     = http.StatusOK
    response.Status   = ""
    return
  }

  // determine subnet
  subnet, err := getInfo(rsp)
  if err != nil {
    subnet = err.Error()
  }

  // update state and endpoint
  response.Configuration = request.Configuration
  response.State         = common.ActiveState
  response.Endpoint      = "Name: " + config.Name + "\nSubnet: " + subnet
  response.Code          = http.StatusOK
  response.Status        = ""
}

//------------------------------------------------------------------------------

// getInfo retrieves subnet information from a response
func getInfo(response *http.Response) (string, error) {
  // retrieve response body
  body, err := ioutil.ReadAll(response.Body)
  if err != nil {
    return "", errors.New("Unable to read response body")
  }

  // unmarshal body to output object
  output := map[string]interface{}{}
  err = yaml.Unmarshal(body, &output )
  if err != nil {
    return "", errors.New("Unable to convert response body into object")
  }

  // determine configuration string
  configurationString, err := mustache.Render("{{spec.config}}", output)
  if err != nil {
    return "", errors.New("Unable to determine configuration string")
  }
  configurationString = strings.Replace(configurationString, "&#34;", "\"", -1)

  // unmarshal configuration string into configuration object
  configuration := map[interface{}]interface{}{}
  err = yaml.Unmarshal([]byte(configurationString), &configuration )
  if err != nil {
    return "", errors.New("Unable to convert configuration string into object: " + configurationString)
  }

  // determine ipam
  subnet, err := mustache.Render("{{ipam.subnet}}", configuration)
  if err != nil {
    return "", errors.New("Unable to determine subnet")
  }

  // success
  return subnet, nil
}

//------------------------------------------------------------------------------
