package evals

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/dataprism/dataprism-commons/utils"
)

type EvaluationRouter struct {
	manager *EvaluationManager
}

func NewRouter(manager *EvaluationManager) (*EvaluationRouter) {
	return &EvaluationRouter{manager:manager}
}

func (router *EvaluationRouter) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	res, err := router.manager.Get(id)

	utils.HandleResponse(w, res, err)
}

func (router *EvaluationRouter) Events(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	res, err := router.manager.Events(id)

	utils.HandleResponse(w, res, err)
}