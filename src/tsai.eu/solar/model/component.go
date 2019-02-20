package model

import (
	"sync"

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
	sync.RWMutex                         `yaml:"mutex,omitempty"` // mutex
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
	Component     string        `yaml:"component"`     // name of the component
	Version       string        `yaml:"version"`       // version of the component
	Configuration string        `yaml:"configuration"` // base configuration of the component
	Dependencies  DependencyMap `yaml:"dependencies"`  // dependencies of component
}

//------------------------------------------------------------------------------

// NewComponent creates a new component
func NewComponent(name string, version string, configuration string) (*Component, error) {
	var component Component

	component.Component = name
	component.Version = version
	component.Configuration = configuration
	component.Dependencies = DependencyMap{Map: map[string]*Dependency{}}

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

	component.Dependencies.RLock()
	for dependency := range component.Dependencies.Map {
		dependencies = append(dependencies, dependency)
	}
	component.Dependencies.RUnlock()

	// success
	return dependencies, nil
}

//------------------------------------------------------------------------------

// GetDependency retrieves a dependency by name
func (component *Component) GetDependency(name string) (*Dependency, error) {
	// determine dependency
	component.Dependencies.RLock()
	dependency, ok := component.Dependencies.Map[name]
	component.Dependencies.RUnlock()

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
	component.Dependencies.RLock()
	_, ok := component.Dependencies.Map[dependency.Dependency]
	component.Dependencies.Unlock()

	if ok {
		return errors.New("variant already exists")
	}

	component.Dependencies.Lock()
	component.Dependencies.Map[dependency.Dependency] = dependency
	component.Dependencies.Unlock()

	// success
	return nil
}

//------------------------------------------------------------------------------

// DeleteDependency deletes a dependency
func (component *Component) DeleteDependency(name string) error {
	// determine dependency
	component.Dependencies.RLock()
	_, ok := component.Dependencies.Map[name]
	component.Dependencies.RUnlock()

	if !ok {
		return errors.New("dependency not found")
	}

	// remove version
	component.Dependencies.Lock()
	delete(component.Dependencies.Map, name)
	component.Dependencies.Unlock()

	// success
	return nil
}

//------------------------------------------------------------------------------