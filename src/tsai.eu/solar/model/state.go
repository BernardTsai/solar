package model

//------------------------------------------------------------------------------

// RelationshipState describes the current state of a relationship
type RelationshipState struct {
  Relationship  string  `yaml:"Relationship"`  // name of relationship
  Dependency    string  `yaml:"Dependency"`    // name of dependency
  Configuration string  `yaml:"Configuration"` // configuration information
  Endpoint      string  `yaml:"Endpoint"`      // endpoint information in yaml format
}

//------------------------------------------------------------------------------

// InstanceState describes the current state of an instance
type InstanceState struct {
  Instance string  `yaml:"Instance"`  // id of an instance
  State    string  `yaml:"State"`     // state of an instance
  Endpoint string  `yaml:"Endpoint"`  // endpoint information in yaml format
}

//------------------------------------------------------------------------------

// TargetState describes the desired state and configuration for an instance
type TargetState struct {
	Domain        string              `yaml:"Domain"`        // name of the domain
  Solution      string              `yaml:"Solution"`      // name of solution
	Version       string              `yaml:"Version"`       // version of solution
  Element       string              `yaml:"Element"`       // name of element
  Cluster       string              `yaml:"Cluster"`       // name of cluster
  Instance      string              `yaml:"Instance"`      // name of instance
  Component     string              `yaml:"Component"`     // name of component
  State         string              `yaml:"State"`         // state of instance
  Min           int                 `yaml:"Min"`           // min. size of the solution element cluster
	Max           int                 `yaml:"Max"`           // max. size of the solution element cluster
	Size          int                 `yaml:"Size"`          // size of the solution element cluster
	Configuration string              `yaml:"Configuration"` // configuration of instance
  Relationships []RelationshipState `yaml:"Relationships"` // current state of all relationships
  Instances     []InstanceState     `yaml:"Instances"`     // current state of all instances
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
  Component     string `yaml:"Component"`             // name of component
  State         string `yaml:"State"`                 // state of instance
	Configuration string `yaml:"Configuration"`         // configuration of instance
	Endpoint      string `yaml:"Endpoint"`              // endpoint of instance
}

//------------------------------------------------------------------------------

// GetTargetState determines the desired state of an element, cluster and instance
func GetTargetState(domainName string, solutionName string,  solutionVersion string, elementName string,  clusterName string, instanceName string) (*TargetState, error) {
	targetState := &TargetState{
    Domain:        domainName,
    Solution:      solutionName,
    Version:       solutionVersion,
    Element:       elementName,
    Cluster:       clusterName,
    Instance:      instanceName,
    Component:     "",
    State:         "initial",
    Min:           0,
    Max:           0,
    Size:          0,
    Configuration: "",
    Relationships: []RelationshipState{},
    Instances:     []InstanceState{},
  }

  // determine domain context
  domain, err := GetDomain(domainName)
  if err != nil {
    return targetState, err
  }

  // determine solution context
  solution, err := domain.GetSolution(solutionName)
  if err != nil {
    return targetState, err
  }

  // determine element context
  element, err := solution.GetElement(elementName)
  if err != nil {
    return targetState, err
  }

  // determine cluster context
  cluster, err := element.GetCluster(clusterName)
  if err != nil {
    return targetState, err
  }

  // determine instance context
  instance, err := cluster.GetInstance(instanceName)
  if err != nil {
    return targetState, err
  }

  // update state of target state
  targetState.State = instance.Target

  // determine architecture context
  architecture, err := domain.GetArchitecture(solutionName, solutionVersion)
  if err != nil {
    return targetState, err
  }

  // determine architecture element context
  architectureElement, err := architecture.GetElement(elementName)
  if err != nil {
    return targetState, err
  }

  // determine architecture cluster context
  architectureCluster, err := architectureElement.GetCluster(clusterName)
  if err != nil {
    return targetState, err
  }

  // determine architecture component context
  architectureComponent, err := domain.GetComponent(element.Component, clusterName)
  if err != nil {
    return targetState, err
  }

  // update component of target state
  targetState.Component = architectureComponent.Component

  // update instance configuration
  instance.Configuration = architectureCluster.Configuration

  // update pool dimension and configuration of target state
  targetState.Min           = cluster.Min
  targetState.Max           = cluster.Max
  targetState.Size          = cluster.Size
  targetState.Configuration = instance.Configuration

  // update relationship information
  relationshipNames, err := cluster.ListRelationships()
  if err != nil {
    return targetState, err
  }

  for _, relationshipName := range(relationshipNames) {
    relationship, err := cluster.GetRelationship(relationshipName)
    if err != nil {
      return targetState, err
    }

    // determine current endpoint information of element
    relationshipElement, err := GetElement(domainName, solutionName, relationship.Element)
    if err != nil {
      return targetState, err
    }

    // add relationship information
    targetState.Relationships = append(targetState.Relationships, RelationshipState{
      Relationship:  relationship.Relationship,
      Dependency:    relationship.Endpoint,
      Configuration: relationship.Configuration,
      Endpoint:      relationshipElement.Endpoint,
    })
  }

  // update instance information
  instanceNames, err := cluster.ListInstances()
  if err != nil {
    return targetState, err
  }

  for _, instanceName := range(instanceNames) {
    clusterInstance, err := cluster.GetInstance(instanceName)
    if err != nil {
      return targetState, err
    }

    // add instance information
    targetState.Instances = append(targetState.Instances, InstanceState{
      Instance: instanceName,
      State:    clusterInstance.State,
      Endpoint: clusterInstance.Endpoint,
    })
  }

  // success
  return targetState, nil
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
	instance.State         = currentState.State
  instance.Configuration = currentState.Configuration
  instance.Endpoint      = currentState.Endpoint

	return nil
}

//------------------------------------------------------------------------------
