package model

import (
	"sync"
	"errors"

	"tsai.eu/solar/util"
)

//------------------------------------------------------------------------------
// Element
// =======
//
// Attributes:
//   - Element
//   - Component
//   - Target
//   - State
//   - Configuration
//   - Endpoint
//   - Clusters
//
// Functions:
//   - NewElement
//
//   - element.Show
//   - element.Load
//   - element.Save
//   - element.Update
//   - element.Reset
//   - element.OK
//   - element.SetState
//
//   - element.ListClusters
//   - element.GetCluster
//   - element.AddCluster
//   - element.DeleteCluster
//------------------------------------------------------------------------------

// Element describes the runtime configuration of a solution element within a domain.
type Element struct {
	Element        string              `yaml:"Element"`             // name of the solution element
	Component      string              `yaml:"Component"`           // type of the solution elmenent
	Target         string              `yaml:"Target"`              // target state of element
	State          string              `yaml:"State"`               // current state of element
	Configuration  string              `yaml:"Configuration"`       // runtime configuration of the solution element
	Endpoint       string              `yaml:"Endpoint"`            // state of the solution element
	Clusters       map[string]*Cluster `yaml:"Clusters"`            // clusters of the solution element
	ClustersX      sync.RWMutex        `yaml:"ClustersX,omitempty"` // mutex for clusters
}

//------------------------------------------------------------------------------

// NewElement creates a new element
func NewElement(name string, component string, configuration string) (*Element, error) {
	var element Element

	element.Element       = name
	element.Component     = component
	element.Configuration = configuration
	element.Endpoint      = ""
	element.Clusters      = map[string]*Cluster{}
	element.ClustersX     = sync.RWMutex{}

	// success
	return &element, nil
}

//------------------------------------------------------------------------------

// Show displays the element information as yaml
func (element *Element) Show() (string, error) {
	return util.ConvertToYAML(element)
}

//------------------------------------------------------------------------------

// Save writes the element as yaml data to a file
func (element *Element) Save(filename string) error {
	return util.SaveYAML(filename, element)
}

//------------------------------------------------------------------------------

// Load reads the element from a file
func (element *Element) Load(filename string) error {
	return util.LoadYAML(filename, element)
}

//------------------------------------------------------------------------------

// ListClusters lists the names of all clusters
func (element *Element) ListClusters() ([]string, error) {
	// collect names
	clusters := []string{}

	element.ClustersX.RLock()
	for cluster := range element.Clusters {
		clusters = append(clusters, cluster)
	}
	element.ClustersX.RUnlock()

	// success
	return clusters, nil
}

//------------------------------------------------------------------------------

// GetCluster retrieves a cluster by name
func (element *Element) GetCluster(name string) (*Cluster, error) {
	// determine dependency
	element.ClustersX.RLock()
	cluster, ok := element.Clusters[name]
	element.ClustersX.RUnlock()

	if !ok {
		return nil, errors.New("cluster not found")
	}

	// success
	return cluster, nil
}

//------------------------------------------------------------------------------

// AddCluster adds a cluster to an element
func (element *Element) AddCluster(cluster *Cluster) {
	element.ClustersX.Lock()
	element.Clusters[cluster.Version] = cluster
	element.ClustersX.Unlock()
}

//------------------------------------------------------------------------------

// DeleteCluster deletes a cluster from an element
func (element *Element) DeleteCluster(version string) {
	element.ClustersX.Lock()
	delete(element.Clusters, version)
	element.ClustersX.Unlock()
}

//------------------------------------------------------------------------------

// Update instantiates/update an element based on an element configuration.
func (element *Element) Update(domainName string, solutionName string, version string, elementConfiguration *ElementConfiguration) error {
	// check if the names are compatible
	if element.Element != elementConfiguration.Element {
		util.LogError("element", "MODEL", "Name of element does match the name of the element configuration")
		return errors.New("Name of element does match the name of the element configuration")
	}

	// check if the components are compatible
	if element.Component != elementConfiguration.Component {
		util.LogError("element", "MODEL", "Type of element does match the type defined in the element configuration")
		return errors.New("Type of element does match the type defined in the element configuration")
	}

	// update target state
	element.Target = ActiveState

	// update configuration
	element.Configuration = elementConfiguration.Configuration

	// update all clusters defined in the element configuration
	clusterNames, _ := elementConfiguration.ListClusters()
	for _, clusterName := range clusterNames {

		cluster, _              := element.GetCluster(clusterName)
		clusterConfiguration, _ := elementConfiguration.GetCluster(clusterName)

		// create new cluster if cluster does not already exist
		if cluster == nil {
			cluster, _ = NewCluster(clusterName, clusterConfiguration.State, clusterConfiguration.Min, clusterConfiguration.Max, clusterConfiguration.Size, "")
			element.AddCluster(cluster)
		}

		// update the element with the configuration information
		if err := cluster.Update(domainName, solutionName, version, element, clusterConfiguration); err != nil {
			util.LogError("element", "MODEL", "Unable to update cluster: '" + clusterName + "' of the element: '" + element.Element + "'\n" + err.Error())
			return err
		}
	}

	// delete all clusters not defined in the element configuration
	clusterNames, _ = element.ListClusters()
	for _, clusterName := range clusterNames {
		cluster, _              := element.GetCluster(clusterName)
		clusterConfiguration, _ := elementConfiguration.GetCluster(clusterName)

		// cluster is not defined in the element configuration
		if clusterConfiguration == nil {
			cluster.Reset()
		}
	}

	// success
	return nil
}

//------------------------------------------------------------------------------

// Reset state of element
func (element *Element) Reset() {
	element.Target = InitialState

	// reset all clusters
	clusterNames, _ := element.ListClusters()
	for _, clusterName := range clusterNames {
		cluster, _ := element.GetCluster(clusterName)

		cluster.Reset()
	}
}

//------------------------------------------------------------------------------

// OK checks if the element has converged to the desired state
func (element *Element) OK() bool {
	// check each cluster
	clusterNames, _ := element.ListClusters()
	for _, clusterName := range clusterNames {
		cluster, _ := element.GetCluster(clusterName)

		if !cluster.OK() {
			return false
		}
	}

	// element is ok
	return true
}

//------------------------------------------------------------------------------

// SetState updates the current state of the element
func (element *Element) SetState(newState string)  {
	if newState == InitialState || newState == InactiveState || newState == ActiveState || newState == FailureState {
		element.State = newState
	}
}

//------------------------------------------------------------------------------
