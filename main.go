package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/sudarshan-reddy/benjerry/configs"
	"github.com/sudarshan-reddy/benjerry/db"
	"github.com/sudarshan-reddy/benjerry/models/postgres"
	"github.com/sudarshan-reddy/benjerry/router"
	"github.com/sudarshan-reddy/benjerry/scripts"
)

const (
	serviceName = "ben-jerry-app"
)

var (
	commitID       = "not assigned"
	buildTimestamp = "not assigned"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s : %s", msg, err)
	}
}

func main() {
	config, err := configs.Load()
	failOnError(err, "error while loading config")
	setupLog(config.LogLevel, config.LogFormat)
	log.Infof("%s built on %s from commit %s", serviceName, buildTimestamp, commitID)

	db.RunMigrateScripts(config.PostgresDBURL, ".")
	postgresDB, err := db.NewPostgresDB(config.PostgresDBURL, config.PostgresDBMaxConnections)
	failOnError(err, "error while connecting to postgresDB")

	iceCreamStore := postgres.NewIceCreamStore(postgresDB)

	if config.LoadData {
		err := scripts.MoveData(iceCreamStore)
		failOnError(err, "error while moving ice cream initial data")
	}

	routerCfg := router.Config{
		IceCreamStore: iceCreamStore,
	}

	apiRouter := router.NewRouter(config.StaticTokens, routerCfg)
	apiRouter.AddRoutes()

	log.Infof("%s running on port %s", serviceName, config.ListenPort)
	http.ListenAndServe(":"+config.ListenPort, apiRouter)
}

func setupLog(logLevel, logFormat string) {
	setLogLevel(logLevel)
	setLogFormat(logFormat)
}

func setLogFormat(logFormat string) {
	switch logFormat {
	case "json":
		log.SetFormatter(&log.JSONFormatter{
			FieldMap: log.FieldMap{log.FieldKeyMsg: "message"},
		})
	default:
		log.SetFormatter(&log.TextFormatter{})
	}
}

func setLogLevel(logLevel string) {
	switch logLevel {
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}
