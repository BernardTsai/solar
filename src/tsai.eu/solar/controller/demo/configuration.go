package demo

import "tsai.eu/solar/util"

//------------------------------------------------------------------------------

// Configuration describes the configuration of a component
type Configuration struct {
	Template string `yaml:"template"`   // template
}

func decodeConfiguration(yaml string) (*Configuration, error) {
	config := Configuration{}

	err := util.ConvertFromYAML(yaml, &config)

	return &config, err
}

func encodeConfiguration(config *Configuration) (string, error) {
	yaml, err := util.ConvertToYAML(config)

	return yaml, err
}

//------------------------------------------------------------------------------
