package model

import (
	"github.com/pkg/errors"
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
//   - domain.Save
//
//   - domain.ListComponents
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
//------------------------------------------------------------------------------

// ComponentMap is a synchronized map for a map of components
type ComponentMap struct {
	Map map[string]*Component `yaml:"map"`             // map of components
}

// MarshalYAML marshals a ComponentMap into yaml
func (m ComponentMap) MarshalYAML() (interface{}, error) {
	return m.Map, nil
}

// UnmarshalYAML unmarshals a ComponentMap from yaml
func (m *ComponentMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	Map := map[string]*Component{}

	err := unmarshal(&Map)
	if err != nil {
		return err
	}

	*m = ComponentMap{Map: Map}

	return nil
}

//------------------------------------------------------------------------------

// ArchitectureMap is a synchronized map for a map of architectures
type ArchitectureMap struct {
	Map map[string]*Architecture `yaml:"map"`             // map of architectures
}

// MarshalYAML marshals a ArchitectureMap into yaml
func (m ArchitectureMap) MarshalYAML() (interface{}, error) {
	return m.Map, nil
}

//------------------------------------------------------------------------------

// SolutionMap is a synchronized map for a map of solutions
type SolutionMap struct {
	Map map[string]*Solution  `yaml:"map"`             // map of solutions
}

// MarshalYAML marshals a SolutionMap into yaml
func (m SolutionMap) MarshalYAML() (interface{}, error) {
	return m.Map, nil
}

// UnmarshalYAML unmarshals an SolutionMap from yaml
func (m *SolutionMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	Map := map[string]*Solution{}

	err := unmarshal(&Map)
	if err != nil {
		return err
	}

	*m = SolutionMap{Map: Map}

	return nil
}

//------------------------------------------------------------------------------

// TaskMap is a synchronized map for a map of tasks
type TaskMap struct {
	Map map[string]*Task  `yaml:"map"`             // map of tasks
}

// MarshalYAML marshals a TaskMap into yaml
func (m TaskMap) MarshalYAML() (interface{}, error) {
	return m.Map, nil
}

//------------------------------------------------------------------------------

// EventMap is a synchronized map for a map of events
type EventMap struct {
	Map map[string]*Event `yaml:"map"`             // map of events
}

// MarshalYAML marshals a EventMap into yaml
func (m EventMap) MarshalYAML() (interface{}, error) {
	return m.Map, nil
}

//------------------------------------------------------------------------------

// Domain describes all artefacts managed with an administrative realm.
type Domain struct {
	Name          string          `yaml:"name"`          // name of the domain
	Components    ComponentMap    `yaml:"components"`    // map of components
	Architectures ArchitectureMap `yaml:"architectures"` // map of architectures
	Solutions     SolutionMap     `yaml:"solutions"`     // list of solutions
	Tasks         TaskMap         `yaml:"tasks"`         // list of tasks
	Events        EventMap        `yaml:"events"`        // list of events
}

//------------------------------------------------------------------------------

// NewDomain creates a new domain
func NewDomain(name string) (*Domain, error) {
	var domain Domain

	domain.Name = name
	domain.Components = ComponentMap{Map: map[string]*Component{}}
	domain.Architectures = ArchitectureMap{Map: map[string]*Architecture{}}
	domain.Solutions = SolutionMap{Map: map[string]*Solution{}}
	domain.Tasks = TaskMap{Map: map[string]*Task{}}
	domain.Events = EventMap{Map: map[string]*Event{}}

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

// ListComponents lists all components of a domain
func (domain *Domain) ListComponents() ([]string, error) {
	// collect names
	components := []string{}

	for component := range domain.Components.Map {
		components = append(components, component)
	}

	// success
	return components, nil
}

//------------------------------------------------------------------------------

// GetComponent retrieves a component by name
func (domain *Domain) GetComponent(name string) (*Component, error) {
	// determine template
	component, ok := domain.Components.Map[name]

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
	_, ok := domain.Components.Map[component.Component]

	if ok {
		return errors.New("component already exists")
	}

	domain.Components.Map[component.Component] = component

	// success
	return nil
}

//------------------------------------------------------------------------------

// DeleteComponent deletes a component
func (domain *Domain) DeleteComponent(name string) error {
	// determine component
	_, ok := domain.Components.Map[name]

	if !ok {
		return errors.New("component not found")
	}

	// remove template
	delete(domain.Components.Map, name)

	// success
	return nil
}

//------------------------------------------------------------------------------

// ListArchitectures lists all architectures of a domain
func (domain *Domain) ListArchitectures() ([]string, error) {
	// collect names
	architectures := []string{}

	for architecture := range domain.Architectures.Map {
		architectures = append(architectures, architecture)
	}

	// success
	return architectures, nil
}

//------------------------------------------------------------------------------

// GetArchitecture get an architecture by name
func (domain *Domain) GetArchitecture(name string) (*Architecture, error) {
	// determine architecture
	architecture, ok := domain.Architectures.Map[name]

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
	_, ok := domain.Architectures.Map[architecture.Architecture]

	if ok {
		return errors.New("architecture already exists")
	}

	domain.Architectures.Map[architecture.Architecture] = architecture

	// success
	return nil
}

//------------------------------------------------------------------------------

// DeleteArchitecture deletes a architecture
func (domain *Domain) DeleteArchitecture(name string) error {
	// determine architecture
	_, ok := domain.Architectures.Map[name]

	if !ok {
		return errors.New("architecture not found")
	}

	// remove architecture
	delete(domain.Architectures.Map, name)

	// success
	return nil
}

//------------------------------------------------------------------------------

// ListSolutions lists all solutions of a domain
func (domain *Domain) ListSolutions() ([]string, error) {
	// collect names
	solutions := []string{}

	for solution := range domain.Solutions.Map {
		solutions = append(solutions, solution)
	}

	// success
	return solutions, nil
}

//------------------------------------------------------------------------------

// GetSolution get a solution by name
func (domain *Domain) GetSolution(name string) (*Solution, error) {
	// determine solution
	solution, ok := domain.Solutions.Map[name]

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
	_, ok := domain.Solutions.Map[solution.Solution]

	if ok {
		return errors.New("solution already exists")
	}

	domain.Solutions.Map[solution.Solution] = solution

	// success
	return nil
}

//------------------------------------------------------------------------------

// DeleteSolution deletes a solution
func (domain *Domain) DeleteSolution(name string) error {
	// determine solution
	_, ok := domain.Solutions.Map[name]

	if !ok {
		return errors.New("solution not found")
	}

	// remove solution
	delete(domain.Solutions.Map, name)

	// success
	return nil
}

//------------------------------------------------------------------------------

// ListTasks all tasks of a domain
func (domain *Domain) ListTasks() ([]string, error) {
	// collect names
	tasks := []string{}

	for task := range domain.Tasks.Map {
		tasks = append(tasks, task)
	}

	// success
	return tasks, nil
}

//------------------------------------------------------------------------------

// GetTask get a task by name
func (domain *Domain) GetTask(name string) (*Task, error) {
	// determine task
	task, ok := domain.Tasks.Map[name]

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
	_, ok := domain.Tasks.Map[task.GetUUID()]

	if ok {
		return errors.New("task already exists")
	}

	domain.Tasks.Map[task.GetUUID()] = task

	// success
	return nil
}

//------------------------------------------------------------------------------

// DeleteTask deletes a task
func (domain *Domain) DeleteTask(uuid string) error {
	// determine task
	_, ok := domain.Tasks.Map[uuid]

	if !ok {
		return errors.New("task not found")
	}

	// remove task
	delete(domain.Tasks.Map, uuid)

	// success
	return nil
}

//------------------------------------------------------------------------------

// ListEvents all events of a domain
func (domain *Domain) ListEvents() ([]string, error) {
	// collect names
	events := []string{}

	for event := range domain.Events.Map {
		events = append(events, event)
	}

	// success
	return events, nil
}

//------------------------------------------------------------------------------

// GetEvent get a event by name
func (domain *Domain) GetEvent(uuid string) (*Event, error) {
	// determine event
	event, ok := domain.Events.Map[uuid]

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
	_, ok := domain.Events.Map[event.UUID]

	if ok {
		return errors.New("event already exists")
	}

	domain.Events.Map[event.UUID] = event

	// success
	return nil
}

//------------------------------------------------------------------------------

// DeleteEvent deletes an event
func (domain *Domain) DeleteEvent(uuid string) error {
	// determine event
	_, ok := domain.Events.Map[uuid]

	if !ok {
		return errors.New("event not found")
	}

	// remove event
	delete(domain.Events.Map, uuid)

	// success
	return nil
}

//------------------------------------------------------------------------------
