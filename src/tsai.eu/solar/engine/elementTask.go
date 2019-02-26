package engine

import (
	"errors"

	"tsai.eu/solar/util"
	"tsai.eu/solar/model"
)

//------------------------------------------------------------------------------

// NewElementTask creates a new task
func NewElementTask(domain string, parent string, solution string, version string, element string) (model.Task, error) {
	var task model.Task

	// TODO: check parameters if context exists
	task.Type         = "ElementTask"
	task.Domain       = domain
	task.Solution     = solution
	task.Version      = version
	task.Element      = element
	task.Cluster      = ""
	task.Instance     = ""
	task.State        = ""
	task.UUID         = util.UUID()
	task.Parent       = parent
	task.Status       = model.TaskStatusInitial
	task.Phase        = 0
	task.Subtasks     = []string{}

	// add handlers
	task.SetExecute(ExecuteElementTask)
	task.SetTerminate(TerminateTask)
	task.SetFailed(FailedTask)
	task.SetTimeout(TimeoutTask)
	task.SetCompleted(CompletedTask)

	// get domain
	d, err := model.GetModel().GetDomain(domain)
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

// ExecuteElementTask is the main task execution routine.
func ExecuteElementTask(task *model.Task) {
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
	element, _ := model.GetElement(task.Domain, task.Solution, task.Element)

	// one by one identify the cluster which may need to be changed
	clusterNames, _ := element.ListClusters()
	for _, clusterName := range clusterNames {
		cluster, _ := element.GetCluster(clusterName)

		// check if the cluster needs to be updated
		if !cluster.OK() {
			// create task to update the cluster
			subtask, _ := NewClusterTask(task.Domain, task.UUID, task.Solution, task.Version,task.Element, clusterName)
			task.AddSubtask(&subtask)

			// trigger the task
			channel <- model.NewEvent(task.Domain, subtask.UUID, model.EventTypeTaskExecution, task.UUID)

			// return and wait for next event
			return
		}
	}

	// execution has completed
	channel <- model.NewEvent(task.Domain, task.UUID, model.EventTypeTaskCompletion, task.UUID)
}

//------------------------------------------------------------------------------
