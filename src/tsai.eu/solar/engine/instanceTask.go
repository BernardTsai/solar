package engine

import (
	"errors"

	ctrl "tsai.eu/solar/controller"
	"tsai.eu/solar/model"
	"tsai.eu/solar/util"
)

//------------------------------------------------------------------------------

// NewInstanceTask creates a new instance task
func NewInstanceTask(domain string, parent string, solution string, version string, element string, cluster string, instance string, state string) (model.Task, error) {
	var task model.Task

	// TODO: check parameters if context exists
	task.Type     = "InstanceTask"
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

// ExecuteInstanceTask is the main task execution routine.
func ExecuteInstanceTask(task *model.Task) {
	// get event channel
	channel := GetEventChannel()

	// check status
	status := task.GetStatus()

	if status != model.TaskStatusInitial && status != model.TaskStatusExecuting {
		return
	}

	// initialize if needed
	if status == model.TaskStatusInitial {
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

	// check if the target state has been reached
	if instance.State == task.State {
		channel <- model.NewEvent(task.Domain, task.UUID, model.EventTypeTaskCompletion, task.UUID)
	}

	// determine required transition
	transition, err := model.GetTransition(instance.State, task.State)
	if err != nil {
		channel <- model.NewEvent(task.Domain, task.UUID, model.EventTypeTaskFailure, task.UUID)
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
		channel <- model.NewEvent(task.Domain, task.UUID, model.EventTypeTaskFailure, task.UUID)
	}

	// execute the required transition
	switch transition {
	case "create":
		_, err = controller.Create(setup)
	case "start":
		_, err = controller.Start(setup)
	case "stop":
		_, err = controller.Stop(setup)
	case "destroy":
		_, err = controller.Destroy(setup)
	case "reset":
		_, err = controller.Reset(setup)
	case "configure":
		_, err = controller.Configure(setup)
	}

	// check for errors and reexecute the task until the desired state has been reached
	if err != nil {
		channel <- model.NewEvent(task.Domain, task.UUID, model.EventTypeTaskFailure, task.UUID)
	} else {
		channel <- model.NewEvent(task.Domain, task.UUID, model.EventTypeTaskExecution, task.UUID)
	}
}

//------------------------------------------------------------------------------
