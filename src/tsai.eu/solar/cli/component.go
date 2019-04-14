package cli

import (
	ishell "gopkg.in/abiosoft/ishell.v2"
	"tsai.eu/solar/model"
	"tsai.eu/solar/util"
)

//------------------------------------------------------------------------------

// ComponentCommand executes the component related subcommands
func ComponentCommand(context *ishell.Context, m *model.Model) {
	// check if the action has been defined
	if len(context.Args) < 1 {
		ComponentUsage(true, context)
		return
	}

	// determine the required action
	action := context.Args[0]

	// handle required action
	switch action {
	case "?":
		ComponentUsage(true, context)
	case _list:
		// check availability of arguments
		if len(context.Args) < 2 || 4 < len(context.Args) {
			ComponentUsage(true, context)
			return
		}

		// set domain name filter
		domainName := context.Args[1]

		// set component name filter
		componentName := ""
		if len(context.Args) >= 3 {
			componentName = context.Args[2]
		}

		// set version filter
		versionName   := ""
		if len(context.Args) == 4 {
			versionName = context.Args[3]
		}

		// determine domain
		domain, err := model.GetDomain(domainName)
		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// determine list of component names
		components := []string{}

		cNameVersions, _ := domain.ListComponents()
		for _, cNameVersion := range cNameVersions {
			if (componentName == "" || componentName == cNameVersion[0]) &&
			   (versionName   == "" || versionName   == cNameVersion[1]) {
				components = append(components, cNameVersion[0] + " - " + cNameVersion[1])
			}
		}

		result, err := util.ConvertToYAML(components)
		handleResult(context, err, "components could not be listed", result)
	case _set:
		// check availability of arguments
		if len(context.Args) != 3 {
			ComponentUsage(true, context)
			return
		}

		// get domain
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// create new component
		component, _ := model.NewComponent("A", "B", "C")

		// load component
		err = component.Load(context.Args[2])

		if err != nil {
			handleResult(context, err, "component could not be loaded", "")
			return
		}

		// add component to domain
		err = domain.AddComponent(component)
		handleResult(context, err, "unable to load component", "")
	case _get:
		// check availability of arguments
		if len(context.Args) != 4 {
			ComponentUsage(true, context)
			return
		}

		// determine component
		component, err := model.GetComponent(context.Args[1], context.Args[2], context.Args[3])

		if err != nil {
			handleResult(context, err, "component can not be identified", "")
			return
		}

		// execute the command
		result, err := component.Show()
		handleResult(context, err, "component can not be displayed", result)
	case _delete:
		// check availability of arguments
		if len(context.Args) != 4 {
			ComponentUsage(true, context)
			return
		}

		// determine domain
		d, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// determine component
		_, err = d.GetComponent(context.Args[2], context.Args[3])

		if err != nil {
			handleResult(context, err, "component can not be identified", "")
			return
		}

		// execute command
		err = d.DeleteComponent(context.Args[2], context.Args[3])
		handleResult(context, err, "component can not be deleted", "component has been deleted")
	default:
		ComponentUsage(true, context)
	}
}

//------------------------------------------------------------------------------

// ComponentUsage describes how to make use of the subcommand
func ComponentUsage(header bool, context *ishell.Context) {
	info := ""
	if header {
		info = _usage
	}
	info += "  component list <domain> <component> <version>\n"
	info += "            set <domain> <filename>\n"
	info += "            get <domain> <component> <version>\n"
	info += "            delete <domain> <component> <version>\n"

  writeInfo(context, info)
}

//------------------------------------------------------------------------------
