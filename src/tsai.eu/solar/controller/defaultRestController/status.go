package defaultRestController

import (
	"net/http"
)

//------------------------------------------------------------------------------

// Status provides the status of an instance
func (c *Controller) Status(w http.ResponseWriter, r *http.Request) {
	targetState, err := readTargetState(r)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Unable to read target state information"))
		return
	}

	// construct current state
	currentState := CurrentState{
		Domain:           targetState.Domain,
		Solution:         targetState.Solution,
		Version:          targetState.Version,
		Element:          targetState.Element,
		ElementEndpoint:  "",
		Cluster:          targetState.Cluster,
		ClusterEndpoint:  "",
		ClusterState:     targetState.ClusterState,
	  Instance:         targetState.Instance,
		InstanceEndpoint: "",
		InstanceState:    targetState.InstanceState,
	}

	// write the result
	writeCurrentState(w, currentState)
}

//------------------------------------------------------------------------------
