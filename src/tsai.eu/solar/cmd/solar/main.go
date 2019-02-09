package main

import (
	"fmt"

	"tsai.eu/solar/engine"
	"tsai.eu/solar/model"
	"tsai.eu/solar/shell"
	"tsai.eu/solar/util"
)

//------------------------------------------------------------------------------

// main entry point for the orchestrator
func main() {
	// initialise command line options
	util.ParseCommandLineOptions()

	// display progam information
	fmt.Println("SOLAR Version 1.0.0")

	// create model
	m := model.GetModel()

	// start the main event loop
	engine.StartDispatcher(m)

	// start the command line interface
	shell.Run(m)
}

//------------------------------------------------------------------------------
