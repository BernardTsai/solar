package model

import (
	"tsai.eu/solar/util"
)

//------------------------------------------------------------------------------
// Controller
// =========
//
// Attributes:
//   - Controller
//   - Version
//   - URL
//   - Types
//   - Status (initial, inactive, active, failure)
//
// Functions:
//   - NewController
//
//   - controller.Show
//   - controller.Load
//   - controller.Load2
//   - controller.Save
//------------------------------------------------------------------------------

// Controller describes a controller for a set of component types.
type Controller struct {
	Controller string      `yaml:"Controller"`  // name of the controller
	Version    string      `yaml:"Version"`     // version of the controller
	URL        string      `yaml:"URL"`         // URL of the controller
	Status     string      `yaml:"Status"`      // status of the controller
}

//------------------------------------------------------------------------------

// NewController creates a new controller
func NewController(controller string, version string) (*Controller, error) {
	var ctrl Controller

	ctrl.Controller = controller
	ctrl.Version    = version
	ctrl.URL        = ""
	ctrl.Status     = InitialState

	// success
	return &ctrl, nil
}

//------------------------------------------------------------------------------

// Show displays the controller information as yaml
func (controller *Controller) Show() (string, error) {
	return util.ConvertToYAML(controller)
}

//------------------------------------------------------------------------------

// Save writes the controller as yaml data to a file
func (controller *Controller) Save(filename string) error {
	return util.SaveYAML(filename, controller)
}

//------------------------------------------------------------------------------

// Load reads the template from a file
func (controller *Controller) Load(filename string) error {
	return util.LoadYAML(filename, controller)
}

//------------------------------------------------------------------------------

// Load2 imports a yaml model
func (controller *Controller) Load2(yaml string) error {
	return util.ConvertFromYAML(yaml, controller)
}

//------------------------------------------------------------------------------
