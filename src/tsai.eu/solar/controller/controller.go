package controller

import (
	"errors"
	"sync"

	"tsai.eu/solar/controller/demo"
	"tsai.eu/solar/controller/dummy"
	"tsai.eu/solar/model"
)

//------------------------------------------------------------------------------

// Controller defines the standard operations of a controller
type Controller interface {
	Status(configuration *model.Setup) (status *model.Status, err error)
	Create(configuration *model.Setup) (status *model.Status, err error)
	Destroy(configuration *model.Setup) (status *model.Status, err error)
	Configure(configuration *model.Setup) (status *model.Status, err error)
	Start(configuration *model.Setup) (status *model.Status, err error)
	Stop(configuration *model.Setup) (status *model.Status, err error)
	Reset(configuration *model.Setup) (status *model.Status, err error)
}

//------------------------------------------------------------------------------

var controllers map[string]Controller

var once sync.Once

//------------------------------------------------------------------------------

// GetController retrieves a controller for a specific component type.
func GetController(componentType string) (Controller, error) {
	// initialise singleton once
	once.Do(func() {
		controllers = map[string]Controller{}

		controllers["Demo"]  = demo.NewController()
		controllers["Dummy"] = dummy.NewController()
	})

	// determine controller
	controller, found := (controllers)[componentType]
	if !found {
		return nil, errors.New("unknown type")
	}

	// success
	return controller, nil
}

//------------------------------------------------------------------------------
