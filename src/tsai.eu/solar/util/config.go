package util

import (
  "github.com/spf13/viper"
)

//------------------------------------------------------------------------------

// MsgConfiguration holds all configuration information for the MSG module
type MsgConfiguration struct {
  Events string
  Status string
}

//------------------------------------------------------------------------------

// Configuration holds all configuration information for the application
type Configuration struct {
  MSG MsgConfiguration
}

//------------------------------------------------------------------------------

// ReadConfiguration reads a file into a Configuration object
func ReadConfiguration() (Configuration, error) {
  var configuration Configuration

  // define the location from which to read the configuration file
  viper.SetConfigName("solar-conf")
  viper.AddConfigPath(".")

  // set default values
  viper.SetDefault("MSG", map[string]string{"Events": "events", "Status": "status"})

  // read configuration
	err := viper.ReadInConfig()
  if err != nil {
    return configuration, err
  }

  // decode the configuration
  err = viper.Unmarshal(&configuration)
	if err != nil {
    return configuration, err
  }

  // success
  return configuration, nil
}

//------------------------------------------------------------------------------
