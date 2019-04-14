package model

import (
	"testing"
	"os"
)

//------------------------------------------------------------------------------

// TestClusterConfiguration01 tests the basic functions of the clusterConfiguration package.
func TestClusterConfiguration01(t *testing.T) {
	filename := "test.yaml"

	// cleanup routine
  defer func() {os.Remove(filename)}()

	clusterConfiguration, _ := NewClusterConfiguration("V1.0.0", "active", 1, 1, 1, "")

	clusterConfiguration.Save(filename)

	clusterConfiguration.Load(filename)

	clusterConfiguration.Show()
}

//------------------------------------------------------------------------------

// TestClusterConfiguration02 tests the relationship related functions of the clusterConfiguration package.
func TestClusterConfiguration02(t *testing.T) {
	model := GetModel()

	model.Load("testdata/testdata1.yaml")

	architecture,         _ := GetArchitecture("demo", "app", "V0.0.0")
	elementConfiguration, _ := architecture.GetElement("oam")
	clusterConfiguration, _ := elementConfiguration.GetCluster("V1.0.0")

	_, err := clusterConfiguration.ListRelationships()
	if err != nil {
		t.Errorf("<clusterConfiguration>.ListRelationships should have returned a list of relationship names")
	}

	_, err = clusterConfiguration.GetRelationship("unknown")
	if err == nil {
		t.Errorf("<clusterConfiguration>.GetRelationship should have complained about a non existing relationship")
	}

	relationship, err := clusterConfiguration.GetRelationship("tenant")
	if err != nil {
		t.Errorf("<clusterConfiguration>.GetCluster should have returned a relationship")
	}

	err = clusterConfiguration.AddRelationship(relationship)
	if err == nil {
		t.Errorf("<clusterConfiguration>.AddCluster should have complained about an already existing relationship")
	}

	err = clusterConfiguration.DeleteRelationship("tenant")
	if err != nil {
		t.Errorf("<clusterConfiguration>.DeleteRelationship should have not have complained when deleting an existing relationship")
	}

	err = clusterConfiguration.DeleteRelationship("tenant")
	if err == nil {
		t.Errorf("<clusterConfiguration>.DeleteRelationship should have have complained about attempting to delete a non-existing cluster")
	}

	err = clusterConfiguration.AddRelationship(relationship)
	if err != nil {
		t.Errorf("<clusterConfiguration>.AddCluster should have not have complained when adding a new relationship")
	}
}

//------------------------------------------------------------------------------
