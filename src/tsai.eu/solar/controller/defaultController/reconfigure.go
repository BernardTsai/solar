package defaultController

import (
	"tsai.eu/solar/model"
)

//------------------------------------------------------------------------------

// Reconfigure reconfigures an instance
func (c *Controller) Reconfigure(targetState *model.TargetState) (currentState *model.CurrentState, err error) {
	return c.Status(targetState)
}

//------------------------------------------------------------------------------
