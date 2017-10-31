package jobs

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.org/dataprism/dataprism-kfunc/utils"
)

type JobsRouter struct {
	manager *JobsManager
}

func NewRouter(manager *JobsManager) (*JobsRouter) {
	return &JobsRouter{manager:manager}
}

func (router *JobsRouter) List(w http.ResponseWriter, r *http.Request) {
	res, err := router.manager.List("")

	utils.HandleResponse(w, res, err)
}

func (router *JobsRouter) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	obj, err := router.manager.Get(id);
	utils.HandleResponse(w, obj, err)
}