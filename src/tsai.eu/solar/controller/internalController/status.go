package internalController

import (
	"tsai.eu/solar/model"
)

//------------------------------------------------------------------------------

// Status provides the status of an instance
func (c *Controller) Status(targetState *model.TargetState) (currentState *model.CurrentState, err error) {
	// construct current state
	currentState = &model.CurrentState{
		Domain:        targetState.Domain,
		Solution:      targetState.Solution,
		Version:       targetState.Version,
		Element:       targetState.Element,
		Cluster:       targetState.Cluster,
	  Instance:      targetState.Instance,
		State:         targetState.State,
		Configuration: targetState.Configuration,
		Endpoint:      "",
	}

	// return results
	return currentState, nil
}

//------------------------------------------------------------------------------
