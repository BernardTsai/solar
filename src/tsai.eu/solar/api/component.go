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

// ComponentListHandler lists the domains of the model.
func ComponentListHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)

  domainName := vars["domain"]

  // determine domain
  domain, err := model.GetDomain(domainName)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  // determine components
  components, err := domain.ListComponents()
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  // convert to yaml
  result, err := util.ConvertToYAML(components)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  // return the result
  io.WriteString(w, result)
}

//------------------------------------------------------------------------------

// ComponentSetHandler handles the uploading of a new component.
func ComponentSetHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)

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

  // create new component
  component, _ := model.NewComponent("", "", "")

  err = component.Load2(string(body))
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  // add component to domain
  err = domain.AddComponent(component)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
}

//------------------------------------------------------------------------------

// ComponentGetHandler retrieves a domain.
func ComponentGetHandler(w http.ResponseWriter, r *http.Request) {
  vars          := mux.Vars(r)
  domainName    := vars["domain"]
  componentName := vars["component"]
  version       := vars["version"]

  // determine component
  component, err := model.GetComponent(domainName, componentName, version)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  // transform domain to string
  yaml, err := component.Show()
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  // write yaml
  // w.Header().Set("Content-Type", "application/x-yaml")
  io.WriteString(w, yaml)
}

//------------------------------------------------------------------------------

// ComponentDeleteHandler deletes a component.
func ComponentDeleteHandler(w http.ResponseWriter, r *http.Request) {
  vars          := mux.Vars(r)
  domainName    := vars["domain"]
  componentName := vars["component"]
  version       := vars["version"]

  // determine domain
  domain, err := model.GetDomain(domainName)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  // delete component
  err = domain.DeleteComponent(componentName, version)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }
}

//------------------------------------------------------------------------------
