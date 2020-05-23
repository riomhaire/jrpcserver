package model

import "net/http"

type Middleware interface {
	Handle(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc)
}

type ServerConfig struct {
	ServiceName string                                                                 `json:"service,omitempty" yaml:"service,omitempty"`
	BaseURI     string                                                                 `json:"baseUri,omitempty" yaml:"baseUri,omitempty"`
	Port        int                                                                    `json:"port,omitempty" yaml:"port,omitempty"`
	Commands    []JRPCCommand                                                          `json:"-" yaml:"-"` // The RPC Commands
	Version     string                                                                 `json:"version,omitempty" yaml:"version,omitempty"`
	Hostname    string                                                                 `json:"host,omitempty" yaml:"host,omitempty"`     // Which host service is bound to - if blank defaults to os.Hostname(), used for consul connection
	Consul      string                                                                 `json:"consul,omitempty" yaml:"consul,omitempty"` // where consul host is located. If blank no consul integration made: its host and port
	Middleware  []func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) `json:"-" yaml:"-"`                               // A list of handlers to use ... logger security etc
}

/**
 Instances allow core access config irrespective of whats in the config
**/
type ServerConfigReader interface {
	ReadServerConfig() (*ServerConfig, error)
}

/**
This is the simplest example of a config
**/
type DefaultConfiguration struct {
	Server ServerConfig `json:"server,omitempty" yaml:"server,omitempty"`
}

/**
Return Server Config
**/
func (config DefaultConfiguration) ReadServerConfig() (*ServerConfig, error) {
	return &config.Server, nil
}
