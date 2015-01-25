package daihinmin

import (
	"encoding/json"
	"log"

	"code.google.com/p/go.net/websocket"
)

const sendQueueSize = 5

type sesh string

type client struct {
	session sesh
	socket  *websocket.Conn
	//	listening map[string]*channel

	user *user

	sendq chan interface{}
}

type seshPair struct {
	sesh
	*client
}

type cmd struct {
	Do   string
	With Cards  `json:",omitempty"`
	To   string `json:",omitempty"`
	What string `json:",omitempty"`
}

func NewClient(socket *websocket.Conn) *client {
	c := &client{
		session: sesh(generateID("c:")),
		socket:  socket,
		user: &user{ // temp
			ID:   generateID("u:"),
			Name: generateID("guest-")[:14],
		},

		sendq: make(chan interface{}, sendQueueSize),
	}
	return c
}

func (c *client) send(msg interface{}) {
	c.sendq <- msg
}

func (c *client) error(msg string) {
	c.send(ErrorReply{
		X:   "error",
		Wrt: "?",
		Msg: msg,
	})
}

func (c *client) Write() {
	for {
		select {
		case msg, ok := <-c.sendq:
			if !ok {
				// our work here is done
				return
			}
			err := websocket.JSON.Send(c.socket, msg)
			if err != nil {
				log.Printf("uh oh: %v", err)
				return
			}
		}
	}
}

func (c *client) Run() {
	for {
		var data []byte
		err := websocket.Message.Receive(c.socket, &data)
		if err != nil {
			break
		}

		var req cmd
		err = json.Unmarshal(data, &req)
		if err != nil {
			c.error("invalid command")
			continue
		}

		log.Printf("Got: %s\n-> %s\n", req.Do, string(data))

		switch req.Do {
		case "login":
			c.send(WelcomeReply{
				X:    "welcome",
				Name: c.user.Name,
				Sesh: c.session,
			})
		default:
			log.Printf("Unknown req %s\n", req.Do)
		}
	}
	// post-disconnect cleanup
	c.socket.Close()
	close(c.sendq)
}

func (c *client) username() string {
	if c == nil {
		return "(nil)"
	}
	if c.user != nil {
		return c.user.Name
	}
	return "Guest"
}
