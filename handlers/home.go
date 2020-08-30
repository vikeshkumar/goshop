package handlers

import (
	"html/template"
	"net/http"
)

var homePageTemplate *template.Template
var home = "templates/index.gohtml"

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	parseAndServe(w, home, homePageTemplate, nil)
}
