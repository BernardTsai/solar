package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"tsai.eu/solar/engine"
	"tsai.eu/solar/model"
	"tsai.eu/solar/cli"
	"tsai.eu/solar/util"
)

//------------------------------------------------------------------------------

const TESTDATA string = "testdata"
const COMMANDS string = "_commands"

//------------------------------------------------------------------------------

// LoadCommands load a list of commands from a file
func LoadCommands() ([][]string, error) {
	list := [][]string{}

	// load data from file
	path := filepath.Join(TESTDATA, COMMANDS)

	data, err := util.LoadFile(path)
	if err != nil {
		return list, err
	}

	// parse commands
	cmds := strings.Split(data, "\n")

	// parse args and add commands to result
	for _, cmd := range cmds {
		args := strings.Split(cmd, " ")

		list = append(list, args)
	}

	// success
	return list, nil
}

//------------------------------------------------------------------------------

// TestCLI verifies the command line interface.
func TestCLI(t *testing.T) {
	// parse configuration file 'solar-conf.yaml' in local directory
	_, err := util.ReadConfiguration()
	if err != nil {
		fmt.Println("unable to read the configuration file")
		fmt.Println(err)
	}

	// initialise command line options
	util.ParseCommandLineOptions()

	// create model
	m := model.GetModel()
 
	// start the main event loop
	engine.StartDispatcher(m)

	// get the command line interface
	shell := cli.Shell(m)

	// load commands from a file
	cmds, err := LoadCommands()
	if err != nil {
		fmt.Println(err)
		t.Errorf("%s", err)
	}

	fmt.Println("Executing tests:")
	fmt.Println("Nr. Command")
	fmt.Println("------------------------------------------------------------")
	for index, cmd := range cmds {
		// skip empty command lines
		if cmd[0] == "" {
			continue
		}

		// construct command line string
		cmdline := strings.Join(cmd, " ")

		// log the command line
		fmt.Printf("%03d %s\n", index, cmdline)

		// process
		err = shell.Process(cmd...)
		if err != nil {
			t.Errorf("Command failed: %s\n%s", cmdline, err)
		}
	}
}

//------------------------------------------------------------------------------
