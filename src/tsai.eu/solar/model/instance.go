package model

import (
	"tsai.eu/solar/util"
)

//------------------------------------------------------------------------------
// Instance
// ========
//
// Attributes:
//   - UUID
//   - State
//   - Configuration
//   - Endpoint
//
// Functions:
//   - NewInstance
//
//   - instance.Show
//   - instance.Load
//   - instance.Save
//------------------------------------------------------------------------------

// Instance describes the runtime configuration of an solution element cluster instance within a domain.
type Instance struct {
	UUID          string  `yaml:"uuid"`          // uuid of the instance
	State         string  `yaml:"state"`         // state of the instance
	Configuration string  `yaml:"configuration"` // runtime configuration of the instance
	Endpoint      string  `yaml:"endpoint"`      // endpoint of the instance
}

//------------------------------------------------------------------------------

// NewInstance creates a new instance
func NewInstance(uuid string, state string, configuration string) (*Instance, error) {
	var instance Instance

	instance.UUID = uuid
	instance.State = state
	instance.Configuration = configuration
	instance.Endpoint = ""

	// success
	return &instance, nil
}

//------------------------------------------------------------------------------

// Show displays the instance information as yaml
func (instance *Instance) Show() (string, error) {
	return util.ConvertToYAML(instance)
}

//------------------------------------------------------------------------------

// Save writes the instance as yaml data to a file
func (instance *Instance) Save(filename string) error {
	return util.SaveYAML(filename, instance)
}

//------------------------------------------------------------------------------

// Load reads the instance from a file
func (instance *Instance) Load(filename string) error {
	return util.LoadYAML(filename, instance)
}

//------------------------------------------------------------------------------
