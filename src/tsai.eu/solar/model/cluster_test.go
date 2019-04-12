package model

import (
	"testing"
	"os"
)

//------------------------------------------------------------------------------

// TestCluster01 tests the basic functions of the cluster package.
func TestCluster01(t *testing.T) {
	filename := "test.yaml"

	// cleanup routine
  defer func() {os.Remove(filename)}()

	cluster, _ := NewCluster("V1.0.0", "active", 1, 1, 1, "")

	cluster.Save(filename)

	cluster.Load(filename)

	cluster.Show()
}

//------------------------------------------------------------------------------

// TestCluster02 tests the instance related functions of the cluster package.
func TestCluster02(t *testing.T) {
	model := GetModel()

	model.Load("testdata/testdata1.yaml")

	cluster, _ := GetCluster("demo", "app", "oam", "V1.0.0")

	_, err := cluster.ListInstances()
	if err != nil {
		t.Errorf("<cluster>.ListInstances should have returned a list of instance names")
	}

	_, err = cluster.GetInstance("unknown")
	if err == nil {
		t.Errorf("<cluster>.GetInstance should have complained about a non existing instance")
	}

	instance, err := cluster.GetInstance("a6c0bea1-ce1a-4fae-b943-1dbcc50cb311")
	if err != nil {
		t.Errorf("<cluster>.GetInstance should have returned a instance")
	}

	err = cluster.AddInstance(instance)
	if err == nil {
		t.Errorf("<cluster>.AddCluster should have complained about an already existing instance")
	}

	err = cluster.DeleteInstance("a6c0bea1-ce1a-4fae-b943-1dbcc50cb311")
	if err != nil {
		t.Errorf("<cluster>.DeleteInstance should have not have complained when deleting an existing instance")
	}

	err = cluster.DeleteInstance("a6c0bea1-ce1a-4fae-b943-1dbcc50cb311")
	if err == nil {
		t.Errorf("<cluster>.DeleteInstance should have have complained about attempting to delete a non-existing instance")
	}

	err = cluster.AddInstance(instance)
	if err != nil {
		t.Errorf("<cluster>.AddInstance should have not have complained when adding a new instance")
	}
}

//------------------------------------------------------------------------------

// TestCluster03 tests the relationship related functions of the cluster package.
func TestCluster03(t *testing.T) {
	model := GetModel()

	model.Load("testdata/testdata1.yaml")

	cluster, _ := GetCluster("demo", "app", "oam", "V1.0.0")

	_, err := cluster.ListRelationships()
	if err != nil {
		t.Errorf("<cluster>.ListRelationships should have returned a list of relationship names")
	}

	_, err = cluster.GetRelationship("unknown")
	if err == nil {
		t.Errorf("<cluster>.GetRelationship should have complained about a non existing relationship")
	}

	relationship, err := cluster.GetRelationship("tenant")
	if err != nil {
		t.Errorf("<cluster>.GetRelationship should have returned a relationship")
	}

	err = cluster.AddRelationship(relationship)
	if err == nil {
		t.Errorf("<cluster>.AddCluster should have complained about an already existing relationship")
	}

	err = cluster.DeleteRelationship("tenant")
	if err != nil {
		t.Errorf("<cluster>.DeleteRelationship should have not have complained when deleting an existing relationship")
	}

	err = cluster.DeleteRelationship("tenant")
	if err == nil {
		t.Errorf("<cluster>.DeleteRelationship should have have complained about attempting to delete a non-existing relationship")
	}

	err = cluster.AddRelationship(relationship)
	if err != nil {
		t.Errorf("<cluster>.AddRelationship should have not have complained when adding a new relationship")
	}
}

//------------------------------------------------------------------------------

// TestCluster04 tests misc. other functions of the cluster package.
func TestCluster04(t *testing.T) {
	model := GetModel()

	model.Load("testdata/testdata1.yaml")

	cluster, _ := GetCluster("demo", "app", "oam", "V1.0.0")

	cluster.Resize(0, 10, 5)

	cluster.Reset()

	cluster.SetState("inactive")
}

//------------------------------------------------------------------------------
