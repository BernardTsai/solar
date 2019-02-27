package api

import (
  "io"
  "net/http"
  "strconv"

  "github.com/gorilla/mux"

  "tsai.eu/solar/model"
  "tsai.eu/solar/util"
  "tsai.eu/solar/cli"
  "tsai.eu/solar/engine"
)

//------------------------------------------------------------------------------

// TaskListHandler lists the tasks of the model.
func TaskListHandler(w http.ResponseWriter, r *http.Request) {
  vars         := mux.Vars(r)
  domainName   := vars["domain"]
  solutionName := vars["solution"]
  elementName  := vars["element"]
  clusterName  := vars["cluster"]
  instanceName := vars["instance"]

  // determine domain
  domain, err := model.GetDomain(domainName)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  // determine list of solution names
  tasks := []string{}

  tNames, _ := domain.ListTasks()
  for _, tName := range tNames {
    task, _ := domain.GetTask(tName)

    if (solutionName == "" || solutionName == task.Solution) &&
       (elementName  == "" || elementName  == task.Element)  &&
       (clusterName  == "" || clusterName  == task.Cluster)  &&
       (instanceName == "" || instanceName == task.Instance) {
      tasks = append(tasks, tName)
    }
  }

  // convert list to yaml
  yaml, err := util.ConvertToYAML(tasks)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  // return the result
  io.WriteString(w, yaml)
}

//------------------------------------------------------------------------------

// TaskGetHandler retrieves a task tree.
func TaskGetHandler(w http.ResponseWriter, r *http.Request) {
  vars         := mux.Vars(r)
  domainName   := vars["domain"]
  taskName     := vars["task"]
  levelNumber  := vars["level"]

  // determine task
  task, err := model.GetTask(domainName, taskName)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  // determine level
  level := 0
  if levelNumber != "" {
    value, err := strconv.Atoi(levelNumber)
    if err != nil {
      w.WriteHeader(http.StatusBadRequest)
      return
    }

    if value < 0 {
      w.WriteHeader(http.StatusBadRequest)
      return
    }
    level = value
  }

  // retrieve the task information
  taskinfo := cli.NewTaskInfo(task, level)
  yaml, err := util.ConvertToYAML(taskinfo)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  // write yaml
  // w.Header().Set("Content-Type", "application/x-yaml")
  io.WriteString(w, yaml)
}

//------------------------------------------------------------------------------

// TaskTerminateHandler retrieves a task tree.
func TaskTerminateHandler(w http.ResponseWriter, r *http.Request) {
  vars         := mux.Vars(r)
  domainName   := vars["domain"]
  taskName     := vars["task"]

  // determine task
  task, err := model.GetTask(domainName, taskName)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  // execute the command
  // get event channel
  channel := engine.GetEventChannel()

  // create event
  channel <- model.NewEvent(domainName, task.UUID, model.EventTypeTaskTermination, "", "initial")
}

//------------------------------------------------------------------------------
