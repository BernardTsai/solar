package util

import (
  "testing"
  "os"
  "bou.ke/monkey"
)

//------------------------------------------------------------------------------

// TestLog01 tests the logging of panic circumstances.
func TestLog01(t *testing.T) {
  defer func(){ if r := recover(); r != nil {} }()

  StartLogging()
  LogLevel("panic")
  LogPanic("test", "UTIL", "TestLog01-A")

  t.Errorf("Logging should have panicked")
}

//------------------------------------------------------------------------------

// TestLog02 tests the logging of fatal circumstances.
func TestLog02(t *testing.T) {
  defer func(){ if r := recover(); r != nil {} }()

  // monkey patch os.Exit(1) scencario
  fakeExit := func(int) {
    panic("os.Exit called")
  }
  patch := monkey.Patch(os.Exit, fakeExit)
  defer patch.Unpatch()

  StartLogging()
  LogLevel("fatal")
  LogFatal("test", "UTIL", "TestLog01-B")

  t.Errorf("Logging should have exited due to a fatal circumstance")
}

//------------------------------------------------------------------------------

// TestLog03 tests all other circumstances.
func TestLog03(t *testing.T) {
  StartLogging()
  LogLevel("error")
  LogError("test", "UTIL", "TestLog01-C")
  LogLevel("warn")
  LogWarn("test", "UTIL", "TestLog01-D")
  LogLevel("info")
  LogInfo("test", "UTIL", "TestLog01-E")
  LogLevel("debug")
  LogDebug("test", "UTIL", "TestLog01-F")
}

//------------------------------------------------------------------------------
