package file

import (
	"tsai.eu/solar/model"
)

//------------------------------------------------------------------------------

// Reset cleans up an instance
func (c Controller) Reset(configuration *model.ComponentConfiguration) (status *model.ComponentStatus, err error) {
	return c.Destroy(configuration)
}

//------------------------------------------------------------------------------
