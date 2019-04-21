package model

import (
	"testing"
)

//------------------------------------------------------------------------------

// TestState01 tests the state related functions of the model package.
func TestState01(t *testing.T) {
	model := GetModel()

	model.Load("testdata/testdata1.yaml")

	_, err :=  GetTargetState("unknown", "app", "V0.0.0", "oam", "V1.0.0", "a6c0bea1-ce1a-4fae-b943-1dbcc50cb311")
	if err == nil {
		t.Errorf("GetTargetState should have reported an unknown domain")
	}

	_, err =  GetTargetState("demo", "unknown", "V0.0.0", "oam", "V1.0.0", "a6c0bea1-ce1a-4fae-b943-1dbcc50cb311")
	if err == nil {
		t.Errorf("GetTargetState should have reported an unknown solution")
	}

	_, err =  GetTargetState("demo", "app", "unknown", "oam", "V1.0.0", "a6c0bea1-ce1a-4fae-b943-1dbcc50cb311")
	if err == nil {
		t.Errorf("GetTargetState should have reported an unknown architecture")
	}

	_, err =  GetTargetState("demo", "app", "V0.0.0", "unknown", "V1.0.0", "a6c0bea1-ce1a-4fae-b943-1dbcc50cb311")
	if err == nil {
		t.Errorf("GetTargetState should have reported an unknown element")
	}

	_, err =  GetTargetState("demo", "app", "V0.0.0", "oam", "unknown", "a6c0bea1-ce1a-4fae-b943-1dbcc50cb311")
	if err == nil {
		t.Errorf("GetTargetState should have reported an unknown cluster")
	}

	_, err =  GetTargetState("demo", "app", "V0.0.0", "oam", "V1.0.0", "unknown")
	if err == nil {
		t.Errorf("GetTargetState should have reported an unknown instance")
	}

	targetState, err := GetTargetState("demo", "app", "V0.0.0", "oam", "V1.0.0", "a6c0bea1-ce1a-4fae-b943-1dbcc50cb311")
	if err != nil {
		t.Errorf("GetTargetState should have returned a target state")
	}

	// update with a derived current state
	currentState := CurrentState {
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

	// update with a correct context
	err = SetCurrentState(&currentState)
	if err != nil {
		t.Errorf("GetTargetState should not have returned an error")
	}

	// update with a wrong context
	currentState.Domain = "unknown"
	err = SetCurrentState(&currentState)
	if err == nil {
		t.Errorf("GetTargetState should returned an error")
	}
}

//------------------------------------------------------------------------------
