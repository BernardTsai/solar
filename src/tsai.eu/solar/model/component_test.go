package model

import (
	"testing"
	"os"
)

//------------------------------------------------------------------------------

// TestComponent01 tests the basic functions of the component package.
func TestComponent01(t *testing.T) {
	filename := "test.yaml"

	// cleanup routine
  defer func() {os.Remove(filename)}()

	component, _ := NewComponent("test", "V1.0.0", "")

	component.Save(filename)

	component.Load(filename)

	yaml, _ := component.Show()

	component.Load2(yaml)
}

//------------------------------------------------------------------------------

// TestComponent02 tests the dependency related functions of the component package.
func TestComponent02(t *testing.T) {
	model := GetModel()

	model.Load("testdata/testdata1.yaml")

	component, _ := GetComponent("demo", "network", "V1.0.0")

	_, err := component.ListDependencies()
	if err != nil {
		t.Errorf("<component>.ListDependencies should have returned a list of dependency names")
	}

	_, err = component.GetDependency("unknown")
	if err == nil {
		t.Errorf("<component>.GetDependency should have complained about a non existing dependency")
	}

	dependency, err := component.GetDependency("tenant")
	if err != nil {
		t.Errorf("<component>.GetDependency should have returned a dependency")
	}

	err = component.AddDependency(dependency)
	if err == nil {
		t.Errorf("<component>.AddComponent should have complained about an already existing dependency")
	}

	err = component.DeleteDependency("tenant")
	if err != nil {
		t.Errorf("<component>.DeleteDependency should have not have complained when deleting an existing dependency")
	}

	err = component.DeleteDependency("tenant")
	if err == nil {
		t.Errorf("<component>.DeleteDependency should have have complained about attempting to delete a non-existing component")
	}

	err = component.AddDependency(dependency)
	if err != nil {
		t.Errorf("<component>.AddDependency should have not have complained when adding a new dependency")
	}
}

//------------------------------------------------------------------------------
