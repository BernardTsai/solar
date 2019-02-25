package dummy

import (
	"tsai.eu/solar/model"
)

//------------------------------------------------------------------------------

// Status provides the status of an instance
func (c Controller) Status(setup *model.Setup) (status *model.Status, err error) {
	// get setups
	elementSetup     := setup.Elements[setup.Element]
	clusterSetup     := elementSetup.Clusters[setup.Cluster]
	instanceSetup    := clusterSetup.Instances[setup.Instance]

	// construct status 
	status = &model.Status{
		Domain:           setup.Domain,
		Solution:         setup.Solution,
		Version:          setup.Version,
		Element:          setup.Element,
		ElementEndpoint:  "",
		Cluster:          setup.Cluster,
		ClusterEndpoint:  "",
		ClusterState:     clusterSetup.Target,
	  Instance:         setup.Instance,
		InstanceEndpoint: instanceSetup.Target,
		InstanceState:    model.UndefinedState,
	}

	// return results
	return status, nil
}

//------------------------------------------------------------------------------