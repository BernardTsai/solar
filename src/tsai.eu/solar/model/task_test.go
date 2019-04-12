package model

import (
	"testing"
	"os"
)

//------------------------------------------------------------------------------

// TestTask01 tests the basic functions of the task package.
func TestTask01(t *testing.T) {
	filename := "test.yaml"

	// cleanup routine
  defer func() {os.Remove(filename)}()

	model := GetModel()

	model.Load("testdata/testdata1.yaml")

	task, _ := GetTask("demo", "0f7bec76-3252-41b0-a7c9-82b75abba9de")

	task.Save(filename)

	task.Load(filename)

	yaml, _ := task.Show()

	task.Load2(yaml)
}

//------------------------------------------------------------------------------

// TestTask02 tests the task related functions of the task package.
func TestTask02(t *testing.T) {
	model := GetModel()

	model.Load("testdata/testdata1.yaml")

	task, _ := GetTask("demo", "0f7bec76-3252-41b0-a7c9-82b75abba9de")

	if task.GetType() != "Cluster" {
		t.Errorf("<task>.GetType returned the wrong value")
	}

	if task.GetDomain() != "demo" {
		t.Errorf("<task>.GetDomain returned the wrong value")
	}

	if task.GetSolution() != "app" {
		t.Errorf("<task>.GetSolution returned the wrong value")
	}

	if task.GetVersion() != "V0.0.0" {
		t.Errorf("<task>.GetVersion returned the wrong value")
	}

	if task.GetElement() != "app-server" {
		t.Errorf("<task>.GetElement returned the wrong value")
	}

	if task.GetCluster() != "V1.0.0" {
		t.Errorf("<task>.GetCluster returned the wrong value")
	}

	if task.GetInstance() != "" {
		t.Errorf("<task>.GetInstance returned the wrong value")
	}

	if task.GetState() != "" {
		t.Errorf("<task>.GetState returned the wrong value")
	}

	if task.GetUUID() != "0f7bec76-3252-41b0-a7c9-82b75abba9de" {
		t.Errorf("<task>.GetUUID returned the wrong value")
	}

	if task.GetParent() != "ffab9264-8f18-4649-ac2d-3f5f8d17c726" {
		t.Errorf("<task>.GetParent returned the wrong value")
	}

	if task.GetStatus() != "completed" {
		t.Errorf("<task>.GetStatus returned the wrong value")
	}

	if task.GetPhase() != 0 {
		t.Errorf("<task>.GetPhase returned the wrong value")
	}

	task.GetTimestamps()

	task.GetSubtasks()

	NewTaskInfo(task, 0)
	NewTaskInfo(task, 1)
	NewTaskInfo(task, 2)

	_, err := task.GetSubtask("unknown")
	if err == nil {
		t.Errorf("<task>.GetSubTask should not have found subtask")
	}

	subtask, err := task.GetSubtask("b51e2ac5-5bf0-44c3-96cc-41de7dbdc7ff")
	if err != nil {
		t.Errorf("<task>.GetSubTask should have found subtask")
	}

	task.AddSubtask(subtask)
}

//------------------------------------------------------------------------------

// TestTask03 tests the handler related functions of the task package.
func TestTask03(t *testing.T) {
	model := GetModel()

	model.Load("testdata/testdata1.yaml")

	task, _ := GetTask("demo", "0f7bec76-3252-41b0-a7c9-82b75abba9de")

	task.SetExecute(nil)
	task.SetCompleted(nil)
	task.SetFailed(nil)
	task.SetTimeout(nil)
	task.SetTerminate(nil)
}

//------------------------------------------------------------------------------
