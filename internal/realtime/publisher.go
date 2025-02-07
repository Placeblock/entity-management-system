package realtime

import (
	"encoding/json"
	"fmt"
	"log"

	zmq "github.com/pebbe/zmq4"
)

type Publisher struct {
	Channel chan interface{}
	socket  *zmq.Socket
}

func NewPublisher(context *zmq.Context) *Publisher {
	channel := make(chan interface{}, 1)
	socket, err := context.NewSocket(zmq.PUB)
	if err != nil {
		log.Fatal("Could not create ZMQ Socket")
		return nil
	}
	socket.Bind("localhost:3008")
	fmt.Println("Created new ZMQ Publisher")
	return &Publisher{channel, socket}
}

func (publisher *Publisher) Listen() {
	for {
		data := <-publisher.Channel
		b, err := json.Marshal(data)
		if err != nil {
			fmt.Printf("Could not send message %s\n", b)
			continue
		}
		publisher.socket.SendBytes(b, zmq.DONTWAIT)
	}
}
