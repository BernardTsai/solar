package util

import (
  "os"
  "testing"
)

//------------------------------------------------------------------------------

// configuration1 represents a valid configuration file
const configuration1 string = `
MSG:
  Notifications:  notifications
  Monitoring:     monitoring
  Address:        127.0.0.1:9092
CORE:
  IDENTIFIER: solar
  LOGLEVEL:   debug
CONTROLLERS:
  - tsai/solar-k8s-controller:V1.0.0
  - tsai/solar-default-controller:V1.0.0
`

//------------------------------------------------------------------------------

// TestConfig01 tests the configuration related functions of the util package.
func TestConfig01(t *testing.T) {
  // cleanup routine
  defer func() {
    os.Remove("./solar-conf.yaml")
  }()

  // create configuration file
  f, err := os.Create("solar-conf.yaml")
  if err != nil {
    t.Fatalf("Unable to create configuration file:\n%s", err )
  }

  // write content
  _, err = f.WriteString(configuration1)
  if err != nil {
    t.Fatalf("Unable to write configuration file:\n%s", err )
  }

  // close configuration file
  err = f.Close()
  if err != nil {
    t.Fatalf("Unable to write configuration file:\n%s", err )
  }

  // read configuration
  configuration, err := GetConfiguration()
  if err != nil {
    t.Errorf("Unable to read and parse configuration file:\n%s", err)
  }

  // vaidate configuration
  if configuration.CORE.LogLevel != "debug" {
    t.Error("Detected inconsistencies when reading configuration")
  }
}

//------------------------------------------------------------------------------

// TestConfig02 tests the configuration related functions of the util package.
func TestConfig02(t *testing.T) {
  // delete configuration file
  os.Remove("./solar-conf.yaml")

  // read configuration
  configuration, err := ReadConfiguration(".")
  if err == nil {
    t.Errorf("ReadConfiguration should have reported an error but has responded with:\n%s", configuration)
  }
}

//------------------------------------------------------------------------------
