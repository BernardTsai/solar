package demo

import (
	"os"
	"path"

	"tsai.eu/solar/model"
)

//------------------------------------------------------------------------------

// Status provides the status of an instance
func (c Controller) Status(setup *model.Setup) (status *model.Status, err error) {
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

	// construct status
	status = &model.Status{
		Domain:           setup.Domain,
		Solution:         setup.Solution,
		Version:          setup.Version,
		Element:          setup.Element,
		Cluster:          setup.Cluster,
		ClusterEndpoint:  "",
		ClusterState:     model.UndefinedState,
	  Instance:         setup.Instance,
		InstanceEndpoint: "",
		InstanceState:    model.UndefinedState,
	}

	// check if parent directory exists
	if _, err = os.Stat(parentPath); os.IsNotExist(err) {
		return status, nil
	}

	// check if element directory exists
	if _, err = os.Stat(elementPath); os.IsNotExist(err) {
		return status, nil
	}

	// check if cluster directory exists
	if _, err = os.Stat(clusterPath); os.IsNotExist(err) {
		return status, nil
	}

	status.ClusterState,  status.ClusterEndpoint  = getClusterState(clusterPath)
	status.InstanceState, status.InstanceEndpoint = getInstanceState(instancePath)

	// return results
	return status, nil
}

//------------------------------------------------------------------------------

// getClusterState determines the state of a cluster
func getClusterState(clusterPath string) (state string, endpoint string) {
	state       = model.InitialState
	endpoint, _ = EncodeEndpoint(NewEndpoint(clusterPath))

	// open directory
	file, err := os.Open(clusterPath)
  if err != nil {
		state = model.FailureState
		return state, endpoint
  }
  defer file.Close()

	// read all files
  list,_ := file.Readdirnames(0) // 0 to read all files and folders
  for _, name := range list {
		instancePath := clusterPath + "/" + name

		instanceState, _ := getInstanceState(instancePath)

		switch instanceState {
		case model.FailureState:
			// demote to failure state if cluster is in initial state
			if state == model.InitialState  {
				state = model.FailureState
			}
		case model.InitialState:
			// does not change the state of the cluster
		case model.InactiveState:
			// promote to inactive state if needed
			if state == model.InitialState || state == model.FailureState {
				state = model.InactiveState
			}
		case model.ActiveState:
			// cluster is in active state
			state = model.ActiveState
		}
  }

	return state, endpoint
}

//------------------------------------------------------------------------------

// getInstanceState determines the state of an instance
func getInstanceState(instancePath string) (state string, endpoint string) {
	state       = model.InitialState
	endpoint, _ = EncodeEndpoint(NewEndpoint(instancePath))

	// check if instance file exists
	if _, err := os.Stat(instancePath); os.IsExist(err) {

		// load information
		if information, err := LoadInformation(instancePath); err != nil {
			// information was probably in the wrong format
			state = model.FailureState
		} else {
			// check state
			switch information.State {
			case "inactive":
				state = model.InactiveState
			case "active":
				state = model.ActiveState
			default:
				state = model.FailureState
			}
		}
	}

	// return the results
	return state, endpoint
}

//------------------------------------------------------------------------------
