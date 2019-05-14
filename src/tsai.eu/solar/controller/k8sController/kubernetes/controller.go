package kubernetes

import (
  "net/http"
  "crypto/tls"
  "gopkg.in/yaml.v2"

  "tsai.eu/solar/controller/k8sController/common"
)

//------------------------------------------------------------------------------

// Configuration holds the information required to address a k8s API endpoint
type Configuration struct {
  URL   string `yaml:"URL"`    // URL of API endpoint
  Token string `yaml:"Token"`  // authorization bearer token
}

// Sample request:
//
// Request: 5cf2625e-abcd-473a-9b8c-5838b014dab4
// Domain: demo
// Solution: app
// Version: V1.0.0
// Element: kubernetes
// Cluster: V1.0.0
// Instance: 5cf2625e-6fb0-473a-9b8c-5838b0140ea3
// Component: k8s-kubernetes
// State: active
// Configuration: |
//   URL: "https://192.168.178.20:6443"
//   Token: "eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImRlbW8tYWRtaW4tdG9rZW4tbHAyY3EiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoiZGVtby1hZG1pbiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6IjA3ODJkZmYyLTY5ZjEtMTFlOS1hNzNkLWYwZGVmMWZjN2FkNSIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZWZhdWx0OmRlbW8tYWRtaW4ifQ.Z6WIRuaqNKhTG2u9cx4eXd8Ado6j1bNwN8JiZFTF9hPLC90i820auALiL9RtHBv-lOo_auSZKJtF9eiLFRxyoon7RsUSie-pZtKHupTKb3YbggEC9ZZxt-NGyKhMhNLRtJxri-FSkSzJqQkkXjZ_qX0CQSgi22-cIbH-ki99C2RvvySkXSJSXeQOhBVDyqSL1Skj042KITlMne20G4AxMB5HSjUsrWdYK8dGCdHW0ohxIRKW1nLvySMTk3G4Em-7Sy-Ug_bVo3KK_ObycmICHaYcbWchHj_6NHYh3ISx0yMpoSsLS3Q3Z6qDABGNVFidMXcvANHNaWtzeiLXNKeVAA"

// Sample response:
//
// Request: 5cf2625e-abcd-473a-9b8c-5838b014dab4
// Action: status
// Code: 200
// Status: ""
// Domain: demo
// Solution: app
// Version: V1.0.0
// Element: kubernetes
// Cluster: V1.0.0
// Instance: 5cf2625e-6fb0-473a-9b8c-5838b0140ea3
// State: active
// Configuration: ""
// Endpoint: |-
//   URL: https://192.168.178.20:6443
//   Token: eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImRlbW8tYWRtaW4tdG9rZW4tbHAyY3EiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoiZGVtby1hZG1pbiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6IjA3ODJkZmYyLTY5ZjEtMTFlOS1hNzNkLWYwZGVmMWZjN2FkNSIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZWZhdWx0OmRlbW8tYWRtaW4ifQ.Z6WIRuaqNKhTG2u9cx4eXd8Ado6j1bNwN8JiZFTF9hPLC90i820auALiL9RtHBv-lOo_auSZKJtF9eiLFRxyoon7RsUSie-pZtKHupTKb3YbggEC9ZZxt-NGyKhMhNLRtJxri-FSkSzJqQkkXjZ_qX0CQSgi22-cIbH-ki99C2RvvySkXSJSXeQOhBVDyqSL1Skj042KITlMne20G4AxMB5HSjUsrWdYK8dGCdHW0ohxIRKW1nLvySMTk3G4Em-7Sy-Ug_bVo3KK_ObycmICHaYcbWchHj_6NHYh3ISx0yMpoSsLS3Q3Z6qDABGNVFidMXcvANHNaWtzeiLXNKeVAA
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

// status determines the status of a kubenetes cluster
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
  req, _    := http.NewRequest("GET", config.URL, nil)
  req.Header.Set("Authorization", "Bearer " + config.Token)
  _, err = client.Do(req)
	if err != nil {
    response.State  = common.FailureState
    response.Code   = http.StatusBadRequest
    response.Status = "Unable to access cluster:\n" + err.Error()
    return
	}

  // update state and endpoint
  response.Configuration = request.Configuration
  response.State         = common.ActiveState
  response.Endpoint      = "URL: " + config.URL + "\nToken: " + config.Token
  response.Code          = http.StatusOK
  response.Status        = ""
}

//------------------------------------------------------------------------------
