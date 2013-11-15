package main

import (
	"bufio"
	"code.google.com/p/go.net/websocket"
	"fmt"
	"github.com/scottferg/goat"
	"os"
)

var (
	g    *goat.Goat
	pool *connectionPool
)

type socketConnection struct {
	socket   *websocket.Conn
	messages chan string
}

type connectionPool struct {
	pool       map[*socketConnection]bool
	register   chan *socketConnection
	unregister chan *socketConnection
	broadcast  chan string
}

func NewConnectionPool() *connectionPool {
	return &connectionPool{
		pool:       make(map[*socketConnection]bool),
		register:   make(chan *socketConnection),
		unregister: make(chan *socketConnection),
		broadcast:  make(chan string),
	}
}

func (c *connectionPool) run() {
	for {
		select {
		case sc := <-c.register:
			c.pool[sc] = true
		case sc := <-c.unregister:
			delete(c.pool, sc)
			close(sc.messages)
		case message := <-c.broadcast:
			for sc := range c.pool {
				sc.write(message)
			}
		}
	}
}

func (s *socketConnection) listen() {
	for {
		var raw string
		err := websocket.Message.Receive(s.socket, &raw)
		if err != nil {
			break
		}
	}
}

func (s *socketConnection) write(message string) {
	err := websocket.Message.Send(s.socket, message)
	if err != nil {
		return
	}
}

func handleSocket(ws *websocket.Conn) {
	s := &socketConnection{
		socket: ws,
	}

	defer func() {
		if err := s.socket.Close(); err != nil {
			fmt.Println("Websocket could not be closed: ", err.Error())
		}
	}()

	pool.register <- s

	client := s.socket.Request().RemoteAddr
	fmt.Println("Client connected: ", client)

	s.listen()
}

func readInput() {
	bio := bufio.NewReader(os.Stdin)

	for {
		line, err := bio.ReadString('\n')
		if err != nil {
			fmt.Println(err.Error())
		}

		pool.broadcast <- line
	}
}

func main() {
	g = goat.NewGoat()

	g.RegisterRoute("/", "socket", goat.GET, websocket.Server{
		Handler: handleSocket,
	})

	pool = NewConnectionPool()
	go pool.run()
	go g.ListenAndServe("5050")

	readInput()
}
