package main

import (
	"github.org/dataprism/dataprism-kfunc/api"
	"github.org/dataprism/dataprism-kfunc/logics"
	consul "github.com/hashicorp/consul/api"
)

func main() {
	API := api.CreateAPI("0.0.0.0:8080")

	client, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		panic(err)
	}

	// -- Profile Providers
	logicsManager := logics.NewManager(client.KV())
	logicsRouter := logics.NewRouter(logicsManager)
	API.RegisterSecureGet("/v1/logics", logicsRouter.ListLogics)
	API.RegisterSecureGet("/v1/logics/{id}", logicsRouter.GetLogic)
	API.RegisterSecurePost("/v1/logics", logicsRouter.SetLogic)
	API.RegisterSecureDelete("/v1/logics/{id}", logicsRouter.RemoveLogic)
}