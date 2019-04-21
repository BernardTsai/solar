package util

import (
  "sync"
  "github.com/spf13/viper"
)

//------------------------------------------------------------------------------

// MsgConfiguration holds all configuration information for the MSG module
type MsgConfiguration struct {
  Notifications string
  Monitoring    string
  Address       string
}

//------------------------------------------------------------------------------

// CoreConfiguration holds all configuration information related to the orchestrator
type CoreConfiguration struct {
  Identifier  string
  LogLevel    string
}

//------------------------------------------------------------------------------

// ControllerConfiguration holds all configuration information related to a controller
type ControllerConfiguration struct {
  Type    string
  Version string
  URL     string
}

//------------------------------------------------------------------------------

// Configuration holds all configuration information for the application
type Configuration struct {
  MSG  MsgConfiguration
  CORE CoreConfiguration
  CTRL []ControllerConfiguration
}

//------------------------------------------------------------------------------

var theConfiguration *Configuration

var configurationInit sync.Once

//------------------------------------------------------------------------------

// GetConfiguration retrieves the configuration.
func GetConfiguration() (*Configuration, error) {
  var err error

	// initialise singleton once
	configurationInit.Do(func() {
    theConfiguration, err = readConfiguration(".")
  })

	// success
	return theConfiguration, err
}

//------------------------------------------------------------------------------

// ReadConfiguration reads a file from a specific path into a Configuration object
func ReadConfiguration(path string) (*Configuration, error) {
  return readConfiguration(path)
}

//------------------------------------------------------------------------------

// readConfiguration reads a file into a Configuration object
func readConfiguration(path string) (*Configuration, error) {
  var configuration Configuration

  // define the location from which to read the configuration file
  viper.SetConfigName("solar-conf")
  viper.AddConfigPath(path)

  // set default values
  viper.SetDefault("MSG",  map[string]string{"Notifications": "notifications", "Monitoring": "monitoring", "Address": "127.0.0.1:9092"})
  viper.SetDefault("CORE", map[string]string{"Identifier": "solar", "LogLevel": ""})
  viper.SetDefault("CTRL", []ControllerConfiguration{})

  // read configuration (ignore any errors)
  if err := viper.ReadInConfig(); err != nil {
    return nil, err
  }

  // decode the configuration
  viper.Unmarshal(&configuration)

  // success
  return &configuration, nil
}

//------------------------------------------------------------------------------
