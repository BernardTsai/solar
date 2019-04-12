package model

import (
	"testing"
	"os"
)

//------------------------------------------------------------------------------

// TestDependency01 tests the basic functions of the dependency package.
func TestDependency01(t *testing.T) {
	filename := "test.yaml"

	// cleanup routine
  defer func() {os.Remove(filename)}()

	dependency, _ := NewDependency("tenant", "context", "tenant", "V1.0.0", "")

	dependency.Save(filename)

	dependency.Load(filename)

	dependency.Show()
}

//------------------------------------------------------------------------------
