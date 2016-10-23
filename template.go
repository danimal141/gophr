package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
)

var layoutFuncs = template.FuncMap{
	"yield": func() (string, error) {
		return "", fmt.Errorf("yield called inappropriately")
	},
}
var layout = template.Must(
	template.New("layout").Funcs(layoutFuncs).ParseFiles("templates/layout.tmpl"),
)

var templates = template.Must(template.New("t").ParseGlob("templates/**/*.tmpl"))

func RenderTemplate(w http.ResponseWriter, r *http.Request, name string, data interface{}) {
	funcs := template.FuncMap{
		"yield": func() (template.HTML, error) {
			buf := bytes.NewBuffer(nil)
			err := templates.ExecuteTemplate(buf, name, data) // Write to buffer instead of ResponseWriter
			return template.HTML(buf.String()), err
		},
	}
	layoutClone, err := layout.Clone()
	if err != nil {
		fmt.Sprintf("Error: ", err)
	}
	layoutClone.Funcs(funcs)

	err = layoutClone.Execute(w, data)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Error: ", err),
			http.StatusInternalServerError,
		)
	}
}
