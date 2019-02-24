package model

import (
	"tsai.eu/solar/util"
)

//------------------------------------------------------------------------------
// RelationshipConfiguration
// =========================
//
// Attributes:
//   - Relationship
//   - Dependency
//   - Type
//   - Element
//   - Version
//   - Configuration
//
// Functions:
//   - NewRelationshipConfiguration
//
//   - relationship.Show
//   - relationship.Load
//   - relationship.Save
//------------------------------------------------------------------------------

// RelationshipConfiguration describes the design time configuration of a relationship between clusters within a domain.
type RelationshipConfiguration struct {
	Relationship  string  `yaml:"relationship"`  // name of the relationship
	Dependency    string  `yaml:"dependency"`    // name of the dependency
	Type          string  `yaml:"type"`          // type of dependency
	Element       string  `yaml:"element"`       // element to which this relationship refers to
	Version       string  `yaml:"version"`       // version of the element to which this relationship refers to
	Configuration string  `yaml:"configuration"` // design time configuration of the relationship
}

//------------------------------------------------------------------------------

// NewRelationshipConfiguration creates a new relationship configuration
func NewRelationshipConfiguration(name string, dependency string, dependencyType string, element string, version string, configuration string) (*RelationshipConfiguration, error) {
	var relationshipConfiguration RelationshipConfiguration

	relationshipConfiguration.Relationship  = name
	relationshipConfiguration.Dependency    = dependency
	relationshipConfiguration.Type          = dependencyType
	relationshipConfiguration.Element       = element
	relationshipConfiguration.Version       = version
	relationshipConfiguration.Configuration = configuration

	// success
	return &relationshipConfiguration, nil
}

//------------------------------------------------------------------------------

// Show displays the relationship configuration information as yaml
func (relationshipConfiguration *RelationshipConfiguration) Show() (string, error) {
	return util.ConvertToYAML(relationshipConfiguration)
}

//------------------------------------------------------------------------------

// Save writes the relationship configuration as yaml data to a file
func (relationshipConfiguration *RelationshipConfiguration) Save(filename string) error {
	return util.SaveYAML(filename, relationshipConfiguration)
}

//------------------------------------------------------------------------------

// Load reads the relationship configuration from a file
func (relationshipConfiguration *RelationshipConfiguration) Load(filename string) error {
	return util.LoadYAML(filename, relationshipConfiguration)
}

//------------------------------------------------------------------------------
