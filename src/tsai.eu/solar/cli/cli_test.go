package cli

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"context"
	"io"
	"os"

	"tsai.eu/solar/engine"
	"tsai.eu/solar/util"
)

//------------------------------------------------------------------------------

const TESTDATA string = "testdata"
const COMMANDS string = "_commands"

//------------------------------------------------------------------------------

// copyFile simply copies an existing file to a not yet existing destination
func copyFile(src string, dest string) {
  srcFile, _ := os.Open(src)
  defer srcFile.Close()

  destFile, _ := os.Create(dest)
  defer destFile.Close()

  io.Copy(destFile, srcFile)
  destFile.Sync()
}

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
	filename := "output.txt"
	// cleanup routine
  defer func() {os.Remove(filename)}()

	// prepare configuration file
  srcConfig  := "testdata/solar-conf.yaml"
  destConfig := "solar-conf.yaml"

  copyFile(srcConfig, destConfig)

  // cleanup routine
  defer func() {os.Remove(destConfig)}()

	// initialise command line options
	util.LogLevel("error")
	util.ParseCommandLineOptions()

	// display progam information
	fmt.Println("SOLAR Version 1.0.0")

	// start the main event loop
	engine.Start(context.Background())

	// get the command line interface
	shell := Shell()

	// load commands from a file
	cmds, err := LoadCommands()
	if err != nil {
		fmt.Println(err)
		t.Errorf("%s", err)
	}

	fmt.Println("Executing tests:")
	fmt.Println("Nr. OK Command")
	fmt.Println("------------------------------------------------------------")
	for index, cmd := range cmds {
		// skip empty command lines
		if cmd[0] == "" {
			continue
		}

		// construct command line string
		cmdline := strings.Join(cmd, " ")

		status  := cmdline[0:2]
		cmdline =  cmdline[3:]
		cmd     = cmd[1:]

		// log the command line
		fmt.Printf("%03d %s %s\n", index, status, cmdline)

		// process
		err = shell.Process(cmd...)
		if err != nil {
			if status != "KO" {
				t.Errorf("Command failed: %s\n%s", cmdline, err)
			}
		} else {
			if status != "OK" {
				t.Errorf("Command should have failed: %s\n%s", cmdline, err)
			}
		}
	}
}

//------------------------------------------------------------------------------
