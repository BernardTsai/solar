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


	// initialise command line options
	util.ParseCommandLineOptions()

	// display progam information
	fmt.Println("SOLAR Version 1.0.0")

	// create model
	m := model.GetModel()

	// start the main event loop
	dispatcher := engine.StartDispatcher(m)
	defer dispatcher.Stop()

	// start the messaging interface listener
	msg, err := msg.StartMSG()
	if err == nil {
		defer msg.Stop()
	} else {
		fmt.Println("Unable to start the messaging interface")
		fmt.Println(err)
	}

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
