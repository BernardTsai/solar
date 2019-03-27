package main

import (
	"fmt"

	"tsai.eu/solar/engine"
	"tsai.eu/solar/model"
	"tsai.eu/solar/api"
	"tsai.eu/solar/monitor"
	"tsai.eu/solar/msg"
	"tsai.eu/solar/cli"
	"tsai.eu/solar/util"
)

//------------------------------------------------------------------------------

// main entry point for the orchestrator
func main() {
	// parse configuration file 'solar-conf.yaml' in local directory
	_, err := util.GetConfiguration()
	if err != nil {
		fmt.Println("Unable to read the configuration file")
		fmt.Println(err)
		return
	}

	// attempt to connect to the message bus
	err = msg.Open()
	if err != nil {
		fmt.Println(err)
	} else {
		defer msg.Close()
	}

	// initialise command line options
	util.ParseCommandLineOptions()

	// display progam information
	fmt.Println("SOLAR Version 1.0.0")

	// create model
	m := model.GetModel()

	// start the main event loop
	dispatcher := engine.StartDispatcher(m)
	defer dispatcher.Stop()

	// start the monitoring loop
	moni := monitor.StartMonitor(m, dispatcher.Channel)
	defer moni.Stop()

	// start the API
	go api.NewRouter()

	// get the command line interface
	shell := cli.Shell(m)

	shell.Run()
}

//------------------------------------------------------------------------------
