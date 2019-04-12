package model

import (
	"testing"
	"os"
)

//------------------------------------------------------------------------------

// TestModel01 tests the top level functions of the model package.
func TestModel01(t *testing.T) {
	filename := "test.yaml"

	// cleanup routine
  defer func() {os.Remove(filename)}()

	model := GetModel()

	model.Reset() 

	model.Save(filename)

	model.Load(filename)

	yaml, _ := model.Show()

	model.Load2(yaml)

	model.ListDomains()
}

//------------------------------------------------------------------------------

// TestModel02 tests the domain related functions of the model package.
func TestModel02(t *testing.T) {
	filename := "test.yaml"

	// cleanup routine
  defer func() {os.Remove(filename)}()

	model := GetModel()

	domain, _ := NewDomain("test")
	model.AddDomain(domain)

	err := domain.Save(filename)
	if err != nil {
		t.Errorf("Domain should have been saved")
	}

	err = domain.Load(filename)
	if err != nil {
		t.Errorf("Domain should have been loaded from a file")
	}

	yaml, err := domain.Show()
	if err != nil {
		t.Errorf("Domain should have been exported as yaml")
	}

	err = domain.Load2(yaml)
	if err != nil {
		t.Errorf("Domain should have been loaded from yaml")
	}

	err = model.AddDomain(domain)
	if err == nil {
		t.Errorf("AddDomain should have complained about an already existing domain")
	}

	GetDomains()

	_, err = GetDomain("unknown")
	if err == nil {
		t.Errorf("GetDomain should have complained about an unknown domain")
	}

	model.DeleteDomain("test")
	err = model.DeleteDomain("test")
	if err == nil {
		t.Errorf("GetDomain should have complained about not being able to delete an unknown domain")
	}



}

//------------------------------------------------------------------------------

// TestModel03 tests the component related functions of the model package.
func TestModel03(t *testing.T) {
}

//------------------------------------------------------------------------------

// TestModelXY tests the solution related functions of the model package.
func TestModelXY(t *testing.T) {
	model := GetModel()

	model.Load("testdata/testdata1.yaml")

	solution, _     := GetSolution("demo", "app")
	architecture, _ := GetArchitecture("demo", "app - V0.0.0")

	solution.Update("demo", architecture)

	// domain related functions
	_, err := GetDomain("demo")
	if err != nil {
		t.Errorf("GetDomain should have found an existing domain")
	}

	_, err = GetDomain("unknown")
	if err == nil {
		t.Errorf("GetDomain should have complained about an non existing domain")
	}

	// component related functions
	_, err = GetComponent("demo", "tenant - V1.0.0")
	if err != nil {
		t.Errorf("GetComponent should have found an existing component")
	}

	_, err = GetComponent("demo", "unknown")
	if err == nil {
		t.Errorf("GetComponent should have complained about an non existing component")
	}

	_, err = GetComponent("unknown", "unknown")
	if err == nil {
		t.Errorf("GetComponent should have complained about an non existing domain")
	}

	// architecture related functions
	_, err = GetArchitecture("demo", "app - V0.0.0")
	if err != nil {
		t.Errorf("GetArchitecture should have found an existing architecture")
	}

	_, err = GetArchitecture("demo", "unknown")
	if err == nil {
		t.Errorf("GetArchitecture should have complained about an non existing architecture")
	}

	_, err = GetArchitecture("unknown", "unknown")
	if err == nil {
		t.Errorf("GetArchitecture should have complained about an non existing domain")
	}

	// solution related functions
	_, err = GetSolution("demo", "app")
	if err != nil {
		t.Errorf("GetSolution should have found an existing solution")
	}

	_, err = GetSolution("demo", "unknown")
	if err == nil {
		t.Errorf("GetSolution should have complained about an non existing solution")
	}

	_, err = GetSolution("unknown", "unknown")
	if err == nil {
		t.Errorf("GetSolution should have complained about an non existing domain")
	}

	// element related functions
	_, err = GetElement("demo", "app", "tenant")
	if err != nil {
		t.Errorf("GetElement should have found an existing element")
	}

	_, err = GetElement("demo", "app", "unknown")
	if err == nil {
		t.Errorf("GetElement should have complained about an non existing element")
	}

	_, err = GetElement("demo", "unknown", "unknown")
	if err == nil {
		t.Errorf("GetElement should have complained about an non existing solution")
	}

	_, err = GetElement("unknown", "unknown", "unknown")
	if err == nil {
		t.Errorf("GetElement should have complained about an non existing domain")
	}

	// cluster related functions
	_, err = GetCluster("demo", "app", "tenant", "V1.0.0")
	if err != nil {
		t.Errorf("GetCluster should have found an existing cluster")
	}

	_, err = GetCluster("demo", "app", "tenant", "unknown")
	if err == nil {
		t.Errorf("GetCluster should have complained about an non existing cluster")
	}

	_, err = GetCluster("demo", "app", "unknown", "unknown")
	if err == nil {
		t.Errorf("GetCluster should have complained about an non existing element")
	}

	_, err = GetCluster("demo", "unknown", "unknown", "unknown")
	if err == nil {
		t.Errorf("GetCluster should have complained about an non existing solution")
	}

	_, err = GetCluster("unknown", "unknown", "unknown", "unknown")
	if err == nil {
		t.Errorf("GetCluster should have complained about an non existing domain")
	}

	// component related functions part 2
	_, err = GetComponent2("demo", "app", "tenant", "V1.0.0")
	if err != nil {
		t.Errorf("GetComponent2 should have found an existing component")
	}

	_, err = GetComponent2("demo", "app", "tenant", "unknown")
	if err == nil {
		t.Errorf("GetComponent2 should have complained about an non existing cluster")
	}

	_, err = GetComponent2("demo", "app", "unknown", "unknown")
	if err == nil {
		t.Errorf("GetComponent2 should have complained about an non existing element")
	}

	_, err = GetComponent2("demo", "unknown", "unknown", "unknown")
	if err == nil {
		t.Errorf("GetComponent2 should have complained about an non existing solution")
	}

	_, err = GetComponent2("unknown", "unknown", "unknown", "unknown")
	if err == nil {
		t.Errorf("GetComponent2 should have complained about an non existing domain")
	}

	// instance related functions
	_, err = GetInstance("demo", "app", "tenant", "V1.0.0", "cda9d59e-4bfc-4eae-bfee-3834b6c11955")
	if err != nil {
		t.Errorf("GetInstance should have found an existing relationship")
	}

	_, err = GetInstance("demo", "app", "tenant", "V1.0.0", "unknown")
	if err == nil {
		t.Errorf("GetInstance should have complained about an non existing instance")
	}

	_, err = GetInstance("demo", "app", "tenant", "unknown", "unknown")
	if err == nil {
		t.Errorf("GetInstance should have complained about an non existing cluster")
	}

	_, err = GetInstance("demo", "app", "unknown", "unknown", "unknown")
	if err == nil {
		t.Errorf("GetInstance should have complained about an non existing element")
	}

	_, err = GetInstance("demo", "unknown", "unknown", "unknown", "unknown")
	if err == nil {
		t.Errorf("GetInstance should have complained about an non existing solution")
	}

	_, err = GetInstance("unknown", "unknown", "unknown", "unknown", "unknown")
	if err == nil {
		t.Errorf("GetInstance should have complained about an non existing domain")
	}

	// relationship related functions
	_, err = GetRelationship("demo", "app", "db", "V1.0.0", "server")
	if err != nil {
		t.Errorf("GetRelationship should have found an existing relationship")
	}

	_, err = GetRelationship("demo", "app", "tenant", "V1.0.0", "unknown")
	if err == nil {
		t.Errorf("GetRelationship should have complained about an non existing relationship")
	}

	_, err = GetRelationship("demo", "app", "tenant", "unknown", "unknown")
	if err == nil {
		t.Errorf("GetRelationship should have complained about an non existing cluster")
	}

	_, err = GetRelationship("demo", "app", "unknown", "unknown", "unknown")
	if err == nil {
		t.Errorf("GetRelationship should have complained about an non existing element")
	}

	_, err = GetRelationship("demo", "unknown", "unknown", "unknown", "unknown")
	if err == nil {
		t.Errorf("GetRelationship should have complained about an non existing solution")
	}

	_, err = GetRelationship("unknown", "unknown", "unknown", "unknown", "unknown")
	if err == nil {
		t.Errorf("GetRelationship should have complained about an non existing domain")
	}

	// task related functions
	_, err = GetTask("demo", "0f7bec76-3252-41b0-a7c9-82b75abba9de")
	if err != nil {
		t.Errorf("GetTask should have found an existing task")
	}

	_, err = GetTask("demo", "unknown")
	if err == nil {
		t.Errorf("GetTask should have complained about an non existing task")
	}

	_, err = GetTask("unknown", "unknown")
	if err == nil {
		t.Errorf("GetTask should have complained about an non existing domain")
	}

	// event related functions
	_, err = GetEvent("demo", "09b2fe2b-b7da-482a-a690-49ab0b24dcc3")
	if err != nil {
		t.Errorf("GetEvent should have found an existing event")
	}

	_, err = GetEvent("demo", "unknown")
	if err == nil {
		t.Errorf("GetEvent should have complained about an non existing event")
	}

	_, err = GetEvent("unknown", "unknown")
	if err == nil {
		t.Errorf("GetEvent should have complained about an non existing domain")
	}
}

//------------------------------------------------------------------------------
