// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: api/api.proto

package api

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "github.com/micro/go-micro/v2/api"
	client "github.com/micro/go-micro/v2/client"
	server "github.com/micro/go-micro/v2/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for IdService service

func NewIdServiceEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for IdService service

type IdService interface {
	GetId(ctx context.Context, in *GetIdRequest, opts ...client.CallOption) (*GetIdResponse, error)
}

type idService struct {
	c    client.Client
	name string
}

func NewIdService(name string, c client.Client) IdService {
	return &idService{
		c:    c,
		name: name,
	}
}

func (c *idService) GetId(ctx context.Context, in *GetIdRequest, opts ...client.CallOption) (*GetIdResponse, error) {
	req := c.c.NewRequest(c.name, "IdService.GetId", in)
	out := new(GetIdResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for IdService service

type IdServiceHandler interface {
	GetId(context.Context, *GetIdRequest, *GetIdResponse) error
}

func RegisterIdServiceHandler(s server.Server, hdlr IdServiceHandler, opts ...server.HandlerOption) error {
	type idService interface {
		GetId(ctx context.Context, in *GetIdRequest, out *GetIdResponse) error
	}
	type IdService struct {
		idService
	}
	h := &idServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&IdService{h}, opts...))
}

type idServiceHandler struct {
	IdServiceHandler
}

func (h *idServiceHandler) GetId(ctx context.Context, in *GetIdRequest, out *GetIdResponse) error {
	return h.IdServiceHandler.GetId(ctx, in, out)
}

// Api Endpoints for UserService service

func NewUserServiceEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for UserService service

type UserService interface {
	CreateUser(ctx context.Context, in *CreateUserRequest, opts ...client.CallOption) (*CreateUserResponse, error)
	GetUser(ctx context.Context, in *GetUserRequest, opts ...client.CallOption) (*GetUserResponse, error)
	GetAllUsers(ctx context.Context, in *GetAllUsersRequest, opts ...client.CallOption) (*GetAllUsersResponse, error)
	HasUser(ctx context.Context, in *HasUserRequest, opts ...client.CallOption) (*HasUserResponse, error)
	Login(ctx context.Context, in *LoginRequest, opts ...client.CallOption) (*LoginResponse, error)
}

type userService struct {
	c    client.Client
	name string
}

func NewUserService(name string, c client.Client) UserService {
	return &userService{
		c:    c,
		name: name,
	}
}

func (c *userService) CreateUser(ctx context.Context, in *CreateUserRequest, opts ...client.CallOption) (*CreateUserResponse, error) {
	req := c.c.NewRequest(c.name, "UserService.CreateUser", in)
	out := new(CreateUserResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userService) GetUser(ctx context.Context, in *GetUserRequest, opts ...client.CallOption) (*GetUserResponse, error) {
	req := c.c.NewRequest(c.name, "UserService.GetUser", in)
	out := new(GetUserResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userService) GetAllUsers(ctx context.Context, in *GetAllUsersRequest, opts ...client.CallOption) (*GetAllUsersResponse, error) {
	req := c.c.NewRequest(c.name, "UserService.GetAllUsers", in)
	out := new(GetAllUsersResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userService) HasUser(ctx context.Context, in *HasUserRequest, opts ...client.CallOption) (*HasUserResponse, error) {
	req := c.c.NewRequest(c.name, "UserService.HasUser", in)
	out := new(HasUserResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userService) Login(ctx context.Context, in *LoginRequest, opts ...client.CallOption) (*LoginResponse, error) {
	req := c.c.NewRequest(c.name, "UserService.Login", in)
	out := new(LoginResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for UserService service

type UserServiceHandler interface {
	CreateUser(context.Context, *CreateUserRequest, *CreateUserResponse) error
	GetUser(context.Context, *GetUserRequest, *GetUserResponse) error
	GetAllUsers(context.Context, *GetAllUsersRequest, *GetAllUsersResponse) error
	HasUser(context.Context, *HasUserRequest, *HasUserResponse) error
	Login(context.Context, *LoginRequest, *LoginResponse) error
}

func RegisterUserServiceHandler(s server.Server, hdlr UserServiceHandler, opts ...server.HandlerOption) error {
	type userService interface {
		CreateUser(ctx context.Context, in *CreateUserRequest, out *CreateUserResponse) error
		GetUser(ctx context.Context, in *GetUserRequest, out *GetUserResponse) error
		GetAllUsers(ctx context.Context, in *GetAllUsersRequest, out *GetAllUsersResponse) error
		HasUser(ctx context.Context, in *HasUserRequest, out *HasUserResponse) error
		Login(ctx context.Context, in *LoginRequest, out *LoginResponse) error
	}
	type UserService struct {
		userService
	}
	h := &userServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&UserService{h}, opts...))
}

type userServiceHandler struct {
	UserServiceHandler
}

func (h *userServiceHandler) CreateUser(ctx context.Context, in *CreateUserRequest, out *CreateUserResponse) error {
	return h.UserServiceHandler.CreateUser(ctx, in, out)
}

func (h *userServiceHandler) GetUser(ctx context.Context, in *GetUserRequest, out *GetUserResponse) error {
	return h.UserServiceHandler.GetUser(ctx, in, out)
}

func (h *userServiceHandler) GetAllUsers(ctx context.Context, in *GetAllUsersRequest, out *GetAllUsersResponse) error {
	return h.UserServiceHandler.GetAllUsers(ctx, in, out)
}

func (h *userServiceHandler) HasUser(ctx context.Context, in *HasUserRequest, out *HasUserResponse) error {
	return h.UserServiceHandler.HasUser(ctx, in, out)
}

func (h *userServiceHandler) Login(ctx context.Context, in *LoginRequest, out *LoginResponse) error {
	return h.UserServiceHandler.Login(ctx, in, out)
}

// Api Endpoints for GameService service

func NewGameServiceEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for GameService service

type GameService interface {
	CreateGame(ctx context.Context, in *CreateGameRequest, opts ...client.CallOption) (*CreateGameResponse, error)
	DeleteGame(ctx context.Context, in *DeleteGameRequest, opts ...client.CallOption) (*DeleteGameResponse, error)
	GetGame(ctx context.Context, in *GetGameRequest, opts ...client.CallOption) (*GetGameResponse, error)
	GetAllGames(ctx context.Context, in *GetAllGamesRequest, opts ...client.CallOption) (*GetAllGamesResponse, error)
}

type gameService struct {
	c    client.Client
	name string
}

func NewGameService(name string, c client.Client) GameService {
	return &gameService{
		c:    c,
		name: name,
	}
}

func (c *gameService) CreateGame(ctx context.Context, in *CreateGameRequest, opts ...client.CallOption) (*CreateGameResponse, error) {
	req := c.c.NewRequest(c.name, "GameService.CreateGame", in)
	out := new(CreateGameResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameService) DeleteGame(ctx context.Context, in *DeleteGameRequest, opts ...client.CallOption) (*DeleteGameResponse, error) {
	req := c.c.NewRequest(c.name, "GameService.DeleteGame", in)
	out := new(DeleteGameResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameService) GetGame(ctx context.Context, in *GetGameRequest, opts ...client.CallOption) (*GetGameResponse, error) {
	req := c.c.NewRequest(c.name, "GameService.GetGame", in)
	out := new(GetGameResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameService) GetAllGames(ctx context.Context, in *GetAllGamesRequest, opts ...client.CallOption) (*GetAllGamesResponse, error) {
	req := c.c.NewRequest(c.name, "GameService.GetAllGames", in)
	out := new(GetAllGamesResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for GameService service

type GameServiceHandler interface {
	CreateGame(context.Context, *CreateGameRequest, *CreateGameResponse) error
	DeleteGame(context.Context, *DeleteGameRequest, *DeleteGameResponse) error
	GetGame(context.Context, *GetGameRequest, *GetGameResponse) error
	GetAllGames(context.Context, *GetAllGamesRequest, *GetAllGamesResponse) error
}

func RegisterGameServiceHandler(s server.Server, hdlr GameServiceHandler, opts ...server.HandlerOption) error {
	type gameService interface {
		CreateGame(ctx context.Context, in *CreateGameRequest, out *CreateGameResponse) error
		DeleteGame(ctx context.Context, in *DeleteGameRequest, out *DeleteGameResponse) error
		GetGame(ctx context.Context, in *GetGameRequest, out *GetGameResponse) error
		GetAllGames(ctx context.Context, in *GetAllGamesRequest, out *GetAllGamesResponse) error
	}
	type GameService struct {
		gameService
	}
	h := &gameServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&GameService{h}, opts...))
}

type gameServiceHandler struct {
	GameServiceHandler
}

func (h *gameServiceHandler) CreateGame(ctx context.Context, in *CreateGameRequest, out *CreateGameResponse) error {
	return h.GameServiceHandler.CreateGame(ctx, in, out)
}

func (h *gameServiceHandler) DeleteGame(ctx context.Context, in *DeleteGameRequest, out *DeleteGameResponse) error {
	return h.GameServiceHandler.DeleteGame(ctx, in, out)
}

func (h *gameServiceHandler) GetGame(ctx context.Context, in *GetGameRequest, out *GetGameResponse) error {
	return h.GameServiceHandler.GetGame(ctx, in, out)
}

func (h *gameServiceHandler) GetAllGames(ctx context.Context, in *GetAllGamesRequest, out *GetAllGamesResponse) error {
	return h.GameServiceHandler.GetAllGames(ctx, in, out)
}

// Api Endpoints for ChatService service

func NewChatServiceEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for ChatService service

type ChatService interface {
	BroadcastMessage(ctx context.Context, in *BroadcastMessageRequest, opts ...client.CallOption) (*BroadcastMessageResponse, error)
	StartWebsocket(ctx context.Context, in *StartWebsocketRequest, opts ...client.CallOption) (*StartWebsocketResponse, error)
}

type chatService struct {
	c    client.Client
	name string
}

func NewChatService(name string, c client.Client) ChatService {
	return &chatService{
		c:    c,
		name: name,
	}
}

func (c *chatService) BroadcastMessage(ctx context.Context, in *BroadcastMessageRequest, opts ...client.CallOption) (*BroadcastMessageResponse, error) {
	req := c.c.NewRequest(c.name, "ChatService.BroadcastMessage", in)
	out := new(BroadcastMessageResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatService) StartWebsocket(ctx context.Context, in *StartWebsocketRequest, opts ...client.CallOption) (*StartWebsocketResponse, error) {
	req := c.c.NewRequest(c.name, "ChatService.StartWebsocket", in)
	out := new(StartWebsocketResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for ChatService service

type ChatServiceHandler interface {
	BroadcastMessage(context.Context, *BroadcastMessageRequest, *BroadcastMessageResponse) error
	StartWebsocket(context.Context, *StartWebsocketRequest, *StartWebsocketResponse) error
}

func RegisterChatServiceHandler(s server.Server, hdlr ChatServiceHandler, opts ...server.HandlerOption) error {
	type chatService interface {
		BroadcastMessage(ctx context.Context, in *BroadcastMessageRequest, out *BroadcastMessageResponse) error
		StartWebsocket(ctx context.Context, in *StartWebsocketRequest, out *StartWebsocketResponse) error
	}
	type ChatService struct {
		chatService
	}
	h := &chatServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&ChatService{h}, opts...))
}

type chatServiceHandler struct {
	ChatServiceHandler
}

func (h *chatServiceHandler) BroadcastMessage(ctx context.Context, in *BroadcastMessageRequest, out *BroadcastMessageResponse) error {
	return h.ChatServiceHandler.BroadcastMessage(ctx, in, out)
}

func (h *chatServiceHandler) StartWebsocket(ctx context.Context, in *StartWebsocketRequest, out *StartWebsocketResponse) error {
	return h.ChatServiceHandler.StartWebsocket(ctx, in, out)
}

// Api Endpoints for QuizContentService service

func NewQuizContentServiceEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for QuizContentService service

type QuizContentService interface {
	GetContentEntry(ctx context.Context, in *GetContentEntryRequest, opts ...client.CallOption) (*GetContentEntryResponse, error)
	GetContentList(ctx context.Context, in *GetContentListRequest, opts ...client.CallOption) (*GetContentListResponse, error)
	GetAllContentLists(ctx context.Context, in *GetAllContentListsRequest, opts ...client.CallOption) (*GetAllContentListsResponse, error)
	HasMatch(ctx context.Context, in *HasMatchRequest, opts ...client.CallOption) (*HasMatchResponse, error)
	CreateNewDataSet(ctx context.Context, in *CreateNewDataSetRequest, opts ...client.CallOption) (*CreateNewDataSetResponse, error)
}

type quizContentService struct {
	c    client.Client
	name string
}

func NewQuizContentService(name string, c client.Client) QuizContentService {
	return &quizContentService{
		c:    c,
		name: name,
	}
}

func (c *quizContentService) GetContentEntry(ctx context.Context, in *GetContentEntryRequest, opts ...client.CallOption) (*GetContentEntryResponse, error) {
	req := c.c.NewRequest(c.name, "QuizContentService.GetContentEntry", in)
	out := new(GetContentEntryResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *quizContentService) GetContentList(ctx context.Context, in *GetContentListRequest, opts ...client.CallOption) (*GetContentListResponse, error) {
	req := c.c.NewRequest(c.name, "QuizContentService.GetContentList", in)
	out := new(GetContentListResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *quizContentService) GetAllContentLists(ctx context.Context, in *GetAllContentListsRequest, opts ...client.CallOption) (*GetAllContentListsResponse, error) {
	req := c.c.NewRequest(c.name, "QuizContentService.GetAllContentLists", in)
	out := new(GetAllContentListsResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *quizContentService) HasMatch(ctx context.Context, in *HasMatchRequest, opts ...client.CallOption) (*HasMatchResponse, error) {
	req := c.c.NewRequest(c.name, "QuizContentService.HasMatch", in)
	out := new(HasMatchResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *quizContentService) CreateNewDataSet(ctx context.Context, in *CreateNewDataSetRequest, opts ...client.CallOption) (*CreateNewDataSetResponse, error) {
	req := c.c.NewRequest(c.name, "QuizContentService.CreateNewDataSet", in)
	out := new(CreateNewDataSetResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for QuizContentService service

type QuizContentServiceHandler interface {
	GetContentEntry(context.Context, *GetContentEntryRequest, *GetContentEntryResponse) error
	GetContentList(context.Context, *GetContentListRequest, *GetContentListResponse) error
	GetAllContentLists(context.Context, *GetAllContentListsRequest, *GetAllContentListsResponse) error
	HasMatch(context.Context, *HasMatchRequest, *HasMatchResponse) error
	CreateNewDataSet(context.Context, *CreateNewDataSetRequest, *CreateNewDataSetResponse) error
}

func RegisterQuizContentServiceHandler(s server.Server, hdlr QuizContentServiceHandler, opts ...server.HandlerOption) error {
	type quizContentService interface {
		GetContentEntry(ctx context.Context, in *GetContentEntryRequest, out *GetContentEntryResponse) error
		GetContentList(ctx context.Context, in *GetContentListRequest, out *GetContentListResponse) error
		GetAllContentLists(ctx context.Context, in *GetAllContentListsRequest, out *GetAllContentListsResponse) error
		HasMatch(ctx context.Context, in *HasMatchRequest, out *HasMatchResponse) error
		CreateNewDataSet(ctx context.Context, in *CreateNewDataSetRequest, out *CreateNewDataSetResponse) error
	}
	type QuizContentService struct {
		quizContentService
	}
	h := &quizContentServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&QuizContentService{h}, opts...))
}

type quizContentServiceHandler struct {
	QuizContentServiceHandler
}

func (h *quizContentServiceHandler) GetContentEntry(ctx context.Context, in *GetContentEntryRequest, out *GetContentEntryResponse) error {
	return h.QuizContentServiceHandler.GetContentEntry(ctx, in, out)
}

func (h *quizContentServiceHandler) GetContentList(ctx context.Context, in *GetContentListRequest, out *GetContentListResponse) error {
	return h.QuizContentServiceHandler.GetContentList(ctx, in, out)
}

func (h *quizContentServiceHandler) GetAllContentLists(ctx context.Context, in *GetAllContentListsRequest, out *GetAllContentListsResponse) error {
	return h.QuizContentServiceHandler.GetAllContentLists(ctx, in, out)
}

func (h *quizContentServiceHandler) HasMatch(ctx context.Context, in *HasMatchRequest, out *HasMatchResponse) error {
	return h.QuizContentServiceHandler.HasMatch(ctx, in, out)
}

func (h *quizContentServiceHandler) CreateNewDataSet(ctx context.Context, in *CreateNewDataSetRequest, out *CreateNewDataSetResponse) error {
	return h.QuizContentServiceHandler.CreateNewDataSet(ctx, in, out)
}

// Api Endpoints for GameSessionService service

func NewGameSessionServiceEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for GameSessionService service

type GameSessionService interface {
	CreateSession(ctx context.Context, in *CreateSessionRequest, opts ...client.CallOption) (*CreateSessionResponse, error)
	StartSession(ctx context.Context, in *StartSessionRequest, opts ...client.CallOption) (*StartSessionResponse, error)
	JoinSession(ctx context.Context, in *JoinSessionRequest, opts ...client.CallOption) (*JoinSessionResponse, error)
	LeaveSession(ctx context.Context, in *LeaveSessionRequest, opts ...client.CallOption) (*LeaveSessionResponse, error)
	IsSessionCreated(ctx context.Context, in *IsSessionCreatedRequest, opts ...client.CallOption) (*IsSessionCreatedResponse, error)
}

type gameSessionService struct {
	c    client.Client
	name string
}

func NewGameSessionService(name string, c client.Client) GameSessionService {
	return &gameSessionService{
		c:    c,
		name: name,
	}
}

func (c *gameSessionService) CreateSession(ctx context.Context, in *CreateSessionRequest, opts ...client.CallOption) (*CreateSessionResponse, error) {
	req := c.c.NewRequest(c.name, "GameSessionService.CreateSession", in)
	out := new(CreateSessionResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameSessionService) StartSession(ctx context.Context, in *StartSessionRequest, opts ...client.CallOption) (*StartSessionResponse, error) {
	req := c.c.NewRequest(c.name, "GameSessionService.StartSession", in)
	out := new(StartSessionResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameSessionService) JoinSession(ctx context.Context, in *JoinSessionRequest, opts ...client.CallOption) (*JoinSessionResponse, error) {
	req := c.c.NewRequest(c.name, "GameSessionService.JoinSession", in)
	out := new(JoinSessionResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameSessionService) LeaveSession(ctx context.Context, in *LeaveSessionRequest, opts ...client.CallOption) (*LeaveSessionResponse, error) {
	req := c.c.NewRequest(c.name, "GameSessionService.LeaveSession", in)
	out := new(LeaveSessionResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameSessionService) IsSessionCreated(ctx context.Context, in *IsSessionCreatedRequest, opts ...client.CallOption) (*IsSessionCreatedResponse, error) {
	req := c.c.NewRequest(c.name, "GameSessionService.IsSessionCreated", in)
	out := new(IsSessionCreatedResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for GameSessionService service

type GameSessionServiceHandler interface {
	CreateSession(context.Context, *CreateSessionRequest, *CreateSessionResponse) error
	StartSession(context.Context, *StartSessionRequest, *StartSessionResponse) error
	JoinSession(context.Context, *JoinSessionRequest, *JoinSessionResponse) error
	LeaveSession(context.Context, *LeaveSessionRequest, *LeaveSessionResponse) error
	IsSessionCreated(context.Context, *IsSessionCreatedRequest, *IsSessionCreatedResponse) error
}

func RegisterGameSessionServiceHandler(s server.Server, hdlr GameSessionServiceHandler, opts ...server.HandlerOption) error {
	type gameSessionService interface {
		CreateSession(ctx context.Context, in *CreateSessionRequest, out *CreateSessionResponse) error
		StartSession(ctx context.Context, in *StartSessionRequest, out *StartSessionResponse) error
		JoinSession(ctx context.Context, in *JoinSessionRequest, out *JoinSessionResponse) error
		LeaveSession(ctx context.Context, in *LeaveSessionRequest, out *LeaveSessionResponse) error
		IsSessionCreated(ctx context.Context, in *IsSessionCreatedRequest, out *IsSessionCreatedResponse) error
	}
	type GameSessionService struct {
		gameSessionService
	}
	h := &gameSessionServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&GameSessionService{h}, opts...))
}

type gameSessionServiceHandler struct {
	GameSessionServiceHandler
}

func (h *gameSessionServiceHandler) CreateSession(ctx context.Context, in *CreateSessionRequest, out *CreateSessionResponse) error {
	return h.GameSessionServiceHandler.CreateSession(ctx, in, out)
}

func (h *gameSessionServiceHandler) StartSession(ctx context.Context, in *StartSessionRequest, out *StartSessionResponse) error {
	return h.GameSessionServiceHandler.StartSession(ctx, in, out)
}

func (h *gameSessionServiceHandler) JoinSession(ctx context.Context, in *JoinSessionRequest, out *JoinSessionResponse) error {
	return h.GameSessionServiceHandler.JoinSession(ctx, in, out)
}

func (h *gameSessionServiceHandler) LeaveSession(ctx context.Context, in *LeaveSessionRequest, out *LeaveSessionResponse) error {
	return h.GameSessionServiceHandler.LeaveSession(ctx, in, out)
}

func (h *gameSessionServiceHandler) IsSessionCreated(ctx context.Context, in *IsSessionCreatedRequest, out *IsSessionCreatedResponse) error {
	return h.GameSessionServiceHandler.IsSessionCreated(ctx, in, out)
}
