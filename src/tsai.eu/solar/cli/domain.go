package cli

import (
	ishell "gopkg.in/abiosoft/ishell.v2"
	"tsai.eu/solar/model"
	"tsai.eu/solar/util"
)

//------------------------------------------------------------------------------

// DomainCommand executes the domain related subcommands
func DomainCommand(context *ishell.Context, m *model.Model) {
	// check if the action has been defined
	if len(context.Args) < 1 {
		DomainUsage(true, context)
		return
	}

	// determine the required action
	action := context.Args[0]

	// handle required action
	switch action {
	case "?":
		DomainUsage(true, context)
	case _list:
		domains, _ := m.ListDomains()
		result, err := util.ConvertToYAML(domains)
		handleResult(context, err, "domains could not be listed", result)
	case _create:
		// check availability of arguments
		if len(context.Args) < 2 {
			DomainUsage(true, context)
			return
		}

		// execute command
		d, _ := model.NewDomain(context.Args[1])
		err := m.AddDomain(d)
    handleResult(context, err, "unable to create domain", "")
	case _delete:
		// check availability of arguments
		if len(context.Args) < 2 {
			DomainUsage(true, context)
			return
		}

		// execute command
		err := m.DeleteDomain(context.Args[1])
		handleResult(context, err, "domain can not be deleted", "")
	case _set:
		// check availability of arguments
		if len(context.Args) < 2 {
			DomainUsage(true, context)
			return
		}

		// execute command
		// create new domain
		d, _ := model.NewDomain("dummy")

		// load domain information
		err := d.Load(context.Args[1])
		if err != nil {
			handleResult(context, err, "domain could not be loaded", "")
			return
		}

		// add domain to model
		err = m.AddDomain(d)
		handleResult(context, err, "domain could not be loaded", "")
	case _get:
		// check availability of arguments
		if len(context.Args) < 2 {
			DomainUsage(true, context)
			return
		}

		// execute command
		d, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain is not known", "")
			return
		}

		// display domain
		result, err := d.Show()
		handleResult(context, err, "domain can not be displayed", result)
	case _reset:
		// check availability of arguments
		if len(context.Args) < 2 {
			DomainUsage(true, context)
			return
		}

		// execute delete command
		m.DeleteDomain(context.Args[1])

		// execute create ommand
		d, _ := model.NewDomain(context.Args[1])
		err := m.AddDomain(d)
		handleResult(context, err, "unable to reset domain", "")
	default:
		DomainUsage(true, context)
	}
}

//------------------------------------------------------------------------------

// DomainUsage describes how to make use of the domain subcommand
func DomainUsage(header bool, context *ishell.Context) {
	info := ""
	if header {
		info = _usage
	}
	info += "  domain list\n"
	info += "         create <domain>\n"
	info += "         delete <domain>\n"
	info += "         set <filename>\n"
	info += "         get <domain>\n"
	info += "         reset <domain>\n"

  writeInfo(context, info)
}

//------------------------------------------------------------------------------
