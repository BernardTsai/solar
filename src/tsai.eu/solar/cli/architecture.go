package cli

import (
	ishell "gopkg.in/abiosoft/ishell.v2"
	"tsai.eu/solar/engine"
	"tsai.eu/solar/model"
	"tsai.eu/solar/util"
)

//------------------------------------------------------------------------------

// ArchitectureCommand executes the architecture related subcommands
func ArchitectureCommand(context *ishell.Context, m *model.Model) {
	// check if the action has been defined
	if len(context.Args) < 1 {
		ArchitectureUsage(true, context)
		return
	}

	// determine the required action
	action := context.Args[0]

	// handle required action
	switch action {
	case "?":
		ArchitectureUsage(true, context)
	case _list:
		// check availability of arguments
		if len(context.Args) < 2 || 4 < len(context.Args) {
			ArchitectureUsage(true, context)
			return
		}

		// set domain name filter
		domainName := context.Args[1]

		// set architecture name filter
		architectureName := ""
		if len(context.Args) >= 3 {
			architectureName = context.Args[2]
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

		// determine list of architecture names
		architectures := []string{}

		aNames, _ := domain.ListArchitectures()
		for _, aName := range aNames {
			architecture, _ := domain.GetArchitecture(aName)

			if (architectureName == "" || architectureName == architecture.Architecture) &&
			   (versionName   == ""    || versionName   == architecture.Version) {
				architectures = append(architectures, aName)
			}
		}

		result, err := util.ConvertToYAML(architectures)
		handleResult(context, err, "architectures could not be listed", result)
	case _set:
		// check availability of arguments
		if len(context.Args) != 3 {
			ArchitectureUsage(true, context)
			return
		}

		// get domain
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// create new architecture
		architecture, _ := model.NewArchitecture("", "", "")

		// load architecture
		err = architecture.Load(context.Args[2])

		if err != nil {
			handleResult(context, err, "architecture could not be loaded", "")
			return
		}

		// add architecture to domain
		err = domain.AddArchitecture(architecture)
		handleResult(context, err, "architecture could not be loaded", "")
	case _get:
		// check availability of arguments
		if len(context.Args) != 3 {
			ArchitectureUsage(true, context)
			return
		}

		// determine architecture
		architecture, err := model.GetArchitecture(context.Args[1], context.Args[2])

		if err != nil {
			handleResult(context, err, "architecture can not be identified", "")
			return
		}

		// execute the command
		result, err := architecture.Show()
		handleResult(context, err, "component can not be displayed", result)
	case _delete:
		// check availability of arguments
		if len(context.Args) != 3 {
			ArchitectureUsage(true, context)
			return
		}

		// determine domain
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// determine architecture
		_, err = domain.GetArchitecture(context.Args[2])

		if err != nil {
			handleResult(context, err, "architecture can not be identified", "")
			return
		}

		// execute command
		err = domain.DeleteArchitecture(context.Args[2])
		handleResult(context, err, "architecture can not be deleted", "architecture has been deleted")
	case _deploy:
		// check availability of arguments
		if len(context.Args) != 3 {
			ArchitectureUsage(true, context)
			return
		}

		// determine domain
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// determine architecture
		architecture, err := domain.GetArchitecture(context.Args[2])

		if err != nil {
			handleResult(context, err, "architecture can not be identified", "")
			return
		}

		// determine solution (create new solution if not found)
		solution, err := domain.GetSolution(architecture.Architecture)
		if err != nil {
			solution, _ = model.NewSolution(architecture.Architecture, architecture.Version, "")

			domain.AddSolution(solution)
		}

		// update the target state of the solution
		if err = solution.Update(domain.Name, architecture); err != nil {
			handleResult(context, err, "unable to create or update the solution", "")
			return
		}

		// create task and start it by signalling an event
		task, err := engine.NewSolutionTask(domain.Name, "", solution)
		if err != nil {
			handleResult(context, err, "task can not be created", "")
			return
		}

		// get event channel
		channel := engine.GetEventChannel()

		// create event
		channel <- model.NewEvent(domain.Name, task.UUID, model.EventTypeTaskExecution, "", "initial")

		handleResult(context, nil, "architecture can not be executed", task.UUID)
	default:
		ArchitectureUsage(true, context)
	}
}

//------------------------------------------------------------------------------

// ArchitectureUsage describes how to make use of the subcommand
func ArchitectureUsage(header bool, context *ishell.Context) {
	info := ""
	if header {
		info = _usage
	}
	info += "  architecture list <domain> <architecture> <version>\n"
	info += "               set <domain> <filename>\n"
	info += "               get <domain> <architecture>\n"
	info += "               delete <domain> <architecture>\n"
	info += "               deploy <domain> <architecture>\n"

  writeInfo(context, info)
}

//------------------------------------------------------------------------------
