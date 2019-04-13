package dummy

import (
	"tsai.eu/solar/model"
)

//------------------------------------------------------------------------------

// Stop deactivates an instance
func (c Controller) Stop(setup *model.Setup) (status *model.Status, err error) {
	return c.Status(setup)
}

//------------------------------------------------------------------------------
