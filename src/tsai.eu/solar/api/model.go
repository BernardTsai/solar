package api

import (
  "io"
  "io/ioutil"
  "net/http"

  "tsai.eu/solar/model"
)

//------------------------------------------------------------------------------

// ModelSetHandler handles the uploading of a new model.
func ModelSetHandler(w http.ResponseWriter, r *http.Request) {
  body, err := ioutil.ReadAll(r.Body)

  // check validity of input
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  // load yaml model
  err = model.GetModel().Load2(string(body))

  // check success of import
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }
}

//------------------------------------------------------------------------------

// ModelGetHandler retrieves the current model.
func ModelGetHandler(w http.ResponseWriter, r *http.Request) {
  result, _ := model.GetModel().Show()

  // w.Header().Set("Content-Type", "application/x-yaml")

  io.WriteString(w, result)
}

//------------------------------------------------------------------------------

// ModelResetHandler resets the current model.
func ModelResetHandler(w http.ResponseWriter, r *http.Request) {
  model.GetModel().Reset()
}

//------------------------------------------------------------------------------
