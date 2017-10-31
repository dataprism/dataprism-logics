package logics

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.org/dataprism/dataprism-kfunc/utils"
	"io/ioutil"
	"encoding/json"
)

type LogicsRouter struct {
	manager *LogicsManager
}

func NewRouter(profileProviderManager *LogicsManager) (*LogicsRouter) {
	return &LogicsRouter{manager:profileProviderManager}
}

func (router *LogicsRouter) ListLogics(w http.ResponseWriter, r *http.Request) {
	intents, err := router.manager.ListLogics(r.Context())
	utils.HandleResponse(w, intents, err)
}

func (router *LogicsRouter) GetLogic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	intents, err := router.manager.GetLogic(r.Context(), id)
	utils.HandleResponse(w, intents, err)
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