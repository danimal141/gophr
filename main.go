package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func NewRouter() *httprouter.Router {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	return router
}

func main() {
	router := NewRouter()
	router.GET("/", HandleHome)
	router.GET("/register", HandleUserNew)
	router.POST("/register", HandleUserCreate)
	router.ServeFiles(
		"/assets/*filepath",
		http.Dir("assets/"),
	)

	m := Middleware{}
	m.Add(router)

	http.ListenAndServe("localhost:3000", m)
}
