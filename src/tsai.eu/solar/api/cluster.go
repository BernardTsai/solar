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

// ClusterUpdateInformation describes desired state a cluster should adopt.
type ClusterUpdateInformation struct {
  Min   int    `yaml:"Min"`   // desired min. size of the solution element cluster
	Max   int    `yaml:"Max"`   // desired max. size of the solution element cluster
	Size  int    `yaml:"Size"`  // desired size of the solution element cluster
	State string `yaml:"State"` // desired state of the solution element cluster
}

//------------------------------------------------------------------------------

// ClusterUpdateHandler updates a cluster.
func ClusterUpdateHandler(w http.ResponseWriter, r *http.Request) {
  vars          := mux.Vars(r)
  domainName    := vars["domain"]
  solutionName  := vars["solution"]
  elementName   := vars["element"]
  clusterName   := vars["cluster"]
  config        := ClusterUpdateInformation{}

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

  // determine cluster
  cluster, err := model.GetCluster(domainName, solutionName, elementName, clusterName)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    io.WriteString(w, "unable to determine instance")
    return
  }

  // update target state and dimensions of cluster
  cluster.Target = config.State

  cluster.Resize(config.Min, config.Max, config.Size)

  // create task and start it by signalling an event
	task, err := engine.NewClusterTask(domainName, "", solutionName, solution.Version, elementName, clusterName)
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
