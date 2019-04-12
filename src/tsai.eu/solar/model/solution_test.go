package model

import (
	"testing"
	"os"
)

//------------------------------------------------------------------------------

// TestSolution01 tests the basic functions of the solution package.
func TestSolution01(t *testing.T) {
	filename := "test.yaml"

	// cleanup routine
  defer func() {os.Remove(filename)}()

	solution, _ := NewSolution("app", "V1.0.0", "")

	solution.Save(filename)

	solution.Load(filename)

	yaml, _ := solution.Show()

	solution.Load2(yaml)
}

//------------------------------------------------------------------------------

// TestSolution02 tests the element related functions of the solution package.
func TestSolution02(t *testing.T) {
	model := GetModel()

	model.Load("testdata/testdata1.yaml")

	solution, _ := GetSolution("demo", "app")

	_, err := solution.ListElements()
	if err != nil {
		t.Errorf("<solution>.ListElements should have returned a list of element names")
	}

	_, err = solution.GetElement("unknown")
	if err == nil {
		t.Errorf("<solution>.GetElement should have complained about a non existing element")
	}

	element, err := solution.GetElement("tenant")
	if err != nil {
		t.Errorf("<solution>.GetElement should have returned a element")
	}

	err = solution.AddElement(element)
	if err == nil {
		t.Errorf("<solution>.AddSolution should have complained about an already existing element")
	}

	err = solution.DeleteElement("tenant")
	if err != nil {
		t.Errorf("<solution>.DeleteElement should have not have complained when deleting an existing element")
	}

	err = solution.DeleteElement("tenant")
	if err == nil {
		t.Errorf("<solution>.DeleteElement should have have complained about attempting to delete a non-existing solution")
	}

	err = solution.AddElement(element)
	if err != nil {
		t.Errorf("<solution>.AddElement should have not have complained when adding a new element")
	}
}

//------------------------------------------------------------------------------

// TestSolution03 tests the state/transition related functions of the solution package.
func TestSolution03(t *testing.T) {
	if IsValidStateOrTransition("abc") {
		t.Errorf("IsValidStateOrTransition should have reported false")
	}

	if !IsValidStateOrTransition("inactive"){
		t.Errorf("IsValidStateOrTransition should have reported true")
	}

	if IsValidState("abc") {
		t.Errorf("IsValidStateOrTransition should have reported false")
	}

	if !IsValidState("inactive") {
		t.Errorf("IsValidStateOrTransition should have reported true")
	}

	if IsValidTransition("abc") {
		t.Errorf("IsValidTransition should have reported false")
	}

	if !IsValidTransition("creating") {
		t.Errorf("IsValidTransition should have reported true")
	}

	_, err := GetTransition("initial", "active")
	if err != nil {
		t.Errorf("GetTransition should have not reported an error")
	}

	_, err = GetTransition("initial", "unkown")
	if err == nil {
		t.Errorf("GetTransition should have reported an error")
	}

	_, err = GetTransition("initial", "failure")
	if err == nil {
		t.Errorf("GetTransition should have reported an error")
	}
}

//------------------------------------------------------------------------------

// TestSolution04 tests the other misc. functions of the solution package.
func TestSolution04(t *testing.T) {
	model := GetModel()

	model.Load("testdata/testdata1.yaml")

	solution, _ := GetSolution("demo", "app")

	ok := solution.OK()
	if !ok {
		t.Errorf("<solution>.OK() should have reported true")
	}

	cluster, _ := GetCluster("demo", "app", "tenant", "V1.0.0")
	cluster.State = InitialState
	ok = solution.OK()
	if ok {
		t.Errorf("<solution>.OK() should have reported false")
	}
}

//------------------------------------------------------------------------------
