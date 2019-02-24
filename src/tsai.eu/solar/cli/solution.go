package cli

import (
	ishell "gopkg.in/abiosoft/ishell.v2"
	"tsai.eu/solar/model"
	"tsai.eu/solar/util"
)

//------------------------------------------------------------------------------

// SolutionCommand executes the solution related subcommands
func SolutionCommand(context *ishell.Context, m *model.Model) {
	// check if the action has been defined
	if len(context.Args) < 1 {
		SolutionUsage(true, context)
		return
	}

	// determine the required action
	action := context.Args[0]

	// handle required action
	switch action {
	case "?":
		SolutionUsage(true, context)
	case _list:
		// check availability of arguments
		if len(context.Args) < 2 || 3 < len(context.Args) {
			SolutionUsage(true, context)
			return
		}

		// set domain name filter
		domainName := context.Args[1]

		// set solution name filter
		solutionName := ""
		if len(context.Args) >= 3 {
			solutionName = context.Args[2]
		}

		// determine domain
		domain, err := model.GetDomain(domainName)
		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// determine list of solution names
		solutions := []string{}

		sNames, _ := domain.ListSolutions()
		for _, sName := range sNames {
			solution, _ := domain.GetSolution(sName)

			if (solutionName == "" || solutionName == solution.Solution) {
				solutions = append(solutions, sName)
			}
		}

		result, err := util.ConvertToJSON(solutions)
		handleResult(context, err, "solutions could not be listed", result)
	case _set:
		// check availability of arguments
		if len(context.Args) != 3 {
			SolutionUsage(true, context)
			return
		}

		// get domain
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// create new solution
		solution, _ := model.NewSolution("", "", "")

		// load solution
		err = solution.Load(context.Args[2])

		if err != nil {
			handleResult(context, err, "solution could not be loaded", "")
		}

		// add solution to domain
		err = domain.AddSolution(solution)
		handleResult(context, err, "unable to load solution", "solution has been loaded")
	case _get:
		// check availability of arguments
		if len(context.Args) != 3 {
			SolutionUsage(true, context)
			return
		}

		// determine solution
		solution, err := model.GetSolution(context.Args[1], context.Args[2])

		if err != nil {
			handleResult(context, err, "solution can not be identified", "")
			return
		}

		// execute the command
		result, err := solution.Show()
		handleResult(context, err, "solution can not be displayed", result)
	case _delete:
		// check availability of arguments
		if len(context.Args) != 3 {
			SolutionUsage(true, context)
			return
		}

		// determine domain
		d, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// determine solution
		_, err = d.GetSolution(context.Args[2])

		if err != nil {
			handleResult(context, err, "solution can not be identified", "")
			return
		}

		// execute command
		err = d.DeleteSolution(context.Args[2])
		handleResult(context, err, "solution can not be deleted", "solution has been deleted")
	default:
		SolutionUsage(true, context)
	}
}

//------------------------------------------------------------------------------

// SolutionUsage describes how to make use of the subcommand
func SolutionUsage(header bool, context *ishell.Context) {
	info := ""
	if header {
		info = _usage
	}
	info += "  solution list <domain> <solution>\n"
	info += "           set <domain> <filename>\n"
	info += "           get <domain> <solution> <version>\n"
	info += "           delete <domain> <solution> <version>\n"

  writeInfo(context, info)
}

//------------------------------------------------------------------------------
