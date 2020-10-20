package main

import (
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/server"
	nats "github.com/micro/go-plugins/broker/nats/v2"
	"github.com/micro/go-plugins/registry/etcdv3/v2"
	"github.com/songquiz/backend/api"
	"github.com/songquiz/backend/evaluation"
)

//Starts a new Evaluation Service
func main() {
	// Uses a new etcd Client Plugin for Service Discovery.
	registry := etcdv3.NewRegistry()
	// Uses a mew NATS Broker Client Plugin for asynchron Messaging.
	broker := nats.NewBroker()

	// A new Service is created.
	service := micro.NewService(
		micro.Name("nq.EvaluationService"),
		micro.Version("latest"),
		micro.Broker(broker),
		micro.Registry(registry),
	)

	service.Init()
	// Register a Evaluation Service Server as Subscriber to the Message Broker.
	err := micro.RegisterSubscriber(
		// The Topic on which will be subscribed to.
		"nq.api.Chat", service.Server(),
		evaluation.New(api.NewQuizContentService("nq.QuizContentService", service.Client()),
			micro.NewEvent("nq.UserMessage", service.Client())), // The Event which will be published.
		server.SubscriberQueue("evaluation")) // Subscribe to the NATS Queue Group called evaluation.
	if err != nil {
		panic(err)
	}

	// Starts the Service.
	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}

}
