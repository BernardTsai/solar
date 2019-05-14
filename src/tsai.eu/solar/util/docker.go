package util

import (
  "net"
  "net/http"
  "encoding/json"
  "context"
  "io/ioutil"
  "strings"
  "strconv"
  "errors"
)

//------------------------------------------------------------------------------

// ImageSummary list details related to an image
type ImageSummary struct {
	ID       string            `yaml:"Id"       json:"Id"`
	Labels   map[string]string `yaml:"Labels"   json:"Labels"`
  RepoTags []string          `yaml:"RepoTags" json:"RepoTags"`
}

// ContainerSummmary has the details of a container
type ContainerSummmary struct {
	ID     string            `yaml:"Id"     json:"Id"`
	Image  string            `yaml:"Image"  json:"Image"`
	Names  []string          `yaml:"Names"  json:"Names"`
  Ports  []ContainerPort   `yaml:"Ports"  json:"Ports"`
	Labels map[string]string `yaml:"Labels" json:"Labels"`
	State  string            `yaml:"State"  json:"State"`
}

// ContainerPort list details related to a container port
type ContainerPort struct {
	IP          string `yaml:"IP"          json:"IP"`
	PrivatePort int    `yaml:"PrivatePort" json:"PrivatePort"`
	PublicPort  int    `yaml:"PublicPort"  json:"PublicPort"`
  Type        string `yaml:"Type"        json:"Type"`
}

type createContainerResponse struct {
  ID       string    `yaml:"Id"       json:"Id"`
  Warnings []string  `yaml:"Warnings" json:"Warnings"`
}

//------------------------------------------------------------------------------

// createContainerTemplate is a template for creating a container
const createContainerTemplate string =
`{
  "Image": "{{IMAGE}}:{{VERSION}}",
  "Tty": true,
  "OpenStdin": true,
  "Labels": {
    "tsai.eu.solar.controller.image":   "{{IMAGE}}",
    "tsai.eu.solar.controller.version": "{{VERSION}}"
  },
  "ExposedPorts": {
    "10000/tcp": {}
  },
  "HostConfig": {
    "PortBindings": {
      "10000/tcp": [
        {
          "HostPort": "{{PORT}}"
        }
      ]
    },
    "RestartPolicy": {
      "Name": "always"
    }
  }
}`

//------------------------------------------------------------------------------

// ListImages lists all available images
func ListImages() ([]ImageSummary, error) {
  result := []ImageSummary{}

  // create client
  httpc := http.Client{
		Transport: &http.Transport{
  		DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
  			return net.Dial("unix", "/var/run/docker.sock")
  		},
		},
  }

  // trigger request
	response, reqErr := httpc.Get("http://v1.39/images/json")
  if reqErr != nil {
    return result, newError("Unable to list images", response, reqErr)
	}

  defer response.Body.Close()

  decodeErr := json.NewDecoder(response.Body).Decode(&result)
  if decodeErr != nil {
    return result, decodeErr
  }

  return result, nil
}

//------------------------------------------------------------------------------

// ListContainers lists all available containers
func ListContainers() ([]ContainerSummmary, error) {
  result := []ContainerSummmary{}

  // create client
  httpc := http.Client{
		Transport: &http.Transport{
  		DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
  			return net.Dial("unix", "/var/run/docker.sock")
  		},
		},
  }

  // trigger request
	response, reqErr := httpc.Get("http://v1.39/containers/json")
	if reqErr != nil {
    return result,  newError("Unable to list containers", response, reqErr)
	}

  defer response.Body.Close()

  decodeErr := json.NewDecoder(response.Body).Decode(&result)
  if decodeErr != nil {
    return result, decodeErr
  }

  return result, nil
}

//------------------------------------------------------------------------------

// PullImage retrieves an image into the local repository
func PullImage(image string, version string) error {
  // create client
  httpc := http.Client{
		Transport: &http.Transport{
  		DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
  			return net.Dial("unix", "/var/run/docker.sock")
  		},
		},
  }

  // trigger request
	response, reqErr := httpc.Post("http://v1.39/images/create?fromImage=" + image + ":" + version, "", nil)
  if reqErr != nil {
    return newError("Unable to pull image '" + image + ":" + version + "'", response, reqErr)
	}

  defer response.Body.Close()

  // wait until the image has been retrieved
  ioutil.ReadAll(response.Body)

  // success
  return nil
}

//------------------------------------------------------------------------------

// StartContainer starts a new container
func StartContainer(image string, version string) (port int, err error) {
  // determine all containers
  containers, err := ListContainers()
  if err != nil{
    return -1, err
  }

  // iterate over all containers
  for _, container := range containers {
    if container.Labels["tsai.eu.solar.controller.image"]   == image &&
       container.Labels["tsai.eu.solar.controller.version"] == version {

      // determine port
      for _, containerPort := range container.Ports {
        return containerPort.PublicPort, nil
      }
    }
  }

  // container has not been found - start a new container
  return startContainer(image, version)
}

//------------------------------------------------------------------------------

// startContainer starts a new container
func startContainer(image string, version string) (int, error) {
  // create client
  httpc := http.Client{
		Transport: &http.Transport{
  		DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
  			return net.Dial("unix", "/var/run/docker.sock")
  		},
		},
  }

  // prepare body of the request
  port, getNewPortError := getNewPort()
  if getNewPortError != nil {
    Print("No new port\n")
    return 0, errors.New("Unable to determine new port:\n" + getNewPortError.Error())
  }

  body := createContainerTemplate
  body =  strings.Replace(body, "{{IMAGE}}",   image,               -1)
  body =  strings.Replace(body, "{{VERSION}}", version,             -1)
  body =  strings.Replace(body, "{{PORT}}",    strconv.Itoa(port),  -1)

  // normalize name
  name := image + "_" + version
  name =  strings.Replace(name, "/", "_", -1)

  // trigger request
	response, reqErr := httpc.Post("http://v1.39/containers/create?name=" + name, "application/json", strings.NewReader(body))
  if reqErr != nil {
    Print("A\n")
    return 0, newError("Unable to create container '" + image + ":" + version + "'", response, reqErr)
	}
  if !strings.HasPrefix(response.Status, "201") {
    Print("%s\n", newError("Failed to create container '" + image + ":" + version + "'", response, reqErr))
    return 0, newError("Failed to create container '" + image + ":" + version + "'", response, reqErr)
  }

  defer response.Body.Close()

  result := createContainerResponse{}

  decodeErr := json.NewDecoder(response.Body).Decode(&result)
  if decodeErr != nil {
    Print("C\n")
    return 0, decodeErr
  }

  // start container
	response2, reqErr2 := httpc.Post("http://v1.39/containers/"+ result.ID + "/start", "application/json", nil)
	if reqErr2 != nil {
    Print("D\n")
    return 0, newError("Unable to start container '" + image + ":" + version + "'", response2, reqErr2)
	}
  if !strings.HasPrefix(response2.Status, "204") {
    Print("E\n")
    return 0, newError("Failed to start container '" + image + ":" + version + "'", response2, reqErr2)
  }

  defer response2.Body.Close()

  ioutil.ReadAll(response2.Body)

  // wait for container
	// response3, reqErr3 := httpc.Post("http://v1.39/containers/"+ result.ID + "/wait?condition=running", "application/json", nil)
	// if reqErr3 != nil {
  //   return 0, reqErr3
	// }
  // if !strings.HasPrefix(response3.Status, "200") {
  //   return 0, errors.New("failed to wait for container: " + response3.Status )
  // }
  // defer response3.Body.Close()

  // ioutil.ReadAll(response3.Body)

  // success
  return port, nil
}

//------------------------------------------------------------------------------

// StopContainer stops a container
func StopContainer(image string, version string) (err error) {
  // create client
  httpc := http.Client{
		Transport: &http.Transport{
  		DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
  			return net.Dial("unix", "/var/run/docker.sock")
  		},
		},
  }

  // normalize name
  name := image + "_" + version
  name =  strings.Replace(name, "/", "_", -1)

  // create request
  request, _ := http.NewRequest("DELETE", "http://v1.39/containers/"+ name + "?force=true", nil)

  // trigger request
	response, reqErr := httpc.Do(request)
  if reqErr != nil {
    return newError("Unable to stop container '" + image + ":" + version + "'", response, reqErr)
	}
  if !strings.HasPrefix(response.Status, "204") {
    return newError("Failed to stop container '" + image + ":" + version + "'", response, reqErr)
  }

  defer response.Body.Close()

  ioutil.ReadAll(response.Body)

  // success
  return nil
}

//------------------------------------------------------------------------------

func getNewPort() (port int, err error) {
  // determine all containers
  containers, err := ListContainers()
  if err != nil{
    return 0, err
  }

  // initialise available ports (starting from 10001)
  // (assumption is that each container has only one public port)
  length    := len(containers) + 1
  available := make([]bool,length, length)

  for i:= 0; i < length; i++ {
    available[i] = true
  }

  // iterate over all containers
  for _, container := range containers {
    _, ok := container.Labels["tsai.eu.solar.controller.image"]
    if ok {
      // remove public port of the container from list of available ports
      for _, containerPort := range container.Ports {
        index := containerPort.PublicPort - 10001
        if 0 <= index && index < length {
          available[index] = false
        }
      }
    }
  }

  // find first free port
  var portNr int
  for portNr = 0; portNr < length; portNr++ {
    if available[portNr] {
      break
    }
  }

  return 10001 + portNr, nil
}

//------------------------------------------------------------------------------

// newError derives a new error from an error and a http response
func newError(message string, response *http.Response, err error) error {
  var txt string

  // construct error message
  if message != "" {
    txt = message + ": "
  } else {
    txt = "Error: "
  }

  // add error message
  if err != nil {
    txt = txt + err.Error()
  }

  // add line break
  txt = txt + "\n"

  // add response body
  if response != nil && response.Body != nil {
    txt = txt + response.Status + "\n"

    defer response.Body.Close()

    bodyBytes, readAllError := ioutil.ReadAll(response.Body)

    if readAllError != nil {
      txt = txt + string(bodyBytes)
    }
  }

  // return new error
  return errors.New(txt)
}

//------------------------------------------------------------------------------
