package main

import (
  "runtime"
  "fmt"
  "net/http"

  "github.com/gorilla/websocket"

  . "slither/types"
)

var (
  Counter = 0
  Clients = make(map[int] Client) // clients list
  Origins = map[string]bool {
    "http://localhost:8000": true, // development
    "http://slither.io": true, // production
  }
)

var bindata = []byte{0, 0, 97, 0, 84, 96, 1, 155, 1, 44, 0, 144, 48, 2, 27, 0, 40, 5, 120, 0, 33, 0, 28, 1, 174, 6} // example data

type Client struct {
  Id int
  Socket *websocket.Conn
  Write chan []byte
  Broadcast chan []byte
}

func NewClient(id int, ws *websocket.Conn) *Client {
  c := Client{
    id,
    ws,
    make(chan []byte),
    make(chan []byte),
  }

  go c.Receiver()
  go c.Writer()

  return &c
}

func (c *Client) Receiver() {
  // Send snake before the loop
  c.Write <- bindata

  for {
    _, p, err := c.Socket.ReadMessage()
    if err != nil {
      println(err.Error())
      break
    }

    byteLength := len(p)
    println("New packet received with a length of", byteLength)
  }

  c.Socket.Close()
}

func (c *Client) Writer() {
  for {
    select {
      case message := <-c.Write:
        c.Socket.WriteMessage(websocket.BinaryMessage, message)

      case message := <-c.Broadcast:
        for client := range Clients {
          Clients[client].Socket.WriteMessage(websocket.BinaryMessage, message)
        }
    }
  }
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
  if !Origins[r.Header.Get("Origin")] {
    http.Error(w, "Origin not allowed", 403)
    return
  }

  // Handshake
  ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
  if _, ok := err.(websocket.HandshakeError); ok {
    http.Error(w, "Not a websocket handshake", 400)
    return
  } else if err != nil {
    return
  }

  // Create a new client
  Counter++
  client := NewClient(Counter, ws)

  fmt.Printf("New user with id: %v\n", client.Id)
}

func main() {
  // Use the maximum available CPU/Cores.
  // GOMAXPROCS is unnecessary when the handlers do not do enough work to justify
  // the time lost communicating between processes.
  runtime.GOMAXPROCS(runtime.NumCPU())

  // Print header
  fmt.Println("*****************************************")
  fmt.Println("**             SLITHER-GO              **")
  fmt.Println("**                                     **")
  fmt.Printf("** Author: Iegor Azuaga   Version: %v **", VERSION)
  fmt.Println()
  fmt.Println("*****************************************")

  // Handle path requests
  http.HandleFunc("/slither", wsHandler)
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "<a href=\"https://github.com/iiegor/slither-go\">Powered by iiegor/slither-go</a>")
  })

  err := http.ListenAndServe(":443", nil)
  if err != nil {
    panic("Can't start the server: " + err.Error())
  }
}