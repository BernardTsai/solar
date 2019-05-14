package controller

import (
  "context"
  "time"
  "net/http"
  "strings"
  "strconv"
  "io/ioutil"

  "tsai.eu/solar/model"
  "tsai.eu/solar/util"
)

//------------------------------------------------------------------------------

// Manager manages controllers
type Manager struct {
  Ticker  *time.Ticker           // ticker
  Active   bool                  // indicates if the monitoring loop should be active
}

//------------------------------------------------------------------------------

// Start creates a process to manage controllers.
func Start(ctx context.Context) (*Manager) {
	// create the manager
	manager := Manager{
    Ticker:  time.NewTicker(10000 * time.Millisecond),
    Active:  false,
	}

	// start the manager
	go manager.Run(ctx)
  manager.Start()

  util.LogInfo("main", "CTL", "controller active")

  // success
  return &manager
}

//------------------------------------------------------------------------------

// Run starts the manager loop validating the status of controllers
func (m *Manager) Run(ctx context.Context) {
  // loop while manager needs to be active
  for {
    select {
    // check if context has expired
    case <-ctx.Done():
      util.LogInfo("main", "CTL", "controller initial")
      m.Ticker.Stop()
      return
    // wait for next tick and monitor solutions
    case <- m.Ticker.C:
      if m.Active {
        checkControllers()
      }
    }
  }
}

//------------------------------------------------------------------------------

// Start will flag the manager to resume execution
func (m *Manager) Start() {
  m.Active = true
  util.LogInfo("main", "CTL", "controller active")
}

//------------------------------------------------------------------------------

// Stop will flag the manager to pause execution
func (m *Manager) Stop() {
  m.Active = false
  util.LogInfo("main", "CTL", "controller inactive")
}

//------------------------------------------------------------------------------

// checkControllers checks the status of the controllers
func checkControllers() {
  // loop over all domains
  domainNames, _ := model.GetDomains()
  for _, domainName := range domainNames {
    domain, _ := model.GetDomain(domainName)

    // loop over all controllers
    controllerNames, _ := domain.ListControllers()
    for _, controllerNameVersion := range controllerNames {
      controller, _ := domain.GetController(controllerNameVersion[0], controllerNameVersion[1])

      // skip internal controller
      if controller == nil || controller.Image == "" {
        continue
      }

      // check status of controller
      checkController(controller)

      // start controller if required
      if controller.Status != model.ActiveState {
        util.LogInfo("main", "CTL", "Controller: " + controller.Controller + ":" + controller.Version + " is NOT activce")
        util.LogInfo("main", "CTL", "Starting controller: " + controller.Controller + ":" + controller.Version)
        stopController(controller)
        startController(controller)
        if controller.Status != model.ActiveState {
          util.LogError("main", "CTL", "Starting controller: " + controller.Controller + ":" + controller.Version + " has failed")
        } else {
          util.LogInfo("main", "CTL", "Controller: " + controller.Controller + ":" + controller.Version + " is activce")
        }
      }
    } // end of loop over all controllers
  } // end of loop over all domains
}

//------------------------------------------------------------------------------

// checkController checks the status of a controller
func checkController(controller *model.Controller) {
  // create client
  httpc := http.Client{}

  // check if the controller is still responding
  if controller.Status == model.ActiveState {
    response, responseError := httpc.Get(controller.URL)
    if responseError != nil {
      controller.Status     = model.InactiveState
      return
    }
    defer response.Body.Close()

    // read response data
    line, readAllErr := ioutil.ReadAll(response.Body)
    if readAllErr != nil {
      controller.Status     = model.InactiveState
      return
    }

    parts := strings.Split(string(line), ":")

    // check if we have a valid response
    if len(parts) != 3 || parts[0] != "SOLAR" {
      controller.Status     = model.InactiveState
      return
    }

    controller.Controller = parts[1]
    controller.Version    = parts[2]
  }
}

//------------------------------------------------------------------------------

// startController starts a controller
func startController(controller *model.Controller) {
  // determine name and version of image
  parts        := strings.Split(controller.Image, ":")
  if len(parts) != 2 {
    return
  }

  imageName    := parts[0]
  imageVersion := parts[1]

  // get all images
  imagesList, listImagesError := util.ListImages()
  if listImagesError != nil {
    util.LogError("main", "CTL", "Unable to list controller images:\n" + listImagesError.Error())
    return
  }

  // pull image if needed
  foundImage := false
  for _, imageItem := range imagesList {
    for _, value := range imageItem.RepoTags {
      if value == controller.Image {
        foundImage = true
        break
      }
      if foundImage {
        break
      }
    }
  }

  if !foundImage {
    pullError := util.PullImage(imageName, imageVersion)
    if pullError != nil {
      util.LogError("main", "CTL", "Unable to pull controller image: " + controller.Image + "\n" + pullError.Error())
      return
    }
  }

  // start container
  util.LogInfo( "main", "CTL", "Starting controller image: " + controller.Image)
  port, startError := util.StartContainer(imageName, imageVersion)
  if startError != nil {
    util.LogError("main", "CTL", "Unable to start controller: " + controller.Image + "due to:\n" + startError.Error() + "\n")
    return
  }

  // update controller information
  controller.URL    = "http://localhost:" + strconv.Itoa(port)
  controller.Status = model.ActiveState
}

//------------------------------------------------------------------------------

// stopController stops a controller
func stopController(controller *model.Controller) {
  // stop local controllers
  if strings.Contains(controller.URL, "//localhost:") {
    // determine name and version of image
    parts        := strings.Split(controller.Image, ":")
    if len(parts) != 2 {
      util.LogError("main", "CTL", "Unable to stop invalid controller: " + controller.Image)
      return
    }

    imageName    := parts[0]
    imageVersion := parts[1]

    // stop container
    util.LogInfo( "main", "CTL", "Stopping controller: " + controller.Image)
    stopError := util.StopContainer(imageName, imageVersion)
    if stopError != nil {
      util.LogError("main", "CTL", "Unable to stop controller: " + controller.Image + "due to:\n" + stopError.Error() + "\n")
    }

    // update controller information
    controller.URL    = ""
    controller.Status = model.InactiveState
  }
}

//------------------------------------------------------------------------------
