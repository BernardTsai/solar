package dummy

import (
	"tsai.eu/solar/model"
)

//------------------------------------------------------------------------------

// Create initialises an instance
func (c Controller) Create(setup *model.Setup) (status *model.Status, err error) {
	return c.Status(setup)
}

//------------------------------------------------------------------------------
