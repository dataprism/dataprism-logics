package logics

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/dataprism/dataprism-commons/utils"
	"io/ioutil"
	"encoding/json"
)

type LogicsRouter struct {
	manager *LogicsManager
}

func NewLogicsRouter(profileProviderManager *LogicsManager) (*LogicsRouter) {
	return &LogicsRouter{manager:profileProviderManager}
}

func (router *LogicsRouter) ListLogics(w http.ResponseWriter, r *http.Request) {
	resp, err := router.manager.ListLogics(r.Context())
	utils.HandleResponse(w, resp, err)
}

func (router *LogicsRouter) GetLogic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	resp, err := router.manager.GetLogic(r.Context(), id)
	utils.HandleResponse(w, resp, err)
}

func (router *LogicsRouter) GetLogicStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	resp, err := router.manager.GetLogicStatus(r.Context(), id)
	utils.HandleResponse(w, resp, err)
}

func (router *LogicsRouter) SetLogic(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var logic Logic
	err = json.Unmarshal(body, &logic)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := router.manager.SetLogic(r.Context(), &logic)
	utils.HandleResponse(w, response, err)
}

func (router *LogicsRouter) RemoveLogic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := router.manager.RemoveLogic(r.Context(), id)
	utils.HandleStatus(w, 200, "Deleted", err)
}