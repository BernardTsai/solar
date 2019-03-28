package monitor

import (
  "time"

  "tsai.eu/solar/model"
  "tsai.eu/solar/engine"
)

//------------------------------------------------------------------------------

// Monitor validates solutions and triggers tasks to converge to the desired target state.
type Monitor struct {
	Model   *model.Model           // repository
  Channel  chan model.Event      // the channel for event notification
  Active   bool                  // indicates if the monitoring loop should be active
}

//------------------------------------------------------------------------------

// StartMonitor creates a process to monitor the consistency of the model.
func StartMonitor(m *model.Model, c  chan model.Event) (*Monitor) {
	// create the monitor
	monitor := Monitor{
		Model:   m,
    Channel: c,
    Active:  true,
	}

	// start the monitor
	go monitor.Run()

  // success
  return &monitor
}

//------------------------------------------------------------------------------

// Run starts the monitor loop validating the model and triggering compensating tasks
func (m *Monitor) Run() {
  channel := m.Channel

  // loop while monitor needs to be active
  for m.Active {
    mismatch := false

    // loop over all domains
    domainNames, _ := m.Model.ListDomains()
    for _, domainName := range domainNames {
      domain, _ := model.GetDomain(domainName)

      // loop over all solutions
      solutionNames, _ := domain.ListSolutions()
      for _, solutionName := range solutionNames {
        solution, _ := domain.GetSolution(solutionName)

        // loop over all elements
        elementNames, _ := solution.ListElements()
        for _, elementName := range elementNames {
          element, _ := solution.GetElement(elementName)

          // check if the element needs to be updated
      		if !element.OK() {

            // only update if no solution related task is currently being executed
            if !runningSolutionTasks(domain, solution) {
              // create task to update the element
        			task, _ := engine.NewSolutionTask(domainName, "", solution)

        			// trigger the task
        			channel <- model.NewEvent(domainName, task.UUID, model.EventTypeTaskExecution, "", "Solution: " + solutionName + "/Element: " + elementName)

              // mark that a mismatch has been found
              mismatch = true

              // exit from elements loop
              break
            }
      		}
        } // end of loop over all elements
      } // end of loop over all solutions
    } // end of loop over all domains

    // sleep a bit if no mismatch between current state and target state has been found
    if !mismatch {
      time.Sleep(100 * time.Millisecond)
    }
  } // end of while active loop
}

//------------------------------------------------------------------------------

// runningSolutionTasks checks if there are any currently running solution tasks
func runningSolutionTasks(domain *model.Domain, solution *model.Solution) (bool) {
  // go through task list and find matching tasks
  taskNames, _ := domain.ListTasks()
  for _, taskName := range taskNames {
    task, _ := domain.GetTask(taskName)

    if task.GetType() == "Solution" && task.GetSolution() == solution.Solution &&
       (task.GetStatus() == model.TaskStatusInitial || task.GetStatus() == model.TaskStatusExecuting) {
      return true
    }
  }

  return false
}

//------------------------------------------------------------------------------

// Stop will flag the monitor to stop execution
func (m *Monitor) Stop() {
  m.Active = false
}

//------------------------------------------------------------------------------
