package api

import (
  "io"
  "net/http"
  "strconv"
  "time"
  "sort"

  "github.com/gorilla/mux"

  "tsai.eu/solar/model"
  "tsai.eu/solar/util"
  "tsai.eu/solar/engine"
)

//------------------------------------------------------------------------------

// TaskSummary captures the most relevant informations related to a task
type TaskSummary struct {
  Type         string   `yaml:"Type"`         // type of task
	Domain       string   `yaml:"Domain"`       // domain of task
	Solution     string   `yaml:"Solution"`     // architecture of entity
	Version      string   `yaml:"Version"`      // architecture version of entity
	Element      string   `yaml:"Element"`      // element of entity
	Cluster      string   `yaml:"Cluster"`      // cluster of entity
	Instance     string   `yaml:"Instance"`     // instance of entity
	State        string   `yaml:"State"`        // desired state of entity
	UUID         string   `yaml:"UUID"`         // uuid of task
	Parent       string   `yaml:"Parent"`       // uuid of parent task
	Status       string   `yaml:"Status"`       // status of task: (execution/completion/failure)
	Phase        int      `yaml:"Phase"`        // phase of task
  Started      string   `yaml:"Started"`      // start date
  Completed    string   `yaml:"Completed"`    // completed date
  Latest       string   `yaml:"Latest"`       // date of latest event
}

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
  tasks := []*TaskSummary{}

  tNames, _ := domain.ListTasks()
  for _, tName := range tNames {
    task, _ := domain.GetTask(tName)

    if (solutionName == task.Solution) &&
       (elementName  == task.Element)  &&
       (clusterName  == task.Cluster)  &&
       (instanceName == task.Instance) {

      // copy task information into summary
      summary := TaskSummary{
        Type:       task.Type,
        Domain:     task.Domain,
        Solution:   task.Solution,
        Version:    task.Version,
        Element:    task.Element,
        Cluster:    task.Cluster,
        Instance:   task.Instance,
        State:      task.State,
        UUID:       task.UUID,
        Parent:     task.Parent,
        Status:     task.Status,
        Phase:      task.Phase,
        Started:    "",
        Completed:  "",
        Latest:     "",
      }

      // derive most relevant event information
      min, max, lst := task.GetTimestamps()

      // copy to summary
      if min != 0 {
        // summary.Started = strconv.FormatInt(min, 10)
        summary.Started = time.Unix(0, min).String()
      }
      if max != 0 {
        // summary.Completed = strconv.FormatInt(max, 10)
        summary.Completed = time.Unix(0, max).String()
      }
      if lst != 0 {
        // summary.Latest = strconv.FormatInt(lst, 10)
        summary.Latest = time.Unix(0, lst).String()
      }

      // append summary to list of tasks
      tasks = append(tasks, &summary)
    }
  }

  // sort tasks along started timestamp
  sort.Slice(tasks, func(i int,j int) bool {
    return tasks[i].Started > tasks[j].Started
  })

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

// TaskTraceHandler retrieves a task trace.
func TaskTraceHandler(w http.ResponseWriter, r *http.Request) {
  vars         := mux.Vars(r)
  domainName   := vars["domain"]
  taskName     := vars["task"]

  // determine task
  task, err := model.GetTask(domainName, taskName)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  // retrieve the task information
  trace := model.NewTrace(task)
  yaml, err := util.ConvertToYAML(trace)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  // write yaml
  // w.Header().Set("Content-Type", "application/x-yaml")
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
  taskinfo := model.NewTaskInfo(task, level)
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
