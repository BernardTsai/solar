package model

import (
	"sync"
	"errors"
	"strings"

	"tsai.eu/solar/util"
)

//------------------------------------------------------------------------------
// Domain
// ======
//
// Attributes:
//   - Name
//   - Templates
//   - Architectures
//   - Components
//   - Tasks
//   - Events
//
// Functions:
//   - NewDomain
//
//   - domain.Show
//   - domain.Load
//   - domain.Load2
//   - domain.Save
//
//   - domain.ListComponents
//   - domain.GetComponents
//   - domain.GetComponent
//   - domain.AddComponent
//   - domain.DeleteComponent
//
//   - domain.ListArchitectures
//   - domain.GetArchitecture
//   - domain.AddArchitecture
//   - domain.DeleteArchitecture
//
//   - domain.ListSolutions
//   - domain.GetSolution
//   - domain.AddSolution
//   - domain.DeleteSolution
//
//   - domain.ListTasks
//   - domain.GetTask
//   - domain.AddTask
//   - domain.DeleteTask
//
//   - domain.ListEvents
//   - domain.GetEvent
//   - domain.AddEvent
//   - domain.DeleteEvent
//
//   - domain.ListControllers
//   - domain.GetController
//   - domain.AddController
//   - domain.DeleteController
//------------------------------------------------------------------------------

// Domain describes all artefacts managed with an administrative realm.
type Domain struct {
	Name           string                   `yaml:"Name"`                     // name of the domain
	Components     map[string]*Component    `yaml:"Components"`               // map of components
	ComponentsX    sync.RWMutex             `yaml:"ComponentsX,omitempty"`    // mutex for components
	Architectures  map[string]*Architecture `yaml:"Architectures"`            // map of architectures
	ArchitecturesX sync.RWMutex             `yaml:"ArchitecturesX,omitempty"` // mutex for architectures
	Solutions      map[string]*Solution     `yaml:"Solutions"`                // list of solutions
	SolutionsX     sync.RWMutex             `yaml:"SolutionsX,omitempty"`     // mutex for solutions
	Tasks          map[string]*Task         `yaml:"Tasks"`                    // list of tasks
	TasksX         sync.RWMutex             `yaml:"TasksX,omitempty"`         // mutex for tasks
	Events         map[string]*Event        `yaml:"Events"`                   // list of events
	EventsX        sync.RWMutex             `yaml:"EventsX,omitempty"`        // mutex for events
	Controllers    map[string]*Controller   `yaml:"Controllers"`              // list of controllers
	ControllersX   sync.RWMutex             `yaml:"ControllersX,omitempty"`   // mutex for controllers
}

//------------------------------------------------------------------------------

// NewDomain creates a new domain
func NewDomain(name string) (*Domain, error) {
	var domain Domain

	domain.Name           = name
	domain.Components     = map[string]*Component{}
	domain.ComponentsX    = sync.RWMutex{}
	domain.Architectures  = map[string]*Architecture{}
	domain.ArchitecturesX = sync.RWMutex{}
	domain.Solutions      = map[string]*Solution{}
	domain.SolutionsX     = sync.RWMutex{}
	domain.Tasks          = map[string]*Task{}
	domain.TasksX         = sync.RWMutex{}
	domain.Events         = map[string]*Event{}
	domain.EventsX        = sync.RWMutex{}
	domain.Controllers    = map[string]*Controller{}
	domain.ControllersX   = sync.RWMutex{}

	// add internal default controller
	ctrl, _ := NewController("Internal", "V1.0.0")
	ctrl.Status = ActiveState

	domain.Controllers["Internal:V1.0.0"] = ctrl

	// add default controllers listed in the configuration file
	configuration, _ := util.GetConfiguration()

	for _, controllerImage := range configuration.CONTROLLERS {
		controller, _ := NewController(util.UUID(),"V0.0.0")
		controller.Image = controllerImage

		domain.AddController(controller)
	}

	// success
	return &domain, nil
}

//------------------------------------------------------------------------------

// Show displays the domain information as json
func (domain *Domain) Show() (string, error) {
	return util.ConvertToYAML(domain)
}

//------------------------------------------------------------------------------

// Save writes the domain as json data to a file
func (domain *Domain) Save(filename string) error {
	return util.SaveYAML(filename, domain)
}

//------------------------------------------------------------------------------

// Load reads the domain from a file
func (domain *Domain) Load(filename string) error {
	return util.LoadYAML(filename, domain)
}

//------------------------------------------------------------------------------

// Load2 imports a yaml model
func (domain *Domain) Load2(yaml string) error {
	return util.ConvertFromYAML(yaml, domain)
}

//------------------------------------------------------------------------------

// ListComponents lists all components of a domain
func (domain *Domain) ListComponents() ([][2]string, error) {
	// collect names
	components := [][2]string{}

	domain.ComponentsX.RLock()
	for _, component := range domain.Components {
		components = append(components, [2]string{component.Component, component.Version})
	}
	domain.ComponentsX.RUnlock()

	// success
	return components, nil
}

//------------------------------------------------------------------------------

// GetComponents retrieves all components
func (domain *Domain) GetComponents() ([]*Component, error) {
	// construct result
	domain.ComponentsX.RLock()

	// iterate over all components
	result := make( []*Component, len(domain.Components) )
	index  := 0
	for _, value := range domain.Components {
		result[index] = value
		index++
	}

	domain.ComponentsX.RUnlock()

	// success
	return result, nil
}

//------------------------------------------------------------------------------

// GetComponent retrieves a component by name and version
func (domain *Domain) GetComponent(name string, version string) (*Component, error) {
	// determine template
	domain.ComponentsX.RLock()
	component, ok := domain.Components[name + " - " + version]
	domain.ComponentsX.RUnlock()

	if !ok {
		return nil, errors.New("component not found")
	}

	// success
	return component, nil
}

//------------------------------------------------------------------------------

// AddComponent adds a component to a domain
func (domain *Domain) AddComponent(component *Component) error {
	// check if component has already been defined
	domain.ComponentsX.RLock()
	_, ok := domain.Components[component.Component + " - " + component.Version]
	domain.ComponentsX.RUnlock()

	if ok {
		return errors.New("component already exists")
	}

	domain.ComponentsX.Lock()
	domain.Components[component.Component + " - " + component.Version] = component
	domain.ComponentsX.Unlock()

	// success
	return nil
}

//------------------------------------------------------------------------------

// DeleteComponent deletes a component
func (domain *Domain) DeleteComponent(name string, version string) error {
	// determine component
	domain.ComponentsX.RLock()
	_, ok := domain.Components[name + " - " + version]
	domain.ComponentsX.RUnlock()

	if !ok {
		return errors.New("component not found")
	}

	// remove template
	domain.ComponentsX.Lock()
	delete(domain.Components, name + " - " + version)
	domain.ComponentsX.Unlock()

	// success
	return nil
}

//------------------------------------------------------------------------------

// ListArchitectures lists all architectures of a domain
func (domain *Domain) ListArchitectures() ([][2]string, error) {
	// collect names
	architectures := [][2]string{}

	domain.ArchitecturesX.RLock()
	for _, architecture := range domain.Architectures {
		architectures = append(architectures, [2]string{architecture.Architecture, architecture.Version})
	}
	domain.ArchitecturesX.RUnlock()

	// success
	return architectures, nil
}

//------------------------------------------------------------------------------

// GetArchitecture get an architecture by name
func (domain *Domain) GetArchitecture(name string, version string) (*Architecture, error) {
	// determine architecture
	domain.ArchitecturesX.RLock()
	architecture, ok := domain.Architectures[name + " - " + version]
	domain.ArchitecturesX.RUnlock()

	if !ok {
		return nil, errors.New("architecture not found")
	}

	// success
	return architecture, nil
}

//------------------------------------------------------------------------------

// AddArchitecture add architecture to a domain
func (domain *Domain) AddArchitecture(architecture *Architecture) error {
	// determine domain
	domain.ArchitecturesX.RLock()
	_, ok := domain.Architectures[architecture.Architecture + " - " + architecture.Version]
	domain.ArchitecturesX.RUnlock()

	if ok {
		// return errors.New("architecture already exists")
	}

	domain.ArchitecturesX.Lock()
	domain.Architectures[architecture.Architecture + " - " + architecture.Version] = architecture
	domain.ArchitecturesX.Unlock()

	// success
	return nil
}

//------------------------------------------------------------------------------

// DeleteArchitecture deletes a architecture
func (domain *Domain) DeleteArchitecture(name string, version string) error {
	// determine architecture
	domain.ArchitecturesX.RLock()
	_, ok := domain.Architectures[name + " - " + version]
	domain.ArchitecturesX.RUnlock()

	if !ok {
		return errors.New("architecture not found")
	}

	// remove architecture
	domain.ArchitecturesX.Lock()
	delete(domain.Architectures, name + " - " + version)
	domain.ArchitecturesX.Unlock()

	// success
	return nil
}

//------------------------------------------------------------------------------

// ListSolutions lists all solutions of a domain
func (domain *Domain) ListSolutions() ([]string, error) {
	// collect names
	solutions := []string{}

	domain.SolutionsX.RLock()
	for solution := range domain.Solutions {
		solutions = append(solutions, solution)
	}
	domain.SolutionsX.RUnlock()

	// success
	return solutions, nil
}

//------------------------------------------------------------------------------

// GetSolution get a solution by name
func (domain *Domain) GetSolution(name string) (*Solution, error) {
	// determine solution
	domain.SolutionsX.RLock()
	solution, ok := domain.Solutions[name]
	domain.SolutionsX.RUnlock()

	if !ok {
		return nil, errors.New("solution not found")
	}

	// success
	return solution, nil
}

//------------------------------------------------------------------------------

// AddSolution adds a solution to a domain
func (domain *Domain) AddSolution(solution *Solution) error {
	// check if solution has already been defined
	domain.SolutionsX.RLock()
	_, ok := domain.Solutions[solution.Solution]
	domain.SolutionsX.RUnlock()

	if ok {
		return errors.New("solution already exists")
	}

	domain.SolutionsX.Lock()
	domain.Solutions[solution.Solution] = solution
	domain.SolutionsX.Unlock()

	// success
	return nil
}

//------------------------------------------------------------------------------

// DeleteSolution deletes a solution
func (domain *Domain) DeleteSolution(name string) error {
	// determine solution
	domain.SolutionsX.RLock()
	_, ok := domain.Solutions[name]
	domain.SolutionsX.RUnlock()

	if !ok {
		return errors.New("solution not found")
	}

	// remove solution
	domain.SolutionsX.Lock()
	delete(domain.Solutions, name)
	domain.SolutionsX.Unlock()

	// success
	return nil
}

//------------------------------------------------------------------------------

// ListTasks all tasks of a domain
func (domain *Domain) ListTasks() ([]string, error) {
	// collect names
	tasks := []string{}

	domain.TasksX.RLock()
	for task := range domain.Tasks {
		tasks = append(tasks, task)
	}
	domain.TasksX.RUnlock()

	// success
	return tasks, nil
}

//------------------------------------------------------------------------------

// GetTask get a task by name
func (domain *Domain) GetTask(name string) (*Task, error) {
	// determine task
	domain.TasksX.RLock()
	task, ok := domain.Tasks[name]
	domain.TasksX.RUnlock()

	if !ok {
		return nil, errors.New("task not found")
	}

	// success
	return task, nil
}

//------------------------------------------------------------------------------

// AddTask adds a task to a domain
func (domain *Domain) AddTask(task *Task) error {
	// check if task has already been defined
	domain.TasksX.RLock()
	_, ok := domain.Tasks[task.GetUUID()]
	domain.TasksX.RUnlock()

	if ok {
		return errors.New("task already exists")
	}

	domain.TasksX.Lock()
	domain.Tasks[task.GetUUID()] = task
	domain.TasksX.Unlock()

	// success
	return nil
}

//------------------------------------------------------------------------------

// DeleteTask deletes a task
func (domain *Domain) DeleteTask(uuid string) error {
	// determine task
	domain.TasksX.RLock()
	_, ok := domain.Tasks[uuid]
	domain.TasksX.RUnlock()

	if !ok {
		return errors.New("task not found")
	}

	// remove task
	domain.TasksX.Lock()
	delete(domain.Tasks, uuid)
	domain.TasksX.Unlock()

	// success
	return nil
}

//------------------------------------------------------------------------------

// ListEvents all events of a domain
func (domain *Domain) ListEvents() ([]string, error) {
	// collect names
	events := []string{}

	domain.EventsX.RLock()
	for event := range domain.Events {
		events = append(events, event)
	}
	domain.EventsX.RUnlock()

	// success
	return events, nil
}

//------------------------------------------------------------------------------

// GetEvent get a event by name
func (domain *Domain) GetEvent(uuid string) (*Event, error) {
	// determine event
	domain.EventsX.RLock()
	event, ok := domain.Events[uuid]
	domain.EventsX.RUnlock()

	if !ok {
		return nil, errors.New("event not found")
	}

	// success
	return event, nil
}

//------------------------------------------------------------------------------

// AddEvent adds a event to a domain
func (domain *Domain) AddEvent(event *Event) error {
	// check if event has already been defined
	domain.EventsX.RLock()
	_, ok := domain.Events[event.UUID]
	domain.EventsX.RUnlock()

	if ok {
		return errors.New("event already exists")
	}

	domain.EventsX.Lock()
	domain.Events[event.UUID] = event
	domain.EventsX.Unlock()

	// register with tasks
	domain.TasksX.Lock()
	task := domain.Tasks[event.Task]
	task.AddEvent(event)
	domain.TasksX.Unlock()

	// register with tasks
	if event.Source != "" {
		domain.TasksX.Lock()
		task = domain.Tasks[event.Source]
		task.AddEvent(event)
		domain.TasksX.Unlock()
	}

	// success
	return nil
}

//------------------------------------------------------------------------------

// DeleteEvent deletes an event
func (domain *Domain) DeleteEvent(uuid string) error {
	// determine event
	domain.EventsX.RLock()
	_, ok := domain.Events[uuid]
	domain.EventsX.RUnlock()

	if !ok {
		return errors.New("event not found")
	}

	// remove event
	domain.EventsX.Lock()
	delete(domain.Events, uuid)
	domain.EventsX.Unlock()

	// success
	return nil
}

//------------------------------------------------------------------------------

// ListControllers all controllers of a domain
func (domain *Domain) ListControllers() ([][2]string, error) {
	// collect names
	controllers := [][2]string{}

	domain.ControllersX.RLock()
	for controllerNameVersion := range domain.Controllers {
		parts := strings.Split(controllerNameVersion, ":")
		if len(parts) == 2 {
			controllers = append(controllers, [2]string{parts[0], parts[1]})
		} else {
			controllers = append(controllers, [2]string{controllerNameVersion, ""})
		}
	}
	domain.ControllersX.RUnlock()

	// success
	return controllers, nil
}

//------------------------------------------------------------------------------

// GetController get a controller by name
func (domain *Domain) GetController(controllerName string, controllerVersion string) (controller *Controller, err error) {
	var ok bool

	// determine controller
	domain.ControllersX.RLock()
	controller, ok = domain.Controllers[controllerName + ":" + controllerVersion]
	domain.ControllersX.RUnlock()

	if !ok {
		return nil, errors.New("controller not found")
	}

	// success
	return controller, nil
}

//------------------------------------------------------------------------------

// AddController adds a controller to a domain
func (domain *Domain) AddController(controller *Controller) error {
	// check if controller has already been defined
	domain.ControllersX.RLock()
	_, ok := domain.Controllers[controller.Controller + ":" + controller.Version]
	domain.ControllersX.RUnlock()

	if ok {
		return errors.New("controller already exists")
	}

	domain.ControllersX.Lock()
	domain.Controllers[controller.Controller + ":" + controller.Version] = controller
	domain.ControllersX.Unlock()

	// success
	return nil
}

//------------------------------------------------------------------------------

// DeleteController deletes an controller
func (domain *Domain) DeleteController(controller string, version string) error {
	// determine controller
	domain.ControllersX.RLock()
	_, ok := domain.Controllers[controller + ":" + version]
	domain.ControllersX.RUnlock()

	if !ok {
		return errors.New("controller not found")
	}

	// remove controller
	domain.ControllersX.Lock()
	delete(domain.Controllers, controller + ":" + version)
	domain.ControllersX.Unlock()

	// success
	return nil
}

//------------------------------------------------------------------------------
