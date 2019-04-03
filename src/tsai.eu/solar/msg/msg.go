package msg

import (
  "fmt"
  "sync"
  "context"
  "errors"
  "strings"

	"github.com/segmentio/kafka-go"

  "tsai.eu/solar/model"
  "tsai.eu/solar/util"
)

//------------------------------------------------------------------------------

// MSG represents the messaging interface to a kafka bus.
type MSG struct {
  Context              context.Context  // runtime context
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

// NewMSG creates a new messaging interface.
func NewMSG(ctx context.Context) (*MSG, error){
  // determine configuration
  configuration, _ := util.GetConfiguration()

  identifier        := configuration.CORE.Identifier
  msgbusAddress     := configuration.MSG.Address
  notificationTopic := configuration.MSG.Notifications
  monitoringTopic   := configuration.MSG.Monitoring

  // construct the messaging interface
  msg := MSG{
    Context:            ctx,
    Address:            msgbusAddress,
    NotificationTopic:  notificationTopic,
    NotificationWriter: nil,
    MonitoringTopic:    monitoringTopic,
    MonitoringReader:   nil,
  }

  // create the notication writer
  msg.NotificationWriter = kafka.NewWriter(kafka.WriterConfig{
    Brokers:  []string{msgbusAddress},
    Topic:    notificationTopic,
    Balancer: &kafka.LeastBytes{},
  })

  // creating the monitoring reader
  msg.MonitoringReader = kafka.NewReader(kafka.ReaderConfig{
    Brokers:   []string{msgbusAddress},
    GroupID:   identifier,
    Topic:     monitoringTopic,
    MinBytes:  10e3, // 10KB
    MaxBytes:  10e6, // 10MB
  })

  // check the result of the initialisation
  if msg.NotificationWriter == nil {
    return nil, errors.New("Unable to connect to the message bus: " + msgbusAddress + " / topic: " + notificationTopic + "\n")
  }

  // check the result of the initialisation
  if msg.MonitoringReader == nil {
    return nil, errors.New("Unable to connect to the message bus: " + msgbusAddress + " / topic: " + monitoringTopic + "\n")
  }

	// start the messaging interface
	go msg.Run()

	return &msg, nil
}

//------------------------------------------------------------------------------

// Run starts the listener of the messaging interface
func (m *MSG) Run() {
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
  for m.MonitoringReader != nil {
    message, err := m.MonitoringReader.ReadMessage(m.Context)
    if err != nil {
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

// Stop shuts down the messaging interface
func (m *MSG) Stop() {
  if m.NotificationWriter != nil {
    m.NotificationWriter.Close()
    m.NotificationWriter = nil
  }

  if m.MonitoringReader != nil {
    m.MonitoringReader.Close()
    m.MonitoringReader = nil
  }
}

//------------------------------------------------------------------------------

// StartMSG creates and return the core new messaging interface.
func StartMSG(ctx context.Context) (*MSG, error){
  // initialise singleton once
	msgInit.Do(func() {
    msgBus, err := NewMSG(ctx)

    if err != nil {
      msg = nil
    } else {
      msg = msgBus
    }
	})

  // the attempt to connect to the message bus has failed before
  if msg == nil {
    return nil, errors.New("Unable to connect to the message bus")
  }

  // return
  return msg, nil
}

//------------------------------------------------------------------------------

// Notify writes data to the message bus
func Notify(key string, value string)  {
  // only publish if a message connection was established
  if msg != nil && msg.NotificationWriter != nil {
    notify(key, value)
  }
}

// notify writes data to the message bus
func notify(key string, value string)  {
  // handle issues with the message bus
  defer func() {
    if r := recover(); r != nil {
      if msg != nil {
        msg.NotificationWriter = nil
      }
    }
  }()

  // only publish if a message connection was established
  go msg.NotificationWriter.WriteMessages(
      msg.Context,
      kafka.Message{
        Key:   []byte(key),
        Value: []byte(value)},
    )
}

//------------------------------------------------------------------------------
