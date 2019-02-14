package model

import (
	"tsai.eu/solar/util"
)

//------------------------------------------------------------------------------
// Dependency
// ==========
//
// Attributes:
//   - Dependency
//   - Type
//   - Component
//   - Version
//   - Configuration
//
// Functions:
//   - NewDependency
//
//   - dependency.Show
//   - dependency.Load
//   - dependency.Save
//------------------------------------------------------------------------------

// Dependency describes what kind of dependency a component within a domain may have.
type Dependency struct {
	Dependency    string `yaml:"dependency"`    // name of the dependency
	Type          string `yaml:"type"`          // type of dependency (service/context)
	Component     string `yaml:"component"`     // component to which the dependency refers to
	Version       string `yaml:"version"`       // version of the component to which the dependency refers to
	Configuration string `yaml:"configuration"` // base configuration of the dependency
}

//------------------------------------------------------------------------------

// NewDependency creates a new dependency
func NewDependency(name string, dtype string, component string, version string, configuration string ) (*Dependency, error) {
	var dependency Dependency

	dependency.Dependency = name
	dependency.Type = dtype
	dependency.Component = component
	dependency.Version = version
	dependency.Configuration = configuration

	// success
	return &dependency, nil
}

//------------------------------------------------------------------------------

// Show displays the dependency information as yaml
func (dependency *Dependency) Show() (string, error) {
	return util.ConvertToYAML(dependency)
}

//------------------------------------------------------------------------------

// Save writes the dependency as yaml data to a file
func (dependency *Dependency) Save(filename string) error {
	return util.SaveYAML(filename, dependency)
}

//------------------------------------------------------------------------------

// Load reads the dependency from a file
func (dependency *Dependency) Load(filename string) error {
	return util.LoadYAML(filename, dependency)
}

//------------------------------------------------------------------------------
