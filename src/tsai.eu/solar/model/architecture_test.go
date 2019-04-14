package model

import (
	"testing"
	"os"
)

//------------------------------------------------------------------------------

// TestArchitecture01 tests the basic functions of the architecture package.
func TestArchitecture01(t *testing.T) {
	filename := "test.yaml"

	// cleanup routine
  defer func() {os.Remove(filename)}()

	architecture, _ := NewArchitecture("test", "V1.0.0", "")

	architecture.Save(filename)

	architecture.Load(filename)

	yaml, _ := architecture.Show()

	architecture.Load2(yaml)
}

//------------------------------------------------------------------------------

// TestArchitecture02 tests the element related functions of the architecture package.
func TestArchitecture02(t *testing.T) {
	model := GetModel()

	model.Load("testdata/testdata1.yaml")

	architecture, _ := GetArchitecture("demo", "app", "V0.0.0")

	_, err := architecture.ListElements()
	if err != nil {
		t.Errorf("<architecture>.ListElements should have returned a list of element names")
	}

	_, err = architecture.GetElement("unknown")
	if err == nil {
		t.Errorf("<architecture>.GetElement should have complained about a non existing element")
	}

	element, err := architecture.GetElement("tenant")
	if err != nil {
		t.Errorf("<architecture>.GetElement should have returned a element")
	}

	err = architecture.AddElement(element)
	if err == nil {
		t.Errorf("<architecture>.AddArchitecture should have complained about an already existing element")
	}

	err = architecture.DeleteElement("tenant")
	if err != nil {
		t.Errorf("<architecture>.DeleteElement should have not have complained when deleting an existing element")
	}

	err = architecture.DeleteElement("tenant")
	if err == nil {
		t.Errorf("<architecture>.DeleteElement should have have complained about attempting to delete a non-existing architecture")
	}

	err = architecture.AddElement(element)
	if err != nil {
		t.Errorf("<architecture>.AddElement should have not have complained when adding a new element")
	}
}

//------------------------------------------------------------------------------
