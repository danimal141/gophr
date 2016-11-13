package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func HandleImageShow(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	img, err := globalImageStore.Find(params.ByName("imageID"))
	if err != nil {
		panic(err)
	}
	if img == nil {
		http.NotFound(w, r)
		return
	}

	user, err := globalUserStore.Find(img.UserID)
	if err != nil {
		panic(err)
	}
	if user == nil {
		panic(fmt.Errorf("Could not find user %s", img.UserID))
	}

	RenderTemplate(w, r, "images/show", map[string]interface{}{
		"Image": img,
		"User":  user,
	})
}

func HandleImageNew(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	RenderTemplate(w, r, "images/new", nil)
}

func HandleImageCreate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.FormValue("url") != "" {
		HandleImageCreateFromURL(w, r)
		return
	}
	HandleImageCreateFromFile(w, r)
}

func HandleImageCreateFromURL(w http.ResponseWriter, r *http.Request) {
	user := RequestUser(r)
	img := NewImage(user)
	img.Description = r.FormValue("description")

	err := img.CreateFromURL(r.FormValue("url"))
	if err != nil {
		if IsValidationError(err) {
			RenderTemplate(w, r, "images/new", map[string]interface{}{
				"Error":    err,
				"ImageURL": r.FormValue("url"),
				"Image":    img,
			})
			return
		}
		panic(err)
	}
	http.Redirect(w, r, "/?flash=Image+Uploaded+Successfully", http.StatusFound)
}

func HandleImageCreateFromFile(w http.ResponseWriter, r *http.Request) {
	user := RequestUser(r)
	img := NewImage(user)
	img.Description = r.FormValue("description")
	file, headers, err := r.FormFile("file")

	if file == nil {
		RenderTemplate(w, r, "images/new", map[string]interface{}{
			"Error": errNoImage,
			"Image": img,
		})
		return
	}
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = img.CreateFromFile(file, headers)
	if err != nil {
		RenderTemplate(w, r, "images/new", map[string]interface{}{
			"Error": err,
			"Image": img,
		})
		return
	}
	http.Redirect(w, r, "/?flash=Image+Uploaded+Successfully", http.StatusFound)
}
