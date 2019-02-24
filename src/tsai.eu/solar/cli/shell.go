package cli

import (
	"os"
	"fmt"
	"strings"

	ishell "gopkg.in/abiosoft/ishell.v2"
	"tsai.eu/solar/model"
	"tsai.eu/solar/util"
)

//------------------------------------------------------------------------------

var output string

//------------------------------------------------------------------------------

// Run executes the main shell event loop
func Run(m *model.Model) {
	// create new shell which by default includes 'exit', 'help' and
	// 'clear' commands
	shell := ishell.New()

	// register a function for the "model" command.
	shell.AddCmd(&ishell.Cmd{
		Name: "usage",
		Help: "usage command",
		Func: func(c *ishell.Context) {
			OutputUsage(true, c)
			ModelUsage(false, c)
			DomainUsage(false, c)
			ComponentUsage(false, c)
			ArchitectureUsage(false, c)
			SolutionUsage(false, c)
		},
	})

	// register a function for the "output" command.
	shell.AddCmd(&ishell.Cmd{
		Name: "output",
		Help: "output commands",
		Func: func(c *ishell.Context) { OutputCommand(c, m) },
	})

	// register a function for the "model" command.
	shell.AddCmd(&ishell.Cmd{
		Name: "model",
		Help: "model commands",
		Func: func(c *ishell.Context) { ModelCommand(c, m) },
	})

	// register a function for the "domain" command.
	shell.AddCmd(&ishell.Cmd{
		Name: "domain",
		Help: "domain commands",
		Func: func(c *ishell.Context) { DomainCommand(c, m) },
	})

	// register a function for the "component" command.
	shell.AddCmd(&ishell.Cmd{
		Name: "component",
		Help: "component commands",
		Func: func(c *ishell.Context) { ComponentCommand(c, m) },
	})

	// register a function for the "architecture" command.
	shell.AddCmd(&ishell.Cmd{
		Name: "architecture",
		Help: "architecture commands",
		Func: func(c *ishell.Context) { ArchitectureCommand(c, m) },
	})

	// register a function for the "solution" command.
	shell.AddCmd(&ishell.Cmd{
		Name: "solution",
		Help: "solution commands",
		Func: func(c *ishell.Context) { SolutionCommand(c, m) },
	})

	// register a function for "#" command.
	shell.AddCmd(&ishell.Cmd{
		Name: "comment",
		Help: "comment",
		Func: func(c *ishell.Context) {
			c.Println(strings.Join(c.Args, " "))
		},
	})

	// run shell
	shell.Run()
}

//------------------------------------------------------------------------------

// handleResult reports error information if present or display success message
func handleResult(context *ishell.Context, err error, fail string, success string) {
	if err != nil {
		if util.Debug() {

			info := fmt.Sprintf("%s\n%+v\n	", fail, err)
			writeError(context, info)
		} else {
			info := fmt.Sprintf("%s\n ", fail)
			writeError(context, info)
		}
	} else {
		writeOutput(context, success)
	}
}

//------------------------------------------------------------------------------

// setOutput defines the name of the output file.
func setOutput(filename string) {
	output = filename
}

//------------------------------------------------------------------------------

// writeInfo prints information to the console.
func writeInfo(context *ishell.Context, info string) {
	context.Println(info)
}

//------------------------------------------------------------------------------

// writeError prints error information to the console.
func writeError(context *ishell.Context, info string) {
	context.Println(info)
}

//------------------------------------------------------------------------------

// writeOutput writes the provided information to the console or a file.
func writeOutput(context *ishell.Context, info string) error {
	// check which channel to use
	if output == "" {
		context.Println(info)
	} else {
		f, err := os.OpenFile(output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
			fmt.Println(err)
      return err
    }
    defer f.Close()

    _, err = f.WriteString(info)
    if err != nil {
			fmt.Println(err)
      return err
    }
	}

	return nil
}

//------------------------------------------------------------------------------
