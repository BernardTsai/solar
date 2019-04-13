package dummy

import (
	"tsai.eu/solar/model"
)

//------------------------------------------------------------------------------

// Reconfigure reconfigures an instance
func (c Controller) Reconfigure(setup *model.Setup) (status *model.Status, err error) {
	return c.Status(setup)
}

//------------------------------------------------------------------------------
