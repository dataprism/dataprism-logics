package nodes

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.org/dataprism/dataprism-kfunc/utils"
)

type NodeRouter struct {
	manager *NodeManager
}

func NewRouter(manager *NodeManager) (*NodeRouter) {
	return &NodeRouter{manager:manager}
}

func (router *NodeRouter) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	res, err := router.manager.Get(id)

	utils.HandleResponse(w, res, err)
}

func (router *NodeRouter) List(w http.ResponseWriter, r *http.Request) {
	res, err := router.manager.List()

	utils.HandleResponse(w, res, err)
}