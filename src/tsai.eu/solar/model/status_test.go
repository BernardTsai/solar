package model

import (
	"testing"
)

//------------------------------------------------------------------------------

// TestStatus01 tests the basic functions of the status package.
func TestStatus01(t *testing.T) {
	model := GetModel()

	model.Load("testdata/testdata1.yaml")

	status, err :=  GetStatus("demo", "app", "oam", "V1.0.0", "a6c0bea1-ce1a-4fae-b943-1dbcc50cb311")
	if err != nil {
		t.Errorf("GetStatus should have returned a status")
	}

	err =  SetStatus(status)
	if err != nil {
		t.Errorf("GetStatus should have set a status")
	}

	_, err =  GetStatus("demo", "app", "oam", "V1.0.0", "unknown")
	if err == nil {
		t.Errorf("GetStatus should have returned an error")
	}

	status, err =  GetStatus("demo", "app", "oam", "V1.0.0", "")
	if err != nil {
		t.Errorf("GetStatus should have returned a status")
	}

	err =  SetStatus(status)
	if err != nil {
		t.Errorf("GetStatus should have set a status")
	}

	_, err =  GetStatus("demo", "app", "oam", "unknown", "")
	if err == nil {
		t.Errorf("GetStatus should have returned an error")
	}

	status, err =  GetStatus("demo", "app", "oam", "", "")
	if err != nil {
		t.Errorf("GetStatus should have returned a status")
	}

	err =  SetStatus(status)
	if err != nil {
		t.Errorf("GetStatus should have set a status")
	}

	_, err =  GetStatus("demo", "app", "unknown", "", "")
	if err == nil {
		t.Errorf("GetStatus should have returned an error")
	}

	_, err =  GetStatus("demo", "unknown", "", "", "")
	if err == nil {
		t.Errorf("GetStatus should have returned an error")
	}

	_, err =  GetStatus("unknown", "", "", "", "")
	if err == nil {
		t.Errorf("GetStatus should have returned an error")
	}
}

//------------------------------------------------------------------------------
