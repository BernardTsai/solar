package model

import (
	"sync"

	"github.com/pkg/errors"
	"tsai.eu/solar/util"
)

//------------------------------------------------------------------------------
// ElementConfiguration
// ====================
//
// Attributes:
//   - Element
//   - Component
//   - Configuration
//   - Clusters
//
// Functions:
//   - NewElementConfiguratin
//
//   - elementConfiguration.Show
//   - elementConfiguration.Load
//   - elementConfiguration.Save
//
//   - elementConfiguration.ListClusters
//   - elementConfiguration.GetCluster
//   - elementConfiguration.AddCluster
//   - elementConfiguration.DeleteCluster
//------------------------------------------------------------------------------

// ElementConfiguration describes the design time configuration of a solution element within a domain.
type ElementConfiguration struct {
	Element       string                           `yaml:"Element"`             // name of the solution element
	Component     string                           `yaml:"Component"`           // type of the solution elmenent
	Configuration string                           `yaml:"Configuration"`       // runtime configuration of the solution element
	Clusters      map[string]*ClusterConfiguration `yaml:"Clusters"`            // cluster configurations of the solution element
	ClustersX     sync.RWMutex                     `yaml:"ClustersX,omitempty"` // mutex for cluster configurations
}

//------------------------------------------------------------------------------

// NewElementConfiguration creates a new element configuration
func NewElementConfiguration(name string, component string, configuration string) (*ElementConfiguration, error) {
	var elementConfiguration ElementConfiguration

	elementConfiguration.Element       = name
	elementConfiguration.Component     = component
	elementConfiguration.Configuration = configuration
	elementConfiguration.Clusters      = map[string]*ClusterConfiguration{}
	elementConfiguration.ClustersX     = sync.RWMutex{}

	// success
	return &elementConfiguration, nil
}

//------------------------------------------------------------------------------

// Show displays the element configuration information as yaml
func (elementConfiguration *ElementConfiguration) Show() (string, error) {
	return util.ConvertToYAML(elementConfiguration)
}

//------------------------------------------------------------------------------

// Save writes the element configuration as yaml data to a file
func (elementConfiguration *ElementConfiguration) Save(filename string) error {
	return util.SaveYAML(filename, elementConfiguration)
}

//------------------------------------------------------------------------------

// Load reads the element configuration from a file
func (elementConfiguration *ElementConfiguration) Load(filename string) error {
	return util.LoadYAML(filename, elementConfiguration)
}

//------------------------------------------------------------------------------

// ListClusters lists the names of all cluster configurations
func (elementConfiguration *ElementConfiguration) ListClusters() ([]string, error) {
	// collect names
	clusterConfigurations := []string{}

  elementConfiguration.ClustersX.RLock()
	for clusterConfiguration := range elementConfiguration.Clusters {
		clusterConfigurations = append(clusterConfigurations, clusterConfiguration)
	}
	elementConfiguration.ClustersX.RUnlock()

	// success
	return clusterConfigurations, nil
}

//------------------------------------------------------------------------------

// GetCluster retrieves a cluster configuration by name
func (elementConfiguration *ElementConfiguration) GetCluster(name string) (*ClusterConfiguration, error) {
	// determine dependency
	elementConfiguration.ClustersX.RLock()
	clusterConfiguration, ok := elementConfiguration.Clusters[name]
	elementConfiguration.ClustersX.RUnlock()

	if !ok {
		return nil, errors.New("cluster configuration not found")
	}

	// success
	return clusterConfiguration, nil
}

//------------------------------------------------------------------------------

// AddCluster adds a cluster configuration to an element
func (elementConfiguration *ElementConfiguration) AddCluster(clusterConfiguration *ClusterConfiguration) {
	elementConfiguration.ClustersX.Lock()
	elementConfiguration.Clusters[clusterConfiguration.Version] = clusterConfiguration
	elementConfiguration.ClustersX.Unlock()
}

//------------------------------------------------------------------------------

// DeleteCluster deletes a cluster configuration from an element
func (elementConfiguration *ElementConfiguration) DeleteCluster(version string) {
	elementConfiguration.ClustersX.Lock()
	delete(elementConfiguration.Clusters, version)
	elementConfiguration.ClustersX.Unlock()
}

//------------------------------------------------------------------------------
