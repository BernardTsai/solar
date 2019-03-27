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
}

//------------------------------------------------------------------------------

// Configuration holds all configuration information for the application
type Configuration struct {
  MSG  MsgConfiguration
  CORE CoreConfiguration
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
    theConfiguration, err = readConfiguration()
  })

	// success
	return theConfiguration, err
}

//------------------------------------------------------------------------------

// ReadConfiguration reads a file into a Configuration object
func readConfiguration() (*Configuration, error) {
  var configuration Configuration

  // define the location from which to read the configuration file
  viper.SetConfigName("solar-conf")
  viper.AddConfigPath(".")

  // set default values
  viper.SetDefault("MSG",  map[string]string{"Notifications": "notifications", "Monitoring": "monitoring", "Address": "127.0.0.1:9092"})
  viper.SetDefault("CORE", map[string]string{"Identifier": "solar"})

  // read configuration (ignore any errors)
	viper.ReadInConfig()

  // decode the configuration
  err := viper.Unmarshal(&configuration)
	if err != nil {
    return &configuration, err
  }

  // success
  return &configuration, nil
}

//------------------------------------------------------------------------------
