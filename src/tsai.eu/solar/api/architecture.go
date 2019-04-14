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

// ArchitectureListHandler lists the architectures of a domain.
func ArchitectureListHandler(w http.ResponseWriter, r *http.Request) {
  vars       := mux.Vars(r)
  domainName := vars["domain"]

  // determine domain
  domain, err := model.GetDomain(domainName)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  // determine architectures
  architectures, err := domain.ListArchitectures()
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  // convert to yaml
  result, err := util.ConvertToYAML(architectures)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  // return the result
  io.WriteString(w, result)
}

//------------------------------------------------------------------------------

// ArchitectureSetHandler handles the uploading of a new architecture.
func ArchitectureSetHandler(w http.ResponseWriter, r *http.Request) {
  vars       := mux.Vars(r)
  domainName := vars["domain"]

  // get yaml
  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    io.WriteString(w, "unable to read architecture:\n" + err.Error())
    return
  }

  // determine domain
  domain, err := model.GetDomain(domainName)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    io.WriteString(w, "domain can not be identified")
    return
  }

  // create new architecture
  architecture, _ := model.NewArchitecture("","","")

  err = architecture.Load2(string(body))
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    io.WriteString(w, "unable to parse architecture:\n" + err.Error())
    return
  }

  // add architecture to domain
  err = domain.AddArchitecture(architecture)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    io.WriteString(w, "unable to add architecture:\n" + err.Error())
    return
  }
}

//------------------------------------------------------------------------------

// ArchitectureGetHandler retrieves an architecture.
func ArchitectureGetHandler(w http.ResponseWriter, r *http.Request) {
  vars                := mux.Vars(r)
  domainName          := vars["domain"]
  architectureName    := vars["architecture"]
  architectureVersion := vars["version"]

  // determine architecture
  architecture, err := model.GetArchitecture(domainName, architectureName, architectureVersion)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  // transform architecture to string
  yaml, err := architecture.Show()
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  // write yaml
  // w.Header().Set("Content-Type", "application/x-yaml")
  io.WriteString(w, yaml)
}

//------------------------------------------------------------------------------

// ArchitectureDeleteHandler deletes an architecture.
func ArchitectureDeleteHandler(w http.ResponseWriter, r *http.Request) {
  vars                := mux.Vars(r)
  domainName          := vars["domain"]
  architectureName    := vars["architecture"]
  architectureVersion := vars["version"]

  // determine domain
  domain, err := model.GetDomain(domainName)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  // delete architecture
  err = domain.DeleteArchitecture(architectureName, architectureVersion)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }
}

//------------------------------------------------------------------------------

// ArchitectureDeployHandler deploys an architecture.
func ArchitectureDeployHandler(w http.ResponseWriter, r *http.Request) {
  vars                := mux.Vars(r)
  domainName          := vars["domain"]
  architectureName    := vars["architecture"]
  architectureVersion := vars["version"]

	// determine domain
	domain, err := model.GetDomain(domainName)

	if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    io.WriteString(w, "domain can not be identified")
    return
	}

	// determine architecture
	architecture, err := domain.GetArchitecture(architectureName, architectureVersion)

	if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    io.WriteString(w, "architecture can not be identified")
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
    io.WriteString(w, "unable to create or update the solution:\n" + err.Error())
    return
	}

	// create task and start it by signalling an event
	task, err := engine.NewSolutionTask(domain.Name, "", solution)
	if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    io.WriteString(w, "task can not be created:\n" + err.Error())
		return
	}

	// get event channel
	channel := engine.GetEventChannel()

	// create event
	channel <- model.NewEvent(domain.Name, task.UUID, model.EventTypeTaskExecution, "", "initial")

  // return the uuid of the task
  io.WriteString(w, task.UUID)
}

//------------------------------------------------------------------------------
