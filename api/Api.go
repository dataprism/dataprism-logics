package api

import (
	"net/http"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/rs/cors"
)

type Rest struct {
	router *mux.Router
	bind string
}

func CreateAPI(bind string) *Rest {
	return &Rest{
		router: mux.NewRouter(),
		bind: bind,
	}
}

func (a *Rest) Start() error {
	handler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"},
		AllowCredentials: true,
		AllowedHeaders: []string{"Authorization", "Content-Type"},
		Debug: false,
	}).Handler(a.router)

	log.Infof("API listening on %s", a.bind);
	return http.ListenAndServe(a.bind, handler)
}

// == GET =====================================================================

func (a *Rest) RegisterGet(path string, h func(w http.ResponseWriter, r *http.Request)) *mux.Route {
	log.WithFields(log.Fields{
		"method": "GET",
		"secured": false,
	}).Info(path)

	return a.router.
		Path(path).
		Methods("GET").
		Handler(http.HandlerFunc(h))
}

// == POST ====================================================================

func (a *Rest) RegisterPost(path string, h func(w http.ResponseWriter, r *http.Request)) *mux.Route {
	log.WithFields(log.Fields{
		"method": "POST",
		"secured": false,
	}).Info(path)

	return a.router.
		Path(path).
		Methods("POST").
		Handler(http.HandlerFunc(h))
}

// == PUT =====================================================================

func (a *Rest) RegisterPut(path string, h func(w http.ResponseWriter, r *http.Request)) *mux.Route {
	log.WithFields(log.Fields{
		"method": "PUT",
		"secured": false,
	}).Info(path)

	return a.router.
		Path(path).
		Methods("PUT").
		Handler(http.HandlerFunc(h))
}

// == DELETE ==================================================================

func (a *Rest) RegisterDelete(path string, h func(w http.ResponseWriter, r *http.Request)) *mux.Route {
	log.WithFields(log.Fields{
		"method": "DELETE",
		"secured": false,
	}).Info(path)

	return a.router.
		Path(path).
		Methods("DELETE").
		Handler(http.HandlerFunc(h))
}