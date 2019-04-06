package engine

import (
	"context"
	"fmt"
	"time"

	"tsai.eu/solar/model"
	"tsai.eu/solar/util"
	"tsai.eu/solar/controller"
)

//------------------------------------------------------------------------------

// Dispatcher receives events from a channel and triggers a task coroutine.
type Dispatcher struct {
	Channel chan model.Event       // the channel for event notification
}

//------------------------------------------------------------------------------

// Start creates a dispatcher and returns a channel for new tasks.
func Start(c context.Context) (*Dispatcher) {
	// create the dispatcher
	dispatcher := Dispatcher{
		Channel: GetEventChannel(),
	}

	// preload the controllers
	controller.GetController("default")

	// start the dispatcher
	go dispatcher.Run(c)

	util.LogInfo("main", "ENG", "engine active")

	return &dispatcher
}

//------------------------------------------------------------------------------

// Run starts the dispatcher loop receiving events and triggering tasks.
func (d *Dispatcher) Run(ctx context.Context) {
	// loop forever until context is shut down
	for {
		select {
		// context has been shut down
		case <-ctx.Done():
			return
		// get next event
		case event := <-d.Channel:
			// terminate if domain is empty = exit request
			if event.Domain == "" {
				return
			}

			// get corresponding domain from the model
			domain, err := model.GetDomain(event.Domain)
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
				// monitor execution of new tasks
				if task.GetStatus() == model.TaskStatusInitial {
					go monitorTask(ctx, task, d.Channel)
				}
				go task.Execute(ctx)

			// handle task completion
			case model.EventTypeTaskCompletion:
				go task.Completed(ctx)

			// handle task failure
			case model.EventTypeTaskFailure:
				go task.Failed(ctx)

			// handle timeout of a task
			case model.EventTypeTaskTimeout:
				go task.Timeout(ctx)

			// handle termination of a task
			case model.EventTypeTaskTermination:
				go task.Terminate(ctx)
			}
		}
	}
}

//------------------------------------------------------------------------------

// monitorTask creates a context for timeout and cancelation of a task
func monitorTask(ctx context.Context, task *model.Task, channel chan model.Event) {
	// derive new timeout context
	monitorCtx, cancel := context.WithTimeout(ctx, 10 * time.Second)
	defer cancel()

	select {
	case <- monitorCtx.Done():
		// check status of task
		status := task.GetStatus()

		if status != model.TaskStatusInitial && status != model.TaskStatusExecuting {
			return
		}

		// task may still be active
		switch monitorCtx.Err().Error() {
		case "context canceled":               // termination of processes
			channel <- model.NewEvent(task.Domain, task.UUID, model.EventTypeTaskTermination, task.UUID, "termination")
		default:                               // timeout
			channel <- model.NewEvent(task.Domain, task.UUID, model.EventTypeTaskTimeout, task.UUID, "timeout")
		}
	}
}

//------------------------------------------------------------------------------
