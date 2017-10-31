package main

import (
	"github.org/dataprism/dataprism-kfunc/api"
	"github.org/dataprism/dataprism-kfunc/logics"
	consul "github.com/hashicorp/consul/api"
	nomad "github.com/hashicorp/nomad/api"
	"github.com/sirupsen/logrus"
	"flag"
	"github.org/dataprism/dataprism-kfunc/evals"
	"github.org/dataprism/dataprism-kfunc/nodes"
)

func main() {
	var jobsDir = flag.String("d", "/tmp", "the directory where job information will be stored")

	API := api.CreateAPI("0.0.0.0:8080")

	client, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		logrus.Error(err)
	}

	nomadClient, err := nomad.NewClient(nomad.DefaultConfig())
	if err != nil {
		logrus.Error(err)
	}

	scheduler := logics.NewScheduler(nomadClient, *jobsDir)

	// -- Profile Providers
	logicsManager := logics.NewManager(client.KV(), nomadClient, scheduler)
	logicsRouter := logics.NewRouter(logicsManager)
	API.RegisterGet("/v1/logics", logicsRouter.ListLogics)
	API.RegisterGet("/v1/logics/{id}", logicsRouter.GetLogic)
	API.RegisterGet("/v1/logics/{id}/status", logicsRouter.GetLogicStatus)
	API.RegisterGet("/v1/logics/{id}/versions", logicsRouter.ListLogicVersions)
	API.RegisterGet("/v1/logics/{id}/versions/latest", logicsRouter.GetLatestLogicVersion)
	API.RegisterGet("/v1/logics/{id}/versions/{version}", logicsRouter.GetLogicVersion)

	API.RegisterPost("/v1/logics/{id}/versions/{version}/schedule", logicsRouter.Schedule)
	API.RegisterDelete("/v1/logics/{id}/versions/{version}/schedule", logicsRouter.Unschedule)

	API.RegisterPost("/v1/logics", logicsRouter.SetLogic)
	API.RegisterPost("/v1/logics/{id}/versions", logicsRouter.SetLogicVersion)

	API.RegisterDelete("/v1/logics/{id}", logicsRouter.RemoveLogic)
	API.RegisterDelete("/v1/logics/{id}/versions/{version}", logicsRouter.RemoveLogicVersion)

	evaluationManager := evals.NewManager(nomadClient)
	evaluationRouter := evals.NewRouter(evaluationManager)
	API.RegisterGet("/v1/evaluations/{id}", evaluationRouter.Get)
	API.RegisterGet("/v1/evaluations/{id}/events", evaluationRouter.Events)

	nodeManager := nodes.NewManager(nomadClient)
	nodeRouter := nodes.NewRouter(nodeManager)
	API.RegisterGet("/v1/nodes", nodeRouter.List)
	API.RegisterGet("/v1/nodes/{id}", nodeRouter.Get)


	err = API.Start()
	if err != nil {
		logrus.Error(err)
	}
}