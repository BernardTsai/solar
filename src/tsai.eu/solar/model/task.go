package model

import (
	"context"
	"sort"
	"errors"

	"tsai.eu/solar/util"
)

//------------------------------------------------------------------------------

// TaskInfo specifies the basic behaviour of a task
type TaskInfo struct {
	Type         string      `yaml:"Type"`         // type of task
	Domain       string      `yaml:"Domain"`       // domain of task
	Solution     string      `yaml:"Solution"`     // architecture of entity
	Version      string      `yaml:"Version"`      // architecture version of entity
	Element      string      `yaml:"Element"`      // element of entity
	Cluster      string      `yaml:"Cluster"`      // cluster of entity
	Instance     string      `yaml:"Instance"`     // instance of entity
	State        string      `yaml:"State"`        // desired state of entity
	Action       string      `yaml:"Action"`       // action to be performed by a controller
	UUID         string      `yaml:"UUID"`         // uuid of task
	Parent       string      `yaml:"Parent"`       // uuid of parent task
	Status       string      `yaml:"Status"`       // status of task: (execution/completion/failure)
	Phase        int         `yaml:"Phase"`        // phase of task
	Subtasks     []*TaskInfo `yaml:"Subtasks"`     // list of subtasks
	Events       []*Event    `yaml:"Events"`       // list of events
}

//------------------------------------------------------------------------------

// NewTaskInfo derives a taskinfo object from a task.
func NewTaskInfo(task *Task, level int) (*TaskInfo) {
	taskinfo := TaskInfo{
		Type:       task.Type,
		Domain:     task.Domain,
		Solution:   task.Solution,
		Version:    task.Version,
		Element:    task.Element,
		Cluster:    task.Cluster,
		Instance:   task.Instance,
		State:      task.State,
		Action:     task.Action,
		UUID:       task.UUID,
		Parent:     task.Parent,
		Status:     task.Status,
		Phase:      task.Phase,
		Subtasks:   []*TaskInfo{},
		Events:     []*Event{},
	}

	// add events
	domain, _ := GetDomain(task.Domain)

	for _, eventUUID := range task.Events {
		event, _ := domain.GetEvent(eventUUID)

		taskinfo.Events = append(taskinfo.Events, event)
	}
	sort.SliceStable(taskinfo.Events, func(i, j int) bool { return taskinfo.Events[i].Time < taskinfo.Events[j].Time })

	// check if details are needed (depends on level)
	if level == 0 || level > 1 {
		sublevel := 0
		if level > 1 {
			sublevel = level - 1
		}

		// add subtasks
		for _, subtaskUUID := range task.Subtasks {
			subtask, _ := GetTask(task.Domain, subtaskUUID)

			taskinfo.Subtasks = append(taskinfo.Subtasks, NewTaskInfo(subtask, sublevel))
		}
	}

	return &taskinfo
}

//------------------------------------------------------------------------------

// TaskHandler is function capable of processing a task related event.
type TaskHandler func(task *Task)

//------------------------------------------------------------------------------

// Task specifies the basic behaviour of a task
type Task struct {
	Type         string     `yaml:"Type"`         // type of task
	Domain       string     `yaml:"Domain"`       // domain of task
	Solution     string     `yaml:"Solution"`     // architecture of entity
	Version      string     `yaml:"Version"`      // architecture version of entity
	Element      string     `yaml:"Element"`      // element of entity
	Cluster      string     `yaml:"Cluster"`      // cluster of entity
	Instance     string     `yaml:"Instance"`     // instance of entity
	State        string     `yaml:"State"`        // desired state of entity
	Action       string     `yaml:"Action"`       // desired action for entity
	UUID         string     `yaml:"UUID"`         // uuid of task
	Parent       string     `yaml:"Parent"`       // uuid of parent task
	Status       string     `yaml:"Status"`       // status of task: (execution/completion/failure)
	Phase        int        `yaml:"Phase"`        // phase of task
	Subtasks     []string   `yaml:"Subtasks"`     // list of subtasks
	Events       []string   `yaml:"Events"`       // list of events
	execute      TaskHandler
	terminate    TaskHandler
	failed       TaskHandler
	timeout      TaskHandler
	completed    TaskHandler
}

//------------------------------------------------------------------------------

// GetType delivers the type of the task.
func (task *Task) GetType() string {
	return task.Type
}

//------------------------------------------------------------------------------

// GetDomain delivers the domain of the task.
func (task *Task) GetDomain() string {
	return task.Domain
}

//------------------------------------------------------------------------------

// GetSolution delivers the architecture of the entity.
func (task *Task) GetSolution() string {
	return task.Solution
}

//------------------------------------------------------------------------------

// GetVersion delivers the architecture version of the entity.
func (task *Task) GetVersion() string {
	return task.Version
}

//------------------------------------------------------------------------------

// GetElement delivers the element of the entity.
func (task *Task) GetElement() string {
	return task.Element
}

//------------------------------------------------------------------------------

// GetCluster delivers the cluster of the entity.
func (task *Task) GetCluster() string {
	return task.Cluster
}

//------------------------------------------------------------------------------

// GetInstance delivers the instance of entity task.
func (task *Task) GetInstance() string {
	return task.Instance
}

//------------------------------------------------------------------------------

// GetState delivers the state of the entity.
func (task *Task) GetState() string {
	return task.State
}

//------------------------------------------------------------------------------

// GetAction delivers the action for the entity.
func (task *Task) GetAction() string {
	return task.Action
}

//------------------------------------------------------------------------------

// GetUUID delivers the universal unique identifier of the task.
func (task *Task) GetUUID() string {
	return task.UUID
}

//------------------------------------------------------------------------------

// GetParent delivers the universal unique identifier of the parent task.
func (task *Task) GetParent() string {
	return task.Parent
}

//------------------------------------------------------------------------------

// GetStatus delivers the status of the task.
func (task *Task) GetStatus() string {
	return task.Status
}

//------------------------------------------------------------------------------

// GetPhase delivers the internal status of the task.
func (task *Task) GetPhase() int {
	return task.Phase
}

//------------------------------------------------------------------------------

// GetSubtask provides the subtask with a given uuid.
func (task *Task) GetSubtask(uuid string) (*Task, error) {
	// check if uuid is in slice of substasks
	found := false
	for _, suuid := range task.Subtasks {
		if suuid == uuid {
			found = true
			break
		}
	}

	if !found {
		return nil, errors.New("unknown subtask")
	}

	// get domain
	domain, _ := GetModel().GetDomain(task.Domain)

	// get subtask
	subtask, err := domain.GetTask(uuid)
	if err != nil {
		return nil, errors.New("unknown subtask")
	}

	// success
	return subtask, nil
}

//------------------------------------------------------------------------------

// GetSubtasks provides a slice of subtask uuids.
func (task *Task) GetSubtasks() []string {
	return task.Subtasks
}

//------------------------------------------------------------------------------

// AddSubtask adds a subtask to the list of subtasks.
func (task *Task) AddSubtask(subtask *Task) {
	task.Subtasks = append(task.Subtasks, subtask.GetUUID())
}

//------------------------------------------------------------------------------

// AddEvent adds an event to the list of events.
func (task *Task) AddEvent(event *Event) {
	task.Events = append(task.Events, event.GetUUID())
}

//------------------------------------------------------------------------------

// Save writes the task as json data to a file
func (task *Task) Save(filename string) error {
	return util.SaveYAML(filename, task)
}

//------------------------------------------------------------------------------

// Load reads the task from a file
func (task *Task) Load(filename string) error {
	return util.LoadYAML(filename, task)
}

//------------------------------------------------------------------------------

// Load2 imports a yaml model
func (task *Task) Load2(yaml string) error {
	return util.ConvertFromYAML(yaml, task)
}

//------------------------------------------------------------------------------

// Show displays the task information as yaml
func (task *Task) Show() (string, error) {
	return util.ConvertToYAML(task)
}

//------------------------------------------------------------------------------

// SetExecute defines the execute event handler of a task
func (task *Task) SetExecute(handler TaskHandler) {
	task.execute = handler
}

//------------------------------------------------------------------------------

// Execute processes a task
func (task *Task) Execute(ctx context.Context) {
	// execute task if appropriate handler has been defined
	if task.execute != nil {
		task.execute(task)
	}
}

//------------------------------------------------------------------------------

// SetTerminate defines the terminate event handler of a task
func (task *Task) SetTerminate(handler TaskHandler) {
	task.terminate = handler
}

//------------------------------------------------------------------------------

// Terminate stops a task
func (task *Task) Terminate(ctx context.Context) {
	// execute task if appropriate handler has been defined
	if task.terminate != nil {
		task.terminate(task)
	}
}

//------------------------------------------------------------------------------

// SetFailed defines the failed event handler of a task
func (task *Task) SetFailed(handler TaskHandler) {
	task.failed = handler
}

//------------------------------------------------------------------------------

// Failed handles the failure of the task
func (task *Task) Failed(ctx context.Context) {
	// execute task if appropriate handler has been defined
	if task.failed != nil {
		task.failed(task)
	}
}

//------------------------------------------------------------------------------

// SetTimeout defines the timeout event handler of a task
func (task *Task) SetTimeout(handler TaskHandler) {
	task.timeout = handler
}

//------------------------------------------------------------------------------

// Timeout handles the timeput of the task
func (task *Task) Timeout(ctx context.Context) {
	// execute task if appropriate handler has been defined
	if task.timeout != nil {
		task.timeout(task)
	}
}

//------------------------------------------------------------------------------

// SetCompleted defines the completed event handler of a task
func (task *Task) SetCompleted(handler TaskHandler) {
	task.completed = handler
}

//------------------------------------------------------------------------------

// Completed handles the completion of the task
func (task *Task) Completed(ctx context.Context) {
	// execute task if appropriate handler has been defined
	if task.completed != nil {
		task.completed(task)
	}
}

//------------------------------------------------------------------------------

// GetTimestamps delivers the timestamps of the task.
func (task *Task) GetTimestamps() (started int64, completed int64, latest int64) {
	var min int64
	var max int64
	var lst int64
	var def int64

	def = 9223372036854775807
	min = def
	max = 0
	lst = 0
	for _, eventUUID := range task.Events {
		event, _ := GetEvent(task.Domain, eventUUID)

		if event.Type == "execution" {
			if event.Time < min {
				min = event.Time
			}
			if lst < event.Time {
				lst = event.Time
			}
		} else {
			if max < event.Time {
				max = event.Time
			}
			if lst < event.Time {
				lst = event.Time
			}
		}
	}

	// copy to results
	started = 0
	if min != def {
		started = min
	}
	completed = 0
	if max != def {
		completed = max
	}
	latest = 0
	if lst != 0 {
		latest = lst
	}

	return started, completed, latest
}

//------------------------------------------------------------------------------
