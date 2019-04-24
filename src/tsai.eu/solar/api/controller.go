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

// ControllerListHandler lists the domains of the model.
func ControllerListHandler(w http.ResponseWriter, r *http.Request) {
  vars       := mux.Vars(r)
  domainName := vars["domain"]

  // determine domain
  domain, err := model.GetDomain(domainName)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  // determine list of controllers
  controllers := []*model.Controller{}

  cNameVersions, err := domain.ListControllers()
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  for _, cNameVersion := range cNameVersions {
    controller, _ := domain.GetController(cNameVersion[0], cNameVersion[1])

    controllers = append(controllers, controller)
  }

  // convert to yaml
  yaml, err := util.ConvertToYAML(controllers)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  // return the result
  io.WriteString(w, yaml)
}

//------------------------------------------------------------------------------

// ControllerSetHandler handles the uploading of a new component.
func ControllerSetHandler(w http.ResponseWriter, r *http.Request) {
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

  // create new controller
  controller, _ := model.NewController("","")

  err = controller.Load2(string(body))
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  // add controller to domain
  err = domain.AddController(controller)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
}

//------------------------------------------------------------------------------

// ControllerGetHandler retrieves a controller.
func ControllerGetHandler(w http.ResponseWriter, r *http.Request) {
  vars              := mux.Vars(r)
  domainName        := vars["domain"]
  controllerName    := vars["controller"]
  controllerVersion := vars["version"]

  // determine controller
  controller, err := model.GetController(domainName, controllerName, controllerVersion)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  // transform controller to string
  yaml, err := controller.Show()
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  // write yaml
  // w.Header().Set("Content-Type", "application/x-yaml")
  io.WriteString(w, yaml)
}

//------------------------------------------------------------------------------

// ControllerDeleteHandler deletes a controller.
func ControllerDeleteHandler(w http.ResponseWriter, r *http.Request) {
  vars              := mux.Vars(r)
  domainName        := vars["domain"]
  controllerName    := vars["controller"]
  controllerVersion := vars["version"]

  // determine domain
  domain, err := model.GetDomain(domainName)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  // delete controller
  err = domain.DeleteController(controllerName, controllerVersion)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }
}

//------------------------------------------------------------------------------

// ControllerResetHandler reset the types of a controller.
func ControllerResetHandler(w http.ResponseWriter, r *http.Request) {
  vars              := mux.Vars(r)
  domainName        := vars["domain"]
  controllerName    := vars["controller"]
  controllerVersion := vars["version"]

  // determine controller
  controller, err := model.GetController(domainName, controllerName, controllerVersion)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  controller.Types = [][2]string{}
}

//------------------------------------------------------------------------------
