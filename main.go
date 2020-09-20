package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/chixm/servertemplate2/config"
	"github.com/chixm/servertemplate2/database"
	"github.com/chixm/servertemplate2/redis"
	"github.com/chixm/servertemplate2/utils"
	"github.com/chixm/servertemplate2/websocket"
	"github.com/gorilla/mux"
)

// instance of the server
var server *http.Server

const useLogFile = false // If this param is true, write out log to /log/application.log file instead of writing to stdout.

func main() {
	initialize()

	launchServer(createServerEndPoints())

	defer terminate()
}

// initialize systems around this server.
func initialize() {
	setupLog(useLogFile)

	config.InitializeConfig()

	utils.InitializeUniqueIDMaker()

	database.InitializeDatabaseConnections(logger)

	websocket.InitializeWebSocket(logger)

	redis.InitializeRedis(logger)

	utils.InitializeWebdriver(logger)

	initializeEmailSender()

	initializeBatch()

	initializeServiceFunctions()
}

func terminate() {
	database.TerminateDatabaseConnections()

	terminateBatch()

	redis.TerminateRedis()

	terminateLog()
}

// settings of endpoints
func createServerEndPoints() *mux.Router {
	r := mux.NewRouter()
	serviceMap, err := LoadServices()
	if err != nil {
		panic(`Coud not define URL of Server service.` + err.Error())
	}

	// load all handlers listed in service.webHandlers.
	for url, handler := range *serviceMap {
		r.HandleFunc(url, handler)
	}

	// load default web socket
	r.HandleFunc("/ws", websocket.Ws)

	// Files under /static can accessed by /static/(filename)...
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(`resources/static`))))
	return r
}

// launch server
func launchServer(r http.Handler) error {
	server := &http.Server{
		Handler:      r,
		Addr:         "localhost:" + strconv.Itoa(config.GetConfig().Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Info("Execute Server ::" + server.Addr)
	return server.ListenAndServe()
}

/** give initialized functions and connections to service. */
func initializeServiceFunctions() {

}
