package demo

import "tsai.eu/solar/util"

//------------------------------------------------------------------------------

// Configuration describes the configuration of a component
type Configuration struct {
	Template string `yaml:"template"`   // template
}

// DecodeConfiguration converts a yaml string into a configuration object
func DecodeConfiguration(yaml string) (*Configuration, error) {
	config := Configuration{}

	err := util.ConvertFromYAML(yaml, &config)

	return &config, err
}

// EncodeConfiguration converts a configuration into a yaml string
func EncodeConfiguration(config *Configuration) (string, error) {
	yaml, err := util.ConvertToYAML(config)

	return yaml, err
}

//------------------------------------------------------------------------------
