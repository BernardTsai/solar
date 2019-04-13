package msg

import (
  "testing"
  "context"
  "io"
  "os"

  "github.com/segmentio/kafka-go"
  "time"

  "bou.ke/monkey"
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

// TestMsg01 tests a missing configuration for the msg package.
// func TestMsg01(t *testing.T) {
//   _, err := Start(context.Background())
//   if err != nil {
//     t.Errorf("msg.Start reported error:\n%s", err)
//   }
// }

//------------------------------------------------------------------------------

// TestMsg02 tests the exposed routines of the msg package.
func TestMsg02(t *testing.T) {
  // prepare configuration file
  srcConfig  := "testdata/solar-conf.yaml"
  destConfig := "solar-conf.yaml"

  copyFile(srcConfig, destConfig)

  // cleanup routine
  defer func() {os.Remove(destConfig)}()

  // monkey patch os.Exit(1) scencario
  msgNr := 0
  var msgs [3][2]string
  msgs[0] = [2]string{"Element",  "demo/app/tenant/active"}
  msgs[1] = [2]string{"Instance", "demo/app/tenant/V1.0.0/uuid/active"}
  msgs[2] = [2]string{"Cluster",  "demo/app/tenant/V1.0.0/active"}
  fakeReadMessage := func(reader *kafka.Reader, ctx context.Context) (kafka.Message, error){
    if msgNr == 3 {
      msgNr = 0
    }
    msgNr = msgNr + 1
    return kafka.Message{
      Key:   []byte(msgs[msgNr-1][0]),
      Value: []byte(msgs[msgNr-1][1]),
    }, nil
  }
  // patch := monkey.Patch(msg.MonitoringReader.ReadMessage, fakeReadMessage)
  patch := monkey.Patch((*kafka.Reader).ReadMessage, fakeReadMessage)
  defer patch.Unpatch()

  // start message interface
  messaging, err := Start(context.Background())
  if err != nil {
    t.Errorf("msg.Start reported error:\n%s", err)
  }

  // notify
  Notify("test", "Hello World")

  // status
  messaging.Status()

  // wait a bit for the messages to be processed
  time.Sleep(100*time.Microsecond)

  // stop
  messaging.Stop()
}

//------------------------------------------------------------------------------
