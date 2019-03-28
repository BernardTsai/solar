package engine

import (
	"errors"

	"tsai.eu/solar/model"
	"tsai.eu/solar/util"
)

//------------------------------------------------------------------------------

// NewSolutionTask creates a new task
func NewSolutionTask(domain string, parent string, solution *model.Solution) (model.Task, error) {
	var task model.Task

	// TODO: check parameters if context exists
	task.Type         = "Solution"
	task.Domain       = domain
	task.Solution     = solution.Solution
	task.Version      = solution.Version
	task.Element      = ""
	task.Cluster      = ""
	task.Instance     = ""
	task.State        = ""
	task.UUID         = util.UUID()
	task.Parent       = parent
	task.Status       = model.TaskStatusInitial
	task.Phase        = 0
	task.Subtasks     = []string{}

	// add handlers
	task.SetExecute(ExecuteSolutionTask)
	task.SetTerminate(TerminateTask)
	task.SetFailed(FailedTask)
	task.SetTimeout(TimeoutTask)
	task.SetCompleted(CompletedTask)

	// get domain
	d, err := model.GetDomain(domain)
	if err != nil {
		return task, errors.New("unknown domain")
	}

	// add task to domain
	err = d.AddTask(&task)
	if err != nil {
		return task, err
	}
	// success
	return task, nil
}

//------------------------------------------------------------------------------

// ExecuteSolutionTask is the main task execution routine.
func ExecuteSolutionTask(task *model.Task) {
	// get event channel
	channel := GetEventChannel()

	// check status
	status := task.GetStatus()

	if status != model.TaskStatusInitial && status != model.TaskStatusExecuting {
		return
	}

	if status == model.TaskStatusInitial {
		task.Status = model.TaskStatusExecuting
	}

	// determine context
	solution, _ := model.GetSolution(task.Domain, task.Solution)

	// one by one identify the element which may need to be changed
	elementNames, _ := solution.ListElements()
	for _, elementName := range elementNames {
		element, _ := solution.GetElement(elementName)

		// check if the element needs to be updated
		if !element.OK() {
			// create task to update the element
			subtask, _ := NewElementTask(task.Domain, task.UUID, task.Solution, task.Version, elementName)
			task.AddSubtask(&subtask)

			// trigger the task
			channel <- model.NewEvent(task.Domain, subtask.UUID, model.EventTypeTaskExecution, task.UUID, "")

			// return and wait for next event
			return
		}
	}

	// execution has completed
	channel <- model.NewEvent(task.Domain, task.UUID, model.EventTypeTaskCompletion, task.UUID, "")
}

//------------------------------------------------------------------------------
