package main

import (
	"log"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-plugins/registry/etcdv3/v2"
	"github.com/songquiz/backend/api"
	"github.com/songquiz/backend/client/contentpreloader"
)

//Starts a new ContentPreloader Service
func main() {
	// Uses a new etcd Client Plugin for Service Discovery.
	registry := etcdv3.NewRegistry()

	// A new Service is created.
	service := micro.NewService(
		micro.Name("nq.ContentPreloader"),
		micro.Version("latest"),
		micro.Registry(registry),
	)

	service.Init()
	// Creates a new ContentPreloader.
	client := contentpreloader.New(
		api.NewQuizContentService("nq.QuizContentService", service.Client()),
	)

	if client == nil {
		log.Fatal(client)
	}
	// Initialize the 3 available Contents.
	client.InitializeData("actors.json")
	client.InitializeData("musicians.json")
	client.InitializeData("scientists.json")

	// Starts the Service.
	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}
}
