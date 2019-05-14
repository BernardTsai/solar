package engine

import (
	"errors"

	"tsai.eu/solar/model"
	"tsai.eu/solar/util"
)

//------------------------------------------------------------------------------

// NewInstanceTask creates a new instance task
func NewInstanceTask(domain string, parent string, solution string, version string, element string, cluster string, instance string, state string) (model.Task, error) {
	var task model.Task

	// TODO: check parameters if context exists
	task.Type     = "Instance"
	task.Domain   = domain
	task.Solution = solution
	task.Version  = version
	task.Element  = element
	task.Cluster  = cluster
	task.Instance = instance
	task.State    = state
	task.Action   = ""
	task.UUID     = util.UUID()
	task.Parent   = parent
	task.Status   = model.TaskStatusInitial
	task.Phase    = 0
	task.Subtasks = []string{}

	// add handlers
	task.SetExecute(ExecuteInstanceTask)
	task.SetTerminate(TerminateTask)
	task.SetFailed(FailedTask)
	task.SetTimeout(TimeoutTask)
	task.SetCompleted(CompletedTask)

	// get domain
	d, err := model.GetModel().GetDomain(domain)
	if err != nil {
		util.LogError(parent, "ENG", "unknown domain")
		return task, errors.New("unknown domain")
	}

	// add task to domain
	err = d.AddTask(&task)
	if err != nil {
		util.LogError(parent, "ENG", "unable to add task")
		return task, err
	}

	// success
	return task, nil
}

//------------------------------------------------------------------------------

// ExecuteInstanceTask is the main task execution routine.
func ExecuteInstanceTask(task *model.Task) {
	// get event channel
	channel := GetEventChannel()

	// check status
	taskStatus := task.GetStatus()

	if taskStatus != model.TaskStatusInitial && taskStatus != model.TaskStatusExecuting {
		return
	}

	// initialize if needed
	if taskStatus == model.TaskStatusInitial {
		// update status
		task.Status = model.TaskStatusExecuting
	}

	// TODO: implement and proper error handling
	// - collect relevant information
	// - determine current state and target state of instance and derive the required transition
	// - trigger controller to execute transition
	// - obtain endpoint
	// - calculate endpoint information
	// - trigger new execution
	// - in case of error trigger failure

	// determine context
	instance, _  := model.GetInstance(task.Domain, task.Solution, task.Element, task.Cluster, task.Instance)

	// update target state of instance
	instance.Target = task.State

	// check if the target state has been reached
	if instance.State == instance.Target {
		channel <- model.NewEvent(task.Domain, task.UUID, model.EventTypeTaskCompletion, task.UUID, "")
		return
	}

	// determine required transition
	transition, err := model.GetTransition(instance.State, task.State)

	if err != nil {
		util.LogError(task.UUID, "ENG", "invalid transition")
		channel <- model.NewEvent(task.Domain, task.UUID, model.EventTypeTaskFailure, task.UUID, "invalid transition")
		return
	}

	// create controller task for the required transition
	subtask, _ := NewControllerTask(task.Domain, task.UUID, task.Solution, task.Version, task.Element, task.Cluster, instance.UUID, task.State, transition )
	task.AddSubtask(&subtask)

	// trigger the task
	channel <- model.NewEvent(task.Domain, subtask.UUID, model.EventTypeTaskExecution, task.UUID, "")
}

//------------------------------------------------------------------------------
