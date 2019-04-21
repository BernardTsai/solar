package controller

import (
	"errors"
	"strings"
	"net/http"
	"io/ioutil"

	"tsai.eu/solar/model"
	"tsai.eu/solar/util"
)

//------------------------------------------------------------------------------

// RestController is an gRPC based implementation of the Controller interface
type RestController struct {
	Type    string  // type of controller
	Version string  // version of the controller
	URL     string  // address to which the controller listens
}

//------------------------------------------------------------------------------

// newRestController creates a gRPC based controller
func newRestController(Type string, Version string, URL string) (*RestController, error) {
	c := RestController{
		Type:       Type,
		Version:    Version,
		URL:        URL,
	}

	// check availability of client
	if !c.Check() {
		util.LogWarn("main", "CTRL", "controller: " + Type + " unavailable")
		return nil, errors.New("controller: " + Type + ":" + Version + " unavailable")
	}

	// success
	util.LogInfo("main", "CTRL", "controller: " + Type + ":" + Version +  " available")
	return &c, nil
}

//------------------------------------------------------------------------------

// process triggers the request towards the controller
func (c *RestController)process(action string, targetState *model.TargetState) (currentState *model.CurrentState, err error) {
	body, _ := util.ConvertToYAML(targetState)

	response, err := http.Post(c.URL + "/" + action, "application/x-yaml", strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	err = util.ConvertFromYAML(string(data), currentState)
	if err != nil {
		return nil, err
	}

	return currentState, nil
}

//------------------------------------------------------------------------------

// Check checks availability of controller
func (c *RestController) Check() bool {
	_, err := http.Get(c.URL + "/check")
	if err != nil {
		return false
	}

	// success
	return true
}

//------------------------------------------------------------------------------

// Status determines the currentState of an instance
func (c *RestController)Status(targetState *model.TargetState) (*model.CurrentState, error) {
	return c.process("status", targetState)
}

//------------------------------------------------------------------------------

// Create instantiates an instance
func (c *RestController)Create(targetState *model.TargetState) (*model.CurrentState, error) {
	return c.process("create", targetState)
}

//------------------------------------------------------------------------------

// Destroy removes an instance
func (c *RestController)Destroy(targetState *model.TargetState) (*model.CurrentState, error) {
	return c.process("destroy", targetState)
}


//------------------------------------------------------------------------------

// Configure configures an instance
func (c *RestController)Configure(targetState *model.TargetState) (*model.CurrentState, error) {
	return c.process("configure", targetState)
}

//------------------------------------------------------------------------------

// Reconfigure reconfigures an instance
func (c *RestController)Reconfigure(targetState *model.TargetState) (*model.CurrentState, error) {
	return c.process("reconfigure", targetState)
}

//------------------------------------------------------------------------------

// Start activates an instance
func (c *RestController)Start(targetState *model.TargetState) (*model.CurrentState, error) {
	return c.process("start", targetState)
}

//------------------------------------------------------------------------------

// Stop deactivates an instance
func (c *RestController)Stop(targetState *model.TargetState) (*model.CurrentState, error) {
	return c.process("stop", targetState)
}

//------------------------------------------------------------------------------

// Reset cleans up a failed instance
func (c *RestController)Reset(targetState *model.TargetState) (*model.CurrentState, error) {
	return c.process("reset", targetState)
}

//------------------------------------------------------------------------------
