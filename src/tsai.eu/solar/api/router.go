package api

import (
  "net/http"
  "github.com/gorilla/mux"
)

//------------------------------------------------------------------------------

// NewRouter creates and starts the API
func NewRouter() {
  router := mux.NewRouter()

  // model
  router.HandleFunc("/model", ModelSetHandler).Methods("POST")
  router.HandleFunc("/model", ModelGetHandler).Methods("GET")
  router.HandleFunc("/model", ModelResetHandler).Methods("PUT")

  // domain
  router.HandleFunc("/domain",          DomainListHandler).Methods("GET")
  router.HandleFunc("/domain/{domain}", DomainCreateHandler).Methods("POST")
  router.HandleFunc("/domain/{domain}", DomainDeleteHandler).Methods("DELETE")
  router.HandleFunc("/domain",          DomainSetHandler).Methods("POST")
  router.HandleFunc("/domain/{domain}", DomainGetHandler).Methods("GET")
  router.HandleFunc("/domain/{domain}", DomainResetHandler).Methods("PUT")

  // component
  router.HandleFunc("/component/{domain}",             ComponentListHandler).Methods("GET")
  router.HandleFunc("/component/{domain}",             ComponentSetHandler).Methods("POST")
  router.HandleFunc("/component/{domain}/{component}", ComponentGetHandler).Methods("GET")
  router.HandleFunc("/component/{domain}/{component}", ComponentDeleteHandler).Methods("DELETE")

  // architecture
  router.HandleFunc("/architecture/{domain}",                ArchitectureListHandler).Methods("GET")
  router.HandleFunc("/architecture/{domain}",                ArchitectureSetHandler).Methods("POST")
  router.HandleFunc("/architecture/{domain}/{architecture}", ArchitectureGetHandler).Methods("GET")
  router.HandleFunc("/architecture/{domain}/{architecture}", ArchitectureDeleteHandler).Methods("DELETE")

  // solution
  router.HandleFunc("/solution/{domain}",                       SolutionListHandler).Methods("GET")
  router.HandleFunc("/solution/{domain}",                       SolutionSetHandler).Methods("POST")
  router.HandleFunc("/solution/{domain}/{solution}",            SolutionGetHandler).Methods("GET")
  router.HandleFunc("/solution/{domain}/{solution}",            SolutionDeleteHandler).Methods("DELETE")
  router.HandleFunc("/solution/{domain}/{solution}/{version}",  SolutionDeployHandler).Methods("POST")

  // task
  router.HandleFunc("/task/{domain}/{solution}/{element}/{cluster}/{instance}", TaskListHandler).Methods("GET")
  router.HandleFunc("/task/{domain}/{task}/{level}",                            TaskGetHandler).Methods("GET")

  // start processing
  http.Handle("/", router)
  http.ListenAndServe(":80", nil)
}

//------------------------------------------------------------------------------