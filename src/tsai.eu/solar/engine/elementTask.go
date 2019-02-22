package engine

import (
	"errors"

	"github.com/google/uuid"
	"tsai.eu/solar/model"
)

//------------------------------------------------------------------------------

// ElementSetup captures all required configurations for an element.
type ElementSetup struct {
	Element  string
	Clusters map[string]ClusterSetup
}

// ClusterSetup captures all required configurations for a cluster of an element.
type ClusterSetup struct {
	Cluster string
	States  map[string]StateSetup
}

// StateSetup captures the sizing of a cluste of an element with a specific state.
type StateSetup struct {
	State     string
	Instances map[string]string
}

//------------------------------------------------------------------------------

func determineCurrentSetup(domain string, element string) ElementSetup {
	// create ElementSetup
	elementSetup := ElementSetup{
		Element:  element,
		Clusters: map[string]ClusterSetup{},
	}

	// loop over all instances of a component/service
	d, _ := model.GetModel().GetDomain(domain) // domain
	c, _ := d.GetElement(element)              // component
	l, _ := c.ListInstances()                  // list of instances
	for n := range l {
		u := l[n]                // uuid
		i, _ := c.GetInstance(u) // instance

		// check if version exists
		versionSetup, found := serviceSetup.Versions[i.Version]
		if !found {
			versionSetup = VersionSetup{
				Version: i.Version,
				States:  map[string]StateSetup{},
			}
		}

		// check if state exists
		stateSetup, found := versionSetup.States[i.State]
		if !found {
			stateSetup = StateSetup{
				State:     i.State,
				Instances: map[string]string{},
			}
		}

		// add instance
		stateSetup.Instances[i.UUID] = i.UUID
	}

	// success
	return serviceSetup
}

//------------------------------------------------------------------------------

func determineTargetSetup(domain string, architecture string, element string) ServiceSetup {
	// create ServiceSetup
	serviceSetup := ServiceSetup{
		Name:     element,
		Versions: map[string]VersionSetup{},
	}

	// loop over all instances of a component/service
	d, _ := model.GetModel().GetDomain(domain) // domain
	a, _ := d.GetArchitecture(architecture)    // architecture
	e, _ := a.GetElement(element)              // service
	l, _ := e.ListClusters()                   // list of clusters
	for i := range l {
		n := l[i]               // cluster name
		t, _ := e.GetCluster(n) // cluster

		// check if version exists
		versionSetup, found := serviceSetup.Versions[t.Version]
		if !found {
			versionSetup = VersionSetup{
				Version: t.Version,
				States:  map[string]StateSetup{},
			}
		}

		// check if state exists
		stateSetup, found := versionSetup.States[t.State]
		if !found {
			stateSetup = StateSetup{
				State:     t.State,
				Instances: map[string]string{},
			}
		}

		// add instances
		for j := 0; j < t.Size; j++ {
			u := uuid.New().String()
			stateSetup.Instances[u] = u
		}
	}

	// success
	return serviceSetup
}

//------------------------------------------------------------------------------

func determineTasks(domain string, architecture string, service string) ([]model.Task, []model.Task, []model.Task) {
	targetSetup := determineTargetSetup(domain, architecture, service)
	currentSetup := determineCurrentSetup(domain, service)
	updateTasks := []model.Task{}
	createTasks := []model.Task{}
	removeTasks := []model.Task{}

	// determine all unchanged instances
	for _, targetVersionSetup := range targetSetup.Versions {
		for _, targetStateSetup := range targetVersionSetup.States {
			for targetInstance := range targetStateSetup.Instances {

				// try to find matching current instance
				currentVersionSetup, found := currentSetup.Versions[targetVersionSetup.Version]
				if !found {
					continue
				}

				currentStateSetup, found := currentVersionSetup.States[targetStateSetup.State]
				if !found {
					continue
				}

				for currentInstance := range currentStateSetup.Instances {
					// instance has been found - now remove instances from the setup
					delete(targetStateSetup.Instances, targetInstance)
					delete(currentStateSetup.Instances, currentInstance)
					break
				}
			}
		}
	}

	// determine all instances which need to be updated
	for targetVersion, targetVersionSetup := range targetSetup.Versions {
		for targetState, targetStateSetup := range targetVersionSetup.States {
			for targetInstance := range targetStateSetup.Instances {

				// try to find matching current instance with matching version
				currentVersionSetup, found := currentSetup.Versions[targetVersionSetup.Version]
				if !found {
					continue
				}

				for _, currentStateSetup := range currentVersionSetup.States {
					for currentInstance := range currentStateSetup.Instances {
						// create new update task
						updateTask, _ := NewInstanceTask(domain, "TODO: unknown", architecture, service, targetVersion, currentInstance, targetState)

						// append new task to set of update tasks
						updateTasks = append(updateTasks, updateTask)

						// instance has been found - now remove instances from the setup
						delete(targetStateSetup.Instances, targetInstance)
						delete(currentStateSetup.Instances, currentInstance)
						break
					}

				}
			}
		}
	}

	// all leftover current instances need to be removed
	for currentVersion, currentVersionSetup := range currentSetup.Versions {
		for _, currentStateSetup := range currentVersionSetup.States {
			for currentInstance := range currentStateSetup.Instances {
				// create new remove task
				removeTask, _ := NewInstanceTask(domain, "TODO: unknown", architecture, service, currentVersion, currentInstance, "initial")

				// append new task to set of remove tasks
				removeTasks = append(removeTasks, removeTask)
			}
		}
	}

	// all leftover target instances need to be created
	for targetVersion, targetVersionSetup := range targetSetup.Versions {
		for targetState, targetStateSetup := range targetVersionSetup.States {
			for targetInstance := range targetStateSetup.Instances {
				// create new create task
				createTask, _ := NewInstanceTask(domain, "TODO: unknown", architecture, service, targetVersion, targetInstance, targetState)

				// append new task to set of create tasks
				createTasks = append(createTasks, createTask)
			}
		}
	}

	// success
	return updateTasks, createTasks, removeTasks
}

//------------------------------------------------------------------------------

// NewElementTask creates a new task
func NewElementTask(domain string, parent string, architecture string, version string, element string) (model.Task, error) {
	var task model.Task

	// TODO: check parameters if context exists
	task.Type         = "ElementTask"
	task.Domain       = domain
	task.Architecture = architecture
	task.Version      = version
	task.Element      = element
	task.Version      = ""
	task.Instance     = ""
	task.State        = ""
	task.UUID         = uuid.New().String()
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

	// initialize if needed
	if status == model.TaskStatusInitial {
		// update status
		task.Status = model.TaskStatusExecuting

		// determine required subtasks
		updateTasks, createTasks, removeTasks := determineTasks(task.Domain, task.Architecture, task.Element)

		// add tasks to domain

		// create task groups
		mainTask, _ := NewParallelTask(task.Domain, task.UUID, []string{})
		task.AddSubtask(&mainTask)

		// add update subtasks
		updateTask, _ := NewParallelTask(task.Domain, mainTask.GetUUID(), []string{})
		mainTask.AddSubtask(&updateTask)
		for _, s := range updateTasks {
			subTask, _ := NewInstanceTask(s.Domain, mainTask.GetUUID(), task.Architecture, s.Element, s.Cluster, s.Instance, s.State)

			updateTask.AddSubtask(&subTask)
		}

		// add create subtasks
		createTask, _ := NewParallelTask(task.Domain, mainTask.GetUUID(), []string{})
		mainTask.AddSubtask(&createTask)
		for _, s := range createTasks {
			subTask, _ := NewInstanceTask(s.Domain, mainTask.GetUUID(), task.Architecture, s.Element, s.Cluster, s.Instance, s.State)

			createTask.AddSubtask(&subTask)
		}

		// add remove subtasks
		removeTask, _ := NewParallelTask(task.Domain, mainTask.GetUUID(), []string{})
		mainTask.AddSubtask(&removeTask)
		for _, s := range removeTasks {
			subTask, _ := NewInstanceTask(s.Domain, mainTask.GetUUID(), task.Architecture, s.Element, s.Cluster, s.Instance, s.State)

			removeTask.AddSubtask(&subTask)
		}

		// trigger execution of main subtask
		channel <- model.NewEvent(task.Domain, mainTask.GetUUID(), model.EventTypeTaskExecution, task.UUID)

		// success
		return
	}

	// success
	return
}

//------------------------------------------------------------------------------
