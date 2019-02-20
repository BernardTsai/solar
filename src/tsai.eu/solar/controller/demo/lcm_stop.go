package demo

import (
	"errors"

	"tsai.eu/solar/model"
)

//------------------------------------------------------------------------------

// Stop deactivates an instance
func (c Controller) Stop(setup *model.Setup) (status *model.Status, err error) {
	// get setups
	elementSetup  := setup.Elements[setup.Element]
	clusterSetup  := elementSetup.Clusters[setup.Cluster]
	parentSetup   := clusterSetup.Relationships["parent"]
	instanceSetup := clusterSetup.Instances[setup.Instance]

	// get paths
	parentPath := c.Root
	if parentSetup != nil {
		parentEndpoint, _ := DecodeEndpoint(parentSetup.Endpoint)

		parentPath = parentPath + parentEndpoint.Path + "/"
	}

	elementPath := parentPath + "/" + elementSetup.Element
	clusterPath := elementPath + "/." + clusterSetup.Cluster
	instancePath := clusterPath + "/" + instanceSetup.Instance

	// load instance information file
	information, _ := LoadInformation(instancePath)

	// set state to active
	information.State = "inactive"

	// save the information
	err = SaveInformation(instancePath, information)
	if err != nil {
		return nil, errors.New("unable to update instance file")
	}

	// success
	return c.Status(setup)
}

//------------------------------------------------------------------------------
