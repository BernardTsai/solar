package engine

import (
	"errors"

	ctrl "tsai.eu/solar/controller"
	"tsai.eu/solar/model"
	"tsai.eu/solar/util"
	"tsai.eu/solar/msg"
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
	component, _ := model.GetComponent2(task.Domain, task.Solution, task.Element, task.Cluster)

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
	}

	// determine setup
	setup, _ := model.GetSetup(task.GetDomain(),
	                           task.GetSolution(),
										      	 task.GetVersion(),
										      	 task.GetElement(),
										      	 task.GetCluster(),
										      	 task.GetInstance() )

	// determine the required controller
	controller, err := ctrl.GetController(component.Component)
	if err != nil {
		util.LogError(task.UUID, "ENG", "unknown controller: " + component.Component)
		channel <- model.NewEvent(task.Domain, task.UUID, model.EventTypeTaskFailure, task.UUID, "unknown controller: " + component.Component)
		return
	}

	// execute the required transition
	var status *model.Status

	switch transition {
	case "create":
		status, err = controller.Create(setup)
	case "start":
		status, err = controller.Start(setup)
	case "stop":
		status, err = controller.Stop(setup)
	case "destroy":
		status, err = controller.Destroy(setup)
	case "reset":
		status, err = controller.Reset(setup)
	case "configure":
		status, err = controller.Configure(setup)
	default:
		util.LogError(task.UUID, "ENG", "invalid transition")
		channel <- model.NewEvent(task.Domain, task.UUID, model.EventTypeTaskFailure, task.UUID, "invalid transition")
		return
	}

	// update status
	if status != nil {
		// remember current state
		currentState := instance.State

		// update status
		model.SetStatus(status)

		// notify if instance state has changed
		if instance.State != currentState {
			msg.Notify( "Instance", status.Domain + "/" + status.Element + "/" + status.Cluster + "/" + status.Instance + "/" + instance.State)
		}
	}

	// check for errors and reexecute the task until the desired state has been reached
	if err != nil {
		util.LogError(task.UUID, "ENG", "controller has reported an error:\n" + err.Error())
		channel <- model.NewEvent(task.Domain, task.UUID, model.EventTypeTaskFailure, task.UUID, err.Error())
	} else {
		channel <- model.NewEvent(task.Domain, task.UUID, model.EventTypeTaskExecution, task.UUID, "")
	}
}

//------------------------------------------------------------------------------
