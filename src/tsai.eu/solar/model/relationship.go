package model

import (
	"tsai.eu/solar/util"
)

//------------------------------------------------------------------------------
// Relationship
// ============
//
// Attributes:
//   - Relationship
//   - Element
//   - Version
//   - Configuration
//   - Endpoint
//
// Functions:
//   - NewRelationship
//
//   - relationship.Show
//   - relationship.Load
//   - relationship.Save
//------------------------------------------------------------------------------

// Relationship describes the runtime configuration of a relationship between clusters within a domain.
type Relationship struct {
	Relationship  string  `yaml:"relationship"`  // name of the relationship
	Element       string  `yaml:"element"`       // element to which this relationship refers to
	Version       string  `yaml:"version"`       // version of the element to which this relationship refers to
	Target        string  `yaml:"target"`        // target state of relationship
	State         string  `yaml:"state"`         // current state of relationship
	Configuration string  `yaml:"configuration"` // runtime configuration of the relationship
	Endpoint      string  `yaml:"endpoint"`      // endpoint of the relationship
}

//------------------------------------------------------------------------------

// NewRelationship creates a new relationship
func NewRelationship(name string, element string, version string, configuration string) (*Relationship, error) {
	var relationship Relationship

	relationship.Relationship  = name
	relationship.Element       = element
	relationship.Version       = version
	relationship.Target        = ActiveState
	relationship.State         = InitialState
	relationship.Configuration = configuration
	relationship.Endpoint      = ""

	// success
	return &relationship, nil
}

//------------------------------------------------------------------------------

// Show displays the relationship information as yaml
func (relationship *Relationship) Show() (string, error) {
	return util.ConvertToYAML(relationship)
}

//------------------------------------------------------------------------------

// Save writes the relationship as yaml data to a file
func (relationship *Relationship) Save(filename string) error {
	return util.SaveYAML(filename, relationship)
}

//------------------------------------------------------------------------------

// Load reads the relationship from a file
func (relationship *Relationship) Load(filename string) error {
	return util.LoadYAML(filename, relationship)
}

//------------------------------------------------------------------------------
