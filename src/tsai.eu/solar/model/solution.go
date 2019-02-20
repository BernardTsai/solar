package model

import (
	"sync"

	"github.com/pkg/errors"
	"tsai.eu/solar/util"
)

//------------------------------------------------------------------------------

// TransitionTable is map of allowed transitions
var transitionTable map[string]map[string]string
var transitionTableInit sync.Once

// IsValidStateOrTransition determines if a string resembles a valid state or transition.
func IsValidStateOrTransition(state string) bool {
	switch state {
	case InitialState, CreatingState, DestroyingState, InactiveState, StartingState, StoppingState, ActiveState, ConfiguringState, FailureState, ResettingState, UndefinedState:
		return true
	}
	return false
}

// IsValidState determines if a string resembles a valid state.
func IsValidState(state string) bool {
	switch state {
	case InitialState, InactiveState, ActiveState, FailureState:
		return true
	}
	return false
}

// IsValidTransition determines if a string resembles a valid transition.
func IsValidTransition(transition string) bool {
	switch transition {
	case CreatingState, DestroyingState, StartingState, StoppingState, ConfiguringState, ResettingState:
		return true
	}
	return false
}

// GetTransition determines the required transition between a current state and a target state.
func GetTransition(currentState string, targetState string) (string, error) {
	// initialise singleton once
	transitionTableInit.Do(func() {
		transitionTable = map[string]map[string]string{}

		transitionTable[InitialState] = map[string]string{
			InitialState:  "none",
			InactiveState: "create",
			ActiveState:   "create",
		}
		transitionTable[InactiveState] = map[string]string{
			InitialState:  "destroy",
			InactiveState: "none",
			ActiveState:   "start",
		}
		transitionTable[ActiveState] = map[string]string{
			InitialState:  "stop",
			InactiveState: "stop",
			ActiveState:   "none",
		}
		transitionTable[FailureState] = map[string]string{
			InitialState:  "reset",
			InactiveState: "reset",
			ActiveState:   "reset",
		}
	})

	// check parameters
	if !IsValidState(currentState) || !IsValidState(targetState) {
		return "", errors.New("invalid state")
	}

	// determine transition
	transition, err := GetTransition(currentState, targetState)

	if err != nil {
		return "", errors.New("invalid transition")
	}

	//success
	return transition, nil
}

//------------------------------------------------------------------------------
// Solution
// ========
//
// Attributes:
//   - Solution
//   - Version
//   - Configuration
//   - Elements
//
// Functions:
//   - NewSolution
//
//   - solution.Show
//   - solution.Load
//   - solution.Save
//
//   - solution.ListElements
//   - solution.GetElement
//   - solution.AddElement
//   - solution.DeleteElement
//------------------------------------------------------------------------------

// ElementMap is a synchronized map for a map of elements
type ElementMap struct {
	sync.RWMutex `yaml:"mutex,omitempty"` // mutex
	Map          map[string]*Element     `yaml:"map"` // map of events
}

// MarshalYAML marshals an ElementMap into yaml
func (m ElementMap) MarshalYAML() (interface{}, error) {
	return m.Map, nil
}

// UnmarshalYAML unmarshals an ElementMap from yaml
func (m *ElementMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	Map := map[string]*Element{}

	err := unmarshal(&Map)
	if err != nil {
		return err
	}

	*m = ElementMap{Map: Map}

	return nil
}

//------------------------------------------------------------------------------

// Solution describes the runtime configuration of a solution within a domain.
type Solution struct {
	Solution      string     `yaml:"solution"`       // name of solution
	Version       string     `yaml:"version"`        // version of solution
	Configuration string     `yaml:"configuration"`  // configuration of solution
	Elements      ElementMap `yaml:"elements"`       // elements of solution
}

//------------------------------------------------------------------------------

// NewSolution creates a new solution
func NewSolution(name string, version string, configuration string) (*Solution, error) {
	var solution Solution

	solution.Solution      = name
	solution.Version       = version
	solution.Configuration = configuration
	solution.Elements      = ElementMap{Map: map[string]*Element{}}

	// success
	return &solution, nil
}

//------------------------------------------------------------------------------

// Show displays the solution information as yaml
func (solution *Solution) Show() (string, error) {
	return util.ConvertToYAML(solution)
}

//------------------------------------------------------------------------------

// Save writes the solution as yaml data to a file
func (solution *Solution) Save(filename string) error {
	return util.SaveYAML(filename, solution)
}

//------------------------------------------------------------------------------

// Load reads the solution from a file
func (solution *Solution) Load(filename string) error {
	return util.LoadYAML(filename, solution)
}

//------------------------------------------------------------------------------

// ListElements lists all elements of a solution
func (solution *Solution) ListElements() ([]string, error) {
	// collect names
	elements := []string{}

	if solution != nil {
		solution.Elements.RLock()
		for element := range solution.Elements.Map {
			elements = append(elements, element)
		}
		solution.Elements.RUnlock()
	}

	// success
	return elements, nil
}

//------------------------------------------------------------------------------

// GetElement retrieves an element by name
func (solution *Solution) GetElement(name string) (*Element, error) {
	// determine instance
	solution.Elements.RLock()
	element, ok := solution.Elements.Map[name]
	solution.Elements.RUnlock()

	if !ok {
		return nil, errors.New("element not found")
	}

	// success
	return element, nil
}

//------------------------------------------------------------------------------

// AddElement adds an element to a solution
func (solution *Solution) AddElement(element *Element) error {
	// check if instance has already been defined
	solution.Elements.RLock()
	_, ok := solution.Elements.Map[element.Element]
	solution.Elements.RUnlock()

	if ok {
		return errors.New("element already exists")
	}

	solution.Elements.Lock()
	solution.Elements.Map[element.Element] = element
	solution.Elements.Unlock()

	// success
	return nil
}

//------------------------------------------------------------------------------

// DeleteElement deletes an element
func (solution *Solution) DeleteElement(uuid string) error {
	// determine element
	solution.Elements.RLock()
	_, ok := solution.Elements.Map[uuid]
	solution.Elements.RUnlock()

	if !ok {
		return errors.New("element not found")
	}

	// remove element
	solution.Elements.Lock()
	delete(solution.Elements.Map, uuid)
	solution.Elements.Unlock()

	// success
	return nil
}

//------------------------------------------------------------------------------