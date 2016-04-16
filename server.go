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
}

func NewClient(id int, ws *websocket.Conn) *Client {
  c := Client{id, ws}

  // Write new client
  c.Socket.WriteMessage(websocket.BinaryMessage, bindata)

  println("Sent!")

  go c.Receiver()

  return &c
}

func (c *Client) Receiver() {
  for {
    _, p, err := c.Socket.ReadMessage()
    if err != nil {
      fmt.Println("breaking..")
      println(err.Error())
      break
    }
    
    println(p)
  }

  c.Socket.Close()
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

  /*
  var err error
  var message []byte

  // Ensure we close the socket at the end
  defer func() {
    if err = ws.Close(); err != nil {
      fmt.Println("Websocket could not be closed", err.Error())
    }
  }()

  Counter++
  client := NewClient(Counter, ws, ws.Request().RemoteAddr)

  // Push client to the list
  Clients[client.id] = client

  // Receiving loop
  for {
    if err = websocket.Message.Receive(ws, &message); err != nil {
      fmt.Println("Can't read the received message: ", err.Error())

      delete(Clients, client.id)
    }

    println(client.id)
  }*/
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
    fmt.Fprintf(w, "<a href=\"http://github.com/iiegor/slither\">Powered by iiegor/slither</a>")
  })

  err := http.ListenAndServe(":8080", nil)
  if err != nil {
    panic("Can't start the server: " + err.Error())
  }
}