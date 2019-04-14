package model

import (
	"testing"
	"os"
)

//------------------------------------------------------------------------------

// TestElementConfiguration01 tests the basic functions of the elementConfiguration package.
func TestElementConfiguration01(t *testing.T) {
	filename := "test.yaml"

	// cleanup routine
  defer func() {os.Remove(filename)}()

	elementConfiguration, _ := NewElementConfiguration("tenant", "V1.0.0", "")

	elementConfiguration.Save(filename)

	elementConfiguration.Load(filename)

	elementConfiguration.Show()
}

//------------------------------------------------------------------------------

// TestElementConfiguration02 tests the cluster related functions of the elementConfiguration package.
func TestElementConfiguration02(t *testing.T) {
	model := GetModel()

	model.Load("testdata/testdata1.yaml")

	architecture, _ := GetArchitecture("demo", "app", "V0.0.0")

	elementConfiguration, _ := architecture.GetElement("tenant")

	_, err := elementConfiguration.ListClusters()
	if err != nil {
		t.Errorf("<elementConfiguration>.ListClusters should have returned a list of cluster names")
	}

	_, err = elementConfiguration.GetCluster("unknown")
	if err == nil {
		t.Errorf("<elementConfiguration>.GetCluster should have complained about a non existing cluster")
	}

	cluster, err := elementConfiguration.GetCluster("V1.0.0")
	if err != nil {
		t.Errorf("<elementConfiguration>.GetCluster should have returned a cluster")
	}

	err = elementConfiguration.AddCluster(cluster)
	if err == nil {
		t.Errorf("<elementConfiguration>.AddCluster should have complained about an already existing cluster")
	}

	err = elementConfiguration.DeleteCluster("V1.0.0")
	if err != nil {
		t.Errorf("<elementConfiguration>.DeleteCluster should have not have complained when deleting an existing cluster")
	}

	err = elementConfiguration.DeleteCluster("V1.0.0")
	if err == nil {
		t.Errorf("<elementConfiguration>.DeleteCluster should have have complained about attempting to delete a non-existing cluster")
	}

	err = elementConfiguration.AddCluster(cluster)
	if err != nil {
		t.Errorf("<elementConfiguration>.AddCluster should have not have complained when adding a new cluster")
	}
}

//------------------------------------------------------------------------------
