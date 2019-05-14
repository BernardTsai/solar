package internalController

import (
	"tsai.eu/solar/model"
)

//------------------------------------------------------------------------------

// Start activates an instance
func (c *Controller) Start(targetState *model.TargetState) (currentState *model.CurrentState, err error) {
	return c.Status(targetState)
}

//------------------------------------------------------------------------------
