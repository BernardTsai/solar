package model

import (
	"testing"
	"os"
)

//------------------------------------------------------------------------------

// TestEvent01 tests the basic functions of the event package.
func TestEvent01(t *testing.T) {
	filename := "test.yaml"

	// cleanup routine
  defer func() {os.Remove(filename)}()

	event := NewEvent("demo", "task-uuid", EventTypeTaskExecution, "", "no comment")

	event.Save(filename)

	event.Load(filename)

	event.Show()
}

//------------------------------------------------------------------------------

// TestEvent02 tests the element related functions of the event package.
func TestEvent02(t *testing.T) {
	event := NewEvent("demo", "task-uuid", EventTypeTaskExecution, "", "no comment")

	if event.GetUUID() != event.UUID {
		t.Errorf("<event>.GetUUID should have returned the uuid of the event")
	}
}

//------------------------------------------------------------------------------
