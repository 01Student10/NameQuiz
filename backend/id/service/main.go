package main

import (
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-plugins/registry/etcdv3/v2"
	"github.com/micro/go-plugins/store/redis/v2"
	"github.com/micro/go-plugins/wrapper/select/roundrobin/v2"
	"github.com/songquiz/backend/api"
	"github.com/songquiz/backend/id"
)

//Starts a new User Service
func main() {
	//Uses a new etcd Client Plugin for Service Discovery.
	registry := etcdv3.NewRegistry()
	// Uses a new redis Client Plugin for Database.
	storeRedis := redis.NewStore()

	// A new Service is created.
	service := micro.NewService(
		micro.Name("nq.IdService"),
		micro.Version("latest"),
		micro.Registry(registry),
		micro.Store(storeRedis),
		// A Client Wrapper Plugin which is responsible for load balancing Id Service Clients.
		micro.WrapClient(roundrobin.NewClientWrapper()),
	)

	service.Init()
	// Register a Id Service Server.
	if err := api.RegisterIdServiceHandler(service.Server(),
		id.New(storeRedis)); err != nil {
		logger.Fatal(err)

	}
	// Starts the Service.
	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}

}
