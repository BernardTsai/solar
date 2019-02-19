package demo

import "tsai.eu/solar/util"

//------------------------------------------------------------------------------

// Endpoint describes the endpoint of the component
type Endpoint struct {
	Path string `yaml:"path"` // path to instance
}

//------------------------------------------------------------------------------

// NewEndpoint creates an endpoint from a path
func NewEndpoint(path string) (endp *Endpoint) {
	endp = &Endpoint{
		Path: path,
	}

	return
}

//------------------------------------------------------------------------------

// DecodeEndpoint decodes endpoint yaml into an object
func DecodeEndpoint(yaml string) (*Endpoint, error) {
	endp := Endpoint{}

	err := util.ConvertFromYAML(yaml, &endp)

	return &endp, err
}

//------------------------------------------------------------------------------

// EncodeEndpoint converts an endpoint into yaml
func EncodeEndpoint(endp *Endpoint) (string, error) {
	yaml, err := util.ConvertToYAML(endp)

	return yaml, err
}

const emptyEndpoint string = ""

//------------------------------------------------------------------------------
