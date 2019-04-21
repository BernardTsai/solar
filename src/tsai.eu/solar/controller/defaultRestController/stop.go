package defaultRestController

import (
	"net/http"
)

//------------------------------------------------------------------------------

// Stop deactivates an instance
func (c *Controller) Stop(w http.ResponseWriter, r *http.Request) {
	c.Status(w, r)
}

//------------------------------------------------------------------------------
