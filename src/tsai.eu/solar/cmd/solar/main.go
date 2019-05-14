package main

import (
	"context"
	"fmt"

	"tsai.eu/solar/engine"
	"tsai.eu/solar/api"
	"tsai.eu/solar/monitor"
	"tsai.eu/solar/controller"
	"tsai.eu/solar/msg"
	"tsai.eu/solar/cli"
	"tsai.eu/solar/util"
)

//------------------------------------------------------------------------------

// Control holds a handle to all running process
type Control struct {
	Cancel     context.CancelFunc      // process context
	Dispatcher *engine.Dispatcher      // the orchestration engine
	MSG        *msg.MSG                // messaging interface
	Monitor    *monitor.Monitor	       // monitoring process
	Controller *controller.Manager     // controller manager
	API        *api.API                // web API
}

//------------------------------------------------------------------------------

// main entry point for the orchestrator
func main() {
	control := Control{
		Cancel:     nil,
		Dispatcher: nil,
		MSG:        nil,
		Monitor:    nil,
	}

	// Create a background context
  ctx := context.Background()

  //Derive a context with cancel
  mainCtx, cancelFunction := context.WithCancel(ctx)
	control.Cancel = cancelFunction

	// defer canceling so that all the resources are freed up for this and the derived contexts
  defer func() { terminate(&control) }()

	// parse configuration file 'solar-conf.yaml' in local directory and initiate logging
	util.StartLogging()

	// initialise command line options
	util.ParseCommandLineOptions()

	// display progam information
	fmt.Println("SOLAR Version 1.0.0")

	// start the controller manager
	control.Controller = controller.Start(mainCtx)

	// start the main event loop
	control.Dispatcher = engine.Start(mainCtx)

	// start the messaging interface listener
	control.MSG, _ = msg.Start(mainCtx)

	// start the monitoring loop
	control.Monitor = monitor.Start(mainCtx)

	// start the API
	control.API = api.Start(mainCtx)

	// get the command line interface
	shell := cli.Shell()

	shell.Run()
}

//------------------------------------------------------------------------------

// terminate frees all resources
func terminate(control *Control) {
	// close the main process context
	control.Cancel()
}

//------------------------------------------------------------------------------
