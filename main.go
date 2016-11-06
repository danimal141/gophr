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
	router.ServeFiles(
		"/assets/*filepath",
		http.Dir("assets/"),
	)

	router.GET("/", HandleHome)
	router.GET("/register", HandleUserNew)
	router.POST("/register", HandleUserCreate)
	router.GET("/login", HandleSessionNew)
	router.POST("/login", HandleSessionCreate)

	secureRouter := NewRouter()
	secureRouter.GET("/account", HandleUserEdit)
	secureRouter.POST("/account", HandleUserUpdate)
	secureRouter.GET("/logout", HandleSessionDestroy)

	m := Middleware{}
	m.Add(router)
	m.Add(http.HandlerFunc(RequireLogin))
	m.Add(secureRouter)

	http.ListenAndServe("localhost:3000", m)
}
