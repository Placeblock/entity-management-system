package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/codelix/ems/internal/realtime"
	"github.com/pebbe/zmq4"
)

func main() {
	zctx, err := zmq4.NewContext()
	socket, err := zctx.NewSocket(zmq4.SUB)
	if err != nil {
		log.Fatal("Could not create ZMQ Socket")
		return
	}
	err = socket.Connect("tcp://127.0.0.1:3008")
	if err != nil {
		log.Fatal("Could not connect to ZMQ Publisher: ", err)
		return
	}
	err = socket.SetSubscribe("")
	if err != nil {
		log.Fatal("Could not set Subscribe Filter: ", err)
		return
	}
	fmt.Println("Connected to ZMQ Socket!")
	for {
		bytes, err := socket.RecvBytes(0)
		fmt.Printf("%s\n", string(bytes[:]))
		if err != nil {
			fmt.Printf("Could not recv data: %v\n", err)
			continue
		}
		var data realtime.Action
		err = json.Unmarshal(bytes, &data)
		if err != nil {
			fmt.Printf("Could not unmarshal data: %v\n", err)
			continue
		}
		fmt.Printf("%+v\n", data)
	}
}
