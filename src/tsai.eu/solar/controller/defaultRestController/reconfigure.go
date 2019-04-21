package defaultRestController

import (
	"net/http"
)

//------------------------------------------------------------------------------

// Reconfigure reconfigures an instance
func (c *Controller) Reconfigure(w http.ResponseWriter, r *http.Request) {
	c.Status(w, r)
}

//------------------------------------------------------------------------------
