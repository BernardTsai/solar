package dummy

import (
	"tsai.eu/solar/model"
)

//------------------------------------------------------------------------------

// Destroy removes instance
func (c Controller) Destroy(setup *model.Setup) (status *model.Status, err error) {
	return c.Status(setup)
}

//------------------------------------------------------------------------------
