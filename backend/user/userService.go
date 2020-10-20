package user

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
	// Not found response of the store should not cause an error.
	notFound = "not found"
	// Secret password string for the Frontend.
	SecretPw = "********"
)

// The User Service struct.
type Service struct {
	userStore store.Store
	idService api.IdService
}

// Ctor of the User Service.
// It receives a Store and an Id Service Client.
// It returns the Service.
func New(store store.Store, idService api.IdService) *Service {
	return &Service{
		userStore: store,
		idService: idService,
	}
}

// Creates a new User for the Namequiz.
// It receives the request and response message defined in the api.proto.
// It always returns nil. Results are passed through the response message.
// The response message contains:
// 	HasCreatedUser - defines weather the creation operation was successful.
//	User - the user struct.
//	Response - if there was an Error it will be displayed in the Response string.
//	Otherwise Response will not be transmitted.
func (us *Service) CreateUser(_ context.Context, req *api.CreateUserRequest, rsp *api.CreateUserResponse) error {
	rsp.HasCreatedUser = 2 // stands for false
	// Check if user already exists
	users, errorMsg := us.getAllRecords()
	if errorMsg == "" {
		alreadyExists := false
		for _, user := range users {
			if user.Username == req.Username && req.Password == user.Password && req.Mail == user.Mail {
				alreadyExists = true
			}
		}
		if !alreadyExists {
			// Calls the Id Service to get a unique ID.
			res, _ := us.idService.GetId(context.Background(), &api.GetIdRequest{})
			if res.Response == "" {
				user := api.User{
					Id:       res.Id,
					Username: req.Username,
					Password: req.Password,
					Mail:     req.Mail,
				}
				// Parse User Struct to bytearray - so the dbuser-database key/value-store can use it.
				value, err := json.Marshal(&user)
				if err != nil {
					rsp.Response = fmt.Sprintf("Parsing Struct to byte Array failed with error: %v", err.Error())
					logger.Info(rsp.Response)
				} else {
					// Create the record for the Store.
					// The Id is also used as Key for the record.
					record := store.Record{
						Key:    fmt.Sprint(res.Id),
						Value:  value,
						Expiry: StoreExpiry,
					}
					// Writes the record to the store.
					writeErr := us.userStore.Write(&record, func(o *store.WriteOptions) { o.Table = StoreTable })

					if writeErr != nil {
						rsp.Response = fmt.Sprintf("error while writing User to store: %v", writeErr.Error())
						logger.Info(rsp.Response)
					} else {
						rsp.HasCreatedUser = 1   // Stands for true.
						user.Password = SecretPw // Defaces the Password in the response.
						rsp.User = &user
					}
				}
			} else {
				rsp.Response = res.Response
			}
		} else {
			rsp.Response = fmt.Sprintf("User: %v with mail: %v already exists,", req.Username, req.Mail)
		}
	} else {
		rsp.Response = errorMsg
	}

	return nil
}

// Calls the Store to get the desired user.
// It receives the request and response message defined in the api.proto.
// It always returns nil. Responses are passed through the response message.
// The response message contains:
//	User - the user struct.
//	Response - if there was an Error it will be displayed in the Response string.
//	Otherwise Response will not be transmitted.
func (us *Service) GetUser(_ context.Context, req *api.GetUserRequest, rsp *api.GetUserResponse) error {

	value, readErr := us.userStore.Read(fmt.Sprint(req.Id), func(o *store.ReadOptions) { o.Table = StoreTable })
	if readErr != nil && readErr.Error() != notFound {
		rsp.Response = fmt.Sprintf("error while reading User from store: %+v", readErr.Error())
		logger.Info(rsp.Response)
	} else {
		if len(value) != 0 {
			var user api.User
			// Parse byte array to user Struct.
			parseErr := json.Unmarshal(value[0].Value, &user)

			if parseErr != nil {
				rsp.Response = fmt.Sprintf("Error while parsing byte array to User struct: %v", parseErr.Error())
				logger.Info(rsp.Response)
			} else {
				rsp.User = &user
			}
		} else {
			rsp.Response = fmt.Sprintf("No User with id : %v found", req.Id)
			logger.Info(rsp.Response)
		}
	}
	return nil
}

// Calls the Store to get all user records.
// It receives the request and response message defined in the api.proto.
// It always returns nil. Responses are passed through the response message.
// The response message contains:
//	User - the user struct.
//	Response - if there was an Error it will be displayed in the Response string.
//	Otherwise Response will not be transmitted.
func (us *Service) GetAllUsers(_ context.Context, _ *api.GetAllUsersRequest, rsp *api.GetAllUsersResponse) error {

	response, errorMsg := us.getAllRecords()
	if errorMsg != "" {
		rsp.Response = errorMsg
	} else {
		rsp.Users = response
	}

	return nil
}

// Checks if the a specific user is stored in the store.
// It receives the request and response message defined in the api.proto.
// It always returns nil. Responses are passed through the response message.
// The response message contains:
//	UserAvailable - tells whether the user is available.
//	Response - if there was an Error it will be displayed in the Response string.
//	Otherwise Response will not be transmitted.
func (us *Service) HasUser(_ context.Context, req *api.HasUserRequest, rsp *api.HasUserResponse) error {
	allUsers, errorMsg := us.getAllRecords()
	var result int32 = 2 // Stands for false.
	if errorMsg == "" {
		for _, v := range allUsers {
			if req.Id == v.Id {
				result = 1 // Stands for true.
			}
		}
	} else {
		rsp.Response = errorMsg
	}
	rsp.UserAvailable = result
	return nil
}

// Checks if the username and password fits to the stored user.
// It receives the request and response message defined in the api.proto.
// It always returns nil. Responses are passed through the response message.
// The response message contains:
//	HasLoggedIn - tells whether the user was able to log in.
//	User - the user struct.
func (us *Service) Login(_ context.Context, req *api.LoginRequest, rsp *api.LoginResponse) error {
	allUsers, _ := us.getAllRecords()
	var result int32 = 2 // Stands for false.
	for _, v := range allUsers {
		if req.UserName == v.Username && req.Password == v.Password {
			result = 1            // Stands for true.
			v.Password = SecretPw // Defaces the Password in the response.
			rsp.User = v
		}
	}

	rsp.HasLoggedIn = result
	return nil
}

// Private function to read a Record from store and parse it to user struct.
// It receives the key of the record to be read.
// It returns the user struct and a possible error message.
func (us *Service) readAndParse(recordKey string) (resUser *api.User, errMsg string) {
	var user api.User
	value, readErr := us.userStore.Read(recordKey, func(o *store.ReadOptions) { o.Table = StoreTable })
	if readErr != nil {
		errMsg = fmt.Sprintf("Error while reading user from store: %v", readErr)
		logger.Infof(errMsg)
	} else {
		// parse byte array to User struct.
		parseErr := json.Unmarshal(value[0].Value, &user)

		if parseErr != nil {
			errMsg = fmt.Sprintf("Error while parsing byte array to user struct: %v", parseErr)
			logger.Infof(errMsg)
		}
	}
	return &user, errMsg
}

// Private function which reads all records of a store and parses them to user structs.
// It returns an user struct array and a possible error message.
func (us *Service) getAllRecords() (userRecordList []*api.User, errorMsg string) {
	userRecords, listErr := us.userStore.List(func(o *store.ListOptions) {})
	response := make([]*api.User, len(userRecords))
	if listErr != nil {
		errorMsg = fmt.Sprintf("Error while reading User Records from store: %v", listErr)
		logger.Infof(errorMsg)
	} else {
		for i, v := range userRecords {
			user, readErr := us.readAndParse(v)
			if readErr == "" {
				response[i] = user
			} else {
				errorMsg = readErr
				break
			}
		}
	}
	return response, errorMsg
}
