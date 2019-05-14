package engine

import (
	"errors"

	ctrl "tsai.eu/solar/controller"
	"tsai.eu/solar/model"
	"tsai.eu/solar/util"
	"tsai.eu/solar/msg"
)

//------------------------------------------------------------------------------

// NewControllerTask creates a new controller task
func NewControllerTask(domain string, parent string, solution string, version string, element string, cluster string, instance string, state string, action string) (model.Task, error) {
	var task model.Task

	// TODO: check parameters if context exists
	task.Type     = "Controller"
	task.Domain   = domain
	task.Solution = solution
	task.Version  = version
	task.Element  = element
	task.Cluster  = cluster
	task.Instance = instance
	task.State    = state
	task.Action   = action
	task.UUID     = util.UUID()
	task.Parent   = parent
	task.Status   = model.TaskStatusInitial
	task.Phase    = 0
	task.Subtasks = []string{}

	// add handlers
	task.SetExecute(ExecuteControllerTask)
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

// ExecuteControllerTask is the main task execution routine.
func ExecuteControllerTask(task *model.Task) {
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

	// determine desired target state
	targetState, _ := model.GetTargetState(
											task.GetDomain(),
	                    task.GetSolution(),
										  task.GetVersion(),
										  task.GetElement(),
										  task.GetCluster(),
										  task.GetInstance() )

	// determine the required controller for the instance
	instance, _     := model.GetInstance(task.Domain, task.Solution, task.Element, task.Cluster, task.Instance)
	component, _    := model.GetComponent2(task.Domain, task.Solution, task.Element, task.Cluster)
	controller, err := ctrl.GetController(component.Controller)
	if err != nil {
		util.LogError(task.UUID, "ENG", "unknown controller: " + component.Component + ":" + task.GetCluster())
		channel <- model.NewEvent(task.Domain, task.UUID, model.EventTypeTaskFailure, task.UUID, "unknown controller: " + component.Component)
		return
	}

	// execute the required transition
	var currentState *model.CurrentState

	switch task.Action {
	case "create":
		currentState, err = controller.Create(targetState)
	case "start":
		currentState, err = controller.Start(targetState)
	case "stop":
		currentState, err = controller.Stop(targetState)
	case "destroy":
		currentState, err = controller.Destroy(targetState)
	case "reset":
		currentState, err = controller.Reset(targetState)
	case "configure":
		currentState, err = controller.Configure(targetState)
	default:
		util.LogError(task.UUID, "ENG", "invalid transition")
		channel <- model.NewEvent(task.Domain, task.UUID, model.EventTypeTaskFailure, task.UUID, "invalid transition")
		return
	}

	// update status
	if currentState != nil {
		// remember current state
		instanceState := instance.State

		// update state
		model.SetCurrentState(currentState)

		// notify if instance state has changed
		if instance.State != instanceState {
			util.LogInfo(task.UUID, "ENG", "Instance: " + currentState.Domain + "/" + currentState.Element + "/" + currentState.Cluster + "/" + currentState.Instance + " has new state:" + instance.State)
			msg.Notify( "Instance", currentState.Domain + "/" + currentState.Element + "/" + currentState.Cluster + "/" + currentState.Instance + "/" + instance.State)
		}
	}

	// check for errors and reexecute the task until the desired state has been reached
	if err != nil {
		util.LogError(task.UUID, "ENG", "controller has reported an error:\n" + err.Error())
		channel <- model.NewEvent(task.Domain, task.UUID, model.EventTypeTaskFailure, task.UUID, err.Error())
		return
	}

	// success
	channel <- model.NewEvent(task.Domain, task.UUID, model.EventTypeTaskCompletion, task.UUID, "")
}

//------------------------------------------------------------------------------
