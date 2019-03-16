package cli

import (
	"strconv"
	"sort"

	ishell "gopkg.in/abiosoft/ishell.v2"
	"tsai.eu/solar/model"
	"tsai.eu/solar/util"
	"tsai.eu/solar/engine"
)

//------------------------------------------------------------------------------

// TaskInfo specifies the basic behaviour of a task
type TaskInfo struct {
	Type         string                  `yaml:"Type"`         // type of task
	Domain       string                  `yaml:"Domain"`       // domain of task
	Solution     string                  `yaml:"Solution"`     // architecture of entity
	Version      string                  `yaml:"Version"`      // architecture version of entity
	Element      string                  `yaml:"Element"`      // element of entity
	Cluster      string                  `yaml:"Cluster"`      // cluster of entity
	Instance     string                  `yaml:"Instance"`     // instance of entity
	State        string                  `yaml:"State"`        // desired state of entity
	UUID         string                  `yaml:"UUID"`         // uuid of task
	Parent       string                  `yaml:"Parent"`       // uuid of parent task
	Status       string                  `yaml:"Status"`       // status of task: (execution/completion/failure)
	Phase        int                     `yaml:"Phase"`        // phase of task
	Subtasks     []*TaskInfo             `yaml:"Subtasks"`     // list of subtasks
	Events       []*model.Event          `yaml:"Events"`       // list of events
}

//------------------------------------------------------------------------------

// TraceElement holds the trace information for an affected element
type TraceElement struct {
	Name     string                    `yaml:"Name"`       // name of element
	Index    int                       `yaml:"Index"`      // index of element
	Clusters map[string]*TraceCluster  `yaml:"Clusters"`   // map of clusters
	Tasks    []*TraceTask              `yaml:"Tasks"`      // map of tasks
}

// TraceCluster holds the trace information for an affected cluster
type TraceCluster struct {
	Name      string                    `yaml:"Name"`       // name of cluster
	Index     int                       `yaml:"Index"`      // index of cluster
	Instances map[string]*TraceInstance `yaml:"Instances"`  // map of instances
	Tasks     []*TraceTask              `yaml:"Tasks"`      // map of tasks
}

// TraceInstance holds the trace information for an affected instance
type TraceInstance struct {
	Name  string                         `yaml:"Name"`       // name of instance
	Index int                            `yaml:"Index"`      // index of instance
	Tasks []*TraceTask                   `yaml:"Tasks"`      // map of tasks
}

// TraceTask holds the trace information for an affected task
type TraceTask struct {
	UUID      string                     `yaml:"UUID"`       // uuid of task
	Started   int64                      `yaml:"Started"`    // start timestamp
	Completed int64                      `yaml:"Completed"`  // end timestamp
	Latest    int64                      `yaml:"Latest"`     // latest timestamp
	Status    string                     `yaml:"Status"`     // status of task
	State     string                     `yaml:"State"`      // desired target state
	Version   string                     `yaml:"Version"`    // desired architecture version
	Layer     int                        `yaml:"Layer"`      // layer
	Layers    int                        `yaml:"Layers"`     // layers
}

// TraceEvent holds the trace information for an affected event
type TraceEvent struct {
	UUID      string                     `yaml:"UUID"`       // uuid of event
	Time      int64                      `yaml:"Time"`       // timestamp of event
	Type      string                     `yaml:"Type"`       // type of event
	Task1     string                     `yaml:"Task1"`      // first task of event
	Element1  string                     `yaml:"Element1"`   // first element of event
	Cluster1  string                     `yaml:"Cluster1"`   // first cluster of event
	Instance1 string                     `yaml:"Instance1"`  // first instance of event
	Index1    int                        `yaml:"Index1"`     // first index of event
	Layer1    int                        `yaml:"Layer1"`     // first layer of event
	Layers1   int                        `yaml:"Layers1"`    // first layers of event
	Task2     string                     `yaml:"Task2"`      // second task of event
	Element2  string                     `yaml:"Element2"`   // second element of event
	Cluster2  string                     `yaml:"Cluster2"`   // second cluster of of event
	Instance2 string                     `yaml:"Instance2"`  // second instance of event
	Index2    int                        `yaml:"Index2"`     // second index of event
	Layer2    int                        `yaml:"Layer2"`     // second layer of event
	Layers2   int                        `yaml:"Layers2"`    // second layers of event
}

// Trace describes the stack trace of a task
type Trace struct {
	Task       string                    `yaml:"Task"`       // UUID of initial task
	Domain     string                    `yaml:"Domain"`     // domain of trace
	Solution   string                    `yaml:"Solution"`   // solution of trace
	Min        int64                     `yaml:"Min"`        // min timestamp
	Max        int64                     `yaml:"Max"`        // max timestamp
	Range      int64                     `yaml:"Range"`      // timestamp range
	Elements   map[string]*TraceElement  `yaml:"Elements"`   // map of affected elements
	Tasks      []*TraceTask              `yaml:"Tasks"`      // map of tasks
  Events     []*TraceEvent             `yaml:"Events"`     // list of affected events
}

//------------------------------------------------------------------------------

// NewTaskInfo derives a taskinfo object from a task.
func NewTaskInfo(task *model.Task, level int) (*TaskInfo) {
	taskinfo := TaskInfo{
		Type:       task.Type,
		Domain:     task.Domain,
		Solution:   task.Solution,
		Version:    task.Version,
		Element:    task.Element,
		Cluster:    task.Cluster,
		Instance:   task.Instance,
		State:      task.State,
		UUID:       task.UUID,
		Parent:     task.Parent,
		Status:     task.Status,
		Phase:      task.Phase,
		Subtasks:   []*TaskInfo{},
		Events:     []*model.Event{},
	}

	// add events
	domain, _ := model.GetDomain(task.Domain)

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
			subtask, _ := model.GetTask(task.Domain, subtaskUUID)

			taskinfo.Subtasks = append(taskinfo.Subtasks, NewTaskInfo(subtask, sublevel))
		}
	}

	return &taskinfo
}

//------------------------------------------------------------------------------

// NewTrace derives a trace object from a task.
func NewTrace(task *model.Task) (*Trace) {
	trace := Trace{
		Task:     task.GetUUID(),
		Domain:   task.GetDomain(),
		Solution: task.GetSolution(),
		Elements: map[string]*TraceElement{},
		Events:   []*TraceEvent{},
	}

	addTaskToTrace(&trace, task)

	sortTrace(&trace)

	return &trace
}

//------------------------------------------------------------------------------

// addTaskToTrace appends task information to a trace
func addTaskToTrace(trace *Trace, task* model.Task) {
	var element     *TraceElement
	var cluster     *TraceCluster
	var instance    *TraceInstance
	var taskEntry   *TraceTask
	var eventEntry  *TraceEvent
	var found       bool
	var added       bool

	// add events
	for _, eventUUID := range task.Events {
		event, _ := model.GetEvent(task.GetDomain(), eventUUID)

		// create new entry for the event
		task2, _ := model.GetTask(task.GetDomain(), event.Task)
		eventEntry = &TraceEvent{
			UUID:       event.UUID,
			Time:       event.Time,
			Type:       event.Type,
			Task1:      "",
			Element1:   "",
			Cluster1:   "",
			Instance1:  "",
			Index1:     0,
			Layer1:     0,
			Layers1:    1,
			Task2:      task2.GetUUID(),
			Element2:   task2.GetElement(),
			Cluster2:   task2.GetCluster(),
			Instance2:  task2.GetInstance(),
			Index2:     0,
			Layer2:     0,
			Layers2:    1,
		}

		task1, _ := model.GetTask(task.GetDomain(), event.Source)
		if task1 != nil {
			eventEntry.Task1     = task1.GetUUID()
			eventEntry.Element1  = task1.GetElement()
			eventEntry.Cluster1  = task1.GetCluster()
			eventEntry.Instance1 = task1.GetInstance()
		}

		// add event entry to trace
		trace.Events = append(trace.Events, eventEntry)
	}

	// create a task task entry
	added = false
	started, completed, latest := task.GetTimestamps()

	taskEntry = &TraceTask{
		UUID:      task.GetUUID(),
		Started:   started,
		Completed: completed,
		Latest:    latest,
		Status:    task.GetStatus(),
		State:     task.GetState(),
		Version:   task.GetVersion(),
	}

	// check if solution element has already been defined
	element = &TraceElement{}
	if task.Element != "" {
		// add missing trace elements
		if element, found = trace.Elements[task.Element]; !found {
			element = &TraceElement{
				Name:     task.Element,
				Index:    -1,
				Clusters: map[string]*TraceCluster{},
			}

			trace.Elements[task.Element] = element
		}
	} else if !added {
		trace.Tasks = append(trace.Tasks, taskEntry)
		added = true
	}

	// check if solution cluster has already been defined
	cluster = &TraceCluster{}
	if task.Cluster != "" {
		// add missing trace clusters
		if cluster, found = element.Clusters[task.Cluster]; !found {
			cluster = &TraceCluster{
				Name:      task.Cluster,
				Index:     -1,
				Instances: map[string]*TraceInstance{},
			}

			element.Clusters[task.Cluster] = cluster
		}
	} else if !added {
		element.Tasks = append(element.Tasks, taskEntry)
		added = true
	}

	// check if solution instance has already been defined
	instance = &TraceInstance{}
	if task.Instance != "" {
		// add missing trace instances
		if instance, found = cluster.Instances[task.Instance]; !found {
			instance = &TraceInstance{
				Name:  task.Instance,
				Index:  -1,
				Tasks: []*TraceTask{},
			}

			cluster.Instances[task.Instance] = instance
		}
		} else if !added {
			cluster.Tasks = append(cluster.Tasks, taskEntry)
			added = true
		}

		if !added {
			instance.Tasks = append(instance.Tasks, taskEntry)
		}


	// add all subtasks
	domainName := task.GetDomain()
	for _, subtaskUUID := range task.Subtasks {
		subtask, _ := model.GetTask(domainName, subtaskUUID)

		addTaskToTrace(trace, subtask)
	}
}

//------------------------------------------------------------------------------

// sortTrace sorts the entities of a trace and adds indices to the entities
func sortTrace(trace *Trace) {
	// vertical index
	index := 0

	// lookup for task and entity indexes
	taskIndices := map[string]int{}
	taskLayer   := map[string]int{}
	taskLayers  := map[string]int{}

	// top level tasks have index 0
	layer := 0
	for _, task := range trace.Tasks {
		taskIndices[task.UUID] = index
		taskLayer[task.UUID]   = layer
		task.Layer             = layer
		layer                  = layer +1
	}
	for _, task := range trace.Tasks {
		task.Layers           = layer
		taskLayers[task.UUID] = layer
	}

	// loop over all elements
	index = index + 1
	elementKeys := make([]string, 0, len(trace.Elements))
	for elementKey := range trace.Elements {
		elementKeys = append(elementKeys, elementKey)
	}
	sort.Slice(elementKeys, func(i, j int) bool {return elementKeys[i] < elementKeys[j]})

	for _, elementKey := range elementKeys {
		element := trace.Elements[elementKey]

		// adjust element index
		element.Index = index

		layer = 0
		for _, task := range element.Tasks {
			taskIndices[task.UUID] = index
			taskLayer[task.UUID]   = layer
			task.Layer             = layer
			layer                  = layer +1
		}
		for _, task := range element.Tasks {
			task.Layers           = layer
			taskLayers[task.UUID] = layer
		}

		index = index + 1

		// loop over all clusters
		clusterKeys := make([]string, 0, len(element.Clusters))
		for clusterKey := range element.Clusters {
			clusterKeys = append(clusterKeys, clusterKey)
		}
		sort.Slice(clusterKeys, func(i, j int) bool {return clusterKeys[i] < clusterKeys[j]})

		for _, clusterKey := range clusterKeys {
			cluster := element.Clusters[clusterKey]

			// adjust cluster index
			cluster.Index = index

			layer = 0
			for _, task := range cluster.Tasks {
				taskIndices[task.UUID] = index
				taskLayer[task.UUID]   = layer
				task.Layer             = layer
				layer                  = layer +1
			}
			for _, task := range cluster.Tasks {
				task.Layers           = layer
				taskLayers[task.UUID] = layer
			}

			index = index + 1

			// loop over all instances
			instanceKeys := make([]string, 0, len(cluster.Instances))
			for instanceKey := range cluster.Instances {
				instanceKeys = append(instanceKeys, instanceKey)
			}

			sort.Slice(instanceKeys, func(i, j int) bool {return instanceKeys[i] < instanceKeys[j]})

			for _, instanceKey := range instanceKeys {
				instance := cluster.Instances[instanceKey]

				// adjust instance index
				instance.Index = index

				layer = 0
				for _, task := range instance.Tasks {
					taskIndices[task.UUID] = index
					taskLayer[task.UUID]   = layer
					task.Layer             = layer
					layer                  = layer +1
				}
				for _, task := range instance.Tasks {
					task.Layers           = layer
					taskLayers[task.UUID] = layer
				}

				index = index + 1
			} // end of instance loop
		} // end of cluster loop
	} // end of element loop

	// adjust indices of events
	trace.Min = 9223372036854775807
	trace.Max = 0
	for _, event := range trace.Events {
		if event.Task1 != "" {
			event.Index1  = taskIndices[event.Task1]
			event.Layer1  = taskLayer[event.Task1]
			event.Layers1 = taskLayers[event.Task1]
		}
		event.Index2  = taskIndices[event.Task2]
		event.Layer2  = taskLayer[event.Task2]
		event.Layers2 = taskLayers[event.Task2]

		if trace.Min > event.Time {
			trace.Min = event.Time
		}
		if trace.Max < event.Time {
			trace.Max = event.Time
		}
	}

	if trace.Min > trace.Max {
		trace.Min = 0
		trace.Max = 0
	}

	trace.Range = trace.Max - trace.Min

	// normalize task events
	for _, event := range trace.Events {
		event.Time = event.Time - trace.Min
	}

	// normalize tasks

	for _, task := range trace.Tasks {
		task.Started   = task.Started   - trace.Min
		task.Completed = task.Completed - trace.Min
		task.Latest    = task.Latest    - trace.Min
	}
	// loop over all elements
	for _, element := range trace.Elements {
		// update tasks
		for _, task := range element.Tasks {
			task.Started   = task.Started   - trace.Min
			task.Completed = task.Completed - trace.Min
			task.Latest    = task.Latest    - trace.Min
		}
		// loop over all clusters
		for _, cluster := range element.Clusters {
			// update tasks
			for _, task := range cluster.Tasks {
				task.Started   = task.Started   - trace.Min
				task.Completed = task.Completed - trace.Min
				task.Latest    = task.Latest    - trace.Min
			}
			// loop over all instances
			for _, instance := range cluster.Instances {
				// update tasks
				for _, task := range instance.Tasks {
					task.Started   = task.Started   - trace.Min
					task.Completed = task.Completed - trace.Min
					task.Latest    = task.Latest    - trace.Min
				}
			} // end instance loop
		} // end cluster loop
	} // end element loop
}

//------------------------------------------------------------------------------

// TaskCommand executes the task related subcommands
func TaskCommand(context *ishell.Context, m *model.Model) {
	// check if the action has been defined
	if len(context.Args) < 1 {
		TaskUsage(true, context)
		return
	}

	// determine the required action
	action := context.Args[0]

	// handle required action
	switch action {
	case "?":
		TaskUsage(true, context)
	case _list:
		// check availability of arguments
		if len(context.Args) < 2 || 6 < len(context.Args) {
			TaskUsage(true, context)
			return
		}

		// set domain name filter
		domainName := context.Args[1]

		// set solution name filter
		solutionName := ""
		if len(context.Args) >= 3 {
			solutionName = context.Args[2]
		}

		// set element name filter
		elementName := ""
		if len(context.Args) >= 4 {
			elementName = context.Args[3]
		}

		// set cluster name filter
		clusterName := ""
		if len(context.Args) >= 5 {
			clusterName = context.Args[4]
		}

		// set instance name filter
		instanceName := ""
		if len(context.Args) >= 6 {
			instanceName = context.Args[5]
		}

		// determine domain
		domain, err := model.GetDomain(domainName)
		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// determine list of solution names
		tasks := []string{}

		tNames, _ := domain.ListTasks()
		for _, tName := range tNames {
			task, _ := domain.GetTask(tName)

			if (solutionName == task.Solution) &&
				 (elementName  == task.Element)  &&
				 (clusterName  == task.Cluster)  &&
				 (instanceName == task.Instance) {
			  tasks = append(tasks, tName)
			}
		}

		result, err := util.ConvertToYAML(tasks)
		handleResult(context, err, "tasks could not be listed", result)
	case _get:
		// check availability of arguments
		if len(context.Args) < 3 || 4 < len(context.Args){
			TaskUsage(true, context)
			return
		}

		// determine task
		task, err := model.GetTask(context.Args[1], context.Args[2])

		if err != nil {
			handleResult(context, err, "task can not be identified", "")
			return
		}

		// determine level
		level := 0
		if len(context.Args) == 4 {
			value, err := strconv.Atoi(context.Args[3])
			if err != nil {
				handleResult(context, err, "invalid level (not an integer)", "")
				return
			}

			if value < 0 {
				handleResult(context, err, "invalid level (must not be negative)", "")
				return
			}
			level = value
		}

		// execute the command
		taskinfo := NewTaskInfo(task, level)
		result, err := util.ConvertToYAML(taskinfo)
		handleResult(context, err, "task can not be displayed", result)
	case _terminate:
		// check availability of arguments
		if len(context.Args) != 3 {
			TaskUsage(true, context)
			return
		}

		// determine task
		task, err := model.GetTask(context.Args[1], context.Args[2])

		if err != nil {
			handleResult(context, err, "task can not be identified", "")
			return
		}

		// execute the command
		// get event channel
		channel := engine.GetEventChannel()

		// create event
		channel <- model.NewEvent(context.Args[1], task.UUID, model.EventTypeTaskTermination, "", "initial")

		handleResult(context, nil, "task can not be terminated", "")
	case _trace:
		// check availability of arguments
		if len(context.Args) != 3 {
			TaskUsage(true, context)
			return
		}

		// determine task
		task, err := model.GetTask(context.Args[1], context.Args[2])

		if err != nil {
			handleResult(context, err, "task can not be identified", "")
			return
		}

		// construct the trace object
		trace := NewTrace(task)

		// finished
		result, _ := util.ConvertToYAML(trace)
		handleResult(context, nil, "trace could not be created", result)
	default:
		TaskUsage(true, context)
	}
}

//------------------------------------------------------------------------------

// TaskUsage describes how to make use of the subcommand
func TaskUsage(header bool, context *ishell.Context) {
	info := ""
	if header {
		info = _usage
	}
	info += "  task list <domain> <solution> <element> <cluster> <instance>\n"
	info += "       get <domain> <task> <level>\n"
	info += "       terminate <domain> <task>\n"
	info += "       trace <domain> <task>\n"

  writeInfo(context, info)
}

//------------------------------------------------------------------------------
