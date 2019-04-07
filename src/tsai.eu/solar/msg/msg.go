package msg

import (
  "fmt"
  "sync"
  "context"
  "errors"
  "strings"
  "time"

	"github.com/segmentio/kafka-go"

  "tsai.eu/solar/model"
  "tsai.eu/solar/util"
)

//------------------------------------------------------------------------------

// MSG represents the messaging interface to a kafka bus.
type MSG struct {
  Context              context.Context  // runtime context
  Identifier           string           // identifier of this node as message source
  Address              string           // address of the kafka bus
  NotificationTopic    string           // topic to which the notifications will be published
  NotificationWriter  *kafka.Writer     // notification connection
  MonitoringTopic      string           // topic from which the message will be published
  MonitoringReader    *kafka.Reader     // monitoring connection
}

//------------------------------------------------------------------------------

var msg     *MSG          // messaging interface singleton
var msgInit  sync.Once    // do once sync for connection initialisation

//------------------------------------------------------------------------------

// Start initialises the messaging interface.
func Start(ctx context.Context) (*MSG, error){
  // initialise singleton once
	msgInit.Do(func() {
    // determine configuration
    configuration, err := util.GetConfiguration()
    if err != nil {
      util.LogError("main", "MSG", "failed to read the configuration: " + err.Error())
      return
    }

    // construct the messaging interface
    msg = &MSG{
      Context:            ctx,
      Identifier:         configuration.CORE.Identifier,
      Address:            configuration.MSG.Address,
      NotificationTopic:  configuration.MSG.Notifications,
      NotificationWriter: nil,
      MonitoringTopic:    configuration.MSG.Monitoring,
      MonitoringReader:   nil,
    }

    // start the messaging interfaces
    err = msg.StartReader()
    if err != nil {
      util.LogError("main", "MSG", "failed to start the reader: " + err.Error())
    }

    err = msg.StartWriter()
    if err != nil {
      util.LogError("main", "MSG", "failed to start the writer: " + err.Error())
    }
	})

  // check if initialisation has been successful
  if msg == nil {
    util.LogError("main", "MSG", "failed to initialise the messaging interface")
  }

  // return
  return msg, nil
}

//------------------------------------------------------------------------------

// Status of the messaging interface.
func (msg *MSG) Status() (messageStatus bool, writerStatus bool, readerStatus bool){
  messageStatus = msg != nil
  writerStatus  = msg.StatusWriter()
  readerStatus  = msg.StatusReader()

  return messageStatus, readerStatus, writerStatus
}

//------------------------------------------------------------------------------

// Stop deactivates the messaging interface.
func (msg *MSG) Stop() {
  msg.StopReader()
  msg.StopWriter()
}

//------------------------------------------------------------------------------

// StartWriter reconnects the writer to the message bus
func (msg *MSG) StartWriter() error {
  // create the notification writer if necessary
  if msg.NotificationWriter == nil {
    msg.NotificationWriter = kafka.NewWriter(kafka.WriterConfig{
      Brokers:  []string{msg.Address},
      Topic:    msg.NotificationTopic,
      Balancer: &kafka.LeastBytes{},
    })
  }

  // check the result of the initialisation
  if msg.NotificationWriter == nil {
    util.LogError("main", "MSG", "failed to start the writer")
    return errors.New("failed to start the writer")
  }

  // success
  util.LogInfo("main", "MSG", "writer active")
  return nil
}

//------------------------------------------------------------------------------

// StopWriter shuts down the writer
func (msg *MSG) StopWriter() {
  if msg.NotificationWriter != nil {
    msg.NotificationWriter.Close()
    msg.NotificationWriter = nil
    util.LogInfo("main", "MSG", "writer inactive")
  }
}

//------------------------------------------------------------------------------

// StatusReader provides the status of the writer
func (msg *MSG) StatusWriter() bool{
  messageStatus := msg != nil
  writerStatus  := false

  if messageStatus {
    writerStatus = msg.NotificationWriter != nil
  }

  return writerStatus
}

//------------------------------------------------------------------------------
//------------------------------------------------------------------------------

// StartReader reconnects the reader to the message bus
func (msg *MSG) StartReader() error {
  // create the notification writer if necessary
  if msg.MonitoringReader == nil {
    msg.MonitoringReader = kafka.NewReader(kafka.ReaderConfig{
      Brokers:   []string{msg.Address},
      GroupID:   msg.Identifier,
      Topic:     msg.MonitoringTopic,
      MinBytes:  10e3, // 10KB
      MaxBytes:  10e6, // 10MB
    })

    // start listening
  	go msg.listen()
  }

  // check the result of the initialisation
  if msg.MonitoringReader == nil {
    util.LogError("main", "MSG", "failed to start the reader")
    return errors.New("failed to start the reader")
  }

  // success
  util.LogInfo("main", "MSG", "reader active")
  return nil
}


//------------------------------------------------------------------------------

// StopReader shuts down the reader
func (msg *MSG) StopReader() {
  if msg.MonitoringReader != nil {
    msg.MonitoringReader.Close()
    msg.MonitoringReader = nil
    util.LogInfo("main", "MSG", "reader inactive")
  }
}

//------------------------------------------------------------------------------

// StatusReader provides the status of the reader
func (msg *MSG) StatusReader() bool{
  messageStatus := msg != nil
  readerStatus  := false

  if messageStatus {
    readerStatus = msg.MonitoringReader != nil
  }

  return readerStatus
}

//------------------------------------------------------------------------------

// listen starts listening for status updates
func (msg *MSG) listen() {
  // handle issues with the message bus
  defer func() {
    if r := recover(); r != nil {
      fmt.Printf("Aborting listener due to error\n")
      if msg != nil {

        msg.MonitoringReader = nil
      }
    }
  }()

  // listen while reader is available
  for msg.MonitoringReader != nil {
    message, err := msg.MonitoringReader.ReadMessage(msg.Context)
    if err != nil {
      util.LogInfo("main", "MSG", err.Error())
      break
    }

    // depending on the message key identify the correct entity to update
    entity := string(message.Key)
    value  := string(message.Value)

    switch entity {
      case "Element":
        names := strings.Split(value, "/")

        if len(names) == 4 {
          element, err := model.GetElement(names[0], names[1], names[2])
          if err == nil {
            element.SetState( names[3] )
          }
        }
      case "Cluster":
        names := strings.Split(value, "/")

        if len(names) == 5 {
          cluster, _ := model.GetCluster(names[0], names[1], names[2], names[3])
          if err == nil {
            cluster.SetState( names[4] )
          }
        }
      case "Instance":
        names := strings.Split(value, "/")

        if len(names) == 6 {
          instance, _ := model.GetInstance(names[0], names[1], names[2], names[3], names[4])
          if err == nil {
            instance.SetState( names[5] )
          }
        }
    }
  }
}

//------------------------------------------------------------------------------

// Notify writes data to the message bus
func Notify(key string, value string)  {
  // only publish if a message connection was established
  if msg != nil && msg.NotificationWriter != nil {
    go notify(key, value)
  }
}

// notify writes data to the message bus
func notify(key string, value string)  {
  // handle issues with the message bus
  defer func() {
    if r := recover(); r != nil {
      util.LogError("main", "MSG", "unable to send notification: " + r.(string))
      if msg != nil {
        msg.NotificationWriter = nil
      }
    }
  }()

  // derive a context with timeout
  ctx, cancel := context.WithTimeout(msg.Context, 1 * time.Second)
  defer cancel()

  // only publish if a message connection was established
  msg.NotificationWriter.WriteMessages(
      ctx,
      kafka.Message{
        Key:   []byte(key),
        Value: []byte(value)},
    )
}

//------------------------------------------------------------------------------
