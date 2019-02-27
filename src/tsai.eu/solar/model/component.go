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
//   - component.Load2
//   - component.Save
//
//   - component.ListDependencies
//   - component.GetDependency
//   - component.AddDependency
//   - component.DeleteDependency
//------------------------------------------------------------------------------

// Component describes a base configuration for a component within a domain.
type Component struct {
	Component     string                 `yaml:"Component"`             // name of the component
	Version       string                 `yaml:"Version"`               // version of the component
	Configuration string                 `yaml:"Configuration"`         // base configuration of the component
	Dependencies  map[string]*Dependency `yaml:"Dependencies"`          // dependencies of component
	DependenciesX sync.RWMutex           `yaml:"ComponentsX,omitempty"` // mutex for dependencies
}

//------------------------------------------------------------------------------

// NewComponent creates a new component
func NewComponent(name string, version string, configuration string) (*Component, error) {
	var component Component

	component.Component     = name
	component.Version       = version
	component.Configuration = configuration
	component.Dependencies  = map[string]*Dependency{}
	component.DependenciesX = sync.RWMutex{}

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

// Load2 imports a yaml model
func (component *Component) Load2(yaml string) error {
	return util.ConvertFromYAML(yaml, component)
}

//------------------------------------------------------------------------------

// ListDependencies lists all dependencies of a template
func (component *Component) ListDependencies() ([]string, error) {
	// collect names
	dependencies := []string{}

	component.DependenciesX.RLock()
	for dependency := range component.Dependencies {
		dependencies = append(dependencies, dependency)
	}
	component.DependenciesX.RUnlock()

	// success
	return dependencies, nil
}

//------------------------------------------------------------------------------

// GetDependency retrieves a dependency by name
func (component *Component) GetDependency(name string) (*Dependency, error) {
	// determine dependency
	component.DependenciesX.RLock()
	dependency, ok := component.Dependencies[name]
	component.DependenciesX.RUnlock()

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
	component.DependenciesX.RLock()
	_, ok := component.Dependencies[dependency.Dependency]
	component.DependenciesX.RUnlock()

	if ok {
		return errors.New("variant already exists")
	}

	component.DependenciesX.Lock()
	component.Dependencies[dependency.Dependency] = dependency
	component.DependenciesX.Unlock()

	// success
	return nil
}

//------------------------------------------------------------------------------

// DeleteDependency deletes a dependency
func (component *Component) DeleteDependency(name string) error {
	// determine dependency
	component.DependenciesX.RLock()
	_, ok := component.Dependencies[name]
	component.DependenciesX.RUnlock()

	if !ok {
		return errors.New("dependency not found")
	}

	// remove version
	component.DependenciesX.Lock()
	delete(component.Dependencies, name)
	component.DependenciesX.Unlock()

	// success
	return nil
}

//------------------------------------------------------------------------------
