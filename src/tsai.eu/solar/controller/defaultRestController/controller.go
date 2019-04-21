package defaultRestController

import (
  "context"
  "errors"
  "io/ioutil"
  "net/http"
  "github.com/gorilla/mux"
  "gopkg.in/yaml.v2"
)

//------------------------------------------------------------------------------

// Controller manages the lifecycle of nothing
type Controller struct {
  Router *mux.Router
  Server *http.Server     // web server
}

//------------------------------------------------------------------------------

// TargetState sent to controller.
type TargetState struct {
	Domain                string `yaml:"Domain"`                // name of the domain
  Solution              string `yaml:"Solution"`              // name of solution
	Version               string `yaml:"Version"`               // version of solution
  Element               string `yaml:"Element"`               // name of element
	ElementConfiguration  string `yaml:"ElementConfiguration"`  // endpoint of element
  Cluster               string `yaml:"Cluster"`               // name of cluster
	ClusterConfiguration  string `yaml:"ClusterConfiguration"`  // endpoint of cluster
  ClusterState          string `yaml:"ClusterState"`     // state of cluster
  Instance              string `yaml:"Instance"`              // name of instance
	InstanceConfiguration string `yaml:"InstanceConfiguration"` // endpoint of instance
  InstanceState         string `yaml:"InstanceState"`    // state of instance
}

//------------------------------------------------------------------------------

// CurrentState received from controller.
type CurrentState struct {
	Domain           string `yaml:"Domain"`           // name of the domain
  Solution         string `yaml:"Solution"`         // name of solution
	Version          string `yaml:"Version"`          // version of solution
  Element          string `yaml:"Element"`          // name of element
	ElementEndpoint  string `yaml:"ElementEndpoint"`  // endpoint of element
  Cluster          string `yaml:"Cluster"`          // name of cluster
	ClusterEndpoint  string `yaml:"ClusterEndpoint"`  // endpoint of cluster
	ClusterState     string `yaml:"ClusterState"`     // state of cluster
  Instance         string `yaml:"Instance"`         // name of instance
	InstanceEndpoint string `yaml:"InstanceEndpoint"` // endpoint of instance
	InstanceState    string `yaml:"InstanceState"`    // state of instance
}

//------------------------------------------------------------------------------

// NewController creates a new controller
func NewController() (*Controller) {
  controller := Controller{}

  // create router
  controller.Router = mux.NewRouter()

  controller.Router.HandleFunc("/",            check).Methods("GET")
  controller.Router.HandleFunc("/create",      controller.Create).Methods("POST")
  controller.Router.HandleFunc("/destroy",     controller.Destroy).Methods("POST")
  controller.Router.HandleFunc("/start",       controller.Start).Methods("POST")
  controller.Router.HandleFunc("/stop",        controller.Stop).Methods("POST")
  controller.Router.HandleFunc("/configure",   controller.Configure).Methods("POST")
  controller.Router.HandleFunc("/reconfigure", controller.Reconfigure).Methods("POST")
  controller.Router.HandleFunc("/reset",       controller.Reset).Methods("POST")
  controller.Router.HandleFunc("/status",      controller.Status).Methods("POST")

  // create server
  controller.Server = &http.Server{Addr: ":10000", Handler: controller.Router}

  // success
  return &controller
}

//------------------------------------------------------------------------------

// Run executes the server
func (c *Controller) Run(ctx context.Context) {
  // start web interface
  go c.Server.ListenAndServe()

  // create a process to check if the server needs to shutdown
  go func() {
    select {
    case <-ctx.Done():
      c.Server.Shutdown(context.Background())
    }
  }()
}

//------------------------------------------------------------------------------

// check allows to ping the server
func check(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(200)
}

//------------------------------------------------------------------------------

// readTargetState reads the target state from a request body
func readTargetState(r *http.Request) (targetState TargetState, err error) {
  targetState = TargetState{}

  // handle any issues with the conversion
  defer func() {
    if r := recover(); r != nil {
      err = errors.New("Unable to retrieve target state")
    }
  }()

  // get body
	body, _ := ioutil.ReadAll(r.Body)

  // convert body into target state
  err = yaml.Unmarshal(body, &targetState)
  if err != nil {
    return targetState, errors.New("Unable to convert body into target state")
  }

  // success
  return targetState, nil
}

//------------------------------------------------------------------------------

// writeCurrentState writes the current state
func writeCurrentState(w http.ResponseWriter, currentState CurrentState) error {
  yaml, err := yaml.Marshal(currentState)
  if err != nil {
    w.WriteHeader(500)
    w.Write([]byte("Unable to write current state"))
  }

  // return results
  w.WriteHeader(200)
  w.Write(yaml)

  // return err status
  return err
}

//------------------------------------------------------------------------------
