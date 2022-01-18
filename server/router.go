package server

import (
	"github.com/gorilla/mux"
)

// RegistHTTPRouter http router for manager api
func (m *MDS) RegistHTTPRouter() *mux.Router {
	router := mux.NewRouter().SkipClean(true)
	peer := router.NewRoute().PathPrefix("/peer").Subrouter()

	// peer router
	peer.Methods("GET").Path("/list").HandlerFunc(nil)
	peer.Methods("PUT").Path("/add").HandlerFunc(nil)
	peer.Methods("PUT").Path("/remove").HandlerFunc(nil)
	peer.Methods("GET").Path("/status").HandlerFunc(m.Status)

	return router
}
