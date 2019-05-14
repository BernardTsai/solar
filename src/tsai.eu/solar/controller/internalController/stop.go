package internalController

import (
	"tsai.eu/solar/model"
)

//------------------------------------------------------------------------------

// Stop deactivates an instance
func (c *Controller) Stop(targetState *model.TargetState) (currentState *model.CurrentState, err error) {
	return c.Status(targetState)
}

//------------------------------------------------------------------------------
