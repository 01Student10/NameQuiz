package quizcontent

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"

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

// The Content is stored into the store. It is not listed in the api.proto file.
type Content struct {
	// the ID of the playlist.
	ListId string
	// The playlist name.
	Name string
	// Contains id string to a ContentEntry struct(which is defined in api.proto)
	Entries map[string]*api.ContentEntry
}

// The QuizContent Service struct.
type Service struct {
	contentStore store.Store
}

// Ctor of the QuizContent Service.
// It receives a Store.
// It returns the Service.
func New(store store.Store) *Service {
	return &Service{
		contentStore: store,
	}
}

// Provides a ContentEntry.
// It receives a listId and an entryId.
// It always returns nil. Results are passed through the response message.
// The response message contains:
// 	Entry - A ContentEntry struct which is defined in the api.proto file and consists of:
//		id - a unique id.
//		name - the name of the person which has to be guessed.
//		path - the path to the picture in the frontend.
//		licence - the licence string of the picture.
//	Response - if there was an Error it will be displayed in the Response string.
//	Otherwise Response will not be transmitted.
func (qc *Service) GetContentEntry(_ context.Context, req *api.GetContentEntryRequest, rsp *api.GetContentEntryResponse) error {
	content, response := qc.readAndParse(req.ListId)
	if response == "" {
		entry, ok := content.Entries[req.EntryId]
		if ok {
			rsp.Entry = entry
		}
	} else {
		rsp.Response = response
	}
	return nil
}

// Gets the content of a playlist and provides an array with the content ids of a playlist.
// It receives a listId.
// It always returns nil. Results are passed through the response message.
// The response message contains:
// 	ListId - the id of the playlist.
//	EntryIds - an array containing the content ids of a playlist.
//	Response - if there was an Error it will be displayed in the Response string.
//	Otherwise Response will not be transmitted.
func (qc *Service) GetContentList(_ context.Context, req *api.GetContentListRequest, rsp *api.GetContentListResponse) error {
	content, response := qc.readAndParse(req.ListId)
	if response == "" {
		entryIds := make([]string, 0)
		for key := range content.Entries {
			entryIds = append(entryIds, key)
		}
		sort.Strings(entryIds)
		rsp.ListId = content.ListId
		rsp.EntryIds = entryIds
	} else {
		rsp.Response = response
	}
	return nil
}

// Provides a mapping of the names of a playlists to the ids of the playlists.
// It always returns nil. Results are passed through the response message.
// The response message contains:
// 	NameToId - a map containing the name of a playlist to the id.
//	Response - if there was an Error it will be displayed in the Response string.
//	Otherwise Response will not be transmitted.
func (qc *Service) GetAllContentLists(_ context.Context, _ *api.GetAllContentListsRequest, rsp *api.GetAllContentListsResponse) error {
	contentRecords, response := qc.getAllRecords()
	if response == "" {
		result := make(map[string]string, len(contentRecords))
		for _, v := range contentRecords {
			result[v.Name] = v.ListId
		}
		rsp.NameToId = result
	} else {
		rsp.Response = response
	}
	return nil
}

// Checks if a certain guess matches the name in the ContentEntry. The guess is not case sensitive.
// It receives a listId, an entryId and a guess.
// It always returns nil. Results are passed through the response message.
// The response message contains:
// 	isRight - tells whether the guess was right or wrong.
//	Response - if there was an Error it will be displayed in the Response string.
//	Otherwise Response will not be transmitted.
func (qc *Service) HasMatch(_ context.Context, req *api.HasMatchRequest, rsp *api.HasMatchResponse) error {
	rsp.IsRight = 2 // Stands for false.
	content, response := qc.readAndParse(req.ListId)
	if response == "" {
		entry, ok := content.Entries[req.EntryId]
		if ok && strings.ToLower(entry.Name) == strings.ToLower(req.Guess) {
			rsp.IsRight = 1 // Stands for true.
		}
	} else {
		rsp.Response = response
	}
	return nil
}

// Creates a new data store entry with from the provided JSON file. The JSON file must have a certain structure.
// It receives an absolute path to the JSON file.
// It always returns nil. Results are passed through the response message.
// The response message contains:
//	Response - if there was an Error it will be displayed in the Response string.
//	Otherwise the response will tell if the playlist is already stored.
func (qc *Service) CreateNewDataSet(_ context.Context, req *api.CreateNewDataSetRequest, rsp *api.CreateNewDataSetResponse) error {

	jsonData, err := ioutil.ReadFile(req.PathToJson)
	if err != nil {
		rsp.Response = fmt.Sprintf("Error while reading file %v with message %v", req.PathToJson, err.Error())
		logger.Info(err)
	} else {
		var content Content
		if parseErr := json.Unmarshal(jsonData, &content); parseErr != nil {
			rsp.Response = fmt.Sprintf("Error while parsing []byte to content struct with msg: %v", parseErr.Error())
			logger.Info(rsp.Response)
		} else {
			// Checks if content is already stored. If so, it will not be stored again.
			records, readErr := qc.contentStore.Read(content.ListId, func(o *store.ReadOptions) { o.Table = StoreTable })

			if readErr != nil && readErr.Error() != "not found" {
				rsp.Response = fmt.Sprintf("Error while reading Content from store: %v", err.Error())
				logger.Info(rsp.Response)
			} else {

				if len(records) == 0 {
					record := store.Record{
						Key:    content.ListId,
						Value:  jsonData,
						Expiry: StoreExpiry,
					}

					writeErr := qc.contentStore.Write(&record, func(o *store.WriteOptions) { o.Table = StoreTable })

					if writeErr != nil {
						rsp.Response = fmt.Sprintf("Error while writing Content to store: %+v", err)
						logger.Info(rsp.Response)
					}

				} else {
					rsp.Response = fmt.Sprintf("Playlist with name : %v already stored..", content.Name)
				}
			}
		}
	}
	return nil
}

// Private funktion which reads a record form the store and parses it to a content struct.
// It receives a recordKey string.
// It Returns a Content struct and an error string. If there was no error the string will be empty.
func (qc *Service) readAndParse(recordKey string) (*Content, string) {
	var errorMsg = ""
	value, readErr := qc.contentStore.Read(recordKey, func(o *store.ReadOptions) { o.Table = StoreTable })
	if readErr != nil && readErr.Error() != "not found" {
		errorMsg = fmt.Sprintf("Error while reading from store with msg: %v", readErr.Error())
		logger.Info(errorMsg)
	}
	var content Content
	// parse byte array to Content Object.
	if len(value) != 0 {
		parseErr := json.Unmarshal(value[0].Value, &content)

		if parseErr != nil {
			errorMsg = fmt.Sprintf("Error while parsing []byte to Content struct with msg: %v", parseErr.Error())
			logger.Info(errorMsg)
		}
	} else {
		errorMsg = fmt.Sprintf("No record with id: %v found", recordKey)
		logger.Info(errorMsg)
	}
	return &content, errorMsg
}

//Private function which gets all contents of the store.
// It returns an array of the Content struct's and an error string if something went wrong.
func (qc *Service) getAllRecords() ([]*Content, string) {
	response := ""
	contentRecords, listErr := qc.contentStore.List(func(o *store.ListOptions) {})
	contents := make([]*Content, len(contentRecords))
	if listErr != nil {
		response = fmt.Sprintf("Error while listing all records from store with msg: %v", listErr.Error())
		logger.Info(response)
	} else {
		for i, v := range contentRecords {
			content, errorMsg := qc.readAndParse(v)
			if errorMsg == "" {
				contents[i] = content
			} else {
				contents = nil
				response = errorMsg
				break
			}
		}
	}
	return contents, response
}
