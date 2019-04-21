package defaultController

import (
	"tsai.eu/solar/model"
)

//------------------------------------------------------------------------------

// Create initialises an instance
func (c *Controller) Create(targetState *model.TargetState) (currentState *model.CurrentState, err error) {
	return c.Status(targetState)
}

//------------------------------------------------------------------------------
