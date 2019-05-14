package common

//------------------------------------------------------------------------------

// UndefinedState indicates a component state is undefined
const UndefinedState string = "undefined"

// FailureState indicates a component related failure has occured
const FailureState string = "failure"

// InitialState indicates a component is in the initial state
const InitialState string = "initial"

// InactiveState indicates a component is in the inactive state
const InactiveState string = "inactive"

// ActiveState indicates a component is in the active state
const ActiveState string = "active"

// CreatingState indicates a component is in the creating state
const CreatingState string = "creating"

// DestroyingState indicates a component is in the destroying state
const DestroyingState string = "destroying"

// StartingState indicates a component is in the starting state
const StartingState string = "starting"

// StoppingState indicates a component is in the stopping state
const StoppingState string = "stopping"

// ConfiguringState indicates a component is in the configuring state
const ConfiguringState string = "configuring"

// ResettingState indicates a component is in the resetting state
const ResettingState string = "resetting"

//------------------------------------------------------------------------------

// CreateAction requests the instantiation of a resource
const CreateAction string = "create"

// DestroyAction requests the instantiation of a resource
const DestroyAction string = "destroy"

// StartAction requests the activation of a resource
const StartAction string = "start"

// StopAction requests the deactivation of a resource
const StopAction string = "stop"

// ConfigureAction requests the configuration of an inactive resource
const ConfigureAction string = "configure"

// ReconfigureAction requests the reconfiguration of an active resource
const ReconfigureAction string = "reconfigure"

// ResetAction requests the removal of a resource
const ResetAction string = "reset"

// StatusAction requests the status of a resource
const StatusAction string = "status"

//------------------------------------------------------------------------------

// Request sent to controller.
type Request struct {
  Request       string `yaml:"Request"`               // request ID
  Domain        string `yaml:"Domain"`                // name of the domain
  Solution      string `yaml:"Solution"`              // name of solution
	Version       string `yaml:"Version"`               // version of solution
  Element       string `yaml:"Element"`               // name of element
  Cluster       string `yaml:"Cluster"`               // name of cluster
  Instance      string `yaml:"Instance"`              // name of instance
  Component     string `yaml:"Component"`             // component type of instance
  State         string `yaml:"State"`                 // state of instance
	Configuration string `yaml:"Configuration"`         // configuration of instance
}

//------------------------------------------------------------------------------

// Response received from controller.
type Response struct {
  Request       string `yaml:"Request"`               // request ID
  Action        string `yaml:"Action"`                // requested action
  Code          int    `yaml:"Code"`                  // response code
  Status        string `yaml:"Status"`                // status information
  Domain        string `yaml:"Domain"`                // name of the domain
  Solution      string `yaml:"Solution"`              // name of solution
	Version       string `yaml:"Version"`               // version of solution
  Element       string `yaml:"Element"`               // name of element
  Cluster       string `yaml:"Cluster"`               // name of cluster
  Instance      string `yaml:"Instance"`              // name of instance
  Component     string `yaml:"Component"`             // component type of instance
  State         string `yaml:"State"`                 // state of instance
	Configuration string `yaml:"Configuration"`         // configuration of instance
	Endpoint      string `yaml:"Endpoint"`              // endpoint of instance
}

//------------------------------------------------------------------------------
