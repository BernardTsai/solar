package model

import (
	"sync"

	"github.com/pkg/errors"
	"tsai.eu/solar/util"
)

//------------------------------------------------------------------------------
// ClusterConfiguration
// ====================
//
// Attributes:
//   - Version
//   - State
//   - Size
//   - Configuration
//   - Relationships
//
// Functions:
//   - NewClusterConfiguration
//
//   - clusterConfiguration.Show
//   - clusterConfiguration.Load
//   - clusterConfiguration.Save
//
//   - clusterConfiguration.ListRelationship
//   - clusterConfiguration.GetRelationship
//   - clusterConfiguration.AddRelationship
//   - clusterConfiguration.DeleteRelationship
//------------------------------------------------------------------------------

// RelationshipConfigurationMap is a synchronized map for a map of relationship configurations
type RelationshipConfigurationMap struct {
	*sync.RWMutex                             `yaml:"mutex,omitempty"` // mutex
	Map map[string]*RelationshipConfiguration `yaml:"map"`             // map of relationships
}

// MarshalYAML marshals a RelationshipConfigurationMap into yaml
func (m RelationshipConfigurationMap) MarshalYAML() (interface{}, error) {
	return m.Map, nil
}

// UnmarshalYAML unmarshals a RelationshipConfigurationMap from yaml
func (m *RelationshipConfigurationMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	Map := map[string]*RelationshipConfiguration{}

	err := unmarshal(&Map)
	if err != nil {
		return err
	}

	*m = RelationshipConfigurationMap{Map: Map}

	return nil
}

//------------------------------------------------------------------------------

// ClusterConfiguration describes the design time configuration of a solution element cluster within a domain.
type ClusterConfiguration struct {
	Version       string                       `yaml:"version"`       // version of the solution element cluster
	State         string                       `yaml:"state"`         // state of the solution element cluster
	Min           int                          `yaml:"min"`           // min. size of the solution element cluster
	Max           int                          `yaml:"max"`           // max. size of the solution element cluster
	Size          int                          `yaml:"size"`          // size of the solution element cluster
	Configuration string                       `yaml:"configuration"` // runtime configuration of the solution element cluster
	Relationships RelationshipConfigurationMap `yaml:"relationships"` // relationships of the solution element cluster
}

//------------------------------------------------------------------------------

// NewClusterConfiguration creates a new cluster configuration
func NewClusterConfiguration(version string, state string, min int, max int, size int, configuration string) (*ClusterConfiguration, error) {
	var clusterConfiguration ClusterConfiguration

	clusterConfiguration.Version       = version
	clusterConfiguration.State         = state
	clusterConfiguration.Min           = min
	clusterConfiguration.Max           = max
	clusterConfiguration.Size          = size
	clusterConfiguration.Configuration = configuration
	clusterConfiguration.Relationships = RelationshipConfigurationMap{Map: map[string]*RelationshipConfiguration{}}

	// success
	return &clusterConfiguration, nil
}

//------------------------------------------------------------------------------

// Show displays the cluster configuration information as yaml
func (clusterConfiguration *ClusterConfiguration) Show() (string, error) {
	return util.ConvertToYAML(clusterConfiguration)
}

//------------------------------------------------------------------------------

// Save writes the cluster configuration as yaml data to a file
func (clusterConfiguration *ClusterConfiguration) Save(filename string) error {
	return util.SaveYAML(filename, clusterConfiguration)
}

//------------------------------------------------------------------------------

// Load reads the cluster configuration from a file
func (clusterConfiguration *ClusterConfiguration) Load(filename string) error {
	return util.LoadYAML(filename, clusterConfiguration)
}

//------------------------------------------------------------------------------

// ListRelationships lists the names of all defined relationship configurations
func (clusterConfiguration *ClusterConfiguration) ListRelationships() ([]string, error) {
	// collect names
	relationships := []string{}

	clusterConfiguration.Relationships.RLock()
	for relationship := range clusterConfiguration.Relationships.Map {
		relationships = append(relationships, relationship)
	}
	clusterConfiguration.Relationships.RUnlock()

	// success
	return relationships, nil
}

//------------------------------------------------------------------------------

// GetRelationship retrieves a relationship configuration by name
func (clusterConfiguration *ClusterConfiguration) GetRelationship(name string) (*RelationshipConfiguration, error) {
	// determine relationship configuration
	clusterConfiguration.Relationships.RLock()
	relationship, ok := clusterConfiguration.Relationships.Map[name]
	clusterConfiguration.Relationships.RUnlock()

	if !ok {
		return nil, errors.New("relationship configuration not found")
	}

	// success
	return relationship, nil
}

//------------------------------------------------------------------------------

// AddRelationship adds a relationship configuration to a cluster
func (clusterConfiguration *ClusterConfiguration) AddRelationship(relationshipConfiguration *RelationshipConfiguration) {
	clusterConfiguration.Relationships.Lock()
	clusterConfiguration.Relationships.Map[relationshipConfiguration.Relationship] = relationshipConfiguration
	clusterConfiguration.Relationships.Unlock()
}

//------------------------------------------------------------------------------

// DeleteRelationship deletes a relationship configuration from a cluster
func (clusterConfiguration *ClusterConfiguration) DeleteRelationship(name string) {
	clusterConfiguration.Relationships.Lock()
	delete(clusterConfiguration.Relationships.Map, name)
	clusterConfiguration.Relationships.Unlock()
}

//------------------------------------------------------------------------------
