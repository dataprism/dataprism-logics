package main

import (
	"github.org/dataprism/dataprism-kfunc/api"
	"github.org/dataprism/dataprism-kfunc/logics"
	consul "github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
)

func main() {
	API := api.CreateAPI("0.0.0.0:8080")

	client, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		logrus.Error(err)
	}

	// -- Profile Providers
	logicsManager := logics.NewManager(client.KV())
	logicsRouter := logics.NewRouter(logicsManager)
	API.RegisterGet("/v1/logics", logicsRouter.ListLogics)
	API.RegisterGet("/v1/logics/{id}", logicsRouter.GetLogic)
	API.RegisterGet("/v1/logics/{id}/versions", logicsRouter.ListLogicVersions)
	API.RegisterGet("/v1/logics/{id}/versions/{version}", logicsRouter.GetLogicVersion)

	API.RegisterPost("/v1/logics", logicsRouter.SetLogic)
	API.RegisterPost("/v1/logics/{id}/versions", logicsRouter.SetLogicVersion)

	API.RegisterDelete("/v1/logics/{id}", logicsRouter.RemoveLogic)
	API.RegisterDelete("/v1/logics/{id}/versions/{version}", logicsRouter.RemoveLogicVersion)

	err = API.Start()
	if err != nil {
		logrus.Error(err)
	}
}