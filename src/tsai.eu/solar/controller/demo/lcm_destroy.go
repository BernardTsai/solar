package demo

import (
	"os"

	"tsai.eu/solar/model"
)

//------------------------------------------------------------------------------

// Destroy removes instance
func (c Controller) Destroy(setup *model.Setup) (status *model.Status, err error) {
	// get setups
	elementSetup     := setup.Elements[setup.Element]
	clusterSetup     := elementSetup.Clusters[setup.Cluster]
	parentSetup      := clusterSetup.Relationships["parent"]
	instanceSetup    := clusterSetup.Instances[setup.Instance]

	// get paths
	parentPath := c.Root
	if parentSetup != nil {
		parentEndpoint, _ := DecodeEndpoint(parentSetup.Endpoint)

		parentPath = parentPath + parentEndpoint.Path + "/"
	}

	elementPath := parentPath + "/" + elementSetup.Element
	clusterPath := elementPath + ".V" + clusterSetup.Cluster
	instancePath := clusterPath + "/" + instanceSetup.Instance

	// delete instance information file
	os.Remove(instancePath)

	// delete cluster directory if empty
	os.Remove(clusterPath)

	// delete element directory if empty
	os.Remove(elementPath)

	// success
	return c.Status(setup)
}

//------------------------------------------------------------------------------
