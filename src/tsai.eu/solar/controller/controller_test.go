package controller

import (
  "testing"
  "io"
  "os"
)

//------------------------------------------------------------------------------

// copyFile simply copies an existing file to a not yet existing destination
func copyFile(src string, dest string) {
  srcFile, _ := os.Open(src)
  defer srcFile.Close()

  destFile, _ := os.Create(dest)
  defer destFile.Close()

  io.Copy(destFile, srcFile)
  destFile.Sync()
}

//------------------------------------------------------------------------------

func TestController01(t *testing.T) {
  // prepare configuration file
  srcConfig  := "testdata/solar-conf.yaml"
  destConfig := "solar-conf.yaml"

  copyFile(srcConfig, destConfig)

  // cleanup routine
  defer func() {os.Remove(destConfig)}()

  // GetController test
  _, err := GetController("dummy")
  if err != nil {
    t.Errorf("GetController is unable to find dummy controller")
  }
}

//------------------------------------------------------------------------------
