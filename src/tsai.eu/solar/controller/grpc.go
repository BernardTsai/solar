package controller

import (
	"context"

	"google.golang.org/grpc"

	"tsai.eu/solar/model"
	"tsai.eu/solar/util"
)

//------------------------------------------------------------------------------

// GRPCController is an gRPC based implementation of the Controller interface
type GRPCController struct {
	Type         string           // type of components which the controller supports
	Version      string           // version of the components which the controller supports
	Address      string           // address to which the controller server listens
	Connection *grpc.ClientConn  // connection to the controller server
	Client       ControllerClient // client for accessing the controller server
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
	util.LogInfo("main", "CTRL", connection.GetState().String())
	c.Connection = connection

	// create the client
	c.Client = NewControllerClient(c.Connection)

	// success
	return &c, nil
}

//------------------------------------------------------------------------------

// Status determines the status of an instance
func (c *GRPCController)Status(setup *model.Setup) (*model.Status, error) {
	setupMessage := SetupMessage{}

	// convert setup to a setup message
	msg, err := util.ConvertToYAML(setup)
	if err != nil {
		util.LogError("status", "CTRL", "unable to convert setup\n" + err.Error())
		return nil, err
	}
	setupMessage.Setup = msg

	// invoke the remote controller
	statusMessage, err := c.Client.Create(context.Background(), &setupMessage )
	if err != nil {
		util.LogError("status", "CTRL", "unable to invoke controller\n" + err.Error())
		return nil, err
	}

	// convert status message to status
	status := model.Status{}

	err = util.ConvertFromYAML(statusMessage.Status, &status)
	if err != nil {
		util.LogError("status", "CTRL", "unable to convert status\n" + err.Error())
		return nil, err
	}

	// success
	return &status, nil
}

//------------------------------------------------------------------------------

// Create instantiates an instance
func (c *GRPCController)Create(setup *model.Setup) (*model.Status, error) {
	setupMessage := SetupMessage{}

	// convert setup to a setup message
	msg, err := util.ConvertToYAML(setup)
	if err != nil {
		util.LogError("create", "CTRL", "unable to convert setup\n" + err.Error())
		return nil, err
	}
	setupMessage.Setup = msg

	// invoke the remote controller
	statusMessage, err := c.Client.Create(context.Background(), &setupMessage )
	if err != nil {
		util.LogError("create", "CTRL", "unable to invoke controller\n" + err.Error())
		return nil, err
	}

	// convert status message to status
	status := model.Status{}

	err = util.ConvertFromYAML(statusMessage.Status, &status)
	if err != nil {
		util.LogError("create", "CTRL", "unable to convert status\n" + err.Error())
		return nil, err
	}

	// success
	return &status, nil
}

//------------------------------------------------------------------------------

// Destroy removes an instance
func (c *GRPCController)Destroy(setup *model.Setup) (*model.Status, error) {
	setupMessage := SetupMessage{}

	// convert setup to a setup message
	msg, err := util.ConvertToYAML(setup)
	if err != nil {
		util.LogError("destroy", "CTRL", "unable to convert setup\n" + err.Error())
		return nil, err
	}
	setupMessage.Setup = msg

	// invoke the remote controller
	statusMessage, err := c.Client.Create(context.Background(), &setupMessage )
	if err != nil {
		util.LogError("destroy", "CTRL", "unable to invoke controller\n" + err.Error())
		return nil, err
	}

	// convert status message to status
	status := model.Status{}

	err = util.ConvertFromYAML(statusMessage.Status, &status)
	if err != nil {
		util.LogError("destroy", "CTRL", "unable to convert status\n" + err.Error())
		return nil, err
	}

	// success
	return &status, nil
}

//------------------------------------------------------------------------------

// Configure configures an instance
func (c *GRPCController)Configure(setup *model.Setup) (*model.Status, error) {
	setupMessage := SetupMessage{}

	// convert setup to a setup message
	msg, err := util.ConvertToYAML(setup)
	if err != nil {
		util.LogError("configure", "CTRL", "unable to convert setup\n" + err.Error())
		return nil, err
	}
	setupMessage.Setup = msg

	// invoke the remote controller
	statusMessage, err := c.Client.Create(context.Background(), &setupMessage )
	if err != nil {
		util.LogError("configure", "CTRL", "unable to invoke controller\n" + err.Error())
		return nil, err
	}

	// convert status message to status
	status := model.Status{}

	err = util.ConvertFromYAML(statusMessage.Status, &status)
	if err != nil {
		util.LogError("configure", "CTRL", "unable to convert status\n" + err.Error())
		return nil, err
	}

	// success
	return &status, nil
}

//------------------------------------------------------------------------------

// Reconfigure reconfigures an instance
func (c *GRPCController)Reconfigure(setup *model.Setup) (*model.Status, error) {
	setupMessage := SetupMessage{}

	// convert setup to a setup message
	msg, err := util.ConvertToYAML(setup)
	if err != nil {
		util.LogError("reconfigure", "CTRL", "unable to convert setup\n" + err.Error())
		return nil, err
	}
	setupMessage.Setup = msg

	// invoke the remote controller
	statusMessage, err := c.Client.Create(context.Background(), &setupMessage )
	if err != nil {
		util.LogError("reconfigure", "CTRL", "unable to invoke controller\n" + err.Error())
		return nil, err
	}

	// convert status message to status
	status := model.Status{}

	err = util.ConvertFromYAML(statusMessage.Status, &status)
	if err != nil {
		util.LogError("reconfigure", "CTRL", "unable to convert status\n" + err.Error())
		return nil, err
	}

	// success
	return &status, nil
}

//------------------------------------------------------------------------------

// Start activates an instance
func (c *GRPCController)Start(setup *model.Setup) (*model.Status, error) {
	setupMessage := SetupMessage{}

	// convert setup to a setup message
	msg, err := util.ConvertToYAML(setup)
	if err != nil {
		util.LogError("start", "CTRL", "unable to convert setup\n" + err.Error())
		return nil, err
	}
	setupMessage.Setup = msg

	// invoke the remote controller
	statusMessage, err := c.Client.Create(context.Background(), &setupMessage )
	if err != nil {
		util.LogError("start", "CTRL", "unable to invoke controller\n" + err.Error())
		return nil, err
	}

	// convert status message to status
	status := model.Status{}

	err = util.ConvertFromYAML(statusMessage.Status, &status)
	if err != nil {
		util.LogError("start", "CTRL", "unable to convert status\n" + err.Error())
		return nil, err
	}

	// success
	return &status, nil
}

//------------------------------------------------------------------------------

// Stop deactivates an instance
func (c *GRPCController)Stop(setup *model.Setup) (*model.Status, error) {
	setupMessage := SetupMessage{}

	// convert setup to a setup message
	msg, err := util.ConvertToYAML(setup)
	if err != nil {
		util.LogError("stop", "CTRL", "unable to convert setup\n" + err.Error())
		return nil, err
	}
	setupMessage.Setup = msg

	// invoke the remote controller
	statusMessage, err := c.Client.Create(context.Background(), &setupMessage )
	if err != nil {
		util.LogError("stop", "CTRL", "unable to invoke controller\n" + err.Error())
		return nil, err
	}

	// convert status message to status
	status := model.Status{}

	err = util.ConvertFromYAML(statusMessage.Status, &status)
	if err != nil {
		util.LogError("stop", "CTRL", "unable to convert status\n" + err.Error())
		return nil, err
	}

	// success
	return &status, nil
}

//------------------------------------------------------------------------------

// Reset cleans up a failed instance
func (c *GRPCController)Reset(setup *model.Setup) (*model.Status, error) {
	setupMessage := SetupMessage{}

	// convert setup to a setup message
	msg, err := util.ConvertToYAML(setup)
	if err != nil {
		util.LogError("reset", "CTRL", "unable to convert setup\n" + err.Error())
		return nil, err
	}
	setupMessage.Setup = msg

	// invoke the remote controller
	statusMessage, err := c.Client.Create(context.Background(), &setupMessage )
	if err != nil {
		util.LogError("reset", "CTRL", "unable to invoke controller\n" + err.Error())
		return nil, err
	}

	// convert status message to status
	status := model.Status{}

	err = util.ConvertFromYAML(statusMessage.Status, &status)
	if err != nil {
		util.LogError("reset", "CTRL", "unable to convert status\n" + err.Error())
		return nil, err
	}

	// success
	return &status, nil
}

//------------------------------------------------------------------------------
