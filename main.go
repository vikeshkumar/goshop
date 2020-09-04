package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v2"
	"io/ioutil"
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

var log = zap.NewExample()

func main() {
	rawJSON := []byte(`{
		"level": "debug",
		"encoding": "json",
		"outputPaths": ["stdout", "./logs"],
		"errorOutputPaths": ["stderr", "./logs"],
		"callerKey": "caller",
		"encoderConfig": {
		  "messageKey": "message",
		  "levelKey": "level",
		  "levelEncoder": "lowercase"
		}
	  }`)
	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}
	cfg.EncoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	log, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer log.Sync()
	config, err := bootstrap(string(c.DefaultConfigurationPath))
	if err != nil {
		log.Error("error reading required configuration file", zap.Error(err))
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
		log.Info("creating handler with path",
			zap.String("pattern", pattern),
			zap.String("function", runtime.FuncForPC(reflect.ValueOf(handlerFunc).Pointer()).Name()))
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
			zap.Any("error", err),
		)
		os.Exit(1)
	}
	if p != nil {
		defer p.Close()
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Error("err", zap.Error(err))
	}
}

// creates a logging middleware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Info(r.RequestURI,
			zap.Any("method", r.Method),
			zap.Any("remote", r.RemoteAddr),
			zap.Any("proto", r.Proto),
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
	log.Info("shutting down", zap.Any("signal", ctx))
	_ = srv.Shutdown(ctx)
	os.Exit(0)
}

// bootstraps the application - only does configuration reading
func bootstrap(configPath string) (c.Configuration, error) {
	b, err := ioutil.ReadFile(configPath)
	config := c.Configuration{}
	if err != nil {
		log.Error("failed to read configuration file, default is config/local.yaml relative to the executable",
			zap.Error(err),
		)
		return config, err
	}
	log.Debug("opened config file for reading",
		zap.String("filename", configPath),
	)
	marshalError := yaml.Unmarshal(b, &config)
	if marshalError != nil {
		log.Error("marshalling error",
			zap.Error(marshalError),
		)
		return config, marshalError
	}
	log.Debug(
		"development set to : ",
		zap.Bool("development", config.ServerConfiguration.Development),
	)
	return config, nil
}
