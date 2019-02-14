package cli

import (
	ishell "gopkg.in/abiosoft/ishell.v2"
  "tsai.eu/solar/model"
)

//------------------------------------------------------------------------------

// OutputCommand defines which output channel to use.
func OutputCommand(context *ishell.Context, m *model.Model) {
	// check if the action has been defined
	if len(context.Args) < 1 {
		OutputUsage(true, context)
		return
	}

	// determine the required action
	filename := context.Args[0]

	// handle required action
	switch filename {
	case "?":
		OutputUsage(true, context)
	case "off":
    setOutput("")
	default:
		setOutput(filename)
	}
}

//------------------------------------------------------------------------------

// OutputUsage describes how to make use of the output subcommand
func OutputUsage(header bool, context *ishell.Context) {
  info := ""
	if header {
		info = "usage:\n"
	}
	info += "  output off\n"
	info += "         <filename>\n"

  writeInfo(context, info)
}

//------------------------------------------------------------------------------
