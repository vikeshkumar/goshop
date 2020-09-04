package handlers

import (
	"go.uber.org/zap"
	"net/http"
)

var home = tc.TemplateDirectory + "/index.gohtml"

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	homeTemplate, err := parse(home)
	if err != nil {
		log.Error("error parsing", zap.Any("error", err))
	}
	execError := homeTemplate.Execute(w, &page{"vikesh.net", "hello"})
	if execError != nil {
		log.Error("error executing template")
	}
}
