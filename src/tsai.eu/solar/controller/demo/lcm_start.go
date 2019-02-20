package demo

import (
	"errors"
	"path"

	"tsai.eu/solar/model"
)

//------------------------------------------------------------------------------

// Start activates an instance
func (c Controller) Start(setup *model.Setup) (status *model.Status, err error) {
	// get setups
	elementSetup     := setup.Elements[setup.Element]
	clusterSetup     := elementSetup.Clusters[setup.Cluster]
	parentSetup      := clusterSetup.Relationships["Parent"]
	instanceSetup    := clusterSetup.Instances[setup.Instance]

	// get paths
	parentPath := c.Root
	if parentSetup != nil {
		parentEndpoint, _ := DecodeEndpoint(parentSetup.Endpoint)

		parentPath = parentPath + parentEndpoint.Path + "/../"
		parentPath = path.Clean(parentPath)
	}

	elementPath := parentPath + "/" + elementSetup.Element
	clusterPath := elementPath + "/." + clusterSetup.Cluster
	instancePath := clusterPath + "/" + instanceSetup.Instance

	// load instance information file
	information, _ := LoadInformation(instancePath)

	// set state to active
	information.State = "active"

	// save the information
	err = SaveInformation(instancePath, information)
	if err != nil {
		return nil, errors.New("unable to update instance file")
	}

	// success
	return c.Status(setup)
}

//------------------------------------------------------------------------------
