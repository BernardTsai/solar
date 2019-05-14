package controller

import (
	"net/http"

	"tsai.eu/solar/controller/defaultController/common"
)

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
  // copy request information to response information
  response.Configuration = request.Configuration
  response.State         = request.State
  response.Endpoint      = ""
  response.Code          = http.StatusOK
  response.Status        = ""
}

//------------------------------------------------------------------------------
