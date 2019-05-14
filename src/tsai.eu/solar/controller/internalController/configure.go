package internalController

import (
	"tsai.eu/solar/model"
)

//------------------------------------------------------------------------------

// Configure reconfigures an instance
func (c *Controller) Configure(targetState *model.TargetState) (currentState *model.CurrentState, err error) {
	return c.Status(targetState)
}

//------------------------------------------------------------------------------
