package defaultRestController

import (
	"net/http"
)

//------------------------------------------------------------------------------

// Create initialises an instance
func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	c.Status(w, r)
}

//------------------------------------------------------------------------------
