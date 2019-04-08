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
	task.Type         = "Cluster"
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

	// check and update status
	status := task.GetStatus()

	if status != model.TaskStatusInitial && status != model.TaskStatusExecuting {
		return
	}

	if status == model.TaskStatusInitial {
		task.Status = model.TaskStatusExecuting
	}

	// determine context
	cluster, _  := model.GetCluster(task.Domain, task.Solution, task.Element, task.Cluster)

	// validate sizing
	if cluster.Min > cluster.Max || cluster.Size < cluster.Min || cluster.Max < cluster.Size {
		channel <- model.NewEvent(task.Domain, task.UUID, model.EventTypeTaskFailure, task.UUID, "inconsistent sizing of cluster: '" + task.Element + " - " + task.Cluster + " of solution: '" + task.Solution + "' within domain:'" + task.Domain)
		return
	}


	// evaluate relationships
	switch cluster.Target {
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
			if refCluster.State != model.ActiveState {
				// check if the desired target state is active (otherwise we have a configuration mismatch)
				if refCluster.Target != model.ActiveState {
					channel <- model.NewEvent(task.Domain, task.UUID, model.EventTypeTaskFailure, task.UUID, "unable to establish context dependency: " + relationship.Element + " - " + relationship.Version)
					return
				}

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
				// check if the desired target state is active (otherwise we have a configuration mismatch)
				if refCluster.Target != model.ActiveState {
					channel <- model.NewEvent(task.Domain, task.UUID, model.EventTypeTaskFailure, task.UUID, "unable to establish service dependency: " + relationship.Element + " - " + relationship.Version)
					return
				}

				// update the related cluster
				triggerClusterTask(task, relationship)

				// return and wait for next event
				return
			}
		}
	}

	// adjust instances
	switch cluster.Target {
	case model.InitialState:
		// one by one identify the instances which need to be reset
		instanceNames, _ := cluster.ListInstances()
		for _, instanceName := range instanceNames {
			instance, _ := cluster.GetInstance(instanceName)
			if instance.State != cluster.Target {
				// update the related instance
				triggerInstanceTask(task, instanceName, cluster.Target)

				// return and wait for next event
				return
			}
		}

		// cluster has reached the desired state
		cluster.State = model.InitialState
	case model.InactiveState:
		_, inactive, active, _, _ := cluster.Pools()

		// one by one identify the instances which need to be reset
		instanceNames, _ := cluster.ListInstances()
		for _, instanceName := range instanceNames {
			instance, _ := cluster.GetInstance(instanceName)

			// cleanup failed instances
			if instance.State == model.FailureState {
				triggerInstanceTask(task, instanceName, model.InitialState)

				return
			}

			// deactivate active instances
			if instance.State == model.ActiveState {
				triggerInstanceTask(task, instanceName, model.InactiveState)

				return
			}

			// ensure that the number of inactive nodes matches the cluster size
			if inactive < cluster.Size {
				// activate the amount of required instances
			  if instance.State != model.InactiveState {
					triggerInstanceTask(task, instanceName, model.InactiveState)

					return
				}
			} else if active > cluster.Size {
				// activate the amount of required instances
			  if instance.State == model.InactiveState {
					triggerInstanceTask(task, instanceName, model.InitialState)

					return
				}
			}

			// cluster has reached the desired state
			cluster.State = model.InactiveState
		}
	case model.ActiveState:
		_, inactive, active, _, _ := cluster.Pools()

		// one by one identify the instances which need to be reset
		instanceNames, _ := cluster.ListInstances()
		for _, instanceName := range instanceNames {
			instance, _ := cluster.GetInstance(instanceName)

			// cleanup failed instances
			if instance.State == model.FailureState {
				triggerInstanceTask(task, instanceName, model.InitialState)

				return
			}

			// ensure that the number of active nodes matches the cluster size
			if active < cluster.Size {
				// activate the amount of required instances
			  if instance.State != cluster.Target {
					triggerInstanceTask(task, instanceName, model.ActiveState)

					return
				}
			} else if active == cluster.Size {
				// remove excess inactive instances
				if inactive > (cluster.Max - cluster.Size) {
					if instance.State == model.InactiveState {
						triggerInstanceTask(task, instanceName, model.InitialState)

						return
					}
				}
			} else if active > cluster.Size {
				// deactivate excess active instances
				if instance.State == model.ActiveState {
					triggerInstanceTask(task, instanceName, model.InactiveState)

					return
				}
			}
		}

		// cluster has reached the desired state
		cluster.State = model.ActiveState
	}

	// execution has completed
	channel <- model.NewEvent(task.Domain, task.UUID, model.EventTypeTaskCompletion, task.UUID, "")
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
	subtask, _ := NewClusterTask(relationship.Domain, task.UUID, task.Solution, task.Version, relationship.Element, relationship.Version)

	task.AddSubtask(&subtask)

	// trigger the task
	channel <- model.NewEvent(task.Domain, subtask.UUID, model.EventTypeTaskExecution, task.UUID, "")
}

//------------------------------------------------------------------------------

// triggerInstanceTask triggers a task to update an instance to a specific state.
func triggerInstanceTask(task *model.Task, instance string, state string)  {
	// get event channel
	channel := GetEventChannel()

	// create task to update the instance
	subtask, _ := NewInstanceTask(task.Domain, task.UUID, task.Solution, task.Version, task.Element, task.Cluster, instance, state )
	task.AddSubtask(&subtask)

	// trigger the task
	channel <- model.NewEvent(task.Domain, subtask.UUID, model.EventTypeTaskExecution, task.UUID, "")
}

//------------------------------------------------------------------------------
