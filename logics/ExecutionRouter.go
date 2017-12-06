package logics

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/dataprism/dataprism-commons/utils"
)

type ExectutionRouter struct {
	manager *ExecutionManager
}

func NewExecutionRouter(manager *ExecutionManager) (*ExectutionRouter) {
	return &ExectutionRouter{manager}
}

func (router *ExectutionRouter) Deploy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	resp, err := router.manager.Deploy(r.Context(), id)
	utils.HandleResponse(w, resp, err)
}

func (router *ExectutionRouter) Undeploy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	resp, err := router.manager.Undeploy(r.Context(), id)
	utils.HandleResponse(w, resp, err)
}