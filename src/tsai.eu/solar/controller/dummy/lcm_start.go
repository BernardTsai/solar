package dummy

import (
	"tsai.eu/solar/model"
)

//------------------------------------------------------------------------------

// Start activates an instance
func (c Controller) Start(setup *model.Setup) (status *model.Status, err error) {
	return c.Status(setup)
}

//------------------------------------------------------------------------------
