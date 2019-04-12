package model

import (
	"testing"
	"os"
)

//------------------------------------------------------------------------------

// TestRelationship01 tests the basic functions of the relationship package.
func TestRelationship01(t *testing.T) {
	filename := "test.yaml"

	// cleanup routine
  defer func() {os.Remove(filename)}()

	relationship, _ := NewRelationship("Tenant", "tenant", "context", "demo", "app", "Tenant", "V1.0.0", "")

	relationship.Save(filename)

	relationship.Load(filename)

	relationship.Show()
}

//------------------------------------------------------------------------------

// TestRelationship02 tests the relationship related functions of the relationship package.
func TestRelationship02(t *testing.T) {
	relationship, _ := NewRelationship("Tenant", "tenant", "context", "demo", "app", "Tenant", "V1.0.0", "")

	relationship.Reset()
}

//------------------------------------------------------------------------------
