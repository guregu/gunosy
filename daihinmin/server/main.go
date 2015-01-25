package main

import (
	"log"
	"net/http"

	"code.google.com/p/go.net/websocket"
	dh "github.com/gophergala/gunosy/daihinmin"
)

func main() {
	bind := "localhost:3000"

	// Create default game
	m := dh.NewMatch("New Game 1")
	log.Printf("Created the initial game: %s\n", m)
	http.Handle("/ws", websocket.Handler(handle))

	log.Println("running:", bind)
	http.ListenAndServe(bind, nil)
}

func handle(ws *websocket.Conn) {
	c := dh.NewClient(ws)
	go c.Write()
	c.Run()
}
