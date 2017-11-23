package main

import (
	"github.com/dataprism/dataprism-commons/api"
	"github.com/dataprism/dataprism-logics/logics"
	consul "github.com/hashicorp/consul/api"
	nomad "github.com/hashicorp/nomad/api"
	"github.com/sirupsen/logrus"
	"flag"
	"github.com/dataprism/dataprism-logics/evals"
	"github.com/dataprism/dataprism-commons/nodes"
	"strconv"
	Nconsul2 "github.com/dataprism/dataprism-commons/consul"
)

func main() {
	var jobsDir = flag.String("d", "/tmp", "the directory where job information will be stored")
	var port = flag.Int("p", 6300, "the port of the dataprism logics rest api")

	API := api.CreateAPI("0.0.0.0:" + strconv.Itoa(*port))

	client, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		logrus.Error(err)
	}

	storage := Nconsul2.NewStorage(client)

	nomadClient, err := nomad.NewClient(nomad.DefaultConfig())
	if err != nil {
		logrus.Error(err)
	}

	scheduler := logics.NewScheduler(nomadClient, *jobsDir)

	// -- Profile Providers
	logicsManager := logics.NewManager(storage, nomadClient, scheduler)
	logicsRouter := logics.NewRouter(logicsManager)
	API.RegisterGet("/v1/logics", logicsRouter.ListLogics)
	API.RegisterPost("/v1/logics", logicsRouter.SetLogic)

	API.RegisterGet("/v1/logics/{id}", logicsRouter.GetLogic)
	API.RegisterDelete("/v1/logics/{id}", logicsRouter.RemoveLogic)

	API.RegisterGet("/v1/logics/{id}/status", logicsRouter.GetLogicStatus)

	API.RegisterGet("/v1/logics/{id}/versions", logicsRouter.ListLogicVersions)
	API.RegisterPost("/v1/logics/{id}/versions", logicsRouter.SetLogicVersion)

	API.RegisterGet("/v1/logics/{id}/versions/latest", logicsRouter.GetLatestLogicVersion)
	API.RegisterGet("/v1/logics/{id}/versions/{version}", logicsRouter.GetLogicVersion)
	API.RegisterDelete("/v1/logics/{id}/versions/{version}", logicsRouter.RemoveLogicVersion)

	API.RegisterPost("/v1/logics/{id}/versions/{version}/schedule", logicsRouter.Schedule)
	API.RegisterDelete("/v1/logics/{id}/versions/{version}/schedule", logicsRouter.Unschedule)

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