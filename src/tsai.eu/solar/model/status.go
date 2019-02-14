package model

import "fmt"

//------------------------------------------------------------------------------

// ComponentStatus object received from controller.
type ComponentStatus struct {
	Domain            string `yaml:"Domain"`            // name of domain
	Solution          string `yaml:"Solution"`          // name of solution
	Element           string `yaml:"Element"`           // name of element
	Version           string `yaml:"Version"`           // version of element
	Instance          string `yaml:"Instance"`          // name of instance
	ComponentEndpoint string `yaml:"ComponentEndpoint"` // endpoint of component
	VersionEndpoint   string `yaml:"VersionEndpoint"`   // endpoint of component version
	InstanceEndpoint  string `yaml:"InstanceEndpoint"`  // endpoint of instance
	InstanceState     string `yaml:"InstanceState"`     // state of instance
	Changed           bool   `yaml:"Changed"`           // indicator if a change occured
}

//------------------------------------------------------------------------------

// DeriveComponentStatus derives a ComponentStatus from a ComponentConfiguration struct
func DeriveComponentStatus(configuration *ComponentConfiguration) (status *ComponentStatus) {
	instance, found := configuration.Instances[configuration.Instance]
	if !found {
		fmt.Println("Hello")
		fmt.Println(configuration.Instance)
		fmt.Println("Hello")
		panic("No instance found")
	}

	status = &ComponentStatus{
		Domain:            configuration.Domain,
		Solution:          configuration.Solution,
		Element:           configuration.Element,
		Version:           instance.Version,
		Instance:          configuration.Instance,
		ComponentEndpoint: configuration.Endpoint,
		VersionEndpoint:   configuration.Endpoints[instance.Version],
		InstanceEndpoint:  instance.Endpoint,
		InstanceState:     instance.State,
		Changed:           false,
	}
	return
}

//------------------------------------------------------------------------------

// SetStatus saves the status received from a controller.
func SetStatus(status ComponentStatus) (err error) {
	if status.Changed {
		// TODO: proper error handling
		domain, _ := GetModel().GetDomain(status.Domain)
		solution, _ := domain.GetSolution(status.Solution)
		element, _ := solution.GetElement(status.Element)
		cluster, _ := element.GetCluster(status.Version)
		instance, _ := cluster.GetInstance(status.Instance)

		// update component
		cluster.Endpoint = status.ComponentEndpoint
		cluster.AddEndpoint(cluster.Version, status.VersionEndpoint)

		// update instance
		instance.Endpoint = status.InstanceEndpoint
		instance.State = status.InstanceState

	}
	// success
	return nil
}
