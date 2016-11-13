package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func HandleImageNew(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	RenderTemplate(w, r, "images/new", nil)
}

// TODO
func HandleImageCreate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
}
