package gRPC

import (
  "testing"
  "context"
  "bou.ke/monkey"
  "os"
)

//------------------------------------------------------------------------------

func TestGRPC01(t *testing.T) {
  defer func(){ if r := recover(); r != nil {} }()

  // monkey patch os.Exit(1) scencario
  fakeExit := func(int) {
    panic("os.Exit called")
  }
  patch := monkey.Patch(os.Exit, fakeExit)
  defer patch.Unpatch()

  // create a gRPC server
	NewController()

  // create controller
  controller := DefaultController{}

  // Test Check
  req1 := VoidMessage{}
  _, err := controller.Check(context.Background(), &req1)
  if err != nil {
    t.Errorf("<controller>.Check should not have returned an error")
  }

  // Test Create
  req2 := SetupMessage{
    Domain:                  "demo",
    Solution:                "app",
    Version:                 "V0.0.0",
    Element:                 "tenant",
    Cluster:                 "V1.0.0",
    Instance:                "first",
    Target:                  "active",
    State:                   "initial",
    DesignTimeConfiguration: "",
    RuntimeConfiguration:    "",
    Elements: map[string]*ElementSetupMessage{
      "tenant": &ElementSetupMessage{
        Element:                 "tenant",
        Component:               "Tenant",
        Target:                  "active",
        State:                   "initial",
        DesignTimeConfiguration: "",
        RuntimeConfiguration:    "",
        Endpoint:                "",
        Clusters: map[string]*ClusterSetupMessage{
          "V1.0.0": &ClusterSetupMessage{
            Cluster:                 "V1.0.0",
            Target:                  "active",
            State:                   "initial",
            Min:                     1,
            Max:                     1,
            Size:                    1,
            BaseConfiguration:       "",
            DesignTimeConfiguration: "",
            RuntimeConfiguration:    "",
            Endpoint:                "",
            Relationships:           map[string]*RelationshipSetupMessage{},
            Instances: map[string]*InstanceSetupMessage{
              "first": &InstanceSetupMessage{
                Instance:                "first",
                Target:                  "active",
                State:                   "initial",
                BaseConfiguration:       "",
                DesignTimeConfiguration: "",
                RuntimeConfiguration:    "",
                Endpoint:                "",
              },
            },
          },
      },
      },
    },
  }
  _, err = controller.Create(context.Background(), &req2)
  if err != nil {
    t.Errorf("<controller>.Create should not have returned an error")
  }

  _, err = controller.Configure(context.Background(), &req2)
  if err != nil {
    t.Errorf("<controller>.Configure should not have returned an error")
  }

  _, err = controller.Start(context.Background(), &req2)
  if err != nil {
    t.Errorf("<controller>.Start should not have returned an error")
  }

  _, err = controller.Stop(context.Background(), &req2)
  if err != nil {
    t.Errorf("<controller>.Stop should not have returned an error")
  }

  _, err = controller.Reconfigure(context.Background(), &req2)
  if err != nil {
    t.Errorf("<controller>.Reconfigure should not have returned an error")
  }

  _, err = controller.Destroy(context.Background(), &req2)
  if err != nil {
    t.Errorf("<controller>.Destroy should not have returned an error")
  }

  _, err = controller.Reset(context.Background(), &req2)
  if err != nil {
    t.Errorf("<controller>.Reset should not have returned an error")
  }

  // Test logs
  Log("error", "test", "log", "error")
  Log("warn",  "test", "log", "warn")
  Log("info",  "test", "log", "info")
  Log("debug", "test", "log", "debug")
  Log("fatal", "test", "log", "fatal")
  Log("panic", "test", "log", "panic")
}

//------------------------------------------------------------------------------
