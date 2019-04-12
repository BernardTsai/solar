package model

import (
	"testing"
)

//------------------------------------------------------------------------------

// TestSetup01 tests the basic functions of the setup package.
func TestSetup01(t *testing.T) {
	model := GetModel()

	model.Load("testdata/testdata1.yaml")

	_, err :=  GetSetup("unknown", "app", "V0.0.0", "oam", "V1.0.0", "a6c0bea1-ce1a-4fae-b943-1dbcc50cb311")
	if err == nil {
		t.Errorf("GetSetup should have reported an unknown domain")
	}

	_, err =  GetSetup("demo", "unknown", "V0.0.0", "oam", "V1.0.0", "a6c0bea1-ce1a-4fae-b943-1dbcc50cb311")
	if err == nil {
		t.Errorf("GetSetup should have reported an unknown solution")
	}

	_, err =  GetSetup("demo", "app", "unknown", "oam", "V1.0.0", "a6c0bea1-ce1a-4fae-b943-1dbcc50cb311")
	if err == nil {
		t.Errorf("GetSetup should have reported an unknown architecture")
	}

	_, err =  GetSetup("demo", "app", "V0.0.0", "unknown", "V1.0.0", "a6c0bea1-ce1a-4fae-b943-1dbcc50cb311")
	if err == nil {
		t.Errorf("GetSetup should have reported an unknown element")
	}

	_, err =  GetSetup("demo", "app", "V0.0.0", "oam", "unknown", "a6c0bea1-ce1a-4fae-b943-1dbcc50cb311")
	if err == nil {
		t.Errorf("GetSetup should have reported an unknown cluster")
	}

		_, err =  GetSetup("demo", "app", "V0.0.0", "oam", "V1.0.0", "unknown")
	if err == nil {
		t.Errorf("GetSetup should have reported an unknown instance")
	}

	setup, err :=  GetSetup("demo", "app", "V0.0.0", "oam", "V1.0.0", "a6c0bea1-ce1a-4fae-b943-1dbcc50cb311")
	if err != nil {
		t.Errorf("GetSetup should have returned a setup")
	}

	err =  SetSetup(setup)
	if err != nil {
		t.Errorf("GetSetup should have set a setup")
	}
}

//------------------------------------------------------------------------------
