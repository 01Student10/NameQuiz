package evaluation

import (
	"context"
	"encoding/json"

	"github.com/micro/go-micro/v2"
	gatewayProto "github.com/micro/go-micro/v2/api/proto"
	"github.com/micro/go-micro/v2/logger"
	"github.com/songquiz/backend/api"
)

// The Evaluation Service struct.
type Evaluation struct {
	quizContent api.QuizContentService
	publisher   micro.Event
}

// Ctor of the Evaluation Service.
// It receives a QuizContent Client and an Event on which will be published.
// It returns the Service.
func New(quizContent api.QuizContentService, publisher micro.Event) *Evaluation {
	return &Evaluation{
		quizContent: quizContent,
		publisher:   publisher,
	}
}

// This Function is called whenever a Message is published to the "np.api.Chat" topic.
// It receives a gatewayProto Event (which is defined in go-micro's api/proto-file)
// In this function, since there is no Caller to return an Error which then can be handled.
// The function throws an Error if something goes wrong.
// It always returns nil.
func (s *Evaluation) Handle(_ context.Context, input *gatewayProto.Event) error {
	var userMessage api.UserMessage
	if parseErr := json.Unmarshal([]byte(input.GetData()), &userMessage); parseErr != nil {
		logger.Fatal("EvaluationService: Error while parsing byte array to UserMessage struct: %v", parseErr)
	}
	res, err := s.quizContent.HasMatch(context.Background(), &api.HasMatchRequest{
		ListId:  userMessage.ListId,
		EntryId: userMessage.EntryId,
		Guess:   userMessage.Guess,
	})
	if err != nil {
		logger.Fatal(err)
	}
	userMessage.WasRight = 2 // Stands for false.
	if res.IsRight == 1 {
		userMessage.WasRight = 1 // Stands for true.
	}
	if pubErr := s.publisher.Publish(context.Background(), &userMessage); pubErr != nil {
		logger.Fatal(pubErr)
	}

	return nil
}
