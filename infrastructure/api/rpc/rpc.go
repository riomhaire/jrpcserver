package rpc

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/riomhaire/jrpcserver/infrastructure/serviceregistry"
	"github.com/riomhaire/jrpcserver/infrastructure/serviceregistry/consulagent"
	"github.com/riomhaire/jrpcserver/infrastructure/serviceregistry/none"
	"github.com/riomhaire/jrpcserver/model"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
	negroniprometheus "github.com/zbindenren/negroni-prometheus"
)

var dispatcher *Dispatcher
var apiconfig *APIConfig

type APIConfig struct {
	ServiceName string
	BaseURI     string
	Port        int
	Commands    []model.JRPCCommand
	Version     func() string
	Hostname    string // Which host service is bound to - if blank defaults to os.Hostname(), used for consul connection
	Consul      string // where consul host is located. If blank no consul integration made: its host and port
}

func StartAPI(config APIConfig) {
	// Create dispatcher for later use
	dispatcher = NewDispatcher(config.Commands)
	consulEnabled := len(config.Consul) != 0
	if len(config.Hostname) == 0 {
		n, _ := os.Hostname()
		config.Hostname = n
	}
	apiconfig = &config

	// Set up registry
	var registryConnector serviceregistry.ServiceRegistry // Default to none
	if consulEnabled {
		registryConnector = consulagent.NewConsulServiceRegistry(config.Consul, config.ServiceName, config.Hostname, config.Port, config.BaseURI, "/health")
	} else {
		registryConnector = none.NewDefaultServiceRegistry() // Default to none

	}

	// Define endpoint
	router := mux.NewRouter()
	// add middleware for a specific route and get params from route
	router.HandleFunc(fmt.Sprintf("%s/{method}", config.BaseURI), rpcHandler)
	router.Handle("/metrics", prometheus.Handler())
	router.HandleFunc("/health", healthHandler)

	// Includes some default middlewares to all routes
	negroni := negroni.Classic()
	// Add some headers
	negroni.UseFunc(AddWorkerHeader)  // Add which instance
	negroni.UseFunc(AddWorkerVersion) // Which version
	// Coors stuff
	handler := cors.New(
		cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"HEAD", "GET", "POST", "PUT", "PATCH", "DELETE"},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: true,
		}).Handler(router) // Add coors
	negroni.UseHandler(handler)

	// Add Server Metrics
	negroni.Use(negroniprometheus.NewMiddleware(config.ServiceName))

	log.Println("Starting JSON RPC Server Version", config.Version(), config.BaseURI, "on port:", config.Port)

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
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", config.Port), negroni))

}

func rpcHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	method := vars["method"]

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

	b, _ := json.MarshalIndent(dipatcherResponse, "", "  ")
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
	version := apiconfig.Version()
	if len(version) == 0 {
		version = "UNKNOWN"
	}
	rw.Header().Add("X-Worker-Version", version)
	if next != nil {
		next(rw, req)
	}
}
