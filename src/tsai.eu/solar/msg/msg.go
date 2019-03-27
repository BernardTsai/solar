package msg

import (
 "time"
  "sync"
  "context"
  "errors"

	"github.com/segmentio/kafka-go"

  "tsai.eu/solar/util"
)

//------------------------------------------------------------------------------

var msgAddress  string      // kafka bus
var msgTopic    string      // topic to which the message will be published
var msgInit     sync.Once   // do once sync for connection initialisation
var msgConn    *kafka.Conn  // message connection

//------------------------------------------------------------------------------

// Open establishes a connection to the message bus
func Open() (error) {
  // initialise once
  msgInit.Do(func() {
    configuration, _ := util.GetConfiguration()

    msgTopic   = configuration.MSG.Events
    msgAddress = configuration.MSG.Address

    conn, err := kafka.DialLeader(context.Background(), "tcp", msgAddress, msgTopic, 0)
    if err == nil {
      msgConn = conn
    } else {
      msgConn = nil
    }
  })

  // check the result of the initialisation
  if msgConn == nil {
    return errors.New("Unable to connect to the message bus: " + msgAddress + " / topic: " + msgTopic)
  }

  // success
  return nil
}

//------------------------------------------------------------------------------

// Publish writes data to the message bus
func Publish(key string, value string)  {
  // handle issues with the message bus
  defer func() {
    if r := recover(); r != nil {
      msgConn = nil
    }
  }()

  // only publish if a message connection was established
  if msgConn != nil {
    msgConn.SetWriteDeadline(time.Now().Add(500*time.Millisecond))
    msgConn.WriteMessages(
        kafka.Message{
          Key:   []byte(key),
          Value: []byte(value)},
    )
  }
}

//------------------------------------------------------------------------------

// Close shuts down the connection to the message bus
func Close() {
  if msgConn != nil {
    msgConn.Close()
    msgConn = nil
  }
}

//------------------------------------------------------------------------------
