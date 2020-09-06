package network

import (
  "tsai.eu/solar/controller/openstackController/common"
)

//------------------------------------------------------------------------------

// Process handles all incoming namespace related requests
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
