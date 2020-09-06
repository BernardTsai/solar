package deployment

import (
  "errors"
  "strings"
  "strconv"
  "io/ioutil"
  "fmt"

  "net/http"
  "crypto/tls"
  "gopkg.in/yaml.v2"
  "github.com/cbroglie/mustache"

  "tsai.eu/solar/controller/openstackController/common"
)

//------------------------------------------------------------------------------

// template resembles an k8s bridge configuration
const template = `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{Name}}
  labels:
    app: {{Name}}
{{Networks}}
spec:
  replicas: {{Size}}
  selector:
    matchLabels:
      app: {{Name}}
  template:
    metadata:
      labels:
        app: {{Name}}
    spec:
      containers:
      - name: {{Name}}
        image: {{Image}}:{{Version}}
        ports:
        - containerPort: 80`

const networkTemplate =
`  annotations:
    k8s.v1.cni.cncf.io/networks: {{Networks}}`

//------------------------------------------------------------------------------

// Configuration holds the information required to address a k8s API endpoint
type Configuration struct {
  URL              string `yaml:"URL"`      // URL of API endpoint
  Token            string `yaml:"Token"`    // authorization bearer token
  Name             string `yaml:"Name"`     // Name of deployment
  Size             int64  `yaml:"Size"`     // Number of replicas
  Image            string `yaml:"Image"`    // Name of image
  Version          string `yaml:"Version"`  // Version of image
  Networks         string `yaml:"Networks"` // Comma seperated list of network names
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
// Component: k8s-deployment
// State: active
// Configuration: |
//   URL: https://192.168.178.20:6443
//   Token: eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImRlbW8tYWRtaW4tdG9rZW4tbHAyY3EiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoiZGVtby1hZG1pbiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6IjA3ODJkZmYyLTY5ZjEtMTFlOS1hNzNkLWYwZGVmMWZjN2FkNSIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZWZhdWx0OmRlbW8tYWRtaW4ifQ.Z6WIRuaqNKhTG2u9cx4eXd8Ado6j1bNwN8JiZFTF9hPLC90i820auALiL9RtHBv-lOo_auSZKJtF9eiLFRxyoon7RsUSie-pZtKHupTKb3YbggEC9ZZxt-NGyKhMhNLRtJxri-FSkSzJqQkkXjZ_qX0CQSgi22-cIbH-ki99C2RvvySkXSJSXeQOhBVDyqSL1Skj042KITlMne20G4AxMB5HSjUsrWdYK8dGCdHW0ohxIRKW1nLvySMTk3G4Em-7Sy-Ug_bVo3KK_ObycmICHaYcbWchHj_6NHYh3ISx0yMpoSsLS3Q3Z6qDABGNVFidMXcvANHNaWtzeiLXNKeVAA
//   Name: nginx
//   Size: 3
//   Image: nginx
//   Version: 1.7.9
//   Networks: oam, m2m

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
// Component: k8s-deployment
// State: active
// Configuration: |
//   URL: https://192.168.178.20:6443
//   Token: eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImRlbW8tYWRtaW4tdG9rZW4tbHAyY3EiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoiZGVtby1hZG1pbiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6IjA3ODJkZmYyLTY5ZjEtMTFlOS1hNzNkLWYwZGVmMWZjN2FkNSIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZWZhdWx0OmRlbW8tYWRtaW4ifQ.Z6WIRuaqNKhTG2u9cx4eXd8Ado6j1bNwN8JiZFTF9hPLC90i820auALiL9RtHBv-lOo_auSZKJtF9eiLFRxyoon7RsUSie-pZtKHupTKb3YbggEC9ZZxt-NGyKhMhNLRtJxri-FSkSzJqQkkXjZ_qX0CQSgi22-cIbH-ki99C2RvvySkXSJSXeQOhBVDyqSL1Skj042KITlMne20G4AxMB5HSjUsrWdYK8dGCdHW0ohxIRKW1nLvySMTk3G4Em-7Sy-Ug_bVo3KK_ObycmICHaYcbWchHj_6NHYh3ISx0yMpoSsLS3Q3Z6qDABGNVFidMXcvANHNaWtzeiLXNKeVAA
//   Name: "nginx"
//   Size: "3"
//   Image: nginx
//   Version: 1.7.9
//   Networks: oam, m2m
// Endpoint: |-
//   Name: nginx
//   Size: 3
//   Image: nginx
//   Version: 1.7.9
//   Networks: oam, m2m
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
  // check if the deployment has been created already
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
  if config.Networks != "" {
    networks, err := mustache.Render(networkTemplate, config)
    if err != nil {
      response.Code   = http.StatusBadRequest
      response.Status = "Unable to render network configuration parameters:\n" + err.Error()
      return
    }
    config.Networks = networks
  }

  body, err := mustache.Render(template, config)
  if err != nil {
    response.Code   = http.StatusInternalServerError
    response.Status = "Unable to render configuration:\n" + err.Error()
		return
  }
  fmt.Println(config)
  fmt.Println(body)

  // create deployment
  transport := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true},}
  client    := &http.Client{Transport: transport}
  req, _    := http.NewRequest("POST", config.URL + "/apis/apps/v1/namespaces/default/deployments", strings.NewReader(body) )
  req.Header.Set("Authorization", "Bearer " + config.Token)
  req.Header.Set("Content-Type", "application/yaml")
  rsp, err := client.Do(req)
	if err != nil {
    response.Code   = http.StatusBadRequest
    response.Status = "Unable to access cluster:\n" + err.Error()
    return
	}

  // check if the deployment has been created
  if rsp.StatusCode != http.StatusCreated {
    rspBody, _ := ioutil.ReadAll(rsp.Body)

    response.Code   = rsp.StatusCode
    response.Status = rsp.Status + "\nRequest:\n" + body + "\n\nResponse:\n" + string(rspBody)
		return
  }

  // determine status
  status(request, response)
}

//------------------------------------------------------------------------------

// destroy removes a network
func destroy(request *common.Request, response *common.Response) {
  // check if the deployment needs to be deleted
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
  req, _    := http.NewRequest("DELETE", config.URL + "/apis/apps/v1/namespaces/default/deployments/" + config.Name, nil )
  req.Header.Set("Authorization", "Bearer " + config.Token)
  rsp, err := client.Do(req)
	if err != nil {
    response.Code   = http.StatusBadRequest
    response.Status = "Unable to access cluster:\n" + err.Error()
    return
	}

  // verify execution
  if rsp.StatusCode != http.StatusOK {
    rspBody, _ := ioutil.ReadAll(rsp.Body)

    response.Code   = rsp.StatusCode
    response.Status = rsp.Status + "\nRequest:\n" + "\n\nResponse:\n" + string(rspBody)
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

// status determines the status of a namespace
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
  req, _    := http.NewRequest("GET", config.URL + "/apis/apps/v1/namespaces/default/deployments/" + config.Name, nil )
  req.Header.Set("Authorization", "Bearer " + config.Token)
  rsp, err := client.Do(req)
	if err != nil {
    response.State  = common.FailureState
    response.Code   = http.StatusBadRequest
    response.Status = "Unable to access cluster:\n" + err.Error()
    return
	}

  // check if the deployment has been found
  if rsp.StatusCode != http.StatusOK {
    response.State    = common.InactiveState
    response.Endpoint = ""
    response.Code     = http.StatusOK
    response.Status   = ""
    return
  }

  // determine size, image, version
  name, size, image, version, networks, err := getInfo(rsp)
  if err != nil {
    name    = ""
    size    = 0
    image   = ""
    version = err.Error()
  }

  // update state and endpoint
  response.State    = common.ActiveState
  response.Endpoint = "Name: " + name + "\nSize: " + string(size) + "\nImage: " + image + "\nVersion: " + version + "\nNetworks: " + networks
  response.Code     = http.StatusOK
  response.Status   = ""
}

//------------------------------------------------------------------------------

type annotations map[string]string;

type metadataStruct struct {
  Annotations map[string]*annotations `yaml:"annotations"`
}

type outputStruct struct {
  MetaData *annotations `yaml:"metadata"`
}

// getInfo retrieves name, size, image and version information from a response
func getInfo(response *http.Response) (name string, size int64, image string, version string, networks string, err error) {
  // retrieve response body
  body, err := ioutil.ReadAll(response.Body)
  if err != nil {
    err = errors.New("Unable to read response body")
    return
  }

  // unmarshal body to output object
  output := map[string]interface{}{}
  err = yaml.Unmarshal(body, &output )
  if err != nil {
    err = errors.New("Unable to convert response body into object")
    return
  }

  // determine name
  name, err = mustache.Render("{{metadata.name}}", output)
  if err != nil {
    err = errors.New("Unable to determine name")
    return
  }

  // determine size
  sizeInfo, err := mustache.Render("{{spec.replicas}}", output)
  if err != nil {
    err = errors.New("Unable to determine size")
    return
  }
  size, err = strconv.ParseInt(sizeInfo, 10, 32)
  if err != nil {
    err = errors.New("Invalid size: " + sizeInfo)
    return
  }

  // determine image
  imageInfo, err := mustache.Render("{{#spec.template.spec.containers}}{{image}}{{/spec.template.spec.containers}}", output)
  if err != nil {
    err = errors.New("Unable to determine image information")
    return
  }

  parts := strings.Split(imageInfo, ":")
  if len(parts) != 2 {
    err = errors.New("Unable to parse image information")
    return
  }
  image   = parts[0]
  version = parts[1]

  // determine networks
  metadata, ok := output["metadata"].(map[interface{}]interface{})
  if !ok {
    networks = ""
    return
  }

  annotations, ok := metadata["annotations"].(map[interface{}]interface{})
  if !ok {
    networks = ""
    return
  }

  networks, ok = annotations["k8s.v1.cni.cncf.io/networks"].(string)
  if !ok {
    networks = ""
    return
  }

  return
}

//------------------------------------------------------------------------------
