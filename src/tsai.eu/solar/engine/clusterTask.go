package engine

import (
	"errors"

	"tsai.eu/solar/util"
	"tsai.eu/solar/model"
)

//------------------------------------------------------------------------------

// NewClusterTask creates a new task
func NewClusterTask(domain string, parent string, solution string, version string, element string, cluster string) (model.Task, error) {
	var task model.Task

	// TODO: check parameters if context exists
	task.Type         = "ElementTask"
	task.Domain       = domain
	task.Solution     = solution
	task.Version      = version
	task.Element      = element
	task.Cluster      = cluster
	task.Instance     = ""
	task.State        = ""
	task.UUID         = util.UUID()
	task.Parent       = parent
	task.Status       = model.TaskStatusInitial
	task.Phase        = 0
	task.Subtasks     = []string{}

	// add handlers
	task.SetExecute(ExecuteClusterTask)
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

// ExecuteClusterTask is the main task execution routine.
func ExecuteClusterTask(task *model.Task) {
	// get event channel
	channel := GetEventChannel()

	// check status
	status := task.GetStatus()

	if status != model.TaskStatusInitial && status != model.TaskStatusExecuting {
		return
	}

	// determine context
	cluster, _   := model.GetCluster(task.Domain, task.Solution, task.Element, task.Cluster)

	// evaluate relationships
	switch cluster.State {
	case model.InactiveState:
		// check if all context relationships are active
		relationshipNames, _ := cluster.ListRelationships()
		for _, relationshipName := range relationshipNames {
			relationship, _ := cluster.GetRelationship(relationshipName)

			if relationship.Type != model.ContextRelationship {
				continue
			}

			// check if the related cluster is in the desired state
			refCluster, _ := model.GetCluster(relationship.Domain, relationship.Solution, relationship.Element, relationship.Version)
			if refCluster.State != model.ActiveState	{
				// update the related cluster
				triggerClusterTask(task, relationship)

				// return and wait for next event
				return
			}
		}
	case model.ActiveState:
		// check if all service and context relationships are active
		relationshipNames, _ := cluster.ListRelationships()
		for _, relationshipName := range relationshipNames {
			relationship, _ := cluster.GetRelationship(relationshipName)

			if relationship.Type != model.ContextRelationship && relationship.Type != model.ServiceRelationship {
				continue
			}

			refCluster, _ := model.GetCluster(relationship.Domain, relationship.Solution, relationship.Element, relationship.Version)
			if refCluster.State != model.ActiveState	{
				// update the related cluster
				triggerClusterTask(task, relationship)

				// return and wait for next event
				return
			}
		}
	}

	// adjust instances
	switch cluster.State {
	case model.InitialState:
		// one by one identify the instances which need to be reset
		instanceNames, _ := cluster.ListInstances()
		for _, instanceName := range instanceNames {
			instance, _ := cluster.GetInstance(instanceName)
			if instance.State != cluster.State {
				// update the related instance
				triggerInstanceTask(task, instanceName, cluster.State)

				// return and wait for next event
				return
			}
		}
	case model.InactiveState:
		count := countInstances(cluster, model.InactiveState)

		// one by one identify the instances which need to be reset
		instanceNames, _ := cluster.ListInstances()
		for _, instanceName := range instanceNames {
			instance, _ := cluster.GetInstance(instanceName)
			if instance.State != cluster.State {
				// update instance to the desired state
				if cluster.Size < count {
					triggerInstanceTask(task, instanceName, cluster.State)
				} else {
					triggerInstanceTask(task, instanceName, model.InitialState)
				}

				// return and wait for next event
				return
			}
		}
	case model.ActiveState:
		count := countInstances(cluster, model.ActiveState)

		// one by one identify the instances which need to be reset
		instanceNames, _ := cluster.ListInstances()
		for _, instanceName := range instanceNames {
			instance, _ := cluster.GetInstance(instanceName)

			// update instance
			if instance.State != cluster.State && cluster.Size < count {
				triggerInstanceTask(task, instanceName, cluster.State)

				// return and wait for next event
				return
			}
		}
	}

	// execution has completed
	channel <- model.NewEvent(task.Domain, task.UUID, model.EventTypeTaskCompletion, task.UUID)
}

//------------------------------------------------------------------------------

// countInstances counts how many instances of a cluster are in a specific state.
func countInstances(cluster *model.Cluster, state string) int {
	count := 0

	// loop over all instances
	instanceNames, _ := cluster.ListInstances()
	for _, instanceName := range instanceNames {
		instance, _ := cluster.GetInstance(instanceName)
		if instance.State == state {
			count++
		}
	}

	return count
}

//------------------------------------------------------------------------------

// triggerClusterTask triggers a task to update a cluster.
func triggerClusterTask(task *model.Task, relationship *model.Relationship)  {
	// get event channel
	channel := GetEventChannel()

	// create task to update the cluster
	subtask, _ := NewClusterTask(relationship.Domain, task.UUID, relationship.Solution, relationship.Version,relationship.Element, relationship.Version)
	task.AddSubtask(&subtask)

	// trigger the task
	channel <- model.NewEvent(task.Domain, subtask.UUID, model.EventTypeTaskExecution, task.UUID)
}

//------------------------------------------------------------------------------

// triggerInstanceTask triggers a task to update an instance to a specific state.
func triggerInstanceTask(task *model.Task, instance string, state string)  {
	// get event channel
	channel := GetEventChannel()

	// create task to update the instance
	subtask, _ := NewInstanceTask(task.Domain, task.UUID, task.Solution, task.Version, task.Element, task.Version, instance, state )
	task.AddSubtask(&subtask)

	// trigger the task
	channel <- model.NewEvent(task.Domain, subtask.UUID, model.EventTypeTaskExecution, task.UUID)
}

//------------------------------------------------------------------------------