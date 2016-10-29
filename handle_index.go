package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func HandleHome(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	RenderTemplate(w, r, "index/home", nil)
}
