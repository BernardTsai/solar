package controller

import (
	"context"
	"errors"

	"google.golang.org/grpc"

	"tsai.eu/solar/model"
	"tsai.eu/solar/util"
)

//------------------------------------------------------------------------------

// GRPCController is an gRPC based implementation of the Controller interface
type GRPCController struct {
	Type        string           // type of components which the controller supports
	Version     string           // version of the components which the controller supports
	Address     string           // address to which the controller server listens
	Connection *grpc.ClientConn  // connection to the controller server
	Client      ControllerClient // client for accessing the controller server
}

//------------------------------------------------------------------------------

// newGRPCController creates a gRPC based controller
func newGRPCController(Type string, Address string) (*GRPCController, error) {
	c := GRPCController{
		Type:       Type,
		Version:    "1.0.0",
		Address:    Address,
		Connection: nil,
		Client:     nil,
	}

	// define the communication options
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	// create the connection
	connection, err := grpc.Dial(Address, opts...)
	if err != nil {
		util.LogError("main", "CTRL", "unable to connect to the controller: " + Type + "\n" + err.Error())
		return nil, err
	}
	c.Connection = connection

	// create the client
	c.Client = NewControllerClient(c.Connection)

	// check availability of client
	if !c.Check() {
		util.LogError("main", "CTRL", "controller: " + Type + " unavailable")
		return nil, errors.New("controller: " + Type + " unavailable")
	}

	// success
	util.LogInfo("main", "CTRL", "controller: " + Type + " available")
	return &c, nil
}

//------------------------------------------------------------------------------

// Status determines the status of an instance
func (c *GRPCController)Status(setup *model.Setup) (*model.Status, error) {
	// convert setup to a setup message
	setupMessage := convertSetup(setup)

	// invoke the remote controller
	statusMessage, err := c.Client.Create(context.Background(), setupMessage )
	if err != nil {
		util.LogError("status", "CTRL", "unable to invoke controller\n" + err.Error())
		return nil, err
	}

	// convert status message to status
	status := convertStatus(statusMessage)

	// success
	return status, nil
}

//------------------------------------------------------------------------------

// Check checks availability of controller
func (c *GRPCController) Check() bool {
	req := VoidMessage{
		Version: "V1.0.0",
	}

	// invoke the remote controller
	_, err := c.Client.Check(context.Background(), &req)
	if err != nil {
		util.LogError("check", "CTRL", "unable to connect to controller " + c.Type + "\n" + err.Error())
		return false
	}

	// success
	return true
}

//------------------------------------------------------------------------------

// Create instantiates an instance
func (c *GRPCController)Create(setup *model.Setup) (*model.Status, error) {
	// convert setup to a setup message
	setupMessage := convertSetup(setup)

	// invoke the remote controller
	statusMessage, err := c.Client.Create(context.Background(), setupMessage )
	if err != nil {
		util.LogError("create", "CTRL", "unable to invoke controller\n" + err.Error())
		return nil, err
	}

	// convert status message to status
	status := convertStatus(statusMessage)

	// success
	return status, nil
}

//------------------------------------------------------------------------------

// Destroy removes an instance
func (c *GRPCController)Destroy(setup *model.Setup) (*model.Status, error) {
	// convert setup to a setup message
	setupMessage := convertSetup(setup)

	// invoke the remote controller
	statusMessage, err := c.Client.Create(context.Background(), setupMessage )
	if err != nil {
		util.LogError("destroy", "CTRL", "unable to invoke controller\n" + err.Error())
		return nil, err
	}

	// convert status message to status
	status := convertStatus(statusMessage)

	// success
	return status, nil
}

//------------------------------------------------------------------------------

// Configure configures an instance
func (c *GRPCController)Configure(setup *model.Setup) (*model.Status, error) {
	// convert setup to a setup message
	setupMessage := convertSetup(setup)

	// invoke the remote controller
	statusMessage, err := c.Client.Create(context.Background(), setupMessage )
	if err != nil {
		util.LogError("configure", "CTRL", "unable to invoke controller\n" + err.Error())
		return nil, err
	}

	// convert status message to status
	status := convertStatus(statusMessage)

	// success
	return status, nil
}

//------------------------------------------------------------------------------

// Reconfigure reconfigures an instance
func (c *GRPCController)Reconfigure(setup *model.Setup) (*model.Status, error) {
	// convert setup to a setup message
	setupMessage := convertSetup(setup)

	// invoke the remote controller
	statusMessage, err := c.Client.Create(context.Background(), setupMessage )
	if err != nil {
		util.LogError("reconfigure", "CTRL", "unable to invoke controller\n" + err.Error())
		return nil, err
	}

	// convert status message to status
	status := convertStatus(statusMessage)

	// success
	return status, nil
}

//------------------------------------------------------------------------------

// Start activates an instance
func (c *GRPCController)Start(setup *model.Setup) (*model.Status, error) {
	// convert setup to a setup message
	setupMessage := convertSetup(setup)

	// invoke the remote controller
	statusMessage, err := c.Client.Start(context.Background(), setupMessage )
	if err != nil {
		util.LogError("start", "CTRL", "unable to invoke controller\n" + err.Error())
		return nil, err
	}

	// convert status message to status
	status := convertStatus(statusMessage)

	// success
	return status, nil
}

//------------------------------------------------------------------------------

// Stop deactivates an instance
func (c *GRPCController)Stop(setup *model.Setup) (*model.Status, error) {
	// convert setup to a setup message
	setupMessage := convertSetup(setup)

	// invoke the remote controller
	statusMessage, err := c.Client.Stop(context.Background(), setupMessage )
	if err != nil {
		util.LogError("stop", "CTRL", "unable to invoke controller\n" + err.Error())
		return nil, err
	}

	// convert status message to status
	status := convertStatus(statusMessage)

	// success
	return status, nil
}

//------------------------------------------------------------------------------

// Reset cleans up a failed instance
func (c *GRPCController)Reset(setup *model.Setup) (*model.Status, error) {
	// convert setup to a setup message
	setupMessage := convertSetup(setup)

	// invoke the remote controller
	statusMessage, err := c.Client.Reset(context.Background(), setupMessage )
	if err != nil {
		util.LogError("reset", "CTRL", "unable to invoke controller\n" + err.Error())
		return nil, err
	}

	// convert status message to status
	status := convertStatus(statusMessage)

	// success
	return status, nil
}

//------------------------------------------------------------------------------

// convertInstanceSetup converts an instance setup to a message
func convertInstanceSetup(setup *model.InstanceSetup) *InstanceSetupMessage {
	instanceSetupMessage := InstanceSetupMessage{
		Instance                : setup.Instance,
		Target                  : setup.Target,
		State                   : setup.State,
		BaseConfiguration       : setup.BaseConfiguration,
		DesignTimeConfiguration : setup.DesignTimeConfiguration,
		RuntimeConfiguration    : setup.RuntimeConfiguration,
		Endpoint                : setup.Endpoint,
	}

	return &instanceSetupMessage
}

//------------------------------------------------------------------------------

// convertRelationshipSetup converts a relationship setup to a message
func convertRelationshipSetup(setup *model.RelationshipSetup) *RelationshipSetupMessage {
	relationshipSetupMessage := RelationshipSetupMessage{
		Relationship            : setup.Relationship,
		Element                 : setup.Element,
		Version                 : setup.Version,
		Target                  : setup.Target,
		State                   : setup.State,
		BaseConfiguration       : setup.BaseConfiguration,
		DesignTimeConfiguration : setup.DesignTimeConfiguration,
		RuntimeConfiguration    : setup.RuntimeConfiguration,
		Endpoint                : setup.Endpoint,
	}

	return &relationshipSetupMessage
}

//------------------------------------------------------------------------------

// convertClusterSetup converts a cluster setup to a message
func convertClusterSetup(setup *model.ClusterSetup) *ClusterSetupMessage {
	clusterSetupMessage := ClusterSetupMessage{
		Cluster                 : setup.Cluster,
		Target                  : setup.Target,
		State                   : setup.State,
		Min                     : int64(setup.Min),
		Max                     : int64(setup.Max),
		Size                    : int64(setup.Size),
		BaseConfiguration       : setup.BaseConfiguration,
		DesignTimeConfiguration : setup.DesignTimeConfiguration,
		RuntimeConfiguration    : setup.RuntimeConfiguration,
		Endpoint                : setup.Endpoint,
		Relationships           : map[string]*RelationshipSetupMessage{},
		Instances               : map[string]*InstanceSetupMessage{},
	}

	for relationshipName, relationshipSetup := range setup.Relationships {
		clusterSetupMessage.Relationships[relationshipName] = convertRelationshipSetup(relationshipSetup)
	}

	for instanceName, instanceSetup := range setup.Instances {
		clusterSetupMessage.Instances[instanceName] = convertInstanceSetup(instanceSetup)
	}

	return &clusterSetupMessage
}

//------------------------------------------------------------------------------

// convertElementSetup converts an element setup to a message
func convertElementSetup(setup *model.ElementSetup) *ElementSetupMessage {
	elementSetupMessage := ElementSetupMessage{
		Element                 : setup.Element,
		Component               : setup.Component,
		Target                  : setup.Target,
		State                   : setup.State,
		DesignTimeConfiguration : setup.DesignTimeConfiguration,
		RuntimeConfiguration    : setup.RuntimeConfiguration,
		Endpoint                : setup.Endpoint,
		Clusters                : map[string]*ClusterSetupMessage{},
	}

	for clusterName, clusterSetup := range setup.Clusters {
		elementSetupMessage.Clusters[clusterName] = convertClusterSetup(clusterSetup)
	}

	return &elementSetupMessage
}

//------------------------------------------------------------------------------

// convertSetup converts a setup to a message
func convertSetup(setup *model.Setup) *SetupMessage {
	setupMessage := SetupMessage{
		Domain                  : setup.Domain,
		Solution                : setup.Solution,
		Version                 : setup.Version,
		Element                 : setup.Element,
		Cluster                 : setup.Cluster,
		Instance                : setup.Instance,
		Target                  : setup.Target,
		State                   : setup.State,
		DesignTimeConfiguration : setup.DesignTimeConfiguration,
		RuntimeConfiguration    : setup.RuntimeConfiguration,
		Elements                : map[string]*ElementSetupMessage{},
	}

	for elementName, elementSetup := range setup.Elements {
		setupMessage.Elements[elementName] = convertElementSetup(elementSetup)
	}

	return &setupMessage
}

//------------------------------------------------------------------------------

// convertStatus converts a message to a status
func convertStatus(message *StatusMessage) *model.Status {
	status := model.Status {
		Domain           : message.Domain,
	  Solution         : message.Solution,
	  Version          : message.Version,
	  Element          : message.Element,
	  ElementEndpoint  : message.ElementEndpoint,
	  Cluster          : message.Cluster,
	  ClusterEndpoint  : message.ClusterEndpoint,
	  ClusterState     : message.ClusterState,
	  Instance         : message.Instance,
	  InstanceEndpoint : message.InstanceEndpoint,
	  InstanceState    : message.InstanceState,
	}

	return &status
}

//------------------------------------------------------------------------------
