package main

import (
	"context"
	"flag"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	c "net.vikesh/goshop/config"
	"net.vikesh/goshop/db"
	"net.vikesh/goshop/handlers"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"time"
)

var logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

func main() {
	config, err := bootstrap(string(c.DefaultConfigurationPath))
	if err != nil {
		logger.Printf("error reading required configuration file %v", err)
		os.Exit(1)
	}
	//configuring URL handlers
	var dir string
	td := config.TemplateConfiguration.TemplateDirectory
	flag.StringVar(&dir, "dir", td, "the directory to serve files from - defaults to the current dir")
	r := mux.NewRouter()
	r.PathPrefix("/_ui/").Handler(http.StripPrefix("/_ui/", http.FileServer(http.Dir(dir))))
	dynamicViews := handlers.CreateHandlers(config.TemplateConfiguration)
	for pattern, handlerFunc := range dynamicViews {
		logger.Printf("creating handler with path = %v, function = %v",
			pattern,
			runtime.FuncForPC(reflect.ValueOf(handlerFunc).Pointer()).Name())
		r.HandleFunc(pattern, handlerFunc)
	}
	r.Schemes("http")

	//configure a logging middleware
	r.Use(loggingMiddleware)
	srv := &http.Server{
		Handler: r,
		Addr:    config.ServerConfiguration.ListenOnPort,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: config.ServerConfiguration.WriteTimeout * time.Second,
		ReadTimeout:  config.ServerConfiguration.ReadTimeout * time.Second,
	}
	// accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	go listenForShutDown(ch, srv, time.Until(time.Now().Add(time.Minute*1)))
	p, err := db.Connect(context.Background(), config.DBConfig.URL)
	if err != nil {
		log.Fatal(
			"failed to connect to database",
			err,
		)
	}
	if p != nil {
		defer p.Close()
	}
	if err := srv.ListenAndServe(); err != nil {
		logger.Printf("err = %v", err)
	}
}

// creates a logging middleware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		logger.Printf("request url = %v, method = %v, remote = %v, protocol =%v",
			r.RequestURI,
			r.Method,
			r.RemoteAddr,
			r.Proto,
		)
		next.ServeHTTP(w, r)
	})
}

// listen for application shutdown and configure timeout to respond to request
// and then terminate
func listenForShutDown(ch chan os.Signal, srv *http.Server, wait time.Duration) {
	<-ch
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	logger.Printf("shutting down, received signal = %v", ctx)
	_ = srv.Shutdown(ctx)
	os.Exit(0)
}

// bootstraps the application - only does configuration reading
func bootstrap(configPath string) (c.Configuration, error) {
	b, err := ioutil.ReadFile(configPath)
	config := c.Configuration{}
	if err != nil {
		log.Println("failed to read configuration file, default is config/local.yaml relative to the executable",
			err,
		)
		return config, err
	}
	logger.Printf("opened config file for reading, configPath = %v",
		configPath,
	)
	marshalError := yaml.Unmarshal(b, &config)
	if marshalError != nil {
		log.Println("marshalling error",
			marshalError,
		)
		return config, marshalError
	}
	logger.Printf(
		"development set to : %v",
		config.ServerConfiguration.Development,
	)
	return config, nil
}
