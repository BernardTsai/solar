package api

import (
  "io"
  "io/ioutil"
  "net/http"

  "github.com/gorilla/mux"

  "tsai.eu/solar/model"
  "tsai.eu/solar/util"
  "tsai.eu/solar/engine"
)

//------------------------------------------------------------------------------

// SolutionListHandler lists the domains of the model.
func SolutionListHandler(w http.ResponseWriter, r *http.Request) {
  vars       := mux.Vars(r)
  domainName := vars["domain"]

  // determine domain
  domain, err := model.GetDomain(domainName)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  // determine solutions
  solutions, err := domain.ListSolutions()
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  // convert to yaml
  yaml, err := util.ConvertToYAML(solutions)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  // return the result
  io.WriteString(w, yaml)
}

//------------------------------------------------------------------------------

// SolutionSetHandler handles the uploading of a new component.
func SolutionSetHandler(w http.ResponseWriter, r *http.Request) {
  vars       := mux.Vars(r)
  domainName := vars["domain"]

  // get yaml
  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  // determine domain
  domain, err := model.GetDomain(domainName)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  // create new solution
  solution, _ := model.NewSolution("","","")

  err = solution.Load2(string(body))
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  // add solution to domain
  err = domain.AddSolution(solution)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
}

//------------------------------------------------------------------------------

// SolutionGetHandler retrieves a solution.
func SolutionGetHandler(w http.ResponseWriter, r *http.Request) {
  vars         := mux.Vars(r)
  domainName   := vars["domain"]
  solutionName := vars["solution"]

  // determine solution
  solution, err := model.GetSolution(domainName, solutionName)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  // transform architecture to string
  yaml, err := solution.Show()
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  // write yaml
  // w.Header().Set("Content-Type", "application/x-yaml")
  io.WriteString(w, yaml)
}

//------------------------------------------------------------------------------

// SolutionDeleteHandler deletes a solution.
func SolutionDeleteHandler(w http.ResponseWriter, r *http.Request) {
  vars         := mux.Vars(r)
  domainName   := vars["domain"]
  solutionName := vars["solution"]

  // determine domain
  domain, err := model.GetDomain(domainName)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  // delete solution
  err = domain.DeleteSolution(solutionName)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }
}

//------------------------------------------------------------------------------

// SolutionDeployHandler deploys an architecture.
func SolutionDeployHandler(w http.ResponseWriter, r *http.Request) {
  vars          := mux.Vars(r)
  domainName    := vars["domain"]
  solutionName  := vars["solution"]
  versionNumber := vars["version"]

  // determine domain
  domain, err := model.GetDomain(domainName)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  // determine architecture
  architecture, err := model.GetArchitecture(domainName, solutionName + " - " + versionNumber)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  // determine solution (create new solution if not found)
  solution, err := domain.GetSolution(architecture.Architecture)
  if err != nil {
    solution, _ = model.NewSolution(architecture.Architecture, architecture.Version, "")

    domain.AddSolution(solution)
  }

  // update the target state of the solution
  if err = solution.Update(domain.Name, architecture); err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  // create task and start it by signalling an event
  task, err := engine.NewSolutionTask(domain.Name, "", solution)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  // get event channel
  channel := engine.GetEventChannel()

  // create event
  channel <- model.NewEvent(domain.Name, task.UUID, model.EventTypeTaskExecution, "", "initial")
}

//------------------------------------------------------------------------------
