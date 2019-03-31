package controller

import (
	"sync"
	"fmt"

	"tsai.eu/solar/model"
	"tsai.eu/solar/util"
	"tsai.eu/solar/controller/dummy"
)

//------------------------------------------------------------------------------

// Controller defines the standard operations of a controller
type Controller interface {
	Status(    setup *model.Setup) (status *model.Status, err error)
	Create(    setup *model.Setup) (status *model.Status, err error)
	Destroy(   setup *model.Setup) (status *model.Status, err error)
	Configure( setup *model.Setup) (status *model.Status, err error)
	Start(     setup *model.Setup) (status *model.Status, err error)
	Stop(      setup *model.Setup) (status *model.Status, err error)
	Reset(     setup *model.Setup) (status *model.Status, err error)
}

//------------------------------------------------------------------------------

var controllers   map[string]Controller
var defController Controller

var once sync.Once

//------------------------------------------------------------------------------

// GetController retrieves a controller for a specific component type.
func GetController(componentType string) (Controller, error) {
	// initialise singleton once
	once.Do(func() {
		// create empty map of controllers
		controllers = map[string]Controller{}

		// initialise the default controller
		defController = dummy.NewController()

		// read the controller configuration
		config, _ := util.GetConfiguration()

		// initialise all gRPC controllers
		for controllerType, controllerAddress := range config.CTRL {
			controller, err := newGRPCController(controllerType, controllerAddress)
			if err != nil {
				fmt.Println("Controller error:\n" + err.Error())
			}
			controllers[controllerType] = controller
		}
	})

	// determine controller (if unknown use the dummy controller)
	controller, found := controllers[componentType]
	if found {
		return controller, nil
	}

	// try to use the dummy controller
	// controller, found = controllers["dummy"]
	// if found {
	// 	return controller, nil
	// }

	// offer the default controller
	return defController, nil


	// TODO: log unknown controller type and error
	// return nil, errors.New("unknown type")
}

//------------------------------------------------------------------------------
