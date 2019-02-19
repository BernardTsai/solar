package demo

import (
	"tsai.eu/solar/model"
)

//------------------------------------------------------------------------------

// Reset cleans up an instance
func (c Controller) Reset(setup *model.Setup) (status *model.Status, err error) {
	return c.Destroy(setup)
}

//------------------------------------------------------------------------------
