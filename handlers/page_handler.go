package handlers

import (
	"html/template"
	"log"
	"net.vikesh/goshop/config"
	"net/http"
	"os"
	"sync"
)

var logger = log.New(os.Stdout, "", log.Llongfile)

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
		log.Println("error reading stat of file", se)
		w.WriteHeader(404)
		return
	}
	f, e := os.Open(tc.Favicon)
	if e != nil {
		log.Println("error reading file", e)
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
		tpl, err := template.ParseGlob(join(path))
		return tpl, err
	}
	if _, ok := templates[path]; !ok {
		tpl, err := template.ParseGlob(join(path))
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

func join(path string) string {
	return tc.TemplateDirectory + path + tc.Suffix
}
