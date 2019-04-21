package gRPC

import (
	"context"

	"google.golang.org/grpc"

  "github.com/rs/zerolog/log"
)

//------------------------------------------------------------------------------

// DefaultController is an implementation of the controller interface
type DefaultController struct {
}

//------------------------------------------------------------------------------

// NewController creates a new controller instance
func NewController() (*grpc.Server, error) {
	// create a gRPC server
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	// register controller
	RegisterControllerServer(grpcServer, &DefaultController{})

	return grpcServer, nil
}

//------------------------------------------------------------------------------

// Check is a keep-alive and version ping
func (c *DefaultController) Check(ctx context.Context, in *VoidMessage) (*VoidMessage, error) {
	Log("info", "check", "gRPC", in.String())
	return in, nil
}

//------------------------------------------------------------------------------

// Create instantiates a component
func (c *DefaultController) Create(ctx context.Context, in *SetupMessage) (*StatusMessage, error) {
	Log("info", "create", "gRPC", in.String())

	// set target state
	setState(in, "inactive")

	return c.Status(ctx, in)
}

//------------------------------------------------------------------------------

// Destroy removes an instance
func (c *DefaultController) Destroy(ctx context.Context, in *SetupMessage) (*StatusMessage, error) {
	Log("info", "destroy", "gRPC", in.String())

	// set target state
	setState(in, "initial")

	return c.Status(ctx, in)
}

//------------------------------------------------------------------------------

// Start activates an instance
func (c *DefaultController) Start(ctx context.Context, in *SetupMessage) (*StatusMessage, error) {
	Log("info", "start", "gRPC", in.String())

	// set target state
	setState(in, "active")

	return c.Status(ctx, in)
}

//------------------------------------------------------------------------------

// Stop activates an instance
func (c *DefaultController) Stop(ctx context.Context, in *SetupMessage) (*StatusMessage, error) {
	Log("info", "stop", "gRPC", in.String())

	// set target state
	setState(in, "inactive")

	return c.Status(ctx, in)
}

//------------------------------------------------------------------------------

// Reset cleans up an instance
func (c *DefaultController) Reset(ctx context.Context, in *SetupMessage) (*StatusMessage, error) {
	Log("info", "reset", "gRPC", in.String())

	// set target state
	setState(in, "initial")

	return c.Status(ctx, in)
}

//------------------------------------------------------------------------------

// Configure reconfigures an instance
func (c *DefaultController) Configure(ctx context.Context, in *SetupMessage) (*StatusMessage, error) {
	Log("info", "configure", "gRPC", in.String())

	// set target state
	setState(in, "active")

	return c.Status(ctx, in)
}

//------------------------------------------------------------------------------

// Reconfigure reconfigures an instance
func (c *DefaultController) Reconfigure(ctx context.Context, in *SetupMessage) (*StatusMessage, error) {
	Log("info", "reconfigure", "gRPC", in.String())

	// set target state
	setState(in, "active")

	return c.Status(ctx, in)
}

//------------------------------------------------------------------------------

// Status provides the status of an instance
func (c *DefaultController) Status(ctx context.Context, in *SetupMessage) (*StatusMessage, error) {
	// get setups
	elementSetup  := in.Elements[in.Element]
	clusterSetup  := elementSetup.Clusters[in.Cluster]
	instanceSetup := clusterSetup.Instances[in.Instance]

	// construct status
	status := StatusMessage{
		Error:            "",
		Domain:           in.Domain,
		Solution:         in.Solution,
		Version:          in.Version,
		Element:          in.Element,
		ElementEndpoint:  "",
		Cluster:          in.Cluster,
		ClusterEndpoint:  "",
		ClusterState:     clusterSetup.Target,
	  Instance:         in.Instance,
		InstanceEndpoint: "",
		InstanceState:    instanceSetup.Target,
	}

	// return results
	return &status, nil
}

//------------------------------------------------------------------------------

// setState tweaks the target state of the setup message
func setState(setup *SetupMessage, state string) {
	// set target state
	elementSetup  := setup.Elements[setup.Element]
	clusterSetup  := elementSetup.Clusters[setup.Cluster]
	instanceSetup := clusterSetup.Instances[setup.Instance]

	instanceSetup.Target = state
}

//------------------------------------------------------------------------------

// Log captures log information
func Log(level string, context string, module string, info string) {
	switch level {
	case "panic":
		log.Panic().Str("Context", context).Str("Module", module).Msg(info)
	case "fatal":
		log.Fatal().Str("Context", context).Str("Module", module).Msg(info)
	case "error":
		log.Error().Str("Context", context).Str("Module", module).Msg(info)
	case "warn":
		log.Warn().Str("Context", context).Str("Module", module).Msg(info)
	case "info":
		log.Info().Str("Context", context).Str("Module", module).Msg(info)
	case "debug":
		log.Debug().Str("Context", context).Str("Module", module).Msg(info)
	}
}

//------------------------------------------------------------------------------