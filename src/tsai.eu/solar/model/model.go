package model

import (
	"sync"
	"errors"

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
//   - GetDomains
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
//   - model.Load2
//   - model.Save
//
//   - model.ListDomains
//   - model.GetDomain
//   - model.AddDomain
//   - model.DeleteDomain
//------------------------------------------------------------------------------

// Model describes all managed artefacts within a model.
type Model struct {
	Schema   string             `yaml:"Schema"`             // schema of the model
	Name     string             `yaml:"Name"`               // name of the model
	Domains  map[string]*Domain `yaml:"Domains"`            // map of domains
	DomainsX sync.RWMutex       `yaml:"DomainsX,omitempty"` // mutex for domains
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
	model.Schema   = "BT V1.0.0"
	model.Name     = "Model"
	model.Domains  = map[string]*Domain{}
	model.DomainsX = sync.RWMutex{}

	// success
	return nil
}

//------------------------------------------------------------------------------

// Show displays the model information on the console as json
func (model *Model) Show() (string, error) {
	return util.ConvertToYAML(model)
}

//------------------------------------------------------------------------------

// Save writes the model as json data to a file
func (model *Model) Save(filename string) error {
	return util.SaveYAML(filename, model)
}

//------------------------------------------------------------------------------

// Load reads the model from a file
func (model *Model) Load(filename string) error {
	return util.LoadYAML(filename, model)
}

//------------------------------------------------------------------------------

// Load2 imports a yaml model
func (model *Model) Load2(yaml string) error {
	return util.ConvertFromYAML(yaml, model)
}

//------------------------------------------------------------------------------

// ListDomains lists all domains of a model
func (model *Model) ListDomains() ([]string, error) {
	domains := []string{}

	model.DomainsX.RLock()
	for domain := range model.Domains {
		domains = append(domains, domain)
	}
 	model.DomainsX.RUnlock()

	// success
	return domains, nil
}

//------------------------------------------------------------------------------

// GetDomain get a domain by name
func (model *Model) GetDomain(name string) (*Domain, error) {
	// determine domain
 	model.DomainsX.RLock()
	domain, ok := model.Domains[name]
 	model.DomainsX.RUnlock()

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
	model.DomainsX.RLock()
	_, ok := model.Domains[domain.Name]
	model.DomainsX.RUnlock()

	if ok {
		return errors.New("domain already exists")
	}

	model.DomainsX.Lock()
	model.Domains[domain.Name] = domain
	model.DomainsX.Unlock()

	// success
	return nil
}

//------------------------------------------------------------------------------

// DeleteDomain deletes a domain
func (model *Model) DeleteDomain(name string) error {
	// determine domain
	model.DomainsX.RLock()
	_, ok := model.Domains[name]
	model.DomainsX.RUnlock()

	if !ok {
		return errors.New("domain not found")
	}

	// remove domain
	model.DomainsX.Lock()
	delete(model.Domains, name)
	model.DomainsX.Unlock()

	// success
	return nil
}

//------------------------------------------------------------------------------

// GetDomain get a domain by name
func GetDomain(domainName string) (*Domain, error) {
	return GetModel().GetDomain(domainName)
}

//------------------------------------------------------------------------------

// GetDomains lists all domains of a model
func GetDomains() ([]string, error) {
	return GetModel().ListDomains()
}

//------------------------------------------------------------------------------

// GetComponent get an component by name
func GetComponent(domainName string, componentName string, version string) (*Component, error) {
	domain, err := GetModel().GetDomain(domainName)
	if err != nil {
		return nil, errors.New("domain not found")
	}

	component, err := domain.GetComponent(componentName, version)
	if err != nil {
		return nil, errors.New("component not found")
	}

	// success
	return component, nil
}

//------------------------------------------------------------------------------

// GetArchitecture get an architecture by name
func GetArchitecture(domainName string, architectureName string, version string) (*Architecture, error) {
	domain, err := GetModel().GetDomain(domainName)
	if err != nil {
		return nil, errors.New("domain not found")
	}

	architecture, err := domain.GetArchitecture(architectureName, version)
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

	component, err := domain.GetComponent(element.Component, cluster.Version)
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

// GetController get a controller by type and version
func GetController(domainName string, controllerType string, controllerVersion string) (*Controller, error) {
	domain, err := GetModel().GetDomain(domainName)
	if err != nil {
		return nil, errors.New("domain not found")
	}

	controller, err := domain.GetController(controllerType, controllerVersion)
	if err != nil {
		return nil, errors.New("controller not found")
	}

	// success
	return controller, nil
}

//------------------------------------------------------------------------------

// GetTask get a task by uuid
func GetTask(domainName string, uuid string) (*Task, error) {
	domain, err := GetModel().GetDomain(domainName)
	if err != nil {
		return nil, errors.New("domain not found")
	}

	task, err := domain.GetTask(uuid)
	if err != nil {
		return nil, errors.New("task not found")
	}

	// success
	return task, nil
}

//------------------------------------------------------------------------------

// GetEvent get an event by uuid
func GetEvent(domainName string, uuid string) (*Event, error) {
	domain, err := GetModel().GetDomain(domainName)
	if err != nil {
		return nil, errors.New("domain not found")
	}

	event, err := domain.GetEvent(uuid)
	if err != nil {
		return nil, errors.New("event not found")
	}

	// success
	return event, nil
}

//------------------------------------------------------------------------------
