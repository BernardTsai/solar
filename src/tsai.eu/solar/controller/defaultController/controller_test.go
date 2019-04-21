package defaultController

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
	State  string `yaml:"state"`      // name of the state file
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
func loadConfiguration(filename string) *model.TargetState {
	state := model.TargetState{}
	path  := filepath.Join(TESTDATA, filename)
	util.LoadYAML(path, &state)

	return &state
}

//------------------------------------------------------------------------------

// TestController verifies the DemoController object.
func TestController(t *testing.T) {
	var state *model.TargetState
	var err    error

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
		fmt.Printf("%03d %10s %s\n", index, entry.Action, entry.State)

		state = loadConfiguration(entry.State)

		switch entry.Action {
		case "create":
			_, err = dc.Create(state)
			if err != nil {
				t.Errorf("Unable to create: %s\n%s", entry, err)
			}
		case "start":
			_, err = dc.Start(state)
			if err != nil {
				t.Errorf("Unable to start: %s\n%s", entry, err)
			}
		case "status":
			_, err = dc.Status(state)
			if err != nil {
				t.Errorf("Unable to get status: %s\n%s", entry, err)
			}
		case "stop":
			_, err = dc.Stop(state)
			if err != nil {
				t.Errorf("Unable to stop: %s\n%s", entry, err)
			}
		case "configure":
			_, err = dc.Configure(state)
			if err != nil {
				t.Errorf("Unable to configure: %s\n%s", entry, err)
			}
		case "reconfigure":
			_, err = dc.Reconfigure(state)
			if err != nil {
				t.Errorf("Unable to configure: %s\n%s", entry, err)
			}
		case "destroy":
			_, err = dc.Destroy(state)
			if err != nil {
				t.Errorf("Unable to destroy: %s\n%s", entry, err)
			}
		case "reset":
			_, err = dc.Reset(state)
			if err != nil {
				t.Errorf("Unable to reset: %s\n%s", entry, err)
			}
		default:
			t.Errorf("Unknown action: %s", entry.Action)
		}
	}
}

//------------------------------------------------------------------------------
