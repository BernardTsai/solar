package defaultRestController

import (
	"net/http"
)

//------------------------------------------------------------------------------

// Start activates an instance
func (c *Controller) Start(w http.ResponseWriter, r *http.Request) {
	c.Status(w, r)
}

//------------------------------------------------------------------------------
