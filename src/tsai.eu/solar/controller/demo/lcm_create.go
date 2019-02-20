package demo

import (
	"errors"
	"os"
	"path"

	"tsai.eu/solar/model"
)

//------------------------------------------------------------------------------

// Create initialises an instance
func (c Controller) Create(setup *model.Setup) (status *model.Status, err error) {
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

	// create path to element directory if it does not exist yet
	if _, err = os.Stat(elementPath); os.IsNotExist(err) {
		if err = os.Mkdir(elementPath, os.ModePerm); err != nil {
			return nil, errors.New("unable to create element directory")
		}
	}

	// create path to cluster directory if it does not exist yet
	if _, err = os.Stat(clusterPath); os.IsNotExist(err) {
		if err = os.Mkdir(clusterPath, os.ModePerm); err != nil {
			return nil, errors.New("unable to create cluster directory")
		}
	}

	// create instance information file
	configuration, _ := DecodeConfiguration(instanceSetup.DesignTimeConfiguration)

	information := Information{
		State:        "inactive",
		Path :        instancePath,
		Template:     configuration.Template,
		Refererences: map[string]string{},
	}

	for relationshipName, relationship := range clusterSetup.Relationships {
		endpoint, _ := DecodeEndpoint(relationship.Endpoint)

		information.Refererences[relationshipName] = endpoint.Path
	}

	err = SaveInformation(instancePath, &information)
	if err != nil {
		return nil, errors.New("unable to create instance file")
	}

	// success
	return c.Status(setup)
}

//------------------------------------------------------------------------------
