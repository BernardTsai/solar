package main

import (
	"fmt"

	"tsai.eu/solar/engine"
	"tsai.eu/solar/model"
	"tsai.eu/solar/cli"
	"tsai.eu/solar/util"
)

//------------------------------------------------------------------------------

// main entry point for the orchestrator
func main() {
	// parse configuration file 'solar-conf.yaml' in local directory
	_, err := util.ReadConfiguration()
	if err != nil {
		fmt.Println("unable to read the configuration file")
		fmt.Println(err)
	}

	// initialise command line options
	util.ParseCommandLineOptions()

	// display progam information
	fmt.Println("SOLAR Version 1.0.0")

	// create model
	m := model.GetModel()

	// start the main event loop
	engine.StartDispatcher(m)

	// get the command line interface
	shell := cli.Shell(m)

	shell.Run()
}

//------------------------------------------------------------------------------
