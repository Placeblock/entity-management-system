package realtime

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Placeblock/nostalgicraft-ems/pkg/realtime"
	zmq "github.com/pebbe/zmq4"
)

type Publisher struct {
	Channel chan realtime.Action
	socket  *zmq.Socket
}

func NewPublisher(context *zmq.Context) *Publisher {
	channel := make(chan realtime.Action, 1)
	socket, err := context.NewSocket(zmq.PUB)
	if err != nil {
		log.Fatal("Could not create ZMQ Socket: ", err)
		return nil
	}
	err = socket.Bind("tcp://127.0.0.1:3008")
	if err != nil {
		log.Fatal("Could not Bind ZMQ Socket: ", err)
		return nil
	}
	fmt.Println("Created new ZMQ Publisher")
	return &Publisher{channel, socket}
}

func (publisher *Publisher) Listen() {
	for {
		data := <-publisher.Channel
		fmt.Printf("Sending Realtime Data: %+v\n", data)
		b, err := json.Marshal(data)
		if err != nil {
			fmt.Printf("Could not send message %s\n", b)
			continue
		}
		publisher.socket.SendBytes(b, zmq.DONTWAIT)
	}
}
