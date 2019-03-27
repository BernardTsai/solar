package engine

import (
	"fmt"

	"tsai.eu/solar/model"
)

//------------------------------------------------------------------------------

// Dispatcher receives events from a channel and triggers a task coroutine.
type Dispatcher struct {
	Model   *model.Model           // repository
	Channel chan model.Event       // the channel for event notification
	Active   bool                  // indicates if the dispatcher loop should be active
}

//------------------------------------------------------------------------------

// StartDispatcher creates a dispatcher and returns a channel for new tasks.
func StartDispatcher(m *model.Model) (*Dispatcher) {

	// create the communication channel
	channel := GetEventChannel()

	// create the dispatcher
	dispatcher := Dispatcher{
		Model:   m,
		Channel: channel,
		Active:  true,
	}

	// start the dispatcher
	go dispatcher.Run()

	return &dispatcher
}

//------------------------------------------------------------------------------

// Run starts the dispatcher loop receiving events and triggering tasks.
func (d *Dispatcher) Run() {
	// loop until exit is requested
	for d.Active {
		// get next event
		event := <-d.Channel

		// terminate if domain is empty = exit request
		if event.Domain == "" {
			return
		}

		// get corresponding domain from the model
		domain, err := d.Model.GetDomain(event.Domain)
		if err != nil {
			// TODO: log unknown domain
			continue
		}

		// save event
		domain.AddEvent(&event)

		// get task
		task, err := domain.GetTask(event.Task)
		if err != nil {
			fmt.Println(err)
			// TODO: log unknown task
			continue
		}

		// determine action by type of event
		// Event types: execute, completed, failed, timeout, terminate
		// Task types can be:
		// - set component state
		// - set instance state
		// - transition component
		// - transition instance
		// - parallel execute tasks
		// - sequentially execute tasks
		switch event.Type {
		// execute the task
		case model.EventTypeTaskExecution:
			go task.Execute()

		// handle task completion
		case model.EventTypeTaskCompletion:
			go task.Completed()

		// handle task failure
		case model.EventTypeTaskFailure:
			go task.Failed()

		// handle timeout of a task
		case model.EventTypeTaskTimeout:
			go task.Timeout()

		// handle termination of a task
		case model.EventTypeTaskTermination:
			go task.Terminate()
		}
	}
}

//------------------------------------------------------------------------------

// Stop will flag the dispatcher to stop execution
func (d *Dispatcher) Stop() {
  d.Active = false
}

//------------------------------------------------------------------------------
