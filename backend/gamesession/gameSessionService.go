package gamesession

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/gorilla/websocket"
	"github.com/micro/go-micro/v2/logger"
	"github.com/rs/cors"
	"github.com/songquiz/backend/api"
)

type Status string

// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Changed by Moritz Laur Copyright 2020.
const (
	Ready    Status = "Ready to play"
	InRound  Status = "In Round"
	Finished Status = "Game Over"
	Error    Status = "Error"

	timePerRound    = 30
	PointsPerAnswer = 20
	NoErrStr        = ""

	SocketaddressPrefix = "ws://localhost"

	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
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
var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// A Struct containing all Information which will be transferred to the Frontend.
type FrontendMessage struct {
	Status          string
	UserNameToScore map[string]int32
	SessionID       int32
	EntryID         string
	Path            string
	PicInfo         string
	Winner          []Pair
	ErrorMsg        string
	GameTitle       string
}

// A Session Struct containing all Information needed in a Session.
type Session struct {
	sessionID      int32
	currentStatus  Status
	userToScore    map[string]int32
	hasScored      map[string]struct{}
	currentGameID  int32
	entryListIndex int32
	hub            Hub
	isRunning      int32
	chatAddress    string
	sessionAddress string
}

// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Copyright 2020 Moritz Laur.
type Client struct {
	// The websocket connection.
	Conn *websocket.Conn
	// Buffered channel of outbound messages.
	Send chan []byte
}

// The Game Session Service struct.
type Service struct {
	// Clients
	userService        api.UserService
	chatService        api.ChatService
	gameService        api.GameService
	quizContentService api.QuizContentService
	// Variables
	sessionIDIndex     int32
	roundTime          int32
	sessions           map[int32]*Session
	chatSocketPorts    map[string]bool
	sessionSocketPorts map[string]bool
}

// Ctor of the Game Session Service.
func New(gameService api.GameService,
	chatService api.ChatService,
	userService api.UserService,
	quizContentService api.QuizContentService) *Service {
	return &Service{
		userService:        userService,
		chatService:        chatService,
		gameService:        gameService,
		quizContentService: quizContentService,
		sessions:           make(map[int32]*Session),
		roundTime:          timePerRound,
		sessionIDIndex:     0,
		// This map tells if a port for the chat Websocket Connection is in use or not.
		chatSocketPorts: map[string]bool{
			":8001": false,
			":8002": false, ":8003": false, ":8004": false,
			":8005": false, ":8006": false, ":8007": false,
			":8008": false, ":8009": false, ":8010": false},
		// This map tells if a port for the session Websocket Connection is in use or not.
		sessionSocketPorts: map[string]bool{
			":8101": false,
			":8102": false, ":8103": false, ":8104": false,
			":8105": false, ":8106": false, ":8107": false,
			":8108": false, ":8109": false, ":8110": false},
	}
}

// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Changed by Moritz Laur Copyright 2020.
//
// Creates a new Session and initializes a Chat Websocket Connection as well as a Session Websocket Connection.
// It receives a Game ID. This is needed to tell of which game the session is to be created.
// It always returns nil. Results are passed through the response message.
// The response message contains:
//	SessionId - The id of the created Session.
//  ChatSocketAddress -the Address of the Chat Websocket the Frontend can connect to.
//  SessionSocketAddress - the Address of the Session Websocket the Frontend can connect to.
//	Response - if there was an Error it will be displayed in the Response string.
//	Otherwise Response will not be transmitted.
func (gs *Service) CreateSession(_ context.Context, req *api.CreateSessionRequest, rsp *api.CreateSessionResponse) error {
	var session Session
	sessionPort, hasSessionPort := selectPort(gs.sessionSocketPorts)
	chatPort, hasChatPort := selectPort(gs.chatSocketPorts)

	if hasSessionPort && hasChatPort {
		gs.sessionIDIndex++
		res, _ := gs.chatService.StartWebsocket(context.Background(), &api.StartWebsocketRequest{
			SessionId:        gs.sessionIDIndex,
			SocketUrlPattern: fmt.Sprintf("/chat/%v", gs.sessionIDIndex),
			SocketPort:       chatPort,
		})

		if res.Successful == 1 {
			hub := newHub()
			go hub.run()
			session.hub = *hub
			errorMsg := gs.initializeSessionWebsocket(sessionPort)
			if errorMsg != "" {
				rsp.Response = fmt.Sprintf("Couldn't create Session cause : %v", errorMsg)
			}
		} else {
			rsp.Response = fmt.Sprintf("Couldn't create Session cause : %v", res.Error)
		}
	} else {
		rsp.Response = "The maximal amount of Game Sessions is reached. Sorry"
	}

	session.chatAddress = fmt.Sprintf("%v%v/chat/%v", SocketaddressPrefix, chatPort, gs.sessionIDIndex)
	session.sessionAddress = fmt.Sprintf("%v%v/session/%v", SocketaddressPrefix, sessionPort, gs.sessionIDIndex)

	session.sessionID = gs.sessionIDIndex
	session.currentGameID = req.GameId
	session.userToScore = make(map[string]int32)
	session.currentStatus = Ready
	session.entryListIndex = -1
	session.isRunning = 2 // stands for false
	gs.sessions[gs.sessionIDIndex] = &session
	rsp.SessionId = session.sessionID
	rsp.ChatSocketAddress = session.chatAddress
	rsp.SessionSocketAddress = session.sessionAddress

	return nil
}

// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Changed by Moritz Laur Copyright 2020.
//
// Private function which handles websocket requests from the peer.
// It registers client to the hub. It starts the writeToSocket goroutine for every Client.
func (gs *Service) initializeSessionWebsocket(sessionPort string) string {
	errorMsg := ""
	mux := http.NewServeMux()
	var sessionId = gs.sessionIDIndex
	mux.HandleFunc(fmt.Sprintf("/session/%v", sessionId), func(w http.ResponseWriter, r *http.Request) {
		r.Header.Add("Connection", "Upgrade")
		r.Header.Add("Upgrade", "websocket")
		r.Header.Add("Sec-WebSocket-Version", "13")
		r.Header.Add("Sec-WebSocket-Key", "1234")

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			msg := fmt.Sprintf("Error while creating socket connection: %v", err)
			logger.Infof(msg)
			errorMsg = msg
		}
		client := &Client{Conn: conn, Send: make(chan []byte, 256)}
		session, ok := gs.sessions[sessionId]
		if ok {
			session.hub.Register <- client
			go gs.writeToSocket(*client)
		} else {
			errorMsg = fmt.Sprintf("Cant find session with id: %v", sessionId)
		}
	})
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	handler := cors.Default().Handler(mux)
	// Decorate existing handler with cors functionality set in c
	handler = c.Handler(handler)
	go startServer(sessionPort, handler)
	return errorMsg
}

// Private function to return a unused port.
// It receives The map from which to choose from.
// Returns a free Port and whether it found one.
func selectPort(mapToSelectFrom map[string]bool) (port string, foundOne bool) {
	foundOne = false
	for key, value := range mapToSelectFrom {
		if !value {
			port = key
			foundOne = true
			break
		}
	}
	if _, ok := mapToSelectFrom[port]; ok {
		mapToSelectFrom[port] = true
	}
	return port, foundOne
}

// Starts a Session.
// It receives Session Id which should be started.
// It always returns nil. Results are passed through the response message.
// The response message contains:
//	HasWorked - tells whether a session has been started.
//	Response - if there was an Error it will be displayed in the Response string.
//	If there was no Error the Response will state that the session has been started.
func (gs *Service) StartSession(_ context.Context, req *api.StartSessionRequest, rsp *api.StartSessionResponse) error {
	session, ok := gs.sessions[req.SessionId]
	if ok {
		session.isRunning = 1 // stands for true
	} else {
		rsp.Response = fmt.Sprintf("Error: no Session with id %v found.", req.SessionId)
	}

	go gs.runGame(req.SessionId)
	rsp.HasWorked = 1 // stands for true
	rsp.Response = "started the session.."

	return nil
}

// Adds a User to a Session.
// It receives User Id of the User to be added and a Session Id which tells to which session the user should be added.
// After the User joined the session the Frontend will be updated.
// It always returns nil. Results are passed through the response message.
// The response message contains:
//	HasWorked - tells whether a user has joined the session.
//	Response - if there was an Error it will be displayed in the Response string.
//	If there was no Error the Response will state that the User joined the session.
func (gs *Service) JoinSession(_ context.Context, req *api.JoinSessionRequest, rsp *api.JoinSessionResponse) error {

	res, _ := gs.userService.GetUser(context.Background(), &api.GetUserRequest{
		Id: req.UserId,
	})
	session, ok := gs.sessions[req.SessionId]
	if ok {
		if res.User != nil {
			game, entryIds, errorStr := gs.getDataFromServices(req.SessionId)
			if len(errorStr) == 0 {
				_, alreadyExists := session.userToScore[res.User.Username]
				if !alreadyExists {
					session.userToScore[res.User.Username] = 0
					rsp.Response = fmt.Sprintf("User %v joined the Session ", res.User.Username)
					rsp.HasWorked = 1 // stands for true
					if session.currentStatus == InRound {
						gs.updateFrontend(session, game, entryIds, "")
					}
				} else {
					rsp.Response = fmt.Sprintf("User %v is already in the Session ", res.User.Username)
					rsp.HasWorked = 2 // stands for false
				}
			} else {
				rsp.Response = errorStr
				rsp.HasWorked = 2 // stands for false
			}
		} else {
			rsp.Response = fmt.Sprintf("No User with id  %v found. ", req.UserId)
		}
	} else {
		rsp.Response = fmt.Sprintf("No Session with id  %v found. ", req.SessionId)
		rsp.HasWorked = 2 // stands for false
	}

	return nil
}

// Removes a User to a Session.
// It receives User Id of the User to be removed and a Session Id which tells from which session the user should be removed.
// After the User left the session the Frontend will be updated.
// It always returns nil. Results are passed through the response message.
// The response message contains:
//	HasWorked - tells whether the user has left the session.
//	Response - if there was an Error it will be displayed in the Response string.
//	If there was no Error the Response will state that the User left the session.
func (gs *Service) LeaveSession(_ context.Context, req *api.LeaveSessionRequest, rsp *api.LeaveSessionResponse) error {
	rsp.HasWorked = 2
	session, ok := gs.sessions[req.SessionId]
	if ok {
		res, _ := gs.userService.GetUser(context.Background(), &api.GetUserRequest{
			Id: req.UserId,
		})
		if res.User != nil {
			game, entryIds, errorStr := gs.getDataFromServices(req.SessionId)
			if len(errorStr) == 0 {
				delete(session.userToScore, res.User.Username)
				delete(session.hasScored, res.User.Username)
				if session.currentStatus == InRound {
					gs.updateFrontend(session, game, entryIds, NoErrStr)
				}
				rsp.Response = fmt.Sprintf("User %v left the Session", res.User.Username)
				rsp.HasWorked = 1
			} else {
				rsp.Response = errorStr
			}
		} else {
			rsp.Response = fmt.Sprintf("There is no User with the following Id: %v ", req.UserId)
		}
	} else {
		rsp.Response = fmt.Sprintf("There is no Session with the following Id: %v ", req.SessionId)
	}

	return nil
}

// This function is called whenever a Message is published to the topic "nq.UserMessage".
// It calls the Chatservice to Broadcast the Message to the Frontend.
// It checks if the user guess was right and updates the Frontend accordingly.
// It receives the Message which will be broadcast.
func (gs *Service) HandleMessage(_ context.Context, message *api.UserMessage) {
	errorMsg := NoErrStr
	session, ok := gs.sessions[message.SessionId]
	if ok {
		resUser, _ := gs.userService.GetUser(context.Background(), &api.GetUserRequest{
			Id: message.UserId,
		})
		if len(resUser.Response) == 0 {
			message.UserName = resUser.User.Username
			resChat, _ := gs.chatService.BroadcastMessage(context.Background(), &api.BroadcastMessageRequest{
				Message: message,
			})
			if resChat.Successful != 1 {
				logger.Infof("couldn't broadcast message to frontend with error: %v", resChat.Error)
				errorMsg = resChat.Error
			}
			if message.WasRight == 1 {

				game, entryIds, errorStr := gs.getDataFromServices(message.SessionId)
				if len(errorStr) == 0 {
					_, alreadyScored := session.hasScored[resUser.User.Username]
					if entryIds[session.entryListIndex] == message.EntryId && !alreadyScored {
						currPoints, ok := session.userToScore[resUser.User.Username]
						if ok {
							var empty struct{}
							session.hasScored[resUser.User.Username] = empty
							session.userToScore[resUser.User.Username] = currPoints + PointsPerAnswer

							gs.updateFrontend(session, game, entryIds, errorMsg)
						}
					}
				} else {
					logger.Info("Error in HandleMessage() occurred with Msg: %v ", errorStr)
					gs.updateFrontend(session, nil, nil, errorStr)
				}
			}

		} else {
			gs.updateFrontend(session, nil, nil, resUser.Response)
		}
	}
}

// Checks if a session has already been created.
// It receives Game Id to identify the Session to be checked on.
// It always returns nil. Results are passed through the response message.
// The response message contains:
//	IsRunning - tells whether the session is running.
//	SessionId - the id of the session.
//  ChatSocketAddress -the Address of the Chat Websocket the Frontend can connect to.
//  SessionSocketAddress - the Address of the Session Websocket the Frontend can connect to.
func (gs *Service) IsSessionCreated(_ context.Context, req *api.IsSessionCreatedRequest, rsp *api.IsSessionCreatedResponse) error {
	rsp.IsRunning = 2 // Stands for false.
	for _, value := range gs.sessions {
		if value.currentGameID == req.GameId {
			rsp.IsRunning = 1
			rsp.SessionId = value.sessionID
			rsp.ChatSocketAddress = value.chatAddress
			rsp.SessionSocketAddress = value.sessionAddress
			break
		}
	}
	return nil
}

// This function prepares the frontendMessage Struct and sends the updated Data to the Frontend.
// It receives the session to of the game what will be updated.
// It receives a game struct to extract data from it.
// It receives an string array of entryIds.
// It receives an errorMsg string. If the string is empty it just updates the Frontend.
// If the errorMsg string is not empty it updates the Frontend by calling an Error State.
func (gs *Service) updateFrontend(session *Session, game *api.Game, entryIds []string, errorMsg string) {
	var status, path, picInfo, resultEntryId string
	winner := make([]Pair, 0)
	if len(errorMsg) == 0 {
		resultEntryId = entryIds[session.entryListIndex]
		if session.entryListIndex < game.AmountOfRounds && session.entryListIndex < int32(len(entryIds)) {
			entry, _ := gs.quizContentService.GetContentEntry(context.Background(), &api.GetContentEntryRequest{
				ListId:  game.PlaylistId,
				EntryId: entryIds[session.entryListIndex],
			})
			if len(entry.Response) == 0 {
				session.currentStatus = InRound
				status = fmt.Sprintf("%v : %v", session.currentStatus, session.entryListIndex+1)
				path = entry.Entry.Path
				picInfo = entry.Entry.Licence
			} else {
				status = fmt.Sprint(Error)
				errorMsg = entry.Response
			}
		} else {
			session.currentStatus = Finished
			winner = gs.getTheWinner(*session)
			status = fmt.Sprint(session.currentStatus)
		}

	} else {
		status = fmt.Sprint(Error)
		entryIds = []string{"default"}
		game = &api.Game{Title: "default"}
		resultEntryId = ""
	}
	feMsg := FrontendMessage{
		SessionID:       session.sessionID,
		Status:          status,
		Path:            path,
		PicInfo:         picInfo,
		UserNameToScore: session.userToScore,
		EntryID:         resultEntryId,
		Winner:          winner,
		ErrorMsg:        errorMsg,
		GameTitle:       game.Title,
	}
	gs.initializeBroadcast(session, feMsg)

}

// Function to call the running hub to which will update the Frontend through the Websocket Connection.
// It receives the session struct in which the WebSocket Connection struct is stored.
// It receives a message Struct with the data to be updated.
func (*Service) initializeBroadcast(session *Session, message FrontendMessage) {

	rawData, err := json.Marshal(message)
	if err != nil {
		logger.Fatal("Error while parsing struct to []byte with message: %v", err.Error())
	}
	rawData = bytes.TrimSpace(bytes.Replace(rawData, newline, space, -1))
	session.hub.Broadcast <- rawData
}

// A goroutine which manages the course of the game. It runs as long as there are rounds left to play or until
//an Error rises.
// After each round it updates the Frontend and then sleeps for a duration of time.
// If an error arises the Session and therefore the Websocket Connection will be closed.
func (gs *Service) runGame(sessionId int32) {
	for {
		session, ok := gs.sessions[sessionId]
		if ok {
			game, entryIds, errorStr := gs.getDataFromServices(sessionId)
			if len(errorStr) == 0 {
				if session.entryListIndex < game.AmountOfRounds && session.entryListIndex < int32(len(entryIds)) {
					session.entryListIndex++
					session.hasScored = make(map[string]struct{})
					gs.updateFrontend(session, game, entryIds, NoErrStr)
				} else {
					logger.Info("about to close the session....!!!!!")
					gs.closeSession(*session)
					return
				}

			} else {
				logger.Infof("Error in runGame() occurred with Msg: %v ", errorStr)
				gs.updateFrontend(session, nil, nil, errorStr)
				time.Sleep(time.Second * 2)
				gs.closeSession(*session)
				return
			}
		} else {
			return
		}

		time.Sleep(time.Second * timePerRound)
	}
}

// Private goroutine which starting the Server for a websocket Connection.
// It receives a port and a http Handler.
func startServer(port string, handler http.Handler) {
	err := http.ListenAndServe(port, handler)
	if err != nil {
		logger.Fatal("ListenAndServe: ", err)
	}
}

// A Function which calls the Game Service, the Quizcontent Service to collect data from them.
// It receives a session id.
// It returns a game struct and an array of entryIds.
// If there was an error calling the other services the error string will contain that error.
// Otherwise the error string will be empty.
func (gs *Service) getDataFromServices(sessionId int32) (game *api.Game, entryIds []string, error string) {
	session, ok := gs.sessions[sessionId]
	if ok {
		res, _ := gs.gameService.GetGame(context.Background(), &api.GetGameRequest{
			Id: session.currentGameID,
		})
		if res.Game == nil {
			error = fmt.Sprintf("Game with id %v is not available", session.currentGameID)
		} else {
			game = res.Game
			quizContentRes, _ := gs.quizContentService.GetContentList(context.Background(), &api.GetContentListRequest{
				ListId: res.Game.PlaylistId,
			})
			if quizContentRes.EntryIds == nil {
				error = fmt.Sprintf("There is no Content with Playlist id %v", res.Game.PlaylistId)
			} else {
				entryIds = quizContentRes.EntryIds
				error = ""
			}
		}
	} else {
		error = fmt.Sprintf("Session with id %v is not available", sessionId)
	}
	return game, entryIds, error
}

// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Changed by Moritz Laur Copyright 2020.
//
// writeToSocket write messages from the hub to the websocket connection.
// A goroutine running writeToSocket is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
// It receives a Client Struct which contains the needed Websocket Connection and a chanel with data to write.
func (gs *Service) writeToSocket(client Client) {
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
			logger.Info("writing message..")
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

type Pair struct {
	Name  string
	Score int32
}

// This function calculates the winner(one or more) of a completed game.
// It receives the session of the game.
// It returns an array of pairs containing the names of the winner and there scores.
func (gs *Service) getTheWinner(session Session) (result []Pair) {
	tempRes := make([]Pair, 0)
	result = make([]Pair, 0)
	for key, value := range session.userToScore {
		tempRes = append(tempRes, Pair{key, value})
	}
	sort.Slice(tempRes, func(first, second int) bool {
		return tempRes[first].Score > tempRes[second].Score
	})
	for index, element := range tempRes { //nolint:wsl
		if index == 0 || result[0].Score == element.Score {
			result = append(result, element)
		}
	}
	return result
}

// This function closes the Websocket Connection and deletes the sessionId from the sessions map.
// It receives the session to be closed and deleted.
func (gs *Service) closeSession(session Session) {
	//TODO: delete Chat connection
	for client := range session.hub.Clients {
		_ = client.Conn.Close()
	}
	delete(gs.sessions, session.sessionID)
}
