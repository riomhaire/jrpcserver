package rpc

import (
	"encoding/json"
	"fmt"
	"log/syslog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gorilla/mux"
	negronilogrus "github.com/meatballhat/negroni-logrus"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	lSyslog "github.com/sirupsen/logrus/hooks/syslog"
	"github.com/urfave/negroni"
	negroniprometheus "github.com/zbindenren/negroni-prometheus"

	"github.com/riomhaire/jrpcserver/infrastructure/serviceregistry"
	"github.com/riomhaire/jrpcserver/infrastructure/serviceregistry/consulagent"
	"github.com/riomhaire/jrpcserver/infrastructure/serviceregistry/none"
)

var dispatcher *Dispatcher

func StartAPI(config interface{}) {
	// Create dispatcher for later use
	dispatcher = NewDispatcher(config)
	consulEnabled := len(dispatcher.serverConfig.Consul) != 0
	if len(dispatcher.serverConfig.Hostname) == 0 {
		n, _ := os.Hostname()
		dispatcher.serverConfig.Hostname = n
	}

	// Set up registry
	var registryConnector serviceregistry.ServiceRegistry // Default to none
	if consulEnabled {
		registryConnector = consulagent.NewConsulServiceRegistry(dispatcher.serverConfig.Consul, dispatcher.serverConfig.ServiceName, dispatcher.serverConfig.Hostname, dispatcher.serverConfig.Port, dispatcher.serverConfig.BaseURI, "/health")
	} else {
		registryConnector = none.NewDefaultServiceRegistry() // Default to none

	}

	// Define endpoint
	router := mux.NewRouter()

	// add middleware for a specific route and get params from route
	router.HandleFunc(fmt.Sprintf("%s/{method}", dispatcher.serverConfig.BaseURI), rpcHandler)
	router.HandleFunc("/health", healthHandler)

	// Includes some default middlewares to all routes
	negroniServer := negroni.New()
	negroniServer.Use(negroni.NewRecovery())

	// add log
	hook, err := lSyslog.NewSyslogHook("", "", syslog.LOG_INFO, "")
	log.StandardLogger().SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	//	log.StandardLogger().
	if err == nil {
		log.StandardLogger().Hooks.Add(hook)
	}

	negroniServer.Use(negronilogrus.NewMiddlewareFromLogger(log.StandardLogger(), dispatcher.serverConfig.ServiceName))

	// Add some useful handlers Add some headers
	negroniServer.UseFunc(AddWorkerHeader)  // Add which instance
	negroniServer.UseFunc(AddWorkerVersion) // Which version
	// If there are any handlers in the config add them
	for _, handler := range dispatcher.serverConfig.Middleware {
		negroniServer.UseFunc(handler)
	}

	// Coors stuff
	handler := cors.New(
		cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"HEAD", "GET", "POST", "PUT", "PATCH", "DELETE"},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: true,
		}).Handler(router) // Add coors
	negroniServer.UseHandler(handler)

	// Add Server Metrics
	negroniServer.Use(negroniprometheus.NewMiddleware(dispatcher.serverConfig.ServiceName))

	log.Println("Starting JSON RPC Server Version", dispatcher.serverConfig.Version, dispatcher.serverConfig.BaseURI, "on port:", dispatcher.serverConfig.Port)

	// Set up shutdown resister
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Set up a way to cleanly shutdown / deregister
	go func() {
		<-c
		registryConnector.Deregister()
		log.Println("Shutting Down")
		os.Exit(0)
	}()

	// Register (if required with consul or other registry)
	registryConnector.Register()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", dispatcher.serverConfig.Port), negroniServer))

}

func rpcHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	method := vars["method"]
	defer r.Body.Close()

	// Get Headers
	headers := make(map[string]string)
	// Loop through headers - we want the last
	for name, values := range r.Header {
		name = strings.ToLower(name)
		for _, h := range values {
			headers[name] = h
		}
	}

	// Make the call
	dipatcherResponse := dispatcher.Execute(method, headers, r.Body)
	var b []byte

	if dipatcherResponse.RawResponse { // non command response
		b, _ = json.MarshalIndent(dipatcherResponse.Value, "", "  ")

	} else { // Encode a command response
		b, _ = json.MarshalIndent(dipatcherResponse, "", "  ")

	}
	// TODO Append schema type
	w.Header().Set("Content-Type", "application/json")
	if dipatcherResponse.Code == 0 {
		w.WriteHeader(http.StatusOK)
		w.Write(b)

	} else {
		w.WriteHeader(dipatcherResponse.Code)
		w.Write(b)

	}
}

func healthHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.Write([]byte("{ \"status\":\"up\"}"))

}

// AddWorkerHeader - adds header of which node actually processed request
func AddWorkerHeader(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	host, err := os.Hostname()
	if err != nil {
		host = "Unknown"
	}
	rw.Header().Add("X-Worker", host)
	if next != nil {
		next(rw, req)
	}
}

// AddWorkerVersion - adds header of which version is installed
func AddWorkerVersion(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if len(dispatcher.serverConfig.Version) == 0 {
		dispatcher.serverConfig.Version = "UNKNOWN"
	}
	rw.Header().Add("X-Worker-Version", dispatcher.serverConfig.Version)
	if next != nil {
		next(rw, req)
	}
}
