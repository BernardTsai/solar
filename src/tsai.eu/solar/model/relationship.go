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
//   - Dependency
//   - Type
//   - Element
//   - Version
//   - Target
//   - State
//   - Configuration
//   - Endpoint
//
// Functions:
//   - NewRelationship
//
//   - relationship.Show
//   - relationship.Load
//   - relationship.Save
//   - relationship.Reset
//------------------------------------------------------------------------------

// Relationship describes the runtime configuration of a relationship between clusters within a domain.
type Relationship struct {
	Relationship  string  `yaml:"relationship"`  // name of the relationship
	Dependency    string  `yaml:"dependency"`    // name of the dependency
	Type          string  `yaml:"type"`          // type of dependency
	Domain        string  `yaml:"domain"`        // domain to which this relationship refers to
	Solution      string  `yaml:"solution"`      // solution to which this relationship refers to
	Element       string  `yaml:"element"`       // element to which this relationship refers to
	Version       string  `yaml:"version"`       // version of the element to which this relationship refers to
	Target        string  `yaml:"target"`        // target state of relationship
	State         string  `yaml:"state"`         // current state of relationship
	Configuration string  `yaml:"configuration"` // runtime configuration of the relationship
	Endpoint      string  `yaml:"endpoint"`      // endpoint of the relationship
}

//------------------------------------------------------------------------------

// NewRelationship creates a new relationship
func NewRelationship(name string, dependency string, dependencyType string, domain string, solution string, element string, version string, configuration string) (*Relationship, error) {
	var relationship Relationship

	relationship.Relationship  = name
	relationship.Dependency    = dependency
	relationship.Type          = dependencyType
	relationship.Domain        = domain
	relationship.Solution      = solution
	relationship.Element       = element
	relationship.Version       = version
	relationship.Target        = InitialState
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

// Reset state of relationship
func (relationship *Relationship) Reset() {
	relationship.Target = InitialState
}

//------------------------------------------------------------------------------
