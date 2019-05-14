package controller

import (
	"sync"
	// "context"

	"tsai.eu/solar/model"
	"tsai.eu/solar/util"
	"tsai.eu/solar/controller/internalController"
)

//------------------------------------------------------------------------------

// Controller defines the standard operations of a controller
type Controller interface {
	Status(      setup *model.TargetState) (status *model.CurrentState, err error)
	Create(      setup *model.TargetState) (status *model.CurrentState, err error)
	Destroy(     setup *model.TargetState) (status *model.CurrentState, err error)
	Configure(   setup *model.TargetState) (status *model.CurrentState, err error)
	Reconfigure( setup *model.TargetState) (status *model.CurrentState, err error)
	Start(       setup *model.TargetState) (status *model.CurrentState, err error)
	Stop(        setup *model.TargetState) (status *model.CurrentState, err error)
	Reset(       setup *model.TargetState) (status *model.CurrentState, err error)
}

//------------------------------------------------------------------------------

var controllers   map[string]Controller   // controller map
var defController Controller              // default controller
var initCtrls     sync.Once               // initialisation guard

var port          int                     // next free port
var initPort      sync.Once               // initialisation guard

//------------------------------------------------------------------------------

// GetController retrieves a controller for a specific version of a controller.
func GetController(controllerVersion string) (Controller, error) {
	// initialise singleton once
	initCtrls.Do(func() {
		// create empty map of controllers
		controllers = map[string]Controller{}

		// initialise the default internal controller
		defController = internalController.NewController()
		controllers["internal:V1.0.0"] = defController
		util.LogInfo("main", "CTL", "internal - controller active")
	})

	// determine controller (if unknown use the default controller)
	controller, found := controllers[controllerVersion]
	if found {
		return controller, nil
	}

	// offer the default internal controller
	return defController, nil


	// TODO: log unknown controller type and error
	// return nil, errors.New("unknown type")
}

//------------------------------------------------------------------------------

// AddController adds a controller to a domain
func AddController(domainName string, controllerName string, controllerVersion string) error {
	// determine domain
	domain, err := model.GetDomain(domainName)
	if err != nil {
		return err
	}

	// determine controller and create a new one if it does not exist
	var ctrl *model.Controller
	ctrl, err = domain.GetController(controllerName, controllerVersion)
	if err != nil {
		ctrl, _ = model.NewController(controllerName, controllerVersion)
		domain.AddController(ctrl)
	}

	// load controller if needed
	// if controller.Status == model.InitialState {
	// 	err = pullImage(image, version)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	controller.Status = model.InactiveState
	// }

	// start controller if needed
	// if controller.Status == model.InactiveState {
	// 	err = startContainer(image, version)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if controller.URL == "" {
	// 	controller.Status = model.ActiveState
	// }

	// success
	return nil
}

//------------------------------------------------------------------------------
