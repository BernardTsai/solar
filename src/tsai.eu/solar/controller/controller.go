package controller

import (
	"sync"
	"context"

	"tsai.eu/solar/model"
	"tsai.eu/solar/util"
	"tsai.eu/solar/controller/defaultController"
	"tsai.eu/solar/controller/defaultRestController"
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

// GetController retrieves a controller for a specific component type.
func GetController(componentType string) (Controller, error) {
	// initialise singleton once
	initCtrls.Do(func() {
		// create empty map of controllers
		controllers = map[string]Controller{}

		// initialise the default controller
		def2Controller := defaultRestController.NewController()
		def2Controller.Run(context.Background())
		// controllers["default2"] = def2Controller
		util.LogInfo("main", "CTRL", "default2 - controller active")

		// initialise the default controller
		defController = defaultController.NewController()
		controllers["default"] = defController
		util.LogInfo("main", "CTRL", "default - controller active")

		// read the controller configuration
		config, _ := util.GetConfiguration()

		// initialise all REST controllers
		for _, info := range config.CTRL {
			controller, err := newRestController(info.Type, info.Version, info.URL)
			if err == nil {
				controllers[info.Type + ":" + info.Version] = controller
				util.LogInfo("main", "CTRL", info.Type + ":" + info.Version + " - controller active")
			}
		}
		// initialise all gRPC controllers
		// for controllerType, controllerAddress := range config.CTRL {
		// 	controller, err := newGRPCController(controllerType, controllerAddress)
		// 	if err == nil {
		// 		controllers[controllerType] = controller
		// 		util.LogInfo("main", "CTRL", controllerType + " - controller active")
		// 	}
		// }
	})

	// determine controller (if unknown use the default controller)
	controller, found := controllers[componentType]
	if found {
		return controller, nil
	}

	// try to use the default controller
	controller, found = controllers["default"]
	if found {
		return controller, nil
	}

	// offer the default controller
	return defController, nil


	// TODO: log unknown controller type and error
	// return nil, errors.New("unknown type")
}

//------------------------------------------------------------------------------

// AddController adds a controller to a domain
func AddController(domainName string, image string, version string) error {
	// determine domain
	domain, err := model.GetDomain(domainName)
	if err != nil {
		return err
	}

	// determine controller and create a new one if it does not exist
	controller, err := domain.GetController(image, version)
	if err != nil {
		controller, _ = model.NewController(image, version)
		domain.AddController(controller)
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
