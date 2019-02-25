package model

import (
	"github.com/pkg/errors"
	"tsai.eu/solar/util"
)

//------------------------------------------------------------------------------
// Architecture
// ============
//
// Attributes:
//   - Architecture
//   - Version
//   - Configuration
//   - Elements
//
// Functions:
//   - NewArchitecture
//
//   - architecture.Show
//   - architecture.Load
//   - architecture.Save
//
//   - architecture.ListElements
//   - architecture.GetElement
//   - architecture.AddElement
//   - architecture.DeleteElement
//------------------------------------------------------------------------------

// ElementConfigurationMap is a synchronized map for a map of element configurations
type ElementConfigurationMap struct {
	Map map[string]*ElementConfiguration `yaml:"map"`             // map of element configurations
}

// MarshalYAML marshals an ElementConfigurationMap into yaml
func (m ElementConfigurationMap) MarshalYAML() (interface{}, error) {
	return m.Map, nil
}

// UnmarshalYAML unmarshals an ElementConfigurationMap from yaml
func (m *ElementConfigurationMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	Map := map[string]*ElementConfiguration{}

	err := unmarshal(&Map)
	if err != nil {
		return err
	}

	*m = ElementConfigurationMap{Map: Map}

	return nil
}

//------------------------------------------------------------------------------

// Architecture describes the design time configuration of a solution within a domain.
type Architecture struct {
	Architecture  string                  `yaml:"Architecture"`   // name of architecture
	Version       string                  `yaml:"Version"`        // type of solution
	Configuration string                  `yaml:"Configuration"`  // configuration of the architecture
	Elements      ElementConfigurationMap `yaml:"Elements"`       // elements configurations of solution
}

//------------------------------------------------------------------------------

// NewArchitecture creates a new architecture
func NewArchitecture(name string, version string, configuration string) (*Architecture, error) {
	var architecture Architecture

	architecture.Architecture  = name
	architecture.Version       = version
	architecture.Configuration = configuration
	architecture.Elements      = ElementConfigurationMap{Map: map[string]*ElementConfiguration{}}

	// success
	return &architecture, nil
}

//------------------------------------------------------------------------------

// Show displays the architecture information as yaml
func (architecture *Architecture) Show() (string, error) {
	return util.ConvertToYAML(architecture)
}

//------------------------------------------------------------------------------

// Save writes the architecture as yaml data to a file
func (architecture *Architecture) Save(filename string) error {
	return util.SaveYAML(filename, architecture)
}

//------------------------------------------------------------------------------

// Load reads the architecture from a file
func (architecture *Architecture) Load(filename string) error {
	return util.LoadYAML(filename, architecture)
}

//------------------------------------------------------------------------------

// ListElements lists all elements of an architecture
func (architecture *Architecture) ListElements() ([]string, error) {
	// collect names
	elementConfigurations := []string{}

	if architecture != nil {
		for elementConfiguration := range architecture.Elements.Map {
			elementConfigurations = append(elementConfigurations, elementConfiguration)
		}
	}

	// success
	return elementConfigurations, nil
}

//------------------------------------------------------------------------------

// GetElement retrieves an element configuration by name
func (architecture *Architecture) GetElement(name string) (*ElementConfiguration, error) {
	// determine instance
	elementConfiguration, ok := architecture.Elements.Map[name]

	if !ok {
		return nil, errors.New("element configuration not found")
	}

	// success
	return elementConfiguration, nil
}

//------------------------------------------------------------------------------

// AddElement adds an element configuration to a component
func (architecture *Architecture) AddElement(elementConfiguration *ElementConfiguration) error {
	// check if instance has already been defined
	_, ok := architecture.Elements.Map[elementConfiguration.Element]

	if ok {
		return errors.New("element configuration already exists")
	}

	architecture.Elements.Map[elementConfiguration.Element] = elementConfiguration

	// success
	return nil
}

//------------------------------------------------------------------------------

// DeleteElement deletes an element configuration
func (architecture *Architecture) DeleteElement(name string) error {
	// determine element
	_, ok := architecture.Elements.Map[name]

	if !ok {
		return errors.New("element configuration not found")
	}

	// remove element
	delete(architecture.Elements.Map, name)

	// success
	return nil
}

//------------------------------------------------------------------------------
