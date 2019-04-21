package defaultRestController

import (
	"fmt"
	"strings"
	"testing"
	"context"
	"path/filepath"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"github.com/gorilla/mux"
)

//------------------------------------------------------------------------------

const TESTDATA string = "testdata"
const COMMANDS string = "_commands"

//------------------------------------------------------------------------------

// LoadCommands load a list of commands from a file
func LoadCommands() ([][]string, error) {
	list := [][]string{}

	// load data from file
	path := filepath.Join(TESTDATA, COMMANDS)

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return list, err
	}

	data := (string)(bytes)

	// parse commands
	cmds := strings.Split(data, "\n")

	// parse args and add commands to result
	for _, cmd := range cmds {
		args := strings.Split(cmd, " ")

		list = append(list, args)
	}

	// success
	return list, nil
}

//------------------------------------------------------------------------------

// TestController verifies the controller.
func TestController(t *testing.T) {
	// create and start controller
	controller := NewController()
	controller.Run(context.Background())

	// load commands from a file
	cmds, err := LoadCommands()
	if err != nil {
		fmt.Println(err)
		t.Errorf("%s", err)
	}

	fmt.Println("Executing tests:")
	fmt.Println("Nr. OK Method Body                 URL")
	fmt.Println("------------------------------------------------------------------------------------")

	for index, cmd := range cmds {
		// skip empty command lines
		if cmd[0] == "" {
			continue
		}

		// construct command line string
		cmdline  := strings.Join(cmd, " ")
		status   := cmdline[0:2]
		method   := strings.TrimSpace(cmdline[3:9])
		filename := strings.TrimSpace(cmdline[10:30])
		url      := cmdline[30:]

		// log the command line
		fmt.Printf("%03d %2s %-6s %-20s %s\n", index, status, method, filename, url)

		// proces the request
		if !process(controller.Router, method, url, filename, status) {
			t.Errorf("Error %s %s %s", status, method, url)
		}
	}
}

//------------------------------------------------------------------------------

// loadFile load contents from filename
func loadFile(filename string) string {
  b, _ := ioutil.ReadFile(TESTDATA + "/" + filename)

	return string(b)
}

//------------------------------------------------------------------------------

// process simulates a request and returns the response
func process(router *mux.Router, method string, url string, filename string, status string) bool {
	var req *http.Request

	if filename != "" {
		body := loadFile(filename)

		req = httptest.NewRequest(method, url, strings.NewReader(body))
	} else {
		req = httptest.NewRequest(method, url, nil)
	}

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code == http.StatusOK {
		return (status == "OK")
	}

	return (status != "OK")
}

//------------------------------------------------------------------------------
