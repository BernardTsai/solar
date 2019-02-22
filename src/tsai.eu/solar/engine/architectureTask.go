package engine

import (
	"errors"
  "fmt"
	"github.com/google/uuid"
	"tsai.eu/solar/model"
)

//------------------------------------------------------------------------------

// NewArchitectureTask creates a new task
func NewArchitectureTask(domain string, parent string, architecture *model.Architecture) (model.Task, error) {
	var task model.Task
 
	// TODO: check parameters if context exists
	task.Type         = "ArchitectureTask"
	task.Domain       = domain
	task.Architecture = architecture.Architecture
	task.Version      = architecture.Version
	task.Element      = ""
	task.Cluster      = ""
	task.Instance     = ""
	task.State        = ""
	task.UUID         = uuid.New().String()
	task.Parent       = parent
	task.Status       = model.TaskStatusInitial
	task.Phase        = 0
	task.Subtasks     = []string{}

	// add handlers
	task.SetExecute(ExecuteParallelTask)
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

	// construct all required subtasks (one for each element)
	architecture.Elements.RLock()
	for element := range architecture.Elements.Map {
		subtask, err := NewElementTask(domain, task.UUID, architecture.Architecture, architecture.Version, element)
		if err != nil {
			return task, errors.New("unable to create subtask for an element")
		}

		task.AddSubtask(&subtask)
	}
	architecture.Elements.RUnlock()

	fmt.Println(task.Subtasks)
	fmt.Println(d.GetTask(task.UUID))

	// success
	return task, nil
}

//------------------------------------------------------------------------------
