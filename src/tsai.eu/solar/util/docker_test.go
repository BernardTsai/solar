package util

import (
  "testing"
)

//------------------------------------------------------------------------------

// TestDockers01 tests the docker related functions of the util package.
func TestDocker01(t *testing.T) {
  _, err := ListImages()
  if err != nil {
    t.Error("Unable to list images:\n" + err.Error())
  }

  _, err = ListContainers()
  if err != nil {
    t.Error("Unable to list containers:\n" + err.Error())
  }

  err = PullImage("alpine", "latest")
  if err != nil {
    t.Error("Unable to pull image:\n" + err.Error())
  }

  // remove any old containers
  StopContainer("alpine", "latest")

  _, err = StartContainer("alpine", "latest")
  if err != nil {
    t.Error("Unable to start container:\n" + err.Error())
  }

  ListContainers()

  err = StopContainer("alpine", "latest")
  if err != nil {
    t.Error("Unable to stop container:\n" + err.Error())
  }
}

//------------------------------------------------------------------------------
