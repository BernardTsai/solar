package main

import (
	"context"
	"net"

	"google.golang.org/grpc"

	"tsai.eu/solar/util"
	"tsai.eu/solar/model"
	"tsai.eu/solar/controller"
)

//------------------------------------------------------------------------------

// ControllerServer is an implementation of the controller interface
type ControllerServer struct {
}

//------------------------------------------------------------------------------

// Create instantiates a component
func (c *ControllerServer) Create(ctx context.Context, in *controller.SetupMessage) (*controller.StatusMessage, error) {
	util.LogInfo("create", "gRPC", in.GetSetup())
	return c.Status(ctx, in)
}

//------------------------------------------------------------------------------

// Destroy removes an instance
func (c *ControllerServer) Destroy(ctx context.Context, in *controller.SetupMessage) (*controller.StatusMessage, error) {
	util.LogInfo("destroy", "gRPC", in.GetSetup())
	return c.Status(ctx, in)
}

//------------------------------------------------------------------------------

// Start activates an instance
func (c *ControllerServer) Start(ctx context.Context, in *controller.SetupMessage) (*controller.StatusMessage, error) {
	util.LogInfo("start", "gRPC", in.GetSetup())
	return c.Status(ctx, in)
}

//------------------------------------------------------------------------------

// Stop activates an instance
func (c *ControllerServer) Stop(ctx context.Context, in *controller.SetupMessage) (*controller.StatusMessage, error) {
	util.LogInfo("stop", "gRPC", in.GetSetup())
	return c.Status(ctx, in)
}

//------------------------------------------------------------------------------

// Reset cleans up an instance
func (c *ControllerServer) Reset(ctx context.Context, in *controller.SetupMessage) (*controller.StatusMessage, error) {
	util.LogInfo("reset", "gRPC", in.GetSetup())
	return c.Status(ctx, in)
}

//------------------------------------------------------------------------------

// Configure reconfigures an instance
func (c *ControllerServer) Configure(ctx context.Context, in *controller.SetupMessage) (*controller.StatusMessage, error) {
	util.LogInfo("configure", "gRPC", in.GetSetup())
	return c.Status(ctx, in)
}

//------------------------------------------------------------------------------

// Reconfigure reconfigures an instance
func (c *ControllerServer) Reconfigure(ctx context.Context, in *controller.SetupMessage) (*controller.StatusMessage, error) {
	util.LogInfo("reconfigure", "gRPC", in.GetSetup())
	return c.Status(ctx, in)
}

//------------------------------------------------------------------------------

// Status provides the status of an instance
func (c *ControllerServer) Status(ctx context.Context, in *controller.SetupMessage) (*controller.StatusMessage, error) {
	setup := model.Setup{}

	// convert message into setup object
	yaml := in.Setup
  err := util.ConvertFromYAML(yaml, &setup)
	if err != nil {
		util.LogError("status", "gRPC", "unable to convert setup: " + err.Error())
		return nil, err
	}

	// get setups
	elementSetup     := setup.Elements[setup.Element]
	clusterSetup     := elementSetup.Clusters[setup.Cluster]
	instanceSetup    := clusterSetup.Instances[setup.Instance]

	// construct status
	status := &model.Status{
		Domain:           setup.Domain,
		Solution:         setup.Solution,
		Version:          setup.Version,
		Element:          setup.Element,
		ElementEndpoint:  "",
		Cluster:          setup.Cluster,
		ClusterEndpoint:  "",
		ClusterState:     clusterSetup.Target,
	  Instance:         setup.Instance,
		InstanceEndpoint: "",
		InstanceState:    instanceSetup.Target,
	}

	// convert to yaml
	out := controller.StatusMessage{}
	out.Status, err = util.ConvertToYAML(status)
	if err != nil {
		util.LogError("status", "gRPC", "unable to convert status: " + err.Error())
		return nil, err
	}

	// return results
	return &out, nil
}

//------------------------------------------------------------------------------

func main() {
	// be verbose
	util.LogLevel("info")

	// open TCP port 10000
	lis, err := net.Listen("tcp", ":10000")
	if err != nil {
		util.LogFatal("main", "gRPC", "failed to listen:" + err.Error())
	}

	// create a gRPC server
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	// register controller and start listening
	controller.RegisterControllerServer(grpcServer, &ControllerServer{})
	grpcServer.Serve(lis)
}

//------------------------------------------------------------------------------
