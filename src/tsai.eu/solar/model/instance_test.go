package model

import (
	"testing"
	"os"
)

//------------------------------------------------------------------------------

// TestInstance01 tests the basic functions of the instance package.
func TestInstance01(t *testing.T) {
	filename := "test.yaml"

	// cleanup routine
  defer func() {os.Remove(filename)}()

	instance, _ := NewInstance("a6c0bea1-ce1a-4fae-b943-1dbcc50cb311", "active", "")

	instance.Save(filename)

	instance.Load(filename)

	instance.Show()
}

//------------------------------------------------------------------------------

// TestInstance02 tests the instance related functions of the instance package.
func TestInstance02(t *testing.T) {
	instance, _ := NewInstance("a6c0bea1-ce1a-4fae-b943-1dbcc50cb311", "active", "")

	instance.SetState("active")

	if !instance.OK() {
		t.Errorf("<instance>.OK should have have complained returned true")
	}

	instance.Reset()

	if instance.OK() {
		t.Errorf("<instance>.OK should have have complained returned false")
	}
}

//------------------------------------------------------------------------------
