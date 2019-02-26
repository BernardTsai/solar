package model

//------------------------------------------------------------------------------

// UndefinedState indicates a component state is undefined
const UndefinedState string = "undefined"

// FailureState indicates a component related failure has occured
const FailureState string = "failure"

// InitialState indicates a component is in the initial state
const InitialState string = "initial"

// InactiveState indicates a component is in the inactive state
const InactiveState string = "inactive"

// ActiveState indicates a component is in the active state
const ActiveState string = "active"

// CreatingState indicates a component is in the creating state
const CreatingState string = "creating"

// DestroyingState indicates a component is in the destroying state
const DestroyingState string = "destroying"

// StartingState indicates a component is in the starting state
const StartingState string = "starting"

// StoppingState indicates a component is in the stopping state
const StoppingState string = "stopping"

// ConfiguringState indicates a component is in the configuring state
const ConfiguringState string = "configuring"

// ResettingState indicates a component is in the resetting state
const ResettingState string = "resetting"

//------------------------------------------------------------------------------

// ContextRelationship indicates that the relationship refers to a runtime context dependency.
const ContextRelationship string = "context"

// ServiceRelationship indicates that the relationship refers to a client/service dependency.
const ServiceRelationship string = "service"

//------------------------------------------------------------------------------

// TaskStatusInitial resembles the initial state of a task
const TaskStatusInitial string = "initial"

// TaskStatusExecuting resembles the execution state of a task
const TaskStatusExecuting string = "executing"

// TaskStatusCompleted resembles the completed state of a task
const TaskStatusCompleted string = "completed"

// TaskStatusFailed resembles the failed state of a task
const TaskStatusFailed string = "failed"
// TaskStatusTimeout resembles the timeout state of a task
const TaskStatusTimeout string = "timeout"
// TaskStatusTerminated resembles the terminated state of a task
const TaskStatusTerminated string = "terminated"

//------------------------------------------------------------------------------
