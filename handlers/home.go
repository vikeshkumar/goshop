package handlers

import (
	"net/http"
)

var home = "index"

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	homeTemplate, err := parse(home)
	if err != nil {
		logger.Printf("error parsing %v", err)
	}
	execError := homeTemplate.Execute(w, &page{"vikesh.net", "hello"})
	if execError != nil {
		logger.Printf("error executing template %v", execError)
	}
}
