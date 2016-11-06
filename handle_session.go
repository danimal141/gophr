package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func HandleSessionNew(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	next := r.URL.Query().Get("next")
	RenderTemplate(w, r, "sessions/new", map[string]interface{}{"Next": next})
}

// TODO: Implement
func HandleSessionCreate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
}

func HandleSessionDestroy(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	session := RequestSession(r)
	if session != nil {
		err := globalSessionStore.Delete(session)
		if err != nil {
			panic(err)
		}
	}
	RenderTemplate(w, r, "sessions/destroy", nil)
}
