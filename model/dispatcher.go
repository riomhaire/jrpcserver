package model

type APIConfig struct {
	ServiceName string
	BaseURI     string
	Port        int
	Commands    []JRPCCommand
	Version     string
	Hostname    string // Which host service is bound to - if blank defaults to os.Hostname(), used for consul connection
	Consul      string // where consul host is located. If blank no consul integration made: its host and port
}
