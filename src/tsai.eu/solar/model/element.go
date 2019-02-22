package model

import (
	"sync"
	"fmt"

	"github.com/pkg/errors"
	"tsai.eu/solar/util"
)

//------------------------------------------------------------------------------
// Element
// =======
//
// Attributes:
//   - Element
//   - Component
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
//
//   - element.ListClusters
//   - element.GetCluster
//   - element.AddCluster
//   - element.DeleteCluster
//------------------------------------------------------------------------------

// ClusterMap is a synchronized map for a map of clusters
type ClusterMap struct {
	*sync.RWMutex `yaml:"mutex,omitempty"`              // mutex
	Map          map[string]*Cluster      `yaml:"map"` // map of clusters
}

// MarshalYAML marshals a ClusterMap into yaml
func (m ClusterMap) MarshalYAML() (interface{}, error) {
	return m.Map, nil
}

// UnmarshalYAML unmarshals a ClusterMap from yaml
func (m *ClusterMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	Map := map[string]*Cluster{}

	err := unmarshal(&Map)
	if err != nil {
		return err
	}

	*m = ClusterMap{Map: Map}

	return nil
}

//------------------------------------------------------------------------------

// Element describes the runtime configuration of a solution element within a domain.
type Element struct {
	Element       string      `yaml:"element"`       // name of the solution element
	Component     string      `yaml:"component"`     // type of the solution elmenent
	Target        string      `yaml:"target"`        // target state of element
	State         string      `yaml:"state"`         // current state of element
	Configuration string      `yaml:"configuration"` // runtime configuration of the solution element
	Endpoint      string      `yaml:"endpoint"`      // state of the solution element
	Clusters      ClusterMap  `yaml:"clusters"`      // clusters of the solution element
}

//------------------------------------------------------------------------------

// NewElement creates a new element
func NewElement(name string, component string, configuration string) (*Element, error) {
	var element Element

	element.Element = name
	element.Component = component
	element.Configuration = configuration
	element.Endpoint = ""
	element.Clusters = ClusterMap{Map: map[string]*Cluster{}}

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

	element.Clusters.RLock()
	for cluster := range element.Clusters.Map {
		clusters = append(clusters, cluster)
	}
	element.Clusters.RUnlock()

	// success
	return clusters, nil
}

//------------------------------------------------------------------------------

// GetCluster retrieves a cluster by name
func (element *Element) GetCluster(name string) (*Cluster, error) {
	// determine dependency
	element.Clusters.RLock()
	cluster, ok := element.Clusters.Map[name]
	element.Clusters.RUnlock()

	if !ok {
		return nil, errors.New("cluster not found")
	}

	// success
	return cluster, nil
}

//------------------------------------------------------------------------------

// AddCluster adds a cluster to an element
func (element *Element) AddCluster(cluster *Cluster) {
	element.Clusters.Lock()
	element.Clusters.Map[cluster.Version] = cluster
	element.Clusters.Unlock()
}

//------------------------------------------------------------------------------

// DeleteCluster deletes a cluster from an element
func (element *Element) DeleteCluster(version string) {
	element.Clusters.Lock()
	delete(element.Clusters.Map, version)
	element.Clusters.Unlock()
}

//------------------------------------------------------------------------------

// Update instantiates/update an element based on an element configuration.
func (element *Element) Update(elementConfiguration *ElementConfiguration) error {
	// check if the names are compatible
	if element.Element != elementConfiguration.Element {
		return errors.New("Name of element does match the name of the element configuration")
	}

	// check if the components are compatible
	if element.Component != elementConfiguration.Component {
		return errors.New("Type of element does match the type defined in the element configuration")
	}

	// update all clusters defined in the element configuration
	clusterNames, _ := elementConfiguration.ListClusters()
	for _, clusterName := range clusterNames {

		cluster, _              := element.GetCluster(clusterName)
		clusterConfiguration, _ := elementConfiguration.GetCluster(clusterName)

		// cluster already exists
		if cluster != nil {
			if err := cluster.Update(clusterConfiguration); err != nil {
				return fmt.Errorf("Unable to update cluster: '%s' of the element '%s'\n%s", clusterName, element.Element, err)
			}
		} else {
			// cluster does not exist
			// create new cluster
			cluster, _ = NewCluster(clusterName, clusterConfiguration.State, clusterConfiguration.Min, clusterConfiguration.Max, clusterConfiguration.Size, "")
			element.AddCluster(cluster)

			// update the element with the configuration information
			if err := cluster.Update(clusterConfiguration); err != nil {
				return fmt.Errorf("Unable to update cluster: '%s' of the element: '%s'\n%s", clusterName, element.Element, err)
			}
		}
	}

	// success
	return nil
}

//------------------------------------------------------------------------------
