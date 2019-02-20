package demo

//------------------------------------------------------------------------------

// ROOTDIR points to the root directory of the file system
const ROOTDIR = "/tmp/demo"

// Controller manages the lifecycle of a file
type Controller struct {
  Root string `yaml:"Root"`   // root directory
}

//------------------------------------------------------------------------------

// NewController creates a new controller for demo elements
func NewController() *Controller {
  controller := Controller{
    Root: ROOTDIR,
  }

  return &controller
}

//------------------------------------------------------------------------------
