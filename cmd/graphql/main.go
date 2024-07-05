package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/ardihikaru/graphql-example-part-1/internal/application"
	"github.com/ardihikaru/graphql-example-part-1/internal/router"

	"github.com/ardihikaru/graphql-example-part-1/pkg/config"
	"github.com/ardihikaru/graphql-example-part-1/pkg/logger"
	e "github.com/ardihikaru/graphql-example-part-1/pkg/utils/error"
)

// Version sets the default build version
var Version = "development"

func main() {
	// loads configuration
	cfg, err := config.Get()
	if err != nil {
		e.FatalOnError(err, "failed to load config")
	}

	// builds private key object
	err = cfg.BuildEncryptionKeys()
	if err != nil {
		e.FatalOnError(err, "failed to build private key object")
	}

	// validates config
	err = cfg.Validate()
	if err != nil {
		e.FatalOnError(err, "config validation occurs")
	}

	// configures logger
	log, err := logger.New(cfg.Log.Level, cfg.Log.Format, cfg.Http.LogHttpRequest, cfg.LogPublisher)
	if err != nil {
		e.FatalOnError(err, "failed to prepare the logger")
	}

	// builds dependencies
	deps := application.BuildDependencies(cfg, log)

	// shows the build version
	name, err := os.Hostname()
	log.Info("starting API service",
		zap.String("Service ID", deps.SvcId),
		zap.String("Hostname", name),
		zap.String("Build Version", Version),
		zap.String("Build Mode", cfg.General.BuildMode),
	)

	// gracefully exit on keyboard interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// starts the api service
	startAPiService(deps)

	// shutdowns the application
	<-c
	log.Info("gracefully shutting down the system")
	os.Exit(0)

	//// creates a new router
	//router := chi.NewRouter()
	//
	//// -- route - config
	//c := cors.New(cors.Options{
	//	AllowedOrigins:   []string{"*"},
	//	AllowedHeaders:   []string{"*"},
	//	AllowCredentials: true,
	//	Debug:            false,
	//})
	//router.Use(c.Handler)

	//// mounts Routes
	//routes.MountRoutes(router)
	//
	//// init DB
	//database.InitDB(log, cfg)
	//defer database.Handle.Close()
	//defer database.HandleSeq.Close()
	//
	//// starts the server
	//log.Info("Connect to http://localhost:" + port + "/ for GraphQL playground")
	//err = http.ListenAndServe(":"+port, router)
	//if err != nil {
	//	e.FatalOnError(err, "failed to start server")
	//}
}

// startAPiService starts the api service
func startAPiService(deps *application.Dependencies) {
	r := router.GetRouter(deps)

	// logs that application is ready
	deps.Log.Info("preparing to serve the request in => " + fmt.Sprintf(":%v", deps.Cfg.Http.Port))

	// builds server params
	address := fmt.Sprintf("%s:%v", deps.Cfg.Http.Address, deps.Cfg.Http.Port)
	server := http.Server{
		ReadTimeout:  deps.Cfg.Http.ReadTimeout,
		WriteTimeout: deps.Cfg.Http.WriteTimeout,
		Handler:      r,
		Addr:         address,
	}

	go func() {
		// stops the application if any error found
		if err := server.ListenAndServe(); err != nil {
			e.FatalOnError(err, "failed to start server")
			os.Exit(1)
		}
	}()
}
