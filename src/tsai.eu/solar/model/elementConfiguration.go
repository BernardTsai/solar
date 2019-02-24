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

// ClusterConfigurationMap is a synchronized map for a map of cluster configurations
type ClusterConfigurationMap struct {
	*sync.RWMutex                        `yaml:"mutex,omitempty"`    // mutex
	Map map[string]*ClusterConfiguration `yaml:"map"`                // map of cluster configurations
}

// MarshalYAML marshals a ClusterConfigurationMap into yaml
func (m ClusterConfigurationMap) MarshalYAML() (interface{}, error) {
	return m.Map, nil
}

// UnmarshalYAML unmarshals a ClusterConfigurationMap from yaml
func (m *ClusterConfigurationMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	Map := map[string]*ClusterConfiguration{}

	err := unmarshal(&Map)
	if err != nil {
		return err
	}

	*m = ClusterConfigurationMap{Map: Map}

	return nil
}

//------------------------------------------------------------------------------

// ElementConfiguration describes the design time configuration of a solution element within a domain.
type ElementConfiguration struct {
	Element       string                  `yaml:"element"`       // name of the solution element
	Component     string                  `yaml:"component"`     // type of the solution elmenent
	Configuration string                  `yaml:"configuration"` // runtime configuration of the solution element
	Clusters      ClusterConfigurationMap `yaml:"clusters"`      // cluster configurations of the solution element
}

//------------------------------------------------------------------------------

// NewElementConfiguratin creates a new element configuration
func NewElementConfiguratin(name string, component string, configuration string) (*ElementConfiguration, error) {
	var elementConfiguration ElementConfiguration

	elementConfiguration.Element = name
	elementConfiguration.Component = component
	elementConfiguration.Configuration = configuration
	elementConfiguration.Clusters = ClusterConfigurationMap{Map: map[string]*ClusterConfiguration{}}

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

	elementConfiguration.Clusters.RLock()
	for clusterConfiguration := range elementConfiguration.Clusters.Map {
		clusterConfigurations = append(clusterConfigurations, clusterConfiguration)
	}
	elementConfiguration.Clusters.RUnlock()

	// success
	return clusterConfigurations, nil
}

//------------------------------------------------------------------------------

// GetCluster retrieves a cluster configuration by name
func (elementConfiguration *ElementConfiguration) GetCluster(name string) (*ClusterConfiguration, error) {
	// determine dependency
	elementConfiguration.Clusters.RLock()
	clusterConfiguration, ok := elementConfiguration.Clusters.Map[name]
	elementConfiguration.Clusters.RUnlock()

	if !ok {
		return nil, errors.New("cluster configuration not found")
	}

	// success
	return clusterConfiguration, nil
}

//------------------------------------------------------------------------------

// AddCluster adds a cluster configuration to an element
func (elementConfiguration *ElementConfiguration) AddCluster(clusterConfiguration *ClusterConfiguration) {
	elementConfiguration.Clusters.Lock()
	elementConfiguration.Clusters.Map[clusterConfiguration.Version] = clusterConfiguration
	elementConfiguration.Clusters.Unlock()
}

//------------------------------------------------------------------------------

// DeleteCluster deletes a cluster configuration from an element
func (elementConfiguration *ElementConfiguration) DeleteCluster(version string) {
	elementConfiguration.Clusters.Lock()
	delete(elementConfiguration.Clusters.Map, version)
	elementConfiguration.Clusters.Unlock()
}

//------------------------------------------------------------------------------
