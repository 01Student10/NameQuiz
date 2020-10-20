package id

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
	// The string at which the record will be stored.
	KEY = "key"
)

// The Id Service struct.
type Service struct {
	idStore store.Store
}

// The id struct which will be stored into the store.
type Id struct {
	IdCounter int32
}

// Ctor of the Id Service.
// It receives a store struct.
// It returns the Service.
func New(store store.Store) *Service {
	return &Service{
		idStore: store,
	}
}

// Provides a unique id.
// It reads the IdCounter value from the Id struct increments it and stores it to the store again.
// It always returns nil. Results are passed through the response message.
// The response message contains:
// 	Id - the unique ID.
//	Response - if there was an Error it will be displayed in the Response string.
//	Otherwise Response will be an empty string.
func (is *Service) GetId(_ context.Context, _ *api.GetIdRequest, rsp *api.GetIdResponse) error {

	records, readErr := is.idStore.Read(KEY, func(o *store.ReadOptions) { o.Table = StoreTable })
	if readErr != nil && readErr.Error() != "not found" {
		rsp.Response = fmt.Sprintf("Error while reading Id from store: %v", readErr.Error())
		logger.Info(rsp.Response)
	}
	var id Id
	if len(records) == 0 {
		id = Id{
			IdCounter: 1, // Starts always with 1. 0 will not be parsed.
		}

	} else {
		// parse byte array to Id struct.
		parseErr := json.Unmarshal(records[0].Value, &id)

		if parseErr != nil {
			rsp.Response = fmt.Sprintf("Error while parsing []byte to id struct with msg: %v", parseErr.Error())
			logger.Info(rsp.Response)
		}

		id.IdCounter++
	}

	rawData, parseErr := json.Marshal(&id)
	if parseErr != nil {
		rsp.Response = fmt.Sprintf("Error while parsing id struct to []byte with msg: %v", parseErr.Error())
		logger.Info(rsp.Response)
	}
	record := store.Record{
		Key:    KEY,
		Value:  rawData,
		Expiry: StoreExpiry,
	}
	writeErr := is.idStore.Write(&record, func(o *store.WriteOptions) { o.Table = StoreTable })

	if writeErr != nil {
		rsp.Response = fmt.Sprintf("error while writing Id to store: %v", writeErr.Error())
		logger.Info(rsp.Response)
	}
	rsp.Id = id.IdCounter

	return nil
}
