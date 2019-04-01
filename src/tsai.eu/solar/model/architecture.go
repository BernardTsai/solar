package model

import (
	"sync"
	"errors"
	
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
//   - architecture.Load2
//   - architecture.Save
//
//   - architecture.ListElements
//   - architecture.GetElement
//   - architecture.AddElement
//   - architecture.DeleteElement
//------------------------------------------------------------------------------

// Architecture describes the design time configuration of a solution within a domain.
type Architecture struct {
	Architecture  string                           `yaml:"Architecture"`        // name of architecture
	Version       string                           `yaml:"Version"`             // type of solution
	Configuration string                           `yaml:"Configuration"`       // configuration of the architecture
	Elements      map[string]*ElementConfiguration `yaml:"Elements"`            // element configurations of solution
	ElementsX     sync.RWMutex                     `yaml:"ElementsX,omitempty"` // mutex for element configurations
}

//------------------------------------------------------------------------------

// NewArchitecture creates a new architecture
func NewArchitecture(name string, version string, configuration string) (*Architecture, error) {
	var architecture Architecture

	architecture.Architecture  = name
	architecture.Version       = version
	architecture.Configuration = configuration
	architecture.Elements      = map[string]*ElementConfiguration{}
	architecture.ElementsX     = sync.RWMutex{}

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

// Load2 imports a yaml model
func (architecture *Architecture) Load2(yaml string) error {
	return util.ConvertFromYAML(yaml, architecture)
}

//------------------------------------------------------------------------------

// ListElements lists all elements of an architecture
func (architecture *Architecture) ListElements() ([]string, error) {
	// collect names
	elementConfigurations := []string{}

  architecture.ElementsX.RLock()
	for elementConfiguration := range architecture.Elements {
		elementConfigurations = append(elementConfigurations, elementConfiguration)
	}
	architecture.ElementsX.RUnlock()

	// success
	return elementConfigurations, nil
}

//------------------------------------------------------------------------------

// GetElement retrieves an element configuration by name
func (architecture *Architecture) GetElement(name string) (*ElementConfiguration, error) {
	// determine instance
	architecture.ElementsX.RLock()
	elementConfiguration, ok := architecture.Elements[name]
	architecture.ElementsX.RUnlock()

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
	architecture.ElementsX.RLock()
	_, ok := architecture.Elements[elementConfiguration.Element]
	architecture.ElementsX.RUnlock()

	if ok {
		return errors.New("element configuration already exists")
	}

	architecture.ElementsX.Lock()
	architecture.Elements[elementConfiguration.Element] = elementConfiguration
	architecture.ElementsX.Unlock()

	// success
	return nil
}

//------------------------------------------------------------------------------

// DeleteElement deletes an element configuration
func (architecture *Architecture) DeleteElement(name string) error {
	// determine element
	architecture.ElementsX.RLock()
	_, ok := architecture.Elements[name]
	architecture.ElementsX.RUnlock()

	if !ok {
		return errors.New("element configuration not found")
	}

	// remove element
	architecture.ElementsX.Lock()
	delete(architecture.Elements, name)
	architecture.ElementsX.Unlock()

	// success
	return nil
}

//------------------------------------------------------------------------------
