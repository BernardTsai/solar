package model

import (
	"sync"

	"github.com/pkg/errors"
	"tsai.eu/solar/util"
)

//------------------------------------------------------------------------------
// Model
// =====
//
// Attributes:
//   - Schema
//   - Name
//   - Domains
//
// Functions:
//   - GetModel
//   - GetDomain
//   - GetComponent
//   - GetArchitecture
//   - GetSolution
//   - GetElement
//   - GetCluster
//   - GetComponent2
//   - GetInstance
//   - GetRelationship
//   - NewModel
//
//   - model.Show
//   - model.Load
//   - model.Save
//
//   - model.ListDomains
//   - model.GetDomain
//   - model.AddDomain
//   - model.DeleteDomain
//------------------------------------------------------------------------------

// DomainMap is a synchronized map for a map of domains
type DomainMap struct {
	Map  map[string]*Domain  `yaml:"map"`             // map of domains
}

// MarshalYAML marshals a EventMap into yaml
func (m DomainMap) MarshalYAML() (interface{}, error) {
	return m.Map, nil
}

// UnmarshalYAML unmarshals a DomainMap from yaml
func (m *DomainMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	Map := map[string]*Domain{}

	err := unmarshal(&Map)
	if err != nil {
		return err
	}

	*m = DomainMap{Map: Map}

	return nil
}

//------------------------------------------------------------------------------

// Model describes all managed artefacts within a model.
type Model struct {
	Schema  string    `yaml:"Schema"`  // schema of the model
	Name    string    `yaml:"Name"`    // name of the model
	Domains DomainMap `yaml:"Domains"` // map of domains
}

var theModel *Model

var modelInit sync.Once

// GetModel retrieves a controller for a specific component type.
func GetModel() *Model {
	// initialise singleton once
	modelInit.Do(func() { theModel, _ = NewModel() })

	// success
	return theModel
}

//------------------------------------------------------------------------------

// NewModel creates a new model
func NewModel() (*Model, error) {
	var model Model

	model.Reset()

	// success
	return &model, nil
}

//------------------------------------------------------------------------------

// Reset resets all model data to its initial values
func (model *Model) Reset() error {
	model.Schema = "BT V1.0.0"
	model.Name = "Model"
	model.Domains = DomainMap{Map: map[string]*Domain{}}

	// success
	return nil
}

//------------------------------------------------------------------------------

// Show displays the model information on the console as json
func (model *Model) Show() (string, error) {
	return util.ConvertToYAML(model)
}

//------------------------------------------------------------------------------

// Save writes the model as json data to a file "model.json"
func (model *Model) Save(filename string) error {
	return util.SaveYAML(filename, model)
}

//------------------------------------------------------------------------------

// Load reads the model from a file "model.json"
func (model *Model) Load(filename string) error {
	return util.LoadYAML(filename, model)
}

//------------------------------------------------------------------------------

// ListDomains lists all domains of a model
func (model *Model) ListDomains() ([]string, error) {
	domains := []string{}

	for domain := range model.Domains.Map {
		domains = append(domains, domain)
	}

	// success
	return domains, nil
}

//------------------------------------------------------------------------------

// GetDomain get a domain by name
func (model *Model) GetDomain(name string) (*Domain, error) {
	// determine domain
	domain, ok := model.Domains.Map[name]

	if !ok {
		return nil, errors.New("domain not found")
	}

	// success
	return domain, nil
}

//------------------------------------------------------------------------------

// AddDomain add a domain to the model
func (model *Model) AddDomain(domain *Domain) error {
	// determine domain
	_, ok := model.Domains.Map[domain.Name]

	if ok {
		return errors.New("domain already exists")
	}

	model.Domains.Map[domain.Name] = domain

	// success
	return nil
}

//------------------------------------------------------------------------------

// DeleteDomain deletes a domain
func (model *Model) DeleteDomain(name string) error {
	// determine domain
	_, ok := model.Domains.Map[name]

	if !ok {
		return errors.New("domain not found")
	}

	// remove domain
	delete(model.Domains.Map, name)

	// success
	return nil
}

//------------------------------------------------------------------------------

// GetDomain get a domain by name
func GetDomain(domainName string) (*Domain, error) {
	domain, err := GetModel().GetDomain(domainName)
	if err != nil {
		return nil, errors.New("domain not found")
	}

	// success
	return domain, nil
}

//------------------------------------------------------------------------------

// GetComponent get an component by name
func GetComponent(domainName string, componentName string) (*Component, error) {
	domain, err := GetModel().GetDomain(domainName)
	if err != nil {
		return nil, errors.New("domain not found")
	}

	component, err := domain.GetComponent(componentName)
	if err != nil {
		return nil, errors.New("component not found")
	}

	// success
	return component, nil
}

//------------------------------------------------------------------------------

// GetArchitecture get an architecture by name
func GetArchitecture(domainName string, architectureName string) (*Architecture, error) {
	domain, err := GetModel().GetDomain(domainName)
	if err != nil {
		return nil, errors.New("domain not found")
	}

	architecture, err := domain.GetArchitecture(architectureName)
	if err != nil {
		return nil, errors.New("architecture not found")
	}

	// success
	return architecture, nil
}

//------------------------------------------------------------------------------

// GetSolution get a solution by name
func GetSolution(domainName string, solutionName string) (*Solution, error) {
	domain, err := GetModel().GetDomain(domainName)
	if err != nil {
		return nil, errors.New("domain not found")
	}

	solution, err := domain.GetSolution(solutionName)
	if err != nil {
		return nil, errors.New("solution not found")
	}

	// success
	return solution, nil
}

//------------------------------------------------------------------------------

// GetElement get an element by name
func GetElement(domainName string, solutionName string, elementName string) (*Element, error) {
	domain, err := GetModel().GetDomain(domainName)
	if err != nil {
		return nil, errors.New("domain not found")
	}

	solution, err := domain.GetSolution(solutionName)
	if err != nil {
		return nil, errors.New("solution not found")
	}

	element, err := solution.GetElement(elementName)
	if err != nil {
		return nil, errors.New("element not found")
	}

	// success
	return element, nil
}

//------------------------------------------------------------------------------

// GetCluster get a cluster by name
func GetCluster(domainName string, solutionName string, elementName string, clusterName string) (*Cluster, error) {
	domain, err := GetModel().GetDomain(domainName)
	if err != nil {
		return nil, errors.New("domain not found")
	}

	solution, err := domain.GetSolution(solutionName)
	if err != nil {
		return nil, errors.New("solution not found")
	}

	element, err := solution.GetElement(elementName)
	if err != nil {
		return nil, errors.New("element not found")
	}

	cluster, err := element.GetCluster(clusterName)
	if err != nil {
		return nil, errors.New("cluster not found")
	}

	// success
	return cluster, nil
}

//------------------------------------------------------------------------------

// GetComponent2 get the component related to a cluster
func GetComponent2(domainName string, solutionName string, elementName string, clusterName string) (*Component, error) {
	domain, err := GetModel().GetDomain(domainName)
	if err != nil {
		return nil, errors.New("domain not found")
	}

	solution, err := domain.GetSolution(solutionName)
	if err != nil {
		return nil, errors.New("solution not found")
	}

	element, err := solution.GetElement(elementName)
	if err != nil {
		return nil, errors.New("element not found")
	}

	cluster, err := element.GetCluster(clusterName)
	if err != nil {
		return nil, errors.New("cluster not found")
	}

	component, err := domain.GetComponent(element.Element + " - " + cluster.Version)
	if err != nil {
		return nil, errors.New("component not found")
	}

	// success
	return component, nil
}

//------------------------------------------------------------------------------

// GetInstance get an instance by name
func GetInstance(domainName string, solutionName string, elementName string, clusterName string, instanceName string) (*Instance, error) {
	domain, err := GetModel().GetDomain(domainName)
	if err != nil {
		return nil, errors.New("domain not found")
	}

	solution, err := domain.GetSolution(solutionName)
	if err != nil {
		return nil, errors.New("solution not found")
	}

	element, err := solution.GetElement(elementName)
	if err != nil {
		return nil, errors.New("element not found")
	}

	cluster, err := element.GetCluster(clusterName)
	if err != nil {
		return nil, errors.New("cluster not found")
	}

	instance, err := cluster.GetInstance(instanceName)
	if err != nil {
		return nil, errors.New("instance not found")
	}

	// success
	return instance, nil
}

//------------------------------------------------------------------------------

// GetRelationship get a relationship by name
func GetRelationship(domainName string, solutionName string, elementName string, clusterName string, relationshipName string) (*Relationship, error) {
	domain, err := GetModel().GetDomain(domainName)
	if err != nil {
		return nil, errors.New("domain not found")
	}

	solution, err := domain.GetSolution(solutionName)
	if err != nil {
		return nil, errors.New("solution not found")
	}

	element, err := solution.GetElement(elementName)
	if err != nil {
		return nil, errors.New("element not found")
	}

	cluster, err := element.GetCluster(clusterName)
	if err != nil {
		return nil, errors.New("cluster not found")
	}

	relationship, err := cluster.GetRelationship(relationshipName)
	if err != nil {
		return nil, errors.New("relationship not found")
	}

	// success
	return relationship, nil
}

//------------------------------------------------------------------------------
