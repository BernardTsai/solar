package model

import (
	"errors"

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
//
//   - controller.ListTypes
//   - controller.AddType
//   - controller.DeleteType
//------------------------------------------------------------------------------

// Controller describes a controller for a set of component types.
type Controller struct {
	Controller string      `yaml:"Controller"`  // name of the controller
	Version    string      `yaml:"Version"`     // version of the controller
	URL        string      `yaml:"URL"`         // URL of the controller
	Status     string      `yaml:"Status"`      // status of the controller
	Types      [][2]string `yaml:"Types"`       // supported component types
}

//------------------------------------------------------------------------------

// NewController creates a new controller
func NewController(controller string, version string) (*Controller, error) {
	var ctrl Controller

	ctrl.Controller = controller
	ctrl.Version    = version
	ctrl.URL        = ""
	ctrl.Status     = InitialState
	ctrl.Types      = [][2]string{}

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

// ListTypes lists all supported component types of a controller
func (controller *Controller) ListTypes() ([][2]string, error) {
	// success
	return controller.Types, nil
}

//------------------------------------------------------------------------------

// AddType adds a component type to a controller
func (controller *Controller) AddType(component string, version string) error {
	// check if component type has already been added
	for _, Type := range controller.Types {
		if Type[0] == component && Type[1] == version {
			return errors.New("component has already been added to the controller")
		}
	}

	controller.Types = append(controller.Types, [2]string{component, version})

	// success
	return nil
}

//------------------------------------------------------------------------------

// DeleteType deletes a component type from a controller
func (controller *Controller) DeleteType(component string, version string) error {
	for index, Type := range controller.Types {
		if Type[0] == component && Type[1] == version {
			controller.Types = append(controller.Types[:index],controller.Types[:index+1]...)

			return nil
		}
	}

	return errors.New("component not found")
}

//------------------------------------------------------------------------------
