package model

import (
	"testing"
	"os"
)

//------------------------------------------------------------------------------

// TestController01 tests the basic functions of the controller.
func TestController01(t *testing.T) {
	filename := "test.yaml"

	// cleanup routine
  defer func() {os.Remove(filename)}()

	controller, _ := NewController("controller", "V1.0.0")

	controller.Save(filename)

	controller.Load(filename)

	yaml, _ := controller.Show()

	controller.Load2(yaml)
}

//------------------------------------------------------------------------------

// TestController02 tests the attribute related functions of the controller.
func TestController02(t *testing.T) {
	controller, _ := NewController("controller", "V1.0.0")

	err := controller.AddType("tenant", "V1.0.0")
	if err != nil {
		t.Errorf("AddType should not have reported a failure")
	}

	err = controller.AddType("tenant", "V1.0.0")
	if err == nil {
		t.Errorf("AddType should have reported a request to create a duplicate entry")
	}

	err = controller.DeleteType("tenant", "V2.0.0")
	if err == nil {
		t.Errorf("DeleteType should  have reported a request to delete an non existing entry")
	}

	err = controller.DeleteType("tenant", "V1.0.0")
	if err != nil {
		t.Errorf("DeleteType should not have reported a failure")
	}

	_, err = controller.ListTypes()
	if err != nil {
		t.Errorf("ListTypes should not have reported a failure")
	}
}

//------------------------------------------------------------------------------
