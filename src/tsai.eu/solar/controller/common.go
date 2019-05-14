package controller

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
