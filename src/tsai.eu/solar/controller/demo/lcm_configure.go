package demo

import (
	"tsai.eu/solar/model"
)

//------------------------------------------------------------------------------

// Configure reconfigures an instance
func (c Controller) Configure(setup *model.Setup) (status *model.Status, err error) {
	return c.Status(setup)
}

//------------------------------------------------------------------------------
