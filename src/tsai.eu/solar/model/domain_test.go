package model

import (
	"testing"
	"os"
)

//------------------------------------------------------------------------------

// TestDomain01 tests the basic functions of a domain.
func TestDomain01(t *testing.T) {
	filename := "test.yaml"

	// cleanup routine
  defer func() {os.Remove(filename)}()

	domain, _ := NewDomain("test")

	domain.Save(filename)

	domain.Load(filename)

	yaml, _ := domain.Show()

	domain.Load2(yaml)
}

//------------------------------------------------------------------------------

// TestDomain02 tests the component related functions of a domain.
func TestDomain02(t *testing.T) {
	model := GetModel()

	model.Load("testdata/testdata1.yaml")

	domain, _ := GetDomain("demo")

	_, err := domain.ListComponents()
	if err != nil {
		t.Errorf("<domain>.ListComponents should have returned a list of component names and versions")
	}

	_, err = domain.GetComponents()
	if err != nil {
		t.Errorf("<domain>.GetComponents should have returned a list of components")
	}

	_, err = domain.GetComponent("tenant", "V2.0.0")
	if err == nil {
		t.Errorf("<domain>.GetComponent should have complained about a non existing component")
	}

	component, err := domain.GetComponent("tenant", "V1.0.0")
	if err != nil {
		t.Errorf("<domain>.GetComponent should have returned a component")
	}

	err = domain.AddComponent(component)
	if err == nil {
		t.Errorf("<domain>.AddComponent should have complained about an already existing component")
	}

	component.Version = "V2.0.0"
	err = domain.AddComponent(component)
	if err != nil {
		t.Errorf("<domain>.AddComponent should have not have complained when adding a new component")
	}

	err = domain.DeleteComponent("tenant", "V2.0.0")
	if err != nil {
		t.Errorf("<domain>.DeleteComponent should have not have complained when deleting an existing component")
	}

	err = domain.DeleteComponent("tenant", "V2.0.0")
	if err == nil {
		t.Errorf("<domain>.DeleteComponent should have have complained about attempting to delete a non-existing component")
	}
}

//------------------------------------------------------------------------------

// TestDomain03 tests the architecture related functions of a domain.
func TestDomain03(t *testing.T) {
	model := GetModel()

	model.Load("testdata/testdata1.yaml")

	domain, _ := GetDomain("demo")

	_, err := domain.ListArchitectures()
	if err != nil {
		t.Errorf("<domain>.ListArchitectures should have returned a list of architecture names")
	}

	_, err = domain.GetArchitecture("app", "V2.0.0")
	if err == nil {
		t.Errorf("<domain>.GetArchitecture should have complained about a non existing architecture")
	}

	architecture, err := domain.GetArchitecture("app", "V0.0.0")
	if err != nil {
		t.Errorf("<domain>.GetArchitecture should have returned an architecture")
	}

	err = domain.AddArchitecture(architecture)
	if err == nil {
		t.Errorf("<domain>.AddArchitecture should have complained about an already existing architecture")
	}

	architecture.Version = "V2.0.0"
	err = domain.AddArchitecture(architecture)
	if err != nil {
		t.Errorf("<domain>.AddArchitecture should have not have complained when adding a new architecture")
	}

	err = domain.DeleteArchitecture("app", "V2.0.0")
	if err != nil {
		t.Errorf("<domain>.DeleteArchitecture should have not have complained when deleting an existing architecture")
	}

	err = domain.DeleteArchitecture("app", "V2.0.0")
	if err == nil {
		t.Errorf("<domain>.DeleteArchitecture should have have complained about attempting to delete a non-existing architecture")
	}
}

//------------------------------------------------------------------------------

// TestDomain04 tests the solution related functions of a domain.
func TestDomain04(t *testing.T) {
	model := GetModel()

	model.Load("testdata/testdata1.yaml")

	domain, _ := GetDomain("demo")

	_, err := domain.ListSolutions()
	if err != nil {
		t.Errorf("<domain>.ListSolutions should have returned a list of solution names")
	}

	_, err = domain.GetSolution("unknown")
	if err == nil {
		t.Errorf("<domain>.GetSolution should have complained about a non existing solution")
	}

	solution, err := domain.GetSolution("app")
	if err != nil {
		t.Errorf("<domain>.GetSolution should have returned a solution")
	}

	err = domain.AddSolution(solution)
	if err == nil {
		t.Errorf("<domain>.AddSolution should have complained about an already existing solution")
	}

	err = domain.DeleteSolution("unknown")
	if err == nil {
		t.Errorf("<domain>.DeleteSolution should have not have complained when deleting an existing solution")
	}

	err = domain.DeleteSolution("app")
	if err != nil {
		t.Errorf("<domain>.DeleteSolution should have have complained about attempting to delete a non-existing solution")
	}

	err = domain.AddSolution(solution)
	if err != nil {
		t.Errorf("<domain>.AddSolution should have not have complained when adding a new solution")
	}
}

//------------------------------------------------------------------------------

// TestDomain05 tests the task related functions of a domain.
func TestDomain05(t *testing.T) {
	model := GetModel()

	model.Load("testdata/testdata1.yaml")

	domain, _ := GetDomain("demo")

	_, err := domain.ListTasks()
	if err != nil {
		t.Errorf("<domain>.ListTasks should have returned a list of task names")
	}

	_, err = domain.GetTask("unknown")
	if err == nil {
		t.Errorf("<domain>.GetTask should have complained about a non existing task")
	}

	task, err := domain.GetTask("0f7bec76-3252-41b0-a7c9-82b75abba9de")
	if err != nil {
		t.Errorf("<domain>.GetTask should have returned a task")
	}

	err = domain.AddTask(task)
	if err == nil {
		t.Errorf("<domain>.AddTask should have complained about an already existing task")
	}

	err = domain.DeleteTask("unknown")
	if err == nil {
		t.Errorf("<domain>.DeleteTask should have not have complained when deleting an existing task")
	}

	err = domain.DeleteTask("0f7bec76-3252-41b0-a7c9-82b75abba9de")
	if err != nil {
		t.Errorf("<domain>.DeleteTask should have have complained about attempting to delete a non-existing task")
	}

	err = domain.AddTask(task)
	if err != nil {
		t.Errorf("<domain>.AddTask should have not have complained when adding a new task")
	}
}

//------------------------------------------------------------------------------

// TestDomain06 tests the event related functions of a domain.
func TestDomain06(t *testing.T) {
	model := GetModel()

	model.Load("testdata/testdata1.yaml")

	domain, _ := GetDomain("demo")

	_, err := domain.ListEvents()
	if err != nil {
		t.Errorf("<domain>.ListEvents should have returned a list of event names")
	}

	_, err = domain.GetEvent("unknown")
	if err == nil {
		t.Errorf("<domain>.GetEvent should have complained about a non existing event")
	}

	event, err := domain.GetEvent("02ed6397-c8fc-4334-983f-f40dd87475cf")
	if err != nil {
		t.Errorf("<domain>.GetEvent should have returned a event")
	}

	err = domain.AddEvent(event)
	if err == nil {
		t.Errorf("<domain>.AddEvent should have complained about an already existing event")
	}

	err = domain.DeleteEvent("unknown")
	if err == nil {
		t.Errorf("<domain>.DeleteEvent should have not have complained when deleting an existing event")
	}

	err = domain.DeleteEvent("02ed6397-c8fc-4334-983f-f40dd87475cf")
	if err != nil {
		t.Errorf("<domain>.DeleteEvent should have have complained about attempting to delete a non-existing event")
	}

	err = domain.AddEvent(event)
	if err != nil {
		t.Errorf("<domain>.AddEvent should have not have complained when adding a new event")
	}
}

//------------------------------------------------------------------------------

// TestDomain07 tests the controller related functions of a domain.
func TestDomain07(t *testing.T) {
	model := GetModel()

	model.Load("testdata/testdata1.yaml")

	domain, _ := GetDomain("demo")

	_, err := domain.GetController("unknown", "unknown")
	if err == nil {
		t.Errorf("<domain>.GetController should have complained about a non existing controller")
	}

	controller, _ := NewController("controller", "V1.0.0")
	err = domain.AddController(controller)
	if err != nil {
		t.Errorf("<domain>.AddController should not have reported a failure")
	}

	err = domain.AddController(controller)
	if err == nil {
		t.Errorf("<domain>.AddController should have complained about an already existing controller")
	}

	_, err = domain.GetController("controller", "V1.0.0")
	if err != nil {
		t.Errorf("<domain>.GetController should have returned a controller")
	}

	_, err = domain.ListControllers()
	if err != nil {
		t.Errorf("<domain>.ListControllers should have returned a list of controller images and versions")
	}

	err = domain.DeleteController("controller", "V1.0.0")
	if err != nil {
		t.Errorf("<domain>.DeleteController should have not have complained when deleting an existing controller")
	}

	err = domain.DeleteController("unknown", "unknown")
	if err == nil {
		t.Errorf("<domain>.DeleteController should have have complained about attempting to delete a non-existing controller")
	}
}

//------------------------------------------------------------------------------
