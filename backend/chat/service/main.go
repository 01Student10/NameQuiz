package main

import (
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-plugins/registry/etcdv3/v2"
	"github.com/micro/go-plugins/wrapper/select/roundrobin/v2"
	"github.com/songquiz/backend/api"
	"github.com/songquiz/backend/chat"
)

//Starts a new Chat Service
func main() {
	// Uses a new etcd Client Plugin for Service Discovery.
	registry := etcdv3.NewRegistry()

	// A new Service is created.
	service := micro.NewService(
		micro.Name("nq.ChatService"),
		micro.Version("latest"),
		micro.Registry(registry),
		micro.WrapClient(roundrobin.NewClientWrapper()),
	)

	service.Init()

	// Register a Chat Service Server.
	if err := api.RegisterChatServiceHandler(service.Server(),
		chat.New()); err != nil {
		logger.Fatal(err)

	}
	// Starts the Service.
	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}

}
