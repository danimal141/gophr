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
	router.GET("/users/new", HandleUserNew)
	router.POST("/users", HandleUserCreate)
	router.GET("/login", HandleSessionNew)
	router.POST("/login", HandleSessionCreate)
	router.ServeFiles(
		"/assets/*filepath",
		http.Dir("assets/"),
	)

	secureRouter := NewRouter()
	secureRouter.GET("/users/edit", HandleUserEdit)
	secureRouter.PATCH("/users", HandleUserUpdate)
	secureRouter.GET("/logout", HandleSessionDestroy)

	m := Middleware{}
	m.Add(router)
	m.Add(http.HandlerFunc(RequireLogin))
	m.Add(secureRouter)

	http.ListenAndServe("localhost:3000", m)
}
