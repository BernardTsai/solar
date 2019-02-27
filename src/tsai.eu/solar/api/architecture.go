package api

import (
  "io"
  "io/ioutil"
  "net/http"

  "github.com/gorilla/mux"

  "tsai.eu/solar/model"
  "tsai.eu/solar/util"
)

//------------------------------------------------------------------------------

// ArchitectureListHandler lists the domains of the model.
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

// ArchitectureSetHandler handles the uploading of a new component.
func ArchitectureSetHandler(w http.ResponseWriter, r *http.Request) {
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

  // create new architecture
  architecture, _ := model.NewArchitecture("","","")

  err = architecture.Load2(string(body))
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  // add architecture to domain
  err = domain.AddArchitecture(architecture)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
}

//------------------------------------------------------------------------------

// ArchitectureGetHandler retrieves a domain.
func ArchitectureGetHandler(w http.ResponseWriter, r *http.Request) {
  vars             := mux.Vars(r)
  domainName       := vars["domain"]
  architectureName := vars["architecture"]

  // determine architecture
  architecture, err := model.GetArchitecture(domainName, architectureName)
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

// ArchitectureDeleteHandler deletes a component.
func ArchitectureDeleteHandler(w http.ResponseWriter, r *http.Request) {
  vars             := mux.Vars(r)
  domainName       := vars["domain"]
  architectureName := vars["architecture"]

  // determine domain
  domain, err := model.GetDomain(domainName)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  // delete architecture
  err = domain.DeleteArchitecture(architectureName)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }
}

//------------------------------------------------------------------------------
