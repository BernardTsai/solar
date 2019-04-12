package model

import (
	"testing"
	"os"
)

//------------------------------------------------------------------------------

// TestRelationshipConfiguration01 tests the basic functions of the relationshipConfiguration package.
func TestRelationshipConfiguration01(t *testing.T) {
	filename := "test.yaml"

	// cleanup routine
  defer func() {os.Remove(filename)}()

	relationshipConfiguration, _ := NewRelationshipConfiguration("tenant", "tenant", "context", "tenant", "V1.0.0", "")

	relationshipConfiguration.Save(filename)

	relationshipConfiguration.Load(filename)

	relationshipConfiguration.Show()
}

//------------------------------------------------------------------------------
