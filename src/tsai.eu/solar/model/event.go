package model

import (
	"time"

	"github.com/google/uuid"
	"tsai.eu/solar/util"
)

//------------------------------------------------------------------------------
// Event
// =====
//
// Attributes:
//   - Domain
//   - UUID
//   - Task
//   - type
//   - Source
//   - Comment
//
// Functions:
//   - NewEvent
//
//   - event.Show
//   - event.Load
//   - event.Save
//------------------------------------------------------------------------------

// Event describes a situation which may trigger further tasks.
type Event struct {
	Domain  string `yaml:"Domain"`  // domain of event
	UUID    string `yaml:"UUID"`    // uuid of event
	Task    string `yaml:"Task"`    // uuid of task
	Type    string `yaml:"Type"`    // type of event: "execution", "completion", "failure"
	Source  string `yaml:"Source"`  // source of the event (uuid of the task or "")
	Time    int64  `yaml:"Time"`    // time since 1.1.1970 in nsecs
	Comment string `yaml:"Comment"` // comments
}

//------------------------------------------------------------------------------

// NewEvent creates a new event
func NewEvent(domain string, task string, etype string, source string, comment string) Event {
	var event Event

	event.Domain  = domain
	event.UUID    = uuid.New().String()
	event.Task    = task
	event.Type    = etype
	event.Source  = source
	event.Time    = time.Now().UnixNano()
	event.Comment = comment

	// success
	return event
}

//------------------------------------------------------------------------------

// Show displays the event information as json
func (event *Event) Show() (string, error) {
	return util.ConvertToYAML(event)
}

//------------------------------------------------------------------------------

// Save writes the event as json data to a file
func (event *Event) Save(filename string) error {
	return util.SaveYAML(filename, event)
}

//------------------------------------------------------------------------------

// Load reads the event from a file
func (event *Event) Load(filename string) error {
	return util.LoadYAML(filename, event)
}

//------------------------------------------------------------------------------

// GetUUID delivers the universal unique identifier of the event
func (event *Event) GetUUID() string {
	return event.UUID
}

//------------------------------------------------------------------------------
