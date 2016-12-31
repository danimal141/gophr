package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func HandleSessionNew(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	next := r.URL.Query().Get("next")
	RenderTemplate(w, r, "sessions/new", map[string]interface{}{"Next": next})
}

func HandleSessionCreate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	next := r.FormValue("next")

	user, err := FindUser(username, password)
	if err != nil {
		if IsValidationError(err) {
			RenderTemplate(w, r, "sessions/new", map[string]interface{}{
				"User":  user,
				"Error": err.Error(),
				"Next":  next,
			})
			return
		}
		log.Fatal(err)
	}

	session := FindOrCreateSession(w, r)
	session.UserID = user.ID
	err = globalSessionStore.Save(session)
	if err != nil {
		log.Fatal(err)
	}
	if next == "" {
		next = "/"
	}
	http.Redirect(w, r, next+"?flash=Logged+in", http.StatusFound)
}

func HandleSessionDestroy(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	session := RequestSession(r)
	if session != nil {
		err := globalSessionStore.Delete(session)
		if err != nil {
			log.Fatal(err)
		}
	}
	RenderTemplate(w, r, "sessions/destroy", nil)
}
