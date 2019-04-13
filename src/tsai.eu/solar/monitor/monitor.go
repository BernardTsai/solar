package monitor

import (
  "context"
  "time"

  "tsai.eu/solar/model"
  "tsai.eu/solar/engine"
  "tsai.eu/solar/util"
)

//------------------------------------------------------------------------------

// Monitor validates solutions and triggers tasks to converge to the desired target state.
type Monitor struct {
  Channel  chan model.Event      // the channel for event notification
  Ticker  *time.Ticker           // ticker
  Active   bool                  // indicates if the monitoring loop should be active
}

//------------------------------------------------------------------------------

// Start creates a process to monitor the consistency of the model.
func Start(ctx context.Context) (*Monitor) {
	// create the monitor
	monitor := Monitor{
    Channel: engine.GetEventChannel(),
    Ticker:  time.NewTicker(100 * time.Millisecond),
    Active:  false,
	}

	// start the monitor
	go monitor.Run(ctx)
  monitor.Start()

  // success
  return &monitor
}

//------------------------------------------------------------------------------

// Run starts the monitor loop validating the model and triggering compensating tasks
func (m *Monitor) Run(ctx context.Context) {
  // loop while monitor needs to be active
  for {
    select {
    // check if context has expired
    case <-ctx.Done():
      util.LogInfo("main", "MON", "monitoring initial")
      m.Ticker.Stop()
      return
    // wait for next tick and monitor solutions
    case <- m.Ticker.C:
      if m.Active {
        checkSolutions()
      }
    }
  }
}

//------------------------------------------------------------------------------

// checkDomains checks if there are any inconsistent domains
func checkSolutions() {
  channel := engine.GetEventChannel()

  // loop over all domains
  domainNames, _ := model.GetDomains()
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

            util.LogInfo(task.UUID, "MON", "starting task to reconcile element: '" + element.Element + "' in solution: '" + solution.Solution + "'")

            // trigger the task
            channel <- model.NewEvent(domainName, task.UUID, model.EventTypeTaskExecution, "", "Solution: " + solutionName + "/Element: " + elementName)

            // exit from elements loop
            break
          }
        }
      } // end of loop over all elements
    } // end of loop over all solutions
  } // end of loop over all domains
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

// Start will flag the monitor to resume execution
func (m *Monitor) Start() {
  m.Active = true
  util.LogInfo("main", "MON", "monitoring active")
}

//------------------------------------------------------------------------------

// Stop will flag the monitor to pause execution
func (m *Monitor) Stop() {
  m.Active = false
  util.LogInfo("main", "MON", "monitoring inactive")
}

//------------------------------------------------------------------------------
