package main

import (
	"context"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/logger"
	nats "github.com/micro/go-plugins/broker/nats/v2"
	"github.com/micro/go-plugins/registry/etcdv3/v2"
	"github.com/songquiz/backend/api"
	"github.com/songquiz/backend/gamesession"
)

//Starts a new Game Session Service
func main() {
	// Uses a new etcd Client Plugin for Service Discovery.
	registry := etcdv3.NewRegistry()
	// Uses a mew NATS Broker Client Plugin for asynchron Messaging.
	broker := nats.NewBroker()

	// A new Service is created.
	service := micro.NewService(
		micro.Name("nq.GameSessionService"),
		micro.Version("latest"),
		micro.Broker(broker),
		micro.Registry(registry),
	)

	service.Init()

	gameSessionVar := gamesession.New(api.NewGameService("nq.GameService", service.Client()),
		api.NewChatService("nq.ChatService", service.Client()),
		api.NewUserService("nq.UserService", service.Client()),
		api.NewQuizContentService("nq.QuizContentService", service.Client()))
	// Register a Game Session Service Server.
	if err := api.RegisterGameSessionServiceHandler(service.Server(), gameSessionVar); err != nil {
		logger.Fatal(err)

	}
	// Register a Game Session Service Server as Subscriber to the Message Broker.
	if err := micro.RegisterSubscriber(
		// The Topic on which will be subscribed to.
		"nq.UserMessage", service.Server(),
		func(ctx context.Context, msg *api.UserMessage) error {
			// Instead of the GameSession Object a function will be called when a Message is published to the topic.
			gameSessionVar.HandleMessage(ctx, msg)
			return nil
		}); err != nil {
		logger.Fatal(err)
	}
	// Starts the Service.
	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}

}
