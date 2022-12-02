//go:build wireinject
// +build wireinject

package main

import (
	"crm-service-go/app/clients"
	"crm-service-go/app/controllers"
	controllerInside "crm-service-go/app/controllers/inside"
	"crm-service-go/app/repositories"
	"crm-service-go/app/services"
	serviceInside "crm-service-go/app/services/inside"
	"crm-service-go/server"
	"github.com/google/wire"
)

func initServer() *server.Server {
	wire.Build(
		repositories.ProviderRepositorySet,
		clients.ProviderHttpClientSet,
		services.ProviderServiceSet,
		serviceInside.ProviderInsideServiceSet,
		controllers.ProviderControllerSet,
		controllerInside.ProviderInsideControllerSet,
		server.ProviderServerSet,
	)
	return &server.Server{}
}
