package api

import (
  "io"
  "io/ioutil"
  "net/http"

  "github.com/gorilla/mux"

  "tsai.eu/solar/model"
  "tsai.eu/solar/engine"
  "tsai.eu/solar/util"
)

//------------------------------------------------------------------------------

// InstanceUpdateInformation describes desired state an instance should adopt.
type InstanceUpdateInformation struct {
	State string `yaml:"State"` // desired state of the instance
}

//------------------------------------------------------------------------------

// InstanceUpdateHandler updates an instance.
func InstanceUpdateHandler(w http.ResponseWriter, r *http.Request) {
  vars          := mux.Vars(r)
  domainName    := vars["domain"]
  solutionName  := vars["solution"]
  elementName   := vars["element"]
  clusterName   := vars["cluster"]
  instanceName  := vars["instance"]
  config        := InstanceUpdateInformation{}

  // determine desired configuration
  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    io.WriteString(w, "unable to read target configuration:\n" + err.Error())
    return
  }

  yaml := string(body)

  util.ConvertFromYAML(yaml, &config)

  // determine solution
  solution, err := model.GetSolution(domainName, solutionName)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    io.WriteString(w, "unable to determine solution")
    return
  }

  // create task and start it by signalling an event
	task, err := engine.NewInstanceTask(domainName, "", solutionName, solution.Version, elementName, clusterName, instanceName, config.State)
	if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    io.WriteString(w, "task can not be created:\n" + err.Error())
		return
	}

	// get event channel
	channel := engine.GetEventChannel()

	// create event
	channel <- model.NewEvent(domainName, task.UUID, model.EventTypeTaskExecution, "", "initial")

  // return the uuid of the task
  io.WriteString(w, task.UUID)
}

//------------------------------------------------------------------------------
