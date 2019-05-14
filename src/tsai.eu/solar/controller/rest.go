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
		util.LogWarn("main", "CTL", "controller: " + Type + " unavailable")
		return nil, errors.New("controller: " + Type + ":" + Version + " unavailable")
	}

	// success
	util.LogInfo("main", "CTL", "controller: " + Type + ":" + Version +  " available")
	return &c, nil
}

//------------------------------------------------------------------------------

// process triggers the request towards the controller
func (c *RestController)process(action string, targetState *model.TargetState) (currentState *model.CurrentState, err error) {
	// convert targetState into a request
	request := Request{
		Request:         util.UUID(),
		Domain:          targetState.Domain,
		Solution:        targetState.Solution,
		Version:         targetState.Version,
		Element:         targetState.Element,
		Cluster:         targetState.Cluster,
		Instance:        targetState.Instance,
		Component:       targetState.Component,
		State:           targetState.State,
		Configuration:   targetState.Configuration,
	}

	// trigger request
	body, _ := util.ConvertToYAML(request)

	rsp, err := http.Post(c.URL + "/" + action, "application/x-yaml", strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	response := &Response{}
	err = util.ConvertFromYAML(string(data), response)
	if err != nil {
		return nil, err
	}

	// convert response into currentState
	currentState = &model.CurrentState{
		Domain         : response.Domain,
		Solution       : response.Solution,
		Version        : response.Version,
		Element        : response.Element,
		Cluster        : response.Cluster,
		Instance       : response.Instance,
		Component      : response.Component,
		State          : response.State,
		Configuration  : response.Configuration,
		Endpoint       : response.Endpoint,
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
