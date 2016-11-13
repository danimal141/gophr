package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

func NewRouter() *httprouter.Router {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	return router
}

func init() {
	// Assign a user store
	userStore, err := NewFileUserStore("./data/users.json")
	if err != nil {
		panic(fmt.Errorf("Error creating user store: %s", err))
	}
	globalUserStore = userStore

	// Assign a session store
	sessionStore, err := NewFileSessionStore("./data/sessions.json")
	if err != nil {
		panic(fmt.Errorf("Error creating session store: %s", err))
	}
	globalSessionStore = sessionStore

	// Assign a sql database
	db, err := NewMySQLDB(os.Getenv("MYSQL_DATA_SOURCE_NAME"))
	if err != nil {
		panic(err)
	}
	globalMySQLDB = db

	// Assign an image store
	globalImageStore = NewDBImageStore()
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
	secureRouter.GET("/images/new", HandleImageNew)
	secureRouter.POST("/images/new", HandleImageCreate)

	m := Middleware{}
	m.Add(router)
	m.Add(http.HandlerFunc(RequireLogin))
	m.Add(secureRouter)

	http.ListenAndServe("localhost:3000", m)
}
