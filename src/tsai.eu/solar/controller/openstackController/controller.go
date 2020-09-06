package main

import (
  "os"
  "fmt"
  "errors"
  "strconv"
  "io/ioutil"
  "net/http"
  "github.com/gorilla/mux"
  "gopkg.in/yaml.v2"

  "tsai.eu/solar/controller/openstackController/common"
  "tsai.eu/solar/controller/openstackController/tenant"
  "tsai.eu/solar/controller/openstackController/network"
  // "tsai.eu/solar/controller/openstackController/server"
)

//------------------------------------------------------------------------------

// Controller manages the lifecycle of nothing
type Controller struct {
  Router *mux.Router
  Server *http.Server     // web server
}

//------------------------------------------------------------------------------

// main entry point for the controller
func main() {
  port := 10000

  // adjust port number if needed
  if len(os.Args) == 2 {
    port, err :=  strconv.Atoi(os.Args[1])

    // check for errors
    if err != nil || port < 1 || port > 65355 {
      fmt.Println("Invalid port number. Usage: openstackController <port>")
      os.Exit(-1)
    }
  }

  controller := NewController(port)

  err := controller.Server.ListenAndServe()
  if err != nil {
    fmt.Println(err.Error())
  }
}

//------------------------------------------------------------------------------

// NewController creates a new controller
func NewController(port int) (*Controller) {
  controller := Controller{}

  // create router
  controller.Router = mux.NewRouter()

  controller.Router.HandleFunc("/",         ping).Methods("GET")
  controller.Router.HandleFunc("/{action}", process).Methods("POST")

  // create server
  controller.Server = &http.Server{Addr: ":" + strconv.Itoa(port), Handler: controller.Router}

  // success
  return &controller
}

//------------------------------------------------------------------------------

// ping allows to ping the server
func ping(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
  w.Write([]byte("SOLAR:OpenStack:V1.0.0"))
}

//------------------------------------------------------------------------------

// process handles all incoming requests
func process(w http.ResponseWriter, r *http.Request) {
  // parse request
  request, err := parseRequest(r)
	if err != nil {
    w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unable to parse request:\n" + err.Error()))
		return
	}

  // construct initial response from request
  response := &common.Response{
    Request:        request.Request,
    Action:         mux.Vars(r)["action"],
    Code:           http.StatusTeapot,
    Status:         "",
    Domain:         request.Domain,
    Solution:       request.Solution,
    Version:        request.Version,
    Element:        request.Element,
    Cluster:        request.Cluster,
    Instance:       request.Instance,
    Component:      request.Component,
    State:          common.UndefinedState,
    Configuration:  request.Configuration,
    Endpoint:       "",
  }

  // determine and verify requested lifecycle transition
  if response.Action != common.CreateAction      &&
     response.Action != common.DestroyAction     &&
     response.Action != common.StartAction       &&
     response.Action != common.StopAction        &&
     response.Action != common.ConfigureAction   &&
     response.Action != common.ReconfigureAction &&
     response.Action != common.ResetAction       &&
     response.Action != common.StatusAction {

    response.Code   = http.StatusBadRequest
    response.Status = "Invalid action: " + response.Action

    writeResponse(w, response)
    return
  }

  // delegate to the corresponding handler
  component := request.Component
  version   := request.Cluster

  switch component + ":" + version {
  case "os-tenant:V1.0.0":
    tenant.Process(request, response)
  case "os-network:V1.0.0":
    network.Process(request, response)
  // case "os-server:V1.0.0":
  //   server.Process(request, response)
  default:
    response.Code   = http.StatusBadRequest
    response.Status = "Unknown component: " + component + ":" + version

    writeResponse(w, response)
    return
  }

  // return results
  writeResponse(w, response)
}

//------------------------------------------------------------------------------

// parseRequest reads the request from a http request body
func parseRequest(r *http.Request) (request *common.Request, err error) {
  request = &common.Request{}

  // handle any issues with the conversion
  defer func() {
    if rec := recover(); rec != nil {
      reason := fmt.Sprintf("Panicked while attempting to parse request:\n%s", rec)
      err = errors.New(reason)
    }
  }()

  // get body
	body, _ := ioutil.ReadAll(r.Body)

  // convert body into target state
  err = yaml.Unmarshal(body, request)
  if err != nil {
    return request, err
  }

  // success
  return request, nil
}

//------------------------------------------------------------------------------

// writeResponse outputs the response structure
func writeResponse(w http.ResponseWriter, response *common.Response) {
  body, _ := yaml.Marshal(response)

  w.WriteHeader(response.Code)
  w.Write(body)
}
