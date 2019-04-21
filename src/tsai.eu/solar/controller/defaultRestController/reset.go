package defaultRestController

import (
	"net/http"
)

//------------------------------------------------------------------------------

// Reset cleans up an instance
func (c *Controller) Reset(w http.ResponseWriter, r *http.Request) {
	c.Status(w, r)
}

//------------------------------------------------------------------------------
