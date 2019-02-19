package demo

import "tsai.eu/solar/util"

//------------------------------------------------------------------------------

// Information describes the runtime configuration of an instance
type Information struct {
	State         string            `yaml:"state"`      // state of instance (active or inactive)
	Path          string            `yaml:"path"`       // path to instance
	Template      string            `yaml:"template"`   // template
	Refererences  map[string]string `yaml:"references"` // path to other references
}

// LoadInformation loads the contents of an instance information file
func LoadInformation(filename string) (info *Information, err error) {
	info = &Information{}

	err = util.LoadYAML(filename, info)

	return
}

// SaveInformation writes an instance information object to a file
func SaveInformation(filename string, info *Information) (err error) {
	err = util.SaveYAML(filename, info)

	return err
}

// DecodeInformation converts yaml into an instance information object
func DecodeInformation(yaml string) (info *Information, err error) {
	info = &Information{}

	err = util.ConvertFromYAML(yaml, info)

	return info, err
}

// EncodeInformation converts an instance information object into yaml
func EncodeInformation(info *Information) (string, error) {
	yaml, err := util.ConvertToYAML(info)

	return yaml, err
}

//------------------------------------------------------------------------------
