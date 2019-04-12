package model

import (
	"testing"
)

//------------------------------------------------------------------------------

// TestTrace01 tests the basic functions of the trace package.
func TestTrace01(t *testing.T) {
	model := GetModel()

	model.Load("testdata/testdata1.yaml")

	task, _ := GetTask("demo", "0f7bec76-3252-41b0-a7c9-82b75abba9de")

	NewTrace(task)

	task, _ = GetTask("demo", "212927e4-cc49-4784-aa35-66430a6bd43b")

	NewTrace(task)
}

//------------------------------------------------------------------------------
