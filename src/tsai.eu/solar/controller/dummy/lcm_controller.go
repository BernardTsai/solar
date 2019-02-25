package dummy

//------------------------------------------------------------------------------

// Controller manages the lifecycle of nothing
type Controller struct {
}

//------------------------------------------------------------------------------

// NewController creates a new controller for dummy elements
func NewController() *Controller {
  controller := Controller{}

  return &controller
}

//------------------------------------------------------------------------------
