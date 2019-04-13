package monitor

import (
  "testing"
  "context"
  "time"

  "tsai.eu/solar/model"
)

//------------------------------------------------------------------------------

// TestMonitor001 tests the basic execution functions
func TestMonitor001(t *testing.T) {
  filename := "testdata/testdata1.yaml"

  // load model
  m := model.GetModel()
  m.Load(filename)

  mon := Start(context.Background())

  mon.Stop()

  mon.Start()

  time.Sleep(time.Second)
}

//------------------------------------------------------------------------------
