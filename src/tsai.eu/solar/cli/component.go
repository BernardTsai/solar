package cli

import (
	ishell "gopkg.in/abiosoft/ishell.v2"
	"tsai.eu/solar/model"
	"tsai.eu/solar/util"
)

//------------------------------------------------------------------------------

// ComponentVersion refers to a specific version of a component in a domain.
type ComponentVersion struct {
	Domain string
	Component string
	Version string
}

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
		// set default filter criteria
		domainName := ""
		componentName := ""
		versionName := ""

		// initialise result
		var result []ComponentVersion

		// check availability of arguments
		if len(context.Args) < 2 || 4 < len(context.Args) {
			ComponentUsage(true, context)
			return
		}
		if len(context.Args) >= 2 {
			domainName = context.Args[1]
		}
		if len(context.Args) >= 3 {
			versionName = context.Args[2]
		}
		if len(context.Args) == 4 {
			versionName = context.Args[3]
		}

		// loop over all domains
		for dName, d in range m.Domains {
			// filter by domain
			if domainName == "" || domainName == dName {

				// loop over all components
				for cName, c in range d.Components {
					// filter by component
					if componentName == "" || componentName == cName {

						// loop over all versions
						for vName, v in range c.Components {
							// filter by component
							if componentName == "" || componentName == cName {






							}
						}
					}
				}
			}
		}

		// execute command
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		components, _ := domain.ListComponents()
		result, err := util.ConvertToJSON(components)
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
		component, _ := model.NewComponent("", "")

		// load component
		err = component.Load(context.Args[2])

		if err != nil {
			handleResult(context, err, "component could not be loaded", "")
		}

		// add component to domain
		err = domain.AddComponent(component)
		handleResult(context, err, "unable to load component", "component has been loaded")
	case _get:
		// check availability of arguments
		if len(context.Args) != 3 {
			ComponentUsage(true, context)
			return
		}

		// determine domain
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// determine component
		t, err := domain.GetComponent(context.Args[2])

		if err != nil {
			handleResult(context, err, "component can not be identified", "")
			return
		}

		// execute the command
		result, err := t.Show()
		handleResult(context, err, "component can not be displayed", result)
	case _delete:
		// check availability of arguments
		if len(context.Args) != 3 {
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
		_, err = d.GetComponent(context.Args[2])

		if err != nil {
			handleResult(context, err, "component can not be identified", "")
			return
		}

		// execute command
		err = d.DeleteComponent(context.Args[2])
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
