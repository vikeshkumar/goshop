package handlers

import (
	"github.com/prometheus/common/log"
	"go.uber.org/zap"
	"html/template"
	"net/http"
)

type TemplateConfiguration struct {
	CacheTemplates bool
}

var templateConfig *TemplateConfiguration

func CreateHandlers(config *TemplateConfiguration) map[string]http.HandlerFunc {
	templateConfig = config
	handlerFunctions := make(map[string]http.HandlerFunc)
	handlerFunctions["/index"] = HomePageHandler
	handlerFunctions["/"] = HomePageHandler
	handlerFunctions["/index.html"] = HomePageHandler
	return handlerFunctions
}

func parseAndServe(w http.ResponseWriter, path string, tpl *template.Template, data *interface{}) {
	var err error
	if !templateConfig.CacheTemplates {
		tpl, err = template.ParseFiles(path)
	} else if tpl == nil {
		tpl, err = template.ParseFiles(path)
	}

	if err != nil {
		log.Error("error parsing template", zap.String("path", path), zap.Error(err))
	}
	exeErr := tpl.Execute(w, data)
	if exeErr != nil {
		log.Error("error executing template", zap.String("path", path), zap.Error(err))
	}
}
