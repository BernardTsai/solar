package dummy

import (
	"fmt"
	"path/filepath"
	"testing"

	"tsai.eu/solar/model"
	"tsai.eu/solar/util"
)

//------------------------------------------------------------------------------

const TESTDATA       string = "testdata"
const CONFIGURATIONS string = "_configurations"

//------------------------------------------------------------------------------

// Task defines an action to be executed with a specific setup
type Task struct {
	Action string `yaml:"action"`     // action to be executed
	Setup  string `yaml:"setup"`      // name of the setup file
}

//------------------------------------------------------------------------------

// Tasks defines a list of tasks to be executed sequentially
type Tasks struct {
	Tasks []Task `yaml:"tasks"` // list of tasks to be executed
}

//------------------------------------------------------------------------------

// loadConfigurations load a list of configuration filenames
func loadConfigurations() (list []Task, err error) {
	path := filepath.Join(TESTDATA, CONFIGURATIONS)
	tasks := Tasks{}

	// retrieve list
	err = util.LoadYAML(path, &tasks)
	if err != nil {
		return list, err
	}

	// success
	list = tasks.Tasks

	return list, nil
}

//------------------------------------------------------------------------------

// loadConfiguration load a configuration file
func loadConfiguration(filename string) *model.Setup {
	setup := model.Setup{}
	path  := filepath.Join(TESTDATA, filename)
	util.LoadYAML(path, &setup)

	return &setup
}

//------------------------------------------------------------------------------

// TestController verifies the DemoController object.
func TestController(t *testing.T) {
	var setup *model.Setup
	// var status *model.ComponentStatus
	var err error

	// create controller
	dc := NewController()

	// Load names of configuration files
	list, err := loadConfigurations()
	if err != nil {
		fmt.Println(err)
		t.Errorf("%s", err)
	}

	fmt.Println("Executing tests:")
	fmt.Println("Nr. Action     Setup          ")
	fmt.Println("------------------------------")
	for index, entry := range list {
		fmt.Printf("%03d %10s %s\n", index, entry.Action, entry.Setup)

		setup = loadConfiguration(entry.Setup)

		switch entry.Action {
		case "create":
			_, err = dc.Create(setup)
			if err != nil {
				t.Errorf("Unable to create: %s\n%s", entry, err)
			}
		case "start":
			_, err = dc.Start(setup)
			if err != nil {
				t.Errorf("Unable to start: %s\n%s", entry, err)
			}
		case "status":
			_, err = dc.Status(setup)
			if err != nil {
				t.Errorf("Unable to get status: %s\n%s", entry, err)
			}
		case "stop":
			_, err = dc.Stop(setup)
			if err != nil {
				t.Errorf("Unable to stop: %s\n%s", entry, err)
			}
		case "configure":
			_, err = dc.Configure(setup)
			if err != nil {
				t.Errorf("Unable to configure: %s\n%s", entry, err)
			}
		case "reconfigure":
			_, err = dc.Reconfigure(setup)
			if err != nil {
				t.Errorf("Unable to configure: %s\n%s", entry, err)
			}
		case "destroy":
			_, err = dc.Destroy(setup)
			if err != nil {
				t.Errorf("Unable to destroy: %s\n%s", entry, err)
			}
		case "reset":
			_, err = dc.Reset(setup)
			if err != nil {
				t.Errorf("Unable to reset: %s\n%s", entry, err)
			}
		default:
			t.Errorf("Unknown action: %s", entry.Action)
		}
	}
}

//------------------------------------------------------------------------------
