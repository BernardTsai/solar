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

// DomainListHandler lists the domains of the model.
func DomainListHandler(w http.ResponseWriter, r *http.Request) {
  domains, err := model.GetModel().ListDomains()

  // check validity of input
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  // load yaml model
  result, err := util.ConvertToYAML(domains)

  // check validity of result
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  // return the result
  io.WriteString(w, result)
}

//------------------------------------------------------------------------------

// DomainCreateHandler creates a new domain.
func DomainCreateHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)

  domainName := vars["domain"]

  domain, _ := model.NewDomain(domainName)

  err := model.GetModel().AddDomain(domain)

  // check validity of result
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }
}

//------------------------------------------------------------------------------

// DomainDeleteHandler deletes a domain.
func DomainDeleteHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)

  domainName := vars["domain"]

  err := model.GetModel().DeleteDomain(domainName)

  // check validity of result
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }
}

//------------------------------------------------------------------------------

// DomainSetHandler handles the uploading of a new domain.
func DomainSetHandler(w http.ResponseWriter, r *http.Request) {
  body, err := ioutil.ReadAll(r.Body)

  // check validity of input
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  // load yaml domain
  domain, _ := model.NewDomain("dummy")

  err = domain.Load2(string(body))

  // check validity of input
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  // add domain to model
  err = model.GetModel().AddDomain(domain)

  // check success of import
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }
}

//------------------------------------------------------------------------------

// DomainGetHandler retrieves a domain.
func DomainGetHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)

  domainName := vars["domain"]

  domain, err := model.GetDomain(domainName)

  // check validity of result
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  // transform domain to string
  result, err := domain.Show()

  // check validity of result
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  // w.Header().Set("Content-Type", "application/x-yaml")

  io.WriteString(w, result)
}

//------------------------------------------------------------------------------

// DomainResetHandler resets a domain.
func DomainResetHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)

  domainName := vars["domain"]

  err := model.GetModel().DeleteDomain(domainName)

  // check validity of result
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  // create new domain
  domain, _ := model.NewDomain(domainName)

  // check if domain has been created
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  // add domain to model
  err = model.GetModel().AddDomain(domain)

  // check success of import
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }
}

//------------------------------------------------------------------------------
