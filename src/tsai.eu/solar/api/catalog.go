package api

import (
  "io"
  "net/http"

  "github.com/gorilla/mux"

  "tsai.eu/solar/model"
  "tsai.eu/solar/util"
)

//------------------------------------------------------------------------------

// CatalogGetHandler retrieves the catalog of a domain.
func CatalogGetHandler(w http.ResponseWriter, r *http.Request) {
  vars       := mux.Vars(r)
  domainName := vars["domain"]

  // determine domain
  domain, err := model.GetDomain(domainName)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  // determine components
  components, err := domain.GetComponents()
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
