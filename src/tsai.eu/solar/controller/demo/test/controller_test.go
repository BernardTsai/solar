package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"tsai.eu/solar/controller/demo"
	"tsai.eu/solar/model"
	"tsai.eu/solar/util"
)

//------------------------------------------------------------------------------

const TESTDATA       string = "testdata"
const CONFIGURATIONS string = "_configurations"

//------------------------------------------------------------------------------

// LoadConfigurations load a list of configuration filenames
func LoadConfigurations() (list []string, err error) {
	path := filepath.Join(TESTDATA, CONFIGURATIONS)

	// retrieve list
	result, err := util.LoadFile(path)
	if err != nil {
		fmt.Println(err)
		return []string{}, err
	}

	// convert to a slice of strings
	list = strings.Split(result, "\n")

	return
}

//------------------------------------------------------------------------------

// LoadConfiguration load a configuration file
func LoadConfiguration(filename string) *model.Setup {
	setup := model.Setup{}
	path  := filepath.Join(TESTDATA, filename)
	util.LoadYAML(path, &setup)

	return &setup
}

//------------------------------------------------------------------------------
// Procedure:
// - create root file component
// - check status
// - start root file component
// - check status
//
// Structure:
//
// <root>
//  nodeA (instances nodeA1(V1.0.0), nodeA2(V1.0.0), nodeA3((V2.0.0), nodeA4((V2.0.0))
//  nodeB (instances nodeB1(V1.0.0), nodeB2(V1.0.0), nodeB3((V2.0.0), nodeB4((V2.0.0))

// TestController verifies the DemoController object.
func TestController(t *testing.T) {
	var setup *model.Setup
	// var status *model.ComponentStatus
	var err error

	// create controller
	fc := demo.Controller{}

	// Load names of configuration files
	list, err := LoadConfigurations()
	if err != nil {
		fmt.Println(err)
		t.Errorf("%s", err)
	}
 
	for _, entry := range list {
		if entry != "" {
			setup = LoadConfiguration(entry)
			_, err = fc.Create(setup)
			if err != nil {
				fmt.Println("Unable to create:" + entry)
				fmt.Println(err)
			}
			_, err = fc.Start(setup)
			if err != nil {
				fmt.Println("Unable to start:" + entry)
				fmt.Println(err)
			}
		}
	}
}

//------------------------------------------------------------------------------
