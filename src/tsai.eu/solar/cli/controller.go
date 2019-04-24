package cli

import (
	ishell "gopkg.in/abiosoft/ishell.v2"
	"tsai.eu/solar/model"
	"tsai.eu/solar/util"
)

//------------------------------------------------------------------------------

// ControllerCommand executes the domain related subcommands
func ControllerCommand(context *ishell.Context, m *model.Model) {
	// check if the action has been defined
	if len(context.Args) < 1 {
		ControllerUsage(true, context)
		return
	}

	// determine the required action
	action := context.Args[0]

	// handle required action
	switch action {
	case "?":
		ControllerUsage(true, context)
	case _list:
    // check availability of arguments
		if len(context.Args) != 2 {
			ControllerUsage(true, context)
			return
		}

		// set domain name filter
		domainName := context.Args[1]

		// determine domain
		domain, err := model.GetDomain(domainName)
		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// determine list of solution names
		controllers := []*model.Controller{}

		cNameVersions, _ := domain.ListControllers()
		for _, cNameVersion := range cNameVersions {
			controller, _ := domain.GetController(cNameVersion[0], cNameVersion[1])

			controllers = append(controllers, controller)
		}

		result, err := util.ConvertToYAML(controllers)
		handleResult(context, err, "solutions could not be listed", result)
	case _delete:
    // check availability of arguments
		if len(context.Args) != 4 {
			ControllerUsage(true, context)
			return
		}

		// determine domain
		d, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// determine controller
		_, err = d.GetController(context.Args[2], context.Args[3])

		if err != nil {
			handleResult(context, err, "controller can not be identified", "")
			return
		}

		// execute command
		err = d.DeleteController(context.Args[2], context.Args[3])
		handleResult(context, err, "controller can not be deleted", "controller has been deleted")
	case _set:
    // check availability of arguments
		if len(context.Args) != 3 {
			ControllerUsage(true, context)
			return
		}

		// get domain
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// create new controller
		controller, _ := model.NewController("", "")

		// load controller
		err = controller.Load(context.Args[2])

		if err != nil {
			handleResult(context, err, "controller could not be loaded", "")
		}

		// add controller to domain
		err = domain.AddController(controller)
		handleResult(context, err, "unable to load controller", "controller has been loaded")
	case _get:
    // check availability of arguments
		if len(context.Args) != 4 {
			ControllerUsage(true, context)
			return
		}

		// determine controller
		controller, err := model.GetController(context.Args[1], context.Args[2], context.Args[3])

		if err != nil {
			handleResult(context, err, "controller can not be identified", "")
			return
		}

		// execute the command
		result, err := controller.Show()
		handleResult(context, err, "controller can not be displayed", result)
	default:
		ControllerUsage(true, context)
	}
}

//------------------------------------------------------------------------------

// ControllerUsage describes how to make use of the domain subcommand
func ControllerUsage(header bool, context *ishell.Context) {
	info := ""
	if header {
		info = _usage
	}
	info += "  controller list <domain>\n"
	info += "             set <domain> <filename>\n"
	info += "             get <domain> <controller> <version>\n"
	info += "             delete <domain> <controller> <version>\n"

  writeInfo(context, info)
}

//------------------------------------------------------------------------------
