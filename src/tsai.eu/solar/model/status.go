package model

//------------------------------------------------------------------------------

// Status object received from controller.
type Status struct {
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

// GetStatus derives the status
func GetStatus(domainName string, solutionName string,  elementName string,  clusterName string, instanceName string) (*Status, error) {
	status := Status{}

	// set status information
  status.Domain                  = domainName
  status.Solution                = solutionName
  status.Element                 = elementName
	status.ElementEndpoint         = ""
  status.Cluster                 = clusterName
	status.ClusterEndpoint         = ""
	status.ClusterState            = ""
  status.Instance                = instanceName
	status.InstanceEndpoint        = ""
	status.InstanceState           = ""

  // determine domain context
  domain, err := GetModel().GetDomain(domainName)
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

	status.ElementEndpoint = element.Endpoint

	// determine cluster context
	if clusterName == "" {
		return &status, nil
	}

	cluster, err := element.GetCluster(clusterName)
  if err != nil {
    return nil, err
  }

	status.ClusterEndpoint = cluster.Endpoint
	status.ClusterState    = cluster.State

	// determine instance context
	if instanceName == "" {
		return &status, nil
	}

	instance, err := cluster.GetInstance(instanceName)
  if err != nil {
    return nil, err
  }

	status.InstanceEndpoint = instance.Endpoint
	status.InstanceState    = instance.State

	return &status, nil
}

//------------------------------------------------------------------------------

// SetStatus updates the model with the provided status information.
func SetStatus(status *Status) (error) {
  // determine domain context
  domain, err := GetModel().GetDomain(status.Domain)
  if err != nil {
    return err
  }

  // determine solution context
  solution, err := domain.GetSolution(status.Solution)
  if err != nil {
    return err
  }

	// determine element context
  element, err := solution.GetElement(status.Element)
  if err != nil {
    return err
  }

	element.Endpoint = status.ElementEndpoint

	// determine cluster context
	if status.Cluster == "" {
		return nil
	}

	cluster, err := element.GetCluster(status.Cluster)
  if err != nil {
    return err
  }

	cluster.Endpoint = status.ClusterEndpoint
	cluster.State    = status.ClusterState

	// determine instance context
	if status.Instance == "" {
		return nil
	}

	instance, err := cluster.GetInstance(status.Instance)
  if err != nil {
    return err
  }

	instance.Endpoint = status.InstanceEndpoint
	instance.State    = status.InstanceState

	return nil
}

//------------------------------------------------------------------------------
