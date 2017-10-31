package logics

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.org/dataprism/dataprism-kfunc/utils"
	"io/ioutil"
	"encoding/json"
	"strconv"
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

func (router *LogicsRouter) ListLogicVersions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	intents, err := router.manager.ListLogicVersions(r.Context(), id)
	utils.HandleResponse(w, intents, err)
}

func (router *LogicsRouter) GetLogic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	intents, err := router.manager.GetLogic(r.Context(), id)
	utils.HandleResponse(w, intents, err)
}

func (router *LogicsRouter) GetLogicVersion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	strVersion := vars["version"]

	version, err := strconv.Atoi(strVersion)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	intents, err := router.manager.GetLogicVersion(r.Context(), id, version)
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

func (router *LogicsRouter) SetLogicVersion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var logicVersion LogicVersion
	err = json.Unmarshal(body, &logicVersion)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := router.manager.SetLogicVersion(r.Context(), id, &logicVersion)
	utils.HandleResponse(w, response, err)
}

func (router *LogicsRouter) RemoveLogic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := router.manager.RemoveLogic(r.Context(), id)
	utils.HandleStatus(w, 200, "Deleted", err)
}

func (router *LogicsRouter) RemoveLogicVersion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	strVersion := vars["version"]

	version, err := strconv.Atoi(strVersion)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = router.manager.RemoveLogicVersion(r.Context(), id, version)
	utils.HandleStatus(w, 200, "Deleted", err)
}