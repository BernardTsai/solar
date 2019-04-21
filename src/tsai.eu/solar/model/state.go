package model

import (
  "github.com/cbroglie/mustache"
  "gopkg.in/yaml.v2"
)

//------------------------------------------------------------------------------

// YAML supports a mapping to an unknown yaml schema
type YAML struct {
}

//------------------------------------------------------------------------------

// TargetState describes the desired state and configuration for an instance
type TargetState struct {
	Domain        string `yaml:"Domain"`                // name of the domain
  Solution      string `yaml:"Solution"`              // name of solution
	Version       string `yaml:"Version"`               // version of solution
  Element       string `yaml:"Element"`               // name of element
  Cluster       string `yaml:"Cluster"`               // name of cluster
  Instance      string `yaml:"Instance"`              // name of instance
  State         string `yaml:"State"`                 // state of instance
	Configuration string `yaml:"Configuration"`         // configuration of instance
}

//------------------------------------------------------------------------------

// CurrentState describes the current state and configuration of an instance
type CurrentState struct {
  Domain        string `yaml:"Domain"`                // name of the domain
  Solution      string `yaml:"Solution"`              // name of solution
	Version       string `yaml:"Version"`               // version of solution
  Element       string `yaml:"Element"`               // name of element
  Cluster       string `yaml:"Cluster"`               // name of cluster
  Instance      string `yaml:"Instance"`              // name of instance
  State         string `yaml:"State"`                 // state of instance
	Configuration string `yaml:"Configuration"`         // configuration of instance
	Endpoint      string `yaml:"Endpoint"`              // endpoint of instance
}

//------------------------------------------------------------------------------

// GetTargetState determines the desired state of an element, cluster and instance
func GetTargetState(domainName string, solutionName string,  solutionVersion string, elementName string,  clusterName string, instanceName string) (*TargetState, error) {
	targetState := TargetState{}

  // determine domain context
  domain, err := GetDomain(domainName)
  if err != nil {
    return nil, err
  }

  // determine solution context
  solution, err := domain.GetSolution(solutionName)
  if err != nil {
    return nil, err
  }

  // determine element context
  element, err := solution.GetElement(elementName)
  if err != nil {
    return nil, err
  }

  // determine cluster context
  cluster, err := element.GetCluster(clusterName)
  if err != nil {
    return nil, err
  }

  // determine instance context
  instance, err := cluster.GetInstance(instanceName)
  if err != nil {
    return nil, err
  }

  // determine architecture context
  architecture, err := domain.GetArchitecture(solutionName, solutionVersion)
  if err != nil {
    return nil, err
  }

  // determine architecture element context
  architectureElement, err := architecture.GetElement(elementName)
  if err != nil {
    return nil, err
  }

  // determine architecture cluster context
  architectureCluster, err := architectureElement.GetCluster(clusterName)
  if err != nil {
    return nil, err
  }

  // determine architecture component context
  architectureComponent, err := domain.GetComponent(element.Component, clusterName)
  if err != nil {
    return nil, err
  }

  // render instance configuration
  template   := architectureComponent.Configuration
  parameters := map[string]interface{}{}

  // convert parameters
  err = yaml.Unmarshal([]byte(architectureCluster.Configuration), &parameters )
  if err != nil {
    return nil, err
  }

  // render
  configuration, err := mustache.Render(template, parameters)
  if err != nil {
    return nil, err
  }

  // update instance configuration
  instance.Configuration = configuration

  // set context information
  targetState.Domain        = domainName
  targetState.Solution      = solutionName
  targetState.Version       = solutionVersion
  targetState.Element       = elementName
  targetState.Cluster       = clusterName
  targetState.Instance      = instanceName
  targetState.State         = instance.Target
  targetState.Configuration = instance.Configuration

  // success
  return &targetState, nil
}

//------------------------------------------------------------------------------

// SetCurrentState updates the model with the provided state information.
func SetCurrentState(currentState *CurrentState) (error) {
	// determine instance context
	instance, err := GetInstance(currentState.Domain, currentState.Solution, currentState.Element, currentState.Cluster, currentState.Instance)
  if err != nil {
    return err
  }

  // update state and endpoint of instance
	instance.State    = currentState.State
  instance.Endpoint = currentState.Endpoint

	return nil
}

//------------------------------------------------------------------------------