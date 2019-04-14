package model

import (
  "errors"

  "tsai.eu/solar/util"
)

//------------------------------------------------------------------------------

// InstanceSetup object describes the setup for an instance of a cluster of a solution element.
type InstanceSetup struct {
  Instance                string  `yaml:"Instance"`                // uuid of the instance
  Target                  string  `yaml:"Target"`                  // target state of the instance
	State                   string  `yaml:"State"`                   // state of the instance
  BaseConfiguration       string  `yaml:"BaseConfiguration"`       // runtime configuration of the instance
  DesignTimeConfiguration string  `yaml:"DesignTimeConfiguration"` // runtime configuration of the instance
	RuntimeConfiguration    string  `yaml:"RuntimeConfiguration"`    // design time configuration of the instance
	Endpoint                string  `yaml:"Endpoint"`                // endpoint of the instance
}

// RelationshipSetup object describes the setup for a relationship.
type RelationshipSetup struct {
  Relationship            string  `yaml:"Relationship"`            // name of the relationship
	Element                 string  `yaml:"Element"`                 // element to which this relationship refers to
	Version                 string  `yaml:"Version"`                 // version of the element to which this relationship refers to
  Target                  string  `yaml:"Target"`                  // target state of the relationship
	State                   string  `yaml:"State"`                   // state of the relationship
  BaseConfiguration       string  `yaml:"BaseConfiguration"`       // runtime configuration of the relationship
  DesignTimeConfiguration string  `yaml:"DesignTimeConfiguration"` // runtime configuration of the relationship
	RuntimeConfiguration    string  `yaml:"RuntimeConfiguration"`    // design time configuration of the relationship
	Endpoint                string  `yaml:"Endpoint"`                // endpoint of the relationship
}

// ClusterSetup object describes the setup for a cluster of a solution element.
type ClusterSetup struct {
  Cluster                 string                        `yaml:"Cluster"`                 // version of the solution element cluster
  Target                  string                        `yaml:"Target"`                  // target state of the cluster
	State                   string                        `yaml:"State"`                   // state of the solution element cluster
  Min                     int                           `yaml:"Min"`                     // min. size of the solution element cluster
  Max                     int                           `yaml:"Max"`                     // max. size of the solution element cluster
	Size                    int                           `yaml:"Size"`                    // size of the solution element cluster
  BaseConfiguration       string                        `yaml:"BaseConfiguration"`       // runtime configuration of the cluster
  DesignTimeConfiguration string                        `yaml:"DesignTimeConfiguration"` // runtime configuration of the cluster
	RuntimeConfiguration    string                        `yaml:"RuntimeConfiguration"`    // design time configuration of the cluster
	Endpoint                string                        `yaml:"Endpoint"`                // endpoint of the solution element cluster
	Relationships           map[string]*RelationshipSetup `yaml:"Relationships"`           // setups for the relationships of the solution element cluster
  Instances               map[string]*InstanceSetup     `yaml:"Instances"`               // setups for the instances of the solution element cluster
}

// ElementSetup object describes the setup for a solution element.
type ElementSetup struct {
  Element                 string                   `yaml:"Element"`                 // name of the solution element
	Component               string                   `yaml:"Component"`               // type of the solution elmenent
  Target                  string                   `yaml:"Target"`                  // target state of the element
	State                   string                   `yaml:"State"`                   // state of the element
  DesignTimeConfiguration string                   `yaml:"DesignTimeConfiguration"` // runtime configuration of the cluster
	RuntimeConfiguration    string                   `yaml:"RuntimeConfiguration"`    // design time configuration of the cluster
	Endpoint                string                   `yaml:"Endpoint"`                // state of the solution element
  Clusters                map[string]*ClusterSetup `yaml:"Clusters"`                // setups for clusters of the solution element
}

// Setup object passed between engine and controller.
type Setup struct {
  Domain                  string                   `yaml:"Domain"`                  // name of the domain
  Solution                string                   `yaml:"Solution"`                // name of solution
	Version                 string                   `yaml:"Version"`                 // version of solution
  Element                 string                   `yaml:"Element"`                 // name of element
  Cluster                 string                   `yaml:"Cluster"`                 // name of cluster
  Instance                string                   `yaml:"Instance"`                // name of instance
  Target                  string                   `yaml:"Target"`                  // target state of the solution
	State                   string                   `yaml:"State"`                   // state of the solution
  DesignTimeConfiguration string                   `yaml:"DesignTimeConfiguration"` // runtime configuration of the solution
	RuntimeConfiguration    string                   `yaml:"RuntimeConfiguration"`    // design time configuration of the solution
	Elements                map[string]*ElementSetup `yaml:"Elements"`                // elements of solution
}

//------------------------------------------------------------------------------

// GetSetup retrieves from the model a setup for the controller.
func GetSetup(domainName string, solutionName string,  solutionVersion string, elementName string,  clusterName string, instanceName string) (*Setup, error) {
	setup := Setup{}

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

  // determine architecture context
  architecture, err := domain.GetArchitecture(solutionName, solutionVersion)
  if err != nil {
    return nil, err
  }

  // set context information
  setup.Domain                  = domainName
  setup.Solution                = solutionName
  setup.Version                 = solutionVersion
  setup.Element                 = elementName
  setup.Cluster                 = clusterName
  setup.Instance                = instanceName
  setup.Target                  = solution.Target
  setup.State                   = solution.State
  setup.DesignTimeConfiguration = architecture.Configuration
  setup.RuntimeConfiguration    = solution.Configuration
  setup.Elements                = map[string]*ElementSetup{}

  // determine element context
  elementNames, _ := solution.ListElements()
  for _, name := range elementNames {
    if elementName == "" || name == elementName {
      setup.Elements[name], err = getElementSetup(&setup, name)
      if err != nil {
        return nil, err
      }
    }
  }

  // check for an undefined element
  if elementName != "" && len(setup.Elements) == 0 {
    return nil, errors.New("unknown element")
  }

  // success
	return &setup, nil
}

//------------------------------------------------------------------------------

// getElementSetup retrieves from the element a setup for the controller.
func getElementSetup(setup *Setup, elementName string) (*ElementSetup, error){
  elementSetup := ElementSetup{}

  // determine domain, solution and architecture context (have been checked before)
  domain, _       := GetDomain(setup.Domain)
  solution, _     := domain.GetSolution(setup.Solution)
  architecture, _ := domain.GetArchitecture(setup.Solution, setup.Version)

  // determine element context
  element, _ := solution.GetElement(elementName)

  // determine element configuration context
  elementConfiguration, err := architecture.GetElement(elementName)
  if err != nil {
    // compensate for undefined configurations
    // e.g. when an element has been removed
    elementConfiguration, _ = NewElementConfiguration(elementName, "", element.Component)
  }

  // set context information
  elementSetup.Element                 = elementName
  elementSetup.Component               = element.Component
  elementSetup.Endpoint                = element.Endpoint
  elementSetup.Target                  = element.Target
  elementSetup.State                   = element.State
  elementSetup.DesignTimeConfiguration = elementConfiguration.Configuration
  elementSetup.RuntimeConfiguration    = element.Configuration
  elementSetup.Clusters                = map[string]*ClusterSetup{}

  // determine cluster context
  clusterNames, _ := element.ListClusters()
  for _, name := range clusterNames {
    if setup.Cluster == "" || name == setup.Cluster {
      elementSetup.Clusters[name], err = getClusterSetup(setup, elementName, name)
      if err != nil {
        return nil, err
      }
    }
  }

  // check for an undefined cluster
  if setup.Cluster != "" && len(elementSetup.Clusters) == 0 {
    return nil, errors.New("unknown cluster")
  }

  // success
	return &elementSetup, nil
}

//------------------------------------------------------------------------------

// getClusterSetup retrieves from the cluster a setup for the controller.
func getClusterSetup(setup *Setup, elementName string, clusterName string) (*ClusterSetup, error){
  clusterSetup := ClusterSetup{}

  // determine domain, solution, element, architecture, elementConfiguration context (have been checked before)
  domain, _               := GetDomain(setup.Domain)
  solution, _             := domain.GetSolution(setup.Solution)
  element, _              := solution.GetElement(elementName)
  architecture, _         := domain.GetArchitecture(setup.Solution, setup.Version)
  elementConfiguration, _ := architecture.GetElement(elementName)

  // determine cluster context
  cluster, err := element.GetCluster(clusterName)
  if err != nil {
    return nil, err
  }

  // determine cluster configuration context
  clusterConfiguration, err := elementConfiguration.GetCluster(clusterName)
  if err != nil {
    // compensate for undefined configurations
    // e.g. when a cluster has been removed
    clusterConfiguration, _ = NewClusterConfiguration(clusterName, InitialState, 0, 0, 0, "")
  }

  // determine component context
  component, err := domain.GetComponent(element.Component, clusterName)
  if err != nil {
    return nil, err
  }

  // set context information
  clusterSetup.Cluster                 = clusterName
  clusterSetup.State                   = cluster.State
  clusterSetup.Min                     = cluster.Min
  clusterSetup.Max                     = cluster.Max
  clusterSetup.Size                    = cluster.Size
  clusterSetup.Target                  = cluster.Target
  clusterSetup.State                   = cluster.State
  clusterSetup.Endpoint                = cluster.Endpoint
  clusterSetup.BaseConfiguration       = component.Configuration
  clusterSetup.DesignTimeConfiguration = clusterConfiguration.Configuration
  clusterSetup.RuntimeConfiguration    = cluster.Configuration
  clusterSetup.Relationships           = map[string]*RelationshipSetup{}
  clusterSetup.Instances               = map[string]*InstanceSetup{}

  // determine relationship context
  relationshipNames, _ := cluster.ListRelationships()
  for _, name := range relationshipNames {
    clusterSetup.Relationships[name], _ = getRelationshipSetup(setup, elementName, clusterName, name)
  }

  // determine instance context
  instanceNames, _ := cluster.ListInstances()
  for _, name := range instanceNames {
    if setup.Instance == "" || name == setup.Instance {
      clusterSetup.Instances[name], _ = getInstanceSetup(setup, elementName, clusterName, name)
    }
  }

  // check for an undefined instance
  if setup.Instance != "" && len(clusterSetup.Instances) == 0 {
    return nil, errors.New("unknown instance")
  }

  // success
	return &clusterSetup, nil
}

//------------------------------------------------------------------------------

// getRelationshipSetup retrieves from the relationship a setup for the controller.
func getRelationshipSetup(setup *Setup, elementName string, clusterName string, relationshipName string) (*RelationshipSetup, error){
  relationshipSetup := RelationshipSetup{}

  // determine domain, solution, element, cluster,
  // architecture, elementConfiguration, clusterConfiguration and
  // component context (have been checked before)
  domain, _               := GetDomain(setup.Domain)
  solution, _             := domain.GetSolution(setup.Solution)
  element, _              := solution.GetElement(elementName)
  cluster, _              := element.GetCluster(clusterName)
  architecture, _         := domain.GetArchitecture(setup.Solution, setup.Version)
  elementConfiguration, _ := architecture.GetElement(elementName)
  clusterConfiguration, _ := elementConfiguration.GetCluster(clusterName)
  component, _            := domain.GetComponent(element.Component, clusterName)

  // determine relationship context
  relationship, _ := cluster.GetRelationship(relationshipName)

  // determine relationship configuration context
  relationshipConfiguration, err := clusterConfiguration.GetRelationship(relationshipName)
  if err != nil {
    // compensate for undefined configurations
    // e.g. when a relationship has been removed
    relationshipConfiguration, _ = NewRelationshipConfiguration(relationshipName, "", "", "", "", "")
  }

  // determine dependency context
  dependency, err := component.GetDependency(relationship.Dependency)
  if err != nil {
    util.LogError("setup", "CORE", "unable to determine dependency: " + relationship.Dependency)
    return nil, err
  }

  // set context information
  relationshipSetup.Relationship            = relationshipName
  relationshipSetup.Element                 = relationship.Element
  relationshipSetup.Version                 = relationship.Version
  relationshipSetup.Target                  = relationship.Target
  relationshipSetup.State                   = relationship.State
  relationshipSetup.Endpoint                = relationship.Endpoint
  relationshipSetup.BaseConfiguration       = dependency.Configuration
  relationshipSetup.DesignTimeConfiguration = relationshipConfiguration.Configuration
  relationshipSetup.RuntimeConfiguration    = relationship.Configuration

  // success
	return &relationshipSetup, nil
}

//------------------------------------------------------------------------------

// getInstanceSetup retrieves from the relationship a setup for the controller.
func getInstanceSetup(setup *Setup, elementName string, clusterName string, instanceName string) (*InstanceSetup, error){
  instanceSetup := InstanceSetup{}

  // determine domain, solution, element, cluster,
  // architecture, elementConfiguration, clusterConfiguration and
  // component context (have been checked before)
  domain, _               := GetDomain(setup.Domain)
  solution, _             := domain.GetSolution(setup.Solution)
  element, _              := solution.GetElement(elementName)
  cluster, _              := element.GetCluster(clusterName)
  architecture, _         := domain.GetArchitecture(setup.Solution, setup.Version)
  elementConfiguration, _ := architecture.GetElement(elementName)
  clusterConfiguration, _ := elementConfiguration.GetCluster(clusterName)
  component, _            := domain.GetComponent(element.Component, clusterName)

  // determine instance context
  instance, _ := cluster.GetInstance(instanceName)

  // set context information
  instanceSetup.Instance                = instanceName
  instanceSetup.Target                  = instance.Target
  instanceSetup.State                   = instance.State
  instanceSetup.Endpoint                = instance.Endpoint
  instanceSetup.BaseConfiguration       = component.Configuration
  instanceSetup.DesignTimeConfiguration = clusterConfiguration.Configuration
  instanceSetup.RuntimeConfiguration    = instance.Configuration

  // success
	return &instanceSetup, nil
}

//------------------------------------------------------------------------------

// SetSetup updates the model with the provided information.
func SetSetup(setup *Setup) (error) {
  // determine element context
  for _, elementSetup := range setup.Elements {
    err := setElementSetup(setup, elementSetup)
    if err != nil {
      return err
    }
  }

  return nil
}

//------------------------------------------------------------------------------

// setElementSetup updates the element setup with the provided information.
func setElementSetup(setup *Setup, elementSetup *ElementSetup) (error) {
  // determine domain context
  domain, err := GetDomain(setup.Domain)
  if err != nil {
    return nil
  }

  // determine solution context
  solution, err := domain.GetSolution(setup.Solution)
  if err != nil {
    return nil
  }

  // determine element context
  element, err := solution.GetElement(elementSetup.Element)
  if err != nil {
    return nil
  }

  // update runtime configuration
  element.State = elementSetup.State

  // determine cluster context
  for _, clusterSetup := range elementSetup.Clusters {
    err := setClusterSetup(setup, elementSetup, clusterSetup)
    if err != nil {
      return err
    }
  }

  return nil
}

//------------------------------------------------------------------------------

// setClusterSetup updates the cluster setup with the provided information.
func setClusterSetup(setup *Setup, elementSetup *ElementSetup, clusterSetup *ClusterSetup) (error) {
  // determine domain, solution and element context
  domain, _   := GetDomain(setup.Domain)
  solution, _ := domain.GetSolution(setup.Solution)
  element, _  := solution.GetElement(elementSetup.Element)

  // determine cluster context
  cluster, err := element.GetCluster(clusterSetup.Cluster)
  if err != nil {
    return nil
  }

  // update runtime configuration
  cluster.State         = clusterSetup.State
  cluster.Configuration = clusterSetup.RuntimeConfiguration
  cluster.Endpoint      = clusterSetup.Endpoint
  cluster.State         = clusterSetup.State
  cluster.Size          = clusterSetup.Size

  // determine relationship context
  for _, relationshipSetup := range clusterSetup.Relationships {
    err := setRelationshipSetup(setup, elementSetup, clusterSetup, relationshipSetup)
    if err != nil {
      return err
    }
  }

  // determine instance context
  for _, instanceSetup := range clusterSetup.Instances {
    err := setInstanceSetup(setup, elementSetup, clusterSetup, instanceSetup)
    if err != nil {
      return err
    }
  }

  return nil
}

//------------------------------------------------------------------------------

// setRelationshipSetup updates the relationship setup with the provided information.
func setRelationshipSetup(setup *Setup, elementSetup *ElementSetup, clusterSetup *ClusterSetup, relationshipSetup *RelationshipSetup) (error) {
  // determine domain, solution, element and cluster context
  domain, _   := GetDomain(setup.Domain)
  solution, _ := domain.GetSolution(setup.Solution)
  element, _  := solution.GetElement(elementSetup.Element)
  cluster, _  := element.GetCluster(clusterSetup.Cluster)

  // determine relationship context
  relationship, err := cluster.GetRelationship(relationshipSetup.Relationship)
  if err != nil {
    return nil
  }

  // update runtime configuration
  relationship.State         = relationshipSetup.State
  relationship.Configuration = relationshipSetup.RuntimeConfiguration
  relationship.Endpoint      = relationshipSetup.Endpoint

  // success
  return nil
}

//------------------------------------------------------------------------------

// setInstanceSetup updates the instance setup with the provided information.
func setInstanceSetup(setup *Setup, elementSetup *ElementSetup, clusterSetup *ClusterSetup, instanceSetup *InstanceSetup) (error) {
  // determine domain, solution, element and cluster context
  domain, _   := GetDomain(setup.Domain)
  solution, _ := domain.GetSolution(setup.Solution)
  element, _  := solution.GetElement(elementSetup.Element)
  cluster, _  := element.GetCluster(clusterSetup.Cluster)

  // determine instance context
  instance, err := cluster.GetInstance(instanceSetup.Instance)
  if err != nil {
    return nil
  }

  // update runtime configuration
  instance.State         = instanceSetup.State
  instance.Configuration = instanceSetup.RuntimeConfiguration
  instance.Endpoint      = instanceSetup.Endpoint

  // success
  return nil
}

//------------------------------------------------------------------------------
