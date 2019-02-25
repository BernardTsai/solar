package cli

import (
	ishell "gopkg.in/abiosoft/ishell.v2"
	"tsai.eu/solar/model"
)

//------------------------------------------------------------------------------

// ModelCommand executes the model related subcommands
func ModelCommand(context *ishell.Context, m *model.Model) {
	// check if the action has been defined
	if len(context.Args) < 1 {
		ModelUsage(true, context)
		return
	}

	// determine the required action
	action := context.Args[0]

	// handle required action
	switch action {
	case "?":
		ModelUsage(true, context)
	case _set:
		err := m.Load(context.Args[1])
		handleResult(context, err, "model could not be loaded", "")
	case _get:
		result, err := m.Show()
		handleResult(context, err, "model can not be displayed", result)
	case _reset:
		err := m.Reset()
		handleResult(context, err, "model could not be reset", "")
	default:
		ModelUsage(true, context)
	}
}

//------------------------------------------------------------------------------

// ModelUsage describes how to make use of the model subcommand
func ModelUsage(header bool, context *ishell.Context) {
	info := ""
	if header {
		info = _usage
	}
	info += "  model set <filename>\n"
	info += "        get\n"
	info += "        reset\n"

  writeInfo(context, info)
}

//------------------------------------------------------------------------------
