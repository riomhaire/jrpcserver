package model

import "net/http"

type Middleware interface {
	Handle(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc)
}

type APIConfig struct {
	ServiceName string                                                                 `json:"service,omitempty"`
	BaseURI     string                                                                 `json:"baseUri,omitempty"`
	Port        int                                                                    `json:"port,omitempty"`
	Commands    []JRPCCommand                                                          `json:"-"` // The RPC Commands
	Version     string                                                                 `json:"version,omitempty"`
	Hostname    string                                                                 `json:"host,omitempty"`   // Which host service is bound to - if blank defaults to os.Hostname(), used for consul connection
	Consul      string                                                                 `json:"consul,omitempty"` // where consul host is located. If blank no consul integration made: its host and port
	Middleware  []func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) `json:"-"`                // A list of handlers to use ... logger security etc
}
