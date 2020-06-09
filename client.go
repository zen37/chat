package main

import (
	"log"

	"github.com/gorilla/websocket"
)

// client represents a single chatting user
type client struct {
	//holds a reference to the websocket that allows communication with the client
	socket *websocket.Conn
	//buffered channel through which the received messages are queued ready to be delivered to the user's browser via the socket
	send chan []byte
	//keeps a reference to the chatting room the user is in, so messages from the user are sent to everyone in the room
	room *room
}

// client reads from the socket, continually sending the received messages
// to the forward channel
//in case of an error the loop will break and the socket will close
func (c *client) read() {
	log.Println("I am in read")
	defer c.socket.Close()
	for {
		_, msg, err := c.socket.ReadMessage()
		log.Println("read msg: ", string(msg))
		if err != nil {
			log.Println("Error : something terrible happen -> ", err)
			return
		}
		c.room.forward <- msg
	}
}

//client continually accepts messages and writes them out of the socket
//in case of an error the loop will break and the socket will close
func (c *client) write() {
	log.Println("I am in write")
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		log.Println("write msg: ", string(msg))
		if err != nil {
			log.Println("Error : something terrible happen -> ", err)
			return
		}
	}
}
