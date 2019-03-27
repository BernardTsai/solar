package util

import (
  "sync"
  "github.com/spf13/viper"
)

//------------------------------------------------------------------------------

// MsgConfiguration holds all configuration information for the MSG module
type MsgConfiguration struct {
  Events  string
  Status  string
  Address string
}

//------------------------------------------------------------------------------

// Configuration holds all configuration information for the application
type Configuration struct {
  MSG MsgConfiguration
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
  viper.SetDefault("MSG", map[string]string{"Events": "events", "Status": "status", "Address": "127.0.0.1:9092"})

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
