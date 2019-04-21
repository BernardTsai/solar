package defaultRestController

import (
	"net/http"
)

//------------------------------------------------------------------------------

// Destroy removes instance
func (c *Controller) Destroy(w http.ResponseWriter, r *http.Request) {
	c.Status(w, r)
}

//------------------------------------------------------------------------------
