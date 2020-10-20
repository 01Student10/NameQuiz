package chat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/micro/go-micro/v2/logger"
	"github.com/rs/cors"
	"github.com/songquiz/backend/api"
)

// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Changed by Moritz Laur Copyright 2020.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Copyright 2020 Moritz Laur.
type Client struct {
	// The websocket connection.
	Conn *websocket.Conn
	// Buffered channel of outbound messages.
	Send chan []byte
}

// The Chat Service struct.
type Service struct {
	socketClients map[int32]Hub
}

// Ctor of the Chat Service.
func New() *Service {
	return &Service{
		socketClients: make(map[int32]Hub),
	}
}

// A UserMessage is send back to all Websocket chat connections in the Frontend.
// It receives the Message Struct which to be send back.
// It always returns nil. Results are passed through the response message.
// The response message contains:
//	 Successful - tells whether the broadcast was successful.
//	 Error - if there was an Error it will be displayed in the Response string.
//	 Otherwise Error will not be transmitted.
func (cs *Service) BroadcastMessage(_ context.Context, req *api.BroadcastMessageRequest, rsp *api.BroadcastMessageResponse) error {
	rawMessage, err := json.Marshal(req.Message)
	if err != nil {
		rsp.Successful = 2 // Stands for false.
		logger.Error(err)
	}
	rawMessage = bytes.TrimSpace(bytes.Replace(rawMessage, newline, space, -1))
	hub, ok := cs.socketClients[req.Message.SessionId]
	if ok {
		hub.Broadcast <- rawMessage
		rsp.Successful = 1 // Stands for True.
	} else {
		rsp.Successful = 2 // Stands for false.
		rsp.Error = fmt.Sprintf("No Clients found for session Id: %v", req.Message.SessionId)
	}
	return nil
}

// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Changed by Moritz Laur Copyright 2020.
//
// creates a new Hub and starts it in a goroutine. It creates a new Websocket Connection and initializes the connection.
// It receives a socket Url, a socket port and a session id to which the connection belongs.
// It always returns nil. Results are passed through the response message.
// The response message contains:
//	Error - if there was an Error it will be displayed in the Response string.
//  Successful - tells whether starting the websocket was successful.
func (cs *Service) StartWebsocket(_ context.Context, req *api.StartWebsocketRequest, rsp *api.StartWebsocketResponse) error {
	rsp.Successful = 2 // Stands for false.
	mux := http.NewServeMux()
	hub := newHub()
	go hub.run()
	cs.socketClients[req.SessionId] = *hub
	mux.HandleFunc(req.SocketUrlPattern, func(w http.ResponseWriter, r *http.Request) {
		cs.serveWs(w, r, req.SessionId)
	})
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	handler := cors.Default().Handler(mux)
	// Decorate existing handler with cors functionality set in c.
	handler = c.Handler(handler)
	go startServer(req.SocketPort, handler)
	rsp.Successful = 1 // Stands for True.
	return nil
}

// Private goroutine which starting the Server for a websocket Connection.
// It receives a port and a http Handler.
func startServer(port string, handler http.Handler) {
	err := http.ListenAndServe(port, handler)
	if err != nil {
		logger.Fatal("ListenAndServe: ", err)
	}
}

// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Changed by Moritz Laur Copyright 2020.
//
// Private function which handles websocket requests from the peer.
// It registers client to the hub. It starts the writeToWebsocket goroutine for every Client.
func (cs *Service) serveWs(w http.ResponseWriter, r *http.Request, sessionId int32) *Client {
	r.Header.Add("Connection", "Upgrade")
	r.Header.Add("Upgrade", "websocket")
	r.Header.Add("Sec-WebSocket-Version", "13")
	r.Header.Add("Sec-WebSocket-Key", "1234")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade went wrong")
		log.Println(err)
		return nil
	}
	client := &Client{Conn: conn, Send: make(chan []byte, 256)}
	hub, ok := cs.socketClients[sessionId]
	if ok {
		hub.Register <- client
	}

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go cs.writeToWebsocket(*client)
	return client
}

// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Changed by Moritz Laur Copyright 2020.
//
// writeToWebsocket write messages from the hub to the websocket connection.
// A goroutine running writeToWebsocket is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
// It receives a Client Struct which contains the needed Websocket Connection and a chanel with data to write.
func (cs *Service) writeToWebsocket(client Client) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = client.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-client.Send:
			_ = client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				_ = client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, _ = w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			_ = client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
