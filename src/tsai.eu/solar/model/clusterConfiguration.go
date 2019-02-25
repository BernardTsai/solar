package model

import (
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
	Version       string                       `yaml:"Version"`       // version of the solution element cluster
	State         string                       `yaml:"State"`         // state of the solution element cluster
	Min           int                          `yaml:"Min"`           // min. size of the solution element cluster
	Max           int                          `yaml:"Max"`           // max. size of the solution element cluster
	Size          int                          `yaml:"Size"`          // size of the solution element cluster
	Configuration string                       `yaml:"Configuration"` // runtime configuration of the solution element cluster
	Relationships RelationshipConfigurationMap `yaml:"Relationships"` // relationships of the solution element cluster
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

	for relationship := range clusterConfiguration.Relationships.Map {
		relationships = append(relationships, relationship)
	}

	// success
	return relationships, nil
}

//------------------------------------------------------------------------------

// GetRelationship retrieves a relationship configuration by name
func (clusterConfiguration *ClusterConfiguration) GetRelationship(name string) (*RelationshipConfiguration, error) {
	// determine relationship configuration
	relationship, ok := clusterConfiguration.Relationships.Map[name]

	if !ok {
		return nil, errors.New("relationship configuration not found")
	}

	// success
	return relationship, nil
}

//------------------------------------------------------------------------------

// AddRelationship adds a relationship configuration to a cluster
func (clusterConfiguration *ClusterConfiguration) AddRelationship(relationshipConfiguration *RelationshipConfiguration) {
	clusterConfiguration.Relationships.Map[relationshipConfiguration.Relationship] = relationshipConfiguration
}

//------------------------------------------------------------------------------

// DeleteRelationship deletes a relationship configuration from a cluster
func (clusterConfiguration *ClusterConfiguration) DeleteRelationship(name string) {
	delete(clusterConfiguration.Relationships.Map, name)
}

//------------------------------------------------------------------------------
