package logics

import (
	"github.com/dataprism/dataprism-commons/core"
	"github.com/dataprism/dataprism-commons/api"
)

func CreateRoutes(platform *core.Platform, API *api.Rest) {
	// -- Profile Providers
	logicsManager := NewLogicsManager(platform)
	logicsRouter := NewLogicsRouter(logicsManager)
	API.RegisterGet("/v1/logics", logicsRouter.ListLogics)
	API.RegisterPost("/v1/logics", logicsRouter.SetLogic)

	API.RegisterGet("/v1/logics/{id}", logicsRouter.GetLogic)
	API.RegisterDelete("/v1/logics/{id}", logicsRouter.RemoveLogic)

	API.RegisterGet("/v1/logics/{id}/status", logicsRouter.GetLogicStatus)

	executionManager := NewExecutionManager(platform, logicsManager)
	executionRouter := NewExecutionRouter(executionManager)
	API.RegisterPost("/v1/logics/{id}/schedule", executionRouter.Deploy)
	API.RegisterDelete("/v1/logics/{id}/schedule", executionRouter.Undeploy)

	//evaluationManager := evals.NewManager(nomadClient)
	//evaluationRouter := evals.NewRouter(evaluationManager)
	//API.RegisterGet("/v1/evaluations/{id}", evaluationRouter.Get)
	//API.RegisterGet("/v1/evaluations/{id}/events", evaluationRouter.Events)

	//nodeManager := nodes.NewManager(nomadClient)
	//nodeRouter := nodes.NewRouter(nodeManager)
	//API.RegisterGet("/v1/nodes", nodeRouter.List)
	//API.RegisterGet("/v1/nodes/{id}", nodeRouter.Get)
}
