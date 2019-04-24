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

	controller, _ := NewController("default", "V1.0.0")

	controller.Save(filename)

	controller.Load(filename)

	yaml, _ := controller.Show()

	controller.Load2(yaml)
}

//------------------------------------------------------------------------------
