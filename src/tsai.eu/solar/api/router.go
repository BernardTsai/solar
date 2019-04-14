package api

import (
  "context"
  "net/http"
  "github.com/gorilla/mux"

  "tsai.eu/solar/util"
)

//------------------------------------------------------------------------------

// API represents the web interface
type API struct {
  Server *http.Server     // web server
}

//------------------------------------------------------------------------------

// redirect forwards requests to main entry URL
func redirect(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "/solar/index.html", 301)
}

//------------------------------------------------------------------------------

// Start starts the web interface
func Start(ctx context.Context) (*API) {
  api := API{}

  // start web interface
  api.Server = NewServer()

  // create a process to check if the server needs to shutdown
  go func() {
    select {
    case <-ctx.Done():
      util.LogInfo("main", "API", "api initial")
      api.Server.Shutdown(context.Background())
    }
  }()

  // success
  return &api
}

//------------------------------------------------------------------------------

// NewServer creates and starts the API
func NewServer() *http.Server{
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

  // catalog
  router.HandleFunc("/catalog/{domain}", CatalogGetHandler).Methods("GET")

  // component
  router.HandleFunc("/component/{domain}",                       ComponentListHandler).Methods("GET")
  router.HandleFunc("/component/{domain}",                       ComponentSetHandler).Methods("POST")
  router.HandleFunc("/component/{domain}/{component}/{version}", ComponentGetHandler).Methods("GET")
  router.HandleFunc("/component/{domain}/{component}/{version}", ComponentDeleteHandler).Methods("DELETE")

  // architecture
  router.HandleFunc("/architecture/{domain}",                          ArchitectureListHandler).Methods("GET")
  router.HandleFunc("/architecture/{domain}",                          ArchitectureSetHandler).Methods("POST")
  router.HandleFunc("/architecture/{domain}/{architecture}/{version}", ArchitectureGetHandler).Methods("GET")
  router.HandleFunc("/architecture/{domain}/{architecture}/{version}", ArchitectureDeleteHandler).Methods("DELETE")
  router.HandleFunc("/architecture/{domain}/{architecture}/{version}", ArchitectureDeployHandler).Methods("POST")

  // solution
  router.HandleFunc("/solution/{domain}",                       SolutionListHandler).Methods("GET")
  router.HandleFunc("/solution/{domain}",                       SolutionSetHandler).Methods("POST")
  router.HandleFunc("/solution/{domain}/{solution}",            SolutionGetHandler).Methods("GET")
  router.HandleFunc("/solution/{domain}/{solution}",            SolutionDeleteHandler).Methods("DELETE")
  router.HandleFunc("/solution/{domain}/{solution}/{version}",  SolutionDeployHandler).Methods("POST")

  // cluster
  router.HandleFunc("/cluster/{domain}/{solution}/{element}/{cluster}", ClusterUpdateHandler).Methods("PUT")

  // instance
  router.HandleFunc("/instance/{domain}/{solution}/{element}/{cluster}/{instance}", InstanceUpdateHandler).Methods("PUT")

  // task
  router.HandleFunc("/tasks/{domain}/{solution}/{element}/{cluster}/{instance}", TaskListHandler).Methods("GET")
  router.HandleFunc("/tasks/{domain}/{solution}/{element}/{cluster}",            TaskListHandler).Methods("GET")
  router.HandleFunc("/tasks/{domain}/{solution}/{element}",                      TaskListHandler).Methods("GET")
  router.HandleFunc("/tasks/{domain}/{solution}",                                TaskListHandler).Methods("GET")
  router.HandleFunc("/tasks/{domain}",                                           TaskListHandler).Methods("GET")

  router.HandleFunc("/task/{domain}/{task}",                                     TaskTraceHandler).Methods("GET")
  router.HandleFunc("/task/{domain}/{task}/{level}",                             TaskGetHandler).Methods("GET")
  router.HandleFunc("/task/{domain}/{task}",                                     TaskTerminateHandler).Methods("DELETE")

  router.HandleFunc("/solar", redirect).Methods("GET")
  router.HandleFunc("/",      redirect).Methods("GET")

  // static files
  router.PathPrefix("/solar/").Handler(http.StripPrefix("/solar/", http.FileServer(http.Dir("./static/"))))

  // start processing
  http.Handle("/", router)

  srv := &http.Server{Addr: ":80"}

  go srv.ListenAndServe()

  // initialisation completed
  util.LogInfo("main", "API", "api active")
  return srv
}

//------------------------------------------------------------------------------
