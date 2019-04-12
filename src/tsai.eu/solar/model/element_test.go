package model

import (
	"testing"
	"os"
)

//------------------------------------------------------------------------------

// TestElement01 tests the basic functions of the element package.
func TestElement01(t *testing.T) {
	filename := "test.yaml"

	// cleanup routine
  defer func() {os.Remove(filename)}()

	element, _ := NewElement("test", "V1.0.0", "")

	element.Save(filename)

	element.Load(filename)

	element.Show()
}

//------------------------------------------------------------------------------

// TestElement02 tests the cluster related functions of the element package.
func TestElement02(t *testing.T) {
	model := GetModel()

	model.Load("testdata/testdata1.yaml")

	element, _ := GetElement("demo", "app", "oam")

	_, err := element.ListClusters()
	if err != nil {
		t.Errorf("<element>.ListClusters should have returned a list of cluster names")
	}

	_, err = element.GetCluster("unknown")
	if err == nil {
		t.Errorf("<element>.GetCluster should have complained about a non existing cluster")
	}

	cluster, err := element.GetCluster("V1.0.0")
	if err != nil {
		t.Errorf("<element>.GetCluster should have returned a cluster")
	}

	err = element.AddCluster(cluster)
	if err == nil {
		t.Errorf("<element>.AddElement should have complained about an already existing cluster")
	}

	err = element.DeleteCluster("V1.0.0")
	if err != nil {
		t.Errorf("<element>.DeleteCluster should have not have complained when deleting an existing cluster")
	}

	err = element.DeleteCluster("V1.0.0")
	if err == nil {
		t.Errorf("<element>.DeleteCluster should have have complained about attempting to delete a non-existing cluster")
	}

	err = element.AddCluster(cluster)
	if err != nil {
		t.Errorf("<element>.AddCluster should have not have complained when adding a new cluster")
	}
}

//------------------------------------------------------------------------------
