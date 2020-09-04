package handlers

import (
	"go.uber.org/zap"
	"html/template"
	"net.vikesh/goshop/config"
	"net/http"
	"os"
	"sync"
)

var log = zap.NewExample()

var templates = make(map[string]*template.Template)
var templateLock sync.Mutex
var tc config.TemplateConfiguration

// CreateHandlers creates all the handlers in the application
func CreateHandlers(c config.TemplateConfiguration) map[string]http.HandlerFunc {
	templateLock.Lock()
	defer templateLock.Unlock()
	tc = c
	handlers := make(map[string]http.HandlerFunc)
	handlers["/index"] = homePageHandler
	handlers["/"] = homePageHandler
	handlers["/index.html"] = homePageHandler
	handlers["/index.html"] = homePageHandler
	handlers["/favicon.ico"] = faviconHandler
	return handlers
}

// handler to write favicon file
func faviconHandler(w http.ResponseWriter, r *http.Request) {
	s, se := os.Stat(tc.Favicon)
	if se != nil {
		log.Error("error reading stat of file", zap.Error(se))
		w.WriteHeader(404)
		return
	}
	f, e := os.Open(tc.Favicon)
	if e != nil {
		log.Error("error reading file", zap.Error(e))
		w.WriteHeader(404)
		return
	}

	http.ServeContent(w, r, r.URL.Path, s.ModTime(), f)
}

type page struct {
	Website string
	Title   string
}

func parse(path string) (*template.Template, error) {
	re := tc.ParseOnce
	if !re {
		tpl, err := template.ParseGlob(path)
		return tpl, err
	}
	if _, ok := templates[path]; !ok {
		tpl, err := template.ParseGlob(path)
		if err != nil {
			return nil, err
		} else {
			templateLock.Lock()
			defer templateLock.Unlock()
			templates[path] = tpl
		}
	}
	return templates[path], nil
}
