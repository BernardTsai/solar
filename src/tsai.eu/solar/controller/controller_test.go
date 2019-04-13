package controller

import (
  "testing"
  "io"
  "os"
  "net"
  "time"

  "tsai.eu/solar/model"
  "tsai.eu/solar/controller/gRPC"
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

// startController creates a new gRPC controller
func startController() {
  // open TCP port 10000
  lis, _ := net.Listen("tcp", ":10000")

  // create a gRPC server
  grpcServer, _ := gRPC.NewController()

  // register controller and start listening
  grpcServer.Serve(lis)
}

//------------------------------------------------------------------------------

// TestController01 evaluates the gRPC functions of the controller package
func TestController01(t *testing.T) {
  // prepare configuration file
  srcConfig  := "testdata/solar-conf.yaml"
  destConfig := "solar-conf.yaml"

  copyFile(srcConfig, destConfig)

  // cleanup routine
  defer func() {os.Remove(destConfig)}()

  // start gRPC controller
  go startController()

  time.Sleep(time.Millisecond)

  // GetController test
  ctrl, err := GetController("dummy")
  if err != nil {
    t.Errorf("GetController is unable to find dummy controller")
  }

  // load model
  m := model.GetModel()

	m.Load("testdata/testdata1.yaml")

  s, _ := model.GetSetup("demo", "app", "V0.0.0", "oam", "V1.0.0", "a6c0bea1-ce1a-4fae-b943-1dbcc50cb311")

  ctrl.Create(s)
  ctrl.Configure(s)
  ctrl.Start(s)
  ctrl.Reconfigure(s)
  ctrl.Stop(s)
  ctrl.Destroy(s)
  ctrl.Reset(s)
  ctrl.Status(s)
}

//------------------------------------------------------------------------------
