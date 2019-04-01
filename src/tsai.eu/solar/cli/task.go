package cli

import (
	"strconv"

	ishell "gopkg.in/abiosoft/ishell.v2"
	"tsai.eu/solar/model"
	"tsai.eu/solar/util"
	"tsai.eu/solar/engine"
)

//------------------------------------------------------------------------------

// TaskCommand executes the task related subcommands
func TaskCommand(context *ishell.Context, m *model.Model) {
	// check if the action has been defined
	if len(context.Args) < 1 {
		TaskUsage(true, context)
		return
	}

	// determine the required action
	action := context.Args[0]

	// handle required action
	switch action {
	case "?":
		TaskUsage(true, context)
	case _list:
		// check availability of arguments
		if len(context.Args) < 2 || 6 < len(context.Args) {
			TaskUsage(true, context)
			return
		}

		// set domain name filter
		domainName := context.Args[1]

		// set solution name filter
		solutionName := ""
		if len(context.Args) >= 3 {
			solutionName = context.Args[2]
		}

		// set element name filter
		elementName := ""
		if len(context.Args) >= 4 {
			elementName = context.Args[3]
		}

		// set cluster name filter
		clusterName := ""
		if len(context.Args) >= 5 {
			clusterName = context.Args[4]
		}

		// set instance name filter
		instanceName := ""
		if len(context.Args) >= 6 {
			instanceName = context.Args[5]
		}

		// determine domain
		domain, err := model.GetDomain(domainName)
		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// determine list of solution names
		tasks := []string{}

		tNames, _ := domain.ListTasks()
		for _, tName := range tNames {
			task, _ := domain.GetTask(tName)

			if (solutionName == task.Solution) &&
				 (elementName  == task.Element)  &&
				 (clusterName  == task.Cluster)  &&
				 (instanceName == task.Instance) {
			  tasks = append(tasks, tName)
			}
		}

		result, err := util.ConvertToYAML(tasks)
		handleResult(context, err, "tasks could not be listed", result)
	case _get:
		// check availability of arguments
		if len(context.Args) < 3 || 4 < len(context.Args){
			TaskUsage(true, context)
			return
		}

		// determine task
		task, err := model.GetTask(context.Args[1], context.Args[2])

		if err != nil {
			handleResult(context, err, "task can not be identified", "")
			return
		}

		// determine level
		level := 0
		if len(context.Args) == 4 {
			value, err := strconv.Atoi(context.Args[3])
			if err != nil {
				handleResult(context, err, "invalid level (not an integer)", "")
				return
			}

			if value < 0 {
				handleResult(context, err, "invalid level (must not be negative)", "")
				return
			}
			level = value
		}

		// execute the command
		taskinfo := model.NewTaskInfo(task, level)
		result, err := util.ConvertToYAML(taskinfo)
		handleResult(context, err, "task can not be displayed", result)
	case _terminate:
		// check availability of arguments
		if len(context.Args) != 3 {
			TaskUsage(true, context)
			return
		}

		// determine task
		task, err := model.GetTask(context.Args[1], context.Args[2])

		if err != nil {
			handleResult(context, err, "task can not be identified", "")
			return
		}

		// execute the command
		// get event channel
		channel := engine.GetEventChannel()

		// create event
		channel <- model.NewEvent(context.Args[1], task.UUID, model.EventTypeTaskTermination, "", "initial")

		handleResult(context, nil, "task can not be terminated", "")
	case _trace:
		// check availability of arguments
		if len(context.Args) != 3 {
			TaskUsage(true, context)
			return
		}

		// determine task
		task, err := model.GetTask(context.Args[1], context.Args[2])

		if err != nil {
			handleResult(context, err, "task can not be identified", "")
			return
		}

		// construct the trace object
		trace := model.NewTrace(task)

		// finished
		result, _ := util.ConvertToYAML(trace)
		handleResult(context, nil, "trace could not be created", result)
	default:
		TaskUsage(true, context)
	}
}

//------------------------------------------------------------------------------

// TaskUsage describes how to make use of the subcommand
func TaskUsage(header bool, context *ishell.Context) {
	info := ""
	if header {
		info = _usage
	}
	info += "  task list <domain> <solution> <element> <cluster> <instance>\n"
	info += "       get <domain> <task> <level>\n"
	info += "       terminate <domain> <task>\n"
	info += "       trace <domain> <task>\n"

  writeInfo(context, info)
}

//------------------------------------------------------------------------------
