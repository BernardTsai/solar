package defaultController

import (
	"tsai.eu/solar/model"
)

//------------------------------------------------------------------------------

// Reset cleans up an instance
func (c *Controller) Reset(targetState *model.TargetState) (currentState *model.CurrentState, err error) {
	return c.Status(targetState)
}

//------------------------------------------------------------------------------
