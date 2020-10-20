package game

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/store"
	"github.com/songquiz/backend/api"
)

const (
	// The table in which the record is stored.
	StoreTable = ""
	// Defines the lifetime of a store record. 0 tells that the record will not expire.
	StoreExpiry = 0
)

// The Game Service struct.
type Service struct {
	gameStore store.Store
	idService api.IdService
}

// Ctor of the Game Service.
// It receives a Store and an Id Service Client.
// It returns the Service.
func New(store store.Store, idService api.IdService) *Service {
	return &Service{
		gameStore: store,
		idService: idService,
	}
}

// Creates a new Game for the Namequiz.
// It receives the title, the amount of songs the playlist id and the user id of the owner of the game.
// It always returns nil. Results are passed through the response message.
// The response message contains:
//	Id - the Id of the game.
//	Response - if there was an Error it will be displayed in the Response string.
//	Otherwise Response will not be transmitted.
func (gs *Service) CreateGame(_ context.Context, req *api.CreateGameRequest, rsp *api.CreateGameResponse) error {
	res, _ := gs.idService.GetId(context.Background(), &api.GetIdRequest{})
	logger.Infof("in create game- id: %v, res: %v!", res.Id, res.Response)
	if res.Response == "" {
		game := api.Game{
			Id:             res.Id,
			Title:          req.Title,
			AmountOfRounds: req.AmountOfSongs,
			PlaylistId:     req.PlaylistId,
			Owner:          req.Owner,
		}
		// parse Game Struct to bytearray - so the dbgame key value store can use it.
		value, err := json.Marshal(&game)
		if err != nil {
			res.Response = fmt.Sprintf("Parsing Game Struct to byte Array failed with error: %v", err.Error())
			logger.Info(res.Response)
		} else {
			record := store.Record{
				Key:    fmt.Sprint(res.Id),
				Value:  value,
				Expiry: StoreExpiry,
			}
			writeErr := gs.gameStore.Write(&record, func(o *store.WriteOptions) { o.Table = StoreTable })

			if writeErr != nil {
				res.Response = fmt.Sprintf("Error while writing Game to store: %+v", err.Error())
				logger.Info(res.Response)
			} else {
				rsp.Id = res.Id
			}
		}
	} else {
		rsp.Response = res.Response
	}
	return nil
}

// Reads a Game from the store and returns it.
// It receives a game id.
// It always returns nil. Results are passed through the response message.
// The response message contains:
//	Id - the Id of the game.
//	Response - if there was an Error it will be displayed in the Response string.
//	Otherwise Response will not be transmitted.
func (gs *Service) GetGame(_ context.Context, req *api.GetGameRequest, rsp *api.GetGameResponse) error {
	value, readErr := gs.gameStore.Read(fmt.Sprint(req.Id), func(o *store.ReadOptions) { o.Table = StoreTable })
	if readErr != nil && readErr.Error() != "not found" {
		logger.Info("error while reading Game to store: %+v", readErr)
	}
	if len(value) != 0 {
		var game api.Game
		// parse byte array to Game Object.
		parseErr := json.Unmarshal(value[0].Value, &game)

		if parseErr != nil {
			rsp.Response = fmt.Sprintf("Error while parsing []byte to Game struct: %v", parseErr.Error())
			logger.Info(rsp.Response)
		}
		rsp.Game = &game
	} else {
		rsp.Response = fmt.Sprintf("No Game with id : %v found", req.Id)
		logger.Info(rsp.Response)
	}
	return nil
}

// Deletes a game from the store.
// It receives a game id.
// It always returns nil. Results are passed through the response message.
// The response message contains:
//	Response - if there was an Error it will be displayed in the Response string.
//	Otherwise Response will not be transmitted.
func (gs *Service) DeleteGame(_ context.Context, req *api.DeleteGameRequest, rsp *api.DeleteGameResponse) error {
	if err := gs.gameStore.Delete(fmt.Sprint(req.Id), func(o *store.DeleteOptions) { o.Table = StoreTable }); err != nil && err.Error() != "not found" {
		logger.Infof("error while deleting game : %v ", err.Error())
		rsp.Response = err.Error()
	}
	return nil
}

// Calls the Store to get all game records and parses them to game struct's.
// It receives the request and response message defined in the api.proto.
// It always returns nil. Responses are passed through the response message.
// The response message contains:
//	Games - an array containing all the games in the store.
//	Response - if there was an Error it will be displayed in the Response string.
//	Otherwise Response will not be transmitted.
func (gs *Service) GetAllGames(_ context.Context, _ *api.GetAllGamesRequest, rsp *api.GetAllGamesResponse) error {

	games, errorMsg := gs.getAllRecords()
	rsp.Games = games
	rsp.Response = errorMsg

	return nil
}

// Private function to read a Record from store and parse it to a game struct.
// It receives the key of the record to be read.
// It returns the user struct and a possible error message.
func (gs *Service) readAndParse(v string) (*api.Game, string) {
	errorMsg := ""
	value, readErr := gs.gameStore.Read(v, func(o *store.ReadOptions) { o.Table = StoreTable })
	if readErr != nil && readErr.Error() != "not found" {
		errorMsg = fmt.Sprintf("Error while reading Game to store: %v", readErr.Error())
		logger.Info(errorMsg)
	}
	var game api.Game
	// parse byte array to Game Object.
	parseErr := json.Unmarshal(value[0].Value, &game)

	if parseErr != nil {
		errorMsg = fmt.Sprintf("Error while parsing byte array to Game struct: %v", parseErr.Error())
		logger.Info(errorMsg)
	}
	return &game, errorMsg
}

// Private function which reads all records of a store and parses them to a game struct.
// It returns an user struct array and a possible error message.
func (gs *Service) getAllRecords() ([]*api.Game, string) {
	errorMsg := ""
	gameRecords, listErr := gs.gameStore.List(func(o *store.ListOptions) {})

	if listErr != nil {
		errorMsg = fmt.Sprintf("Error while writing Game to store: %v", listErr.Error())
		logger.Info(errorMsg)
	}
	games := make([]*api.Game, len(gameRecords))

	for i, v := range gameRecords {
		game, err := gs.readAndParse(v)
		if err == "" {
			games[i] = game
		} else {
			errorMsg = err
		}
	}
	return games, errorMsg
}
