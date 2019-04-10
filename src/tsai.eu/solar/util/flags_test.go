package util

import (
  "testing"
)

//------------------------------------------------------------------------------

// TestFlags01 tests the command line options related functions of the util package.
func TestFlags01(t *testing.T) {
  ParseCommandLineOptions()
  if Debug() == true {
    t.Errorf("Default debug flag should be false")
  }
}

//------------------------------------------------------------------------------
