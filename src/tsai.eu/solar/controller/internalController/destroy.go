package internalController

import (
	"tsai.eu/solar/model"
)

//------------------------------------------------------------------------------

// Destroy removes instance
func (c *Controller) Destroy(targetState *model.TargetState) (currentState *model.CurrentState, err error) {
	return c.Status(targetState)
}

//------------------------------------------------------------------------------
