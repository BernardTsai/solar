package defaultRestController

import (
	"net/http"
)

//------------------------------------------------------------------------------

// Configure reconfigures an instance
func (c *Controller) Configure(w http.ResponseWriter, r *http.Request) {
	c.Status(w, r)
}

//------------------------------------------------------------------------------
