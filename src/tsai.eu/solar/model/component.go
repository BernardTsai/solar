package model

import (
	"github.com/pkg/errors"
	"tsai.eu/solar/util"
)

//------------------------------------------------------------------------------
// Component
// =========
//
// Attributes:
//   - Component
//   - Version
//   - Configuration
//   - Dependencies
//
// Functions:
//   - NewComponent
//
//   - component.Show
//   - component.Load
//   - component.Save
//
//   - component.ListDependencies
//   - component.GetDependency
//   - component.AddDependency
//   - component.DeleteDependency
//------------------------------------------------------------------------------

// DependencyMap is a synchronized map for a map of dependencies
type DependencyMap struct {
	Map          map[string]*Dependency  `yaml:"map"`             // map of dependencies
}

// MarshalYAML marshals a DependencyMap into yaml
func (m DependencyMap) MarshalYAML() (interface{}, error) {
	return m.Map, nil
}

// UnmarshalYAML unmarshals a DependencyMap from yaml
func (m *DependencyMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	Map := map[string]*Dependency{}

	err := unmarshal(&Map)
	if err != nil {
		return err
	}

	*m = DependencyMap{Map: Map}

	return nil
}

//------------------------------------------------------------------------------

// Component describes a base configuration for a component within a domain.
type Component struct {
	Component     string        `yaml:"Component"`     // name of the component
	Version       string        `yaml:"Version"`       // version of the component
	Configuration string        `yaml:"Configuration"` // base configuration of the component
	Dependencies  DependencyMap `yaml:"Dependencies"`  // dependencies of component
}

//------------------------------------------------------------------------------

// NewComponent creates a new component
func NewComponent(name string, version string, configuration string) (*Component, error) {
	var component Component

	component.Component     = name
	component.Version       = version
	component.Configuration = configuration
	component.Dependencies  = DependencyMap{Map: map[string]*Dependency{}}

	// success
	return &component, nil
}

//------------------------------------------------------------------------------

// Show displays the component information as yaml
func (component *Component) Show() (string, error) {
	return util.ConvertToYAML(component)
}

//------------------------------------------------------------------------------

// Save writes the component as yaml data to a file
func (component *Component) Save(filename string) error {
	return util.SaveYAML(filename, component)
}

//------------------------------------------------------------------------------

// Load reads the template from a file
func (component *Component) Load(filename string) error {
	return util.LoadYAML(filename, component)
}

//------------------------------------------------------------------------------

// ListDependencies lists all dependencies of a template
func (component *Component) ListDependencies() ([]string, error) {
	// collect names
	dependencies := []string{}

	for dependency := range component.Dependencies.Map {
		dependencies = append(dependencies, dependency)
	}

	// success
	return dependencies, nil
}

//------------------------------------------------------------------------------

// GetDependency retrieves a dependency by name
func (component *Component) GetDependency(name string) (*Dependency, error) {
	// determine dependency
	dependency, ok := component.Dependencies.Map[name]

	if !ok {
		return nil, errors.New("dependency not found")
	}

	// success
	return dependency, nil
}

//------------------------------------------------------------------------------

// AddDependency adds a dependency to a component
func (component *Component) AddDependency(dependency *Dependency) error {
	// check if dependency has already been defined
	_, ok := component.Dependencies.Map[dependency.Dependency]

	if ok {
		return errors.New("variant already exists")
	}

	component.Dependencies.Map[dependency.Dependency] = dependency

	// success
	return nil
}

//------------------------------------------------------------------------------

// DeleteDependency deletes a dependency
func (component *Component) DeleteDependency(name string) error {
	// determine dependency
	_, ok := component.Dependencies.Map[name]

	if !ok {
		return errors.New("dependency not found")
	}

	// remove version
	delete(component.Dependencies.Map, name)

	// success
	return nil
}

//------------------------------------------------------------------------------
