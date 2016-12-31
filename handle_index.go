package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func HandleHome(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	imgs, err := globalImageStore.FindAll(0)
	if err != nil {
		log.Fatal(err)
	}
	RenderTemplate(w, r, "index/home", map[string]interface{}{
		"Images": imgs,
	})
}
