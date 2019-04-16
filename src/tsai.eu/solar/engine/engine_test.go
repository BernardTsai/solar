package engine

import (
  "testing"
  "context"
  "os"
  "io"
  "time"

  "tsai.eu/solar/model"
)

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

// TestEngine001 tests the basic execution functions
func TestEngine001(t *testing.T) {
  filename := "testdata/testdata1.yaml"

  // prepare configuration file
  srcConfig  := "testdata/solar-conf.yaml"
  destConfig := "solar-conf.yaml"

  copyFile(srcConfig, destConfig)

  // cleanup routine
  defer func() {os.Remove(destConfig)}()

  Start(context.Background())

  // load model
  m := model.GetModel()
  m.Load(filename)

  // trigger execution
  // create new solution
  domain, _       := model.GetDomain("demo")
  solution, _     := model.NewSolution("app", "V0.0.0", "")
  architecture, _ := model.GetArchitecture("demo", "app", "V0.0.0")

	domain.AddSolution(solution)

	// update the target state of the solution
	err := solution.Update("demo", architecture)
  if err != nil {
    t.Errorf("unable to create or update the solution:\n" + err.Error())
	}

  task, err := NewSolutionTask("demo", "", solution)
  if err != nil {
    t.Errorf("task can not be created")
  }

  // get event channel
  channel := GetEventChannel()

  // create event
  channel <- model.NewEvent("demo", task.UUID, model.EventTypeTaskExecution, "", "initial")

  time.Sleep(time.Millisecond)
}

//------------------------------------------------------------------------------
