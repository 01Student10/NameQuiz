package contentpreloader

import (
	"context"

	"github.com/micro/go-micro/v2/logger"
	"github.com/songquiz/backend/api"
)

// The ContentPreLoader struct.
type ContentPreLoader struct {
	quizContent api.QuizContentService
}

// Ctor of the Game Service.
// It receives a QuizContent Service Client.
// It returns the Service
func New(quizContent api.QuizContentService) *ContentPreLoader {
	return &ContentPreLoader{
		quizContent: quizContent,
	}
}

// Calls the QuizContent Service to Create a new Data Set with the given path to the .json file.
// Prints the response.
func (ci *ContentPreLoader) InitializeData(path string) {
	response, err := ci.quizContent.CreateNewDataSet(context.Background(),
		&api.CreateNewDataSetRequest{
			PathToJson: path,
		})
	if err != nil {
		logger.Error(err)
	}
	logger.Infof("new DataSet created with Response : %v", response)
}
