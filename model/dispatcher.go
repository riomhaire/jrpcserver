package model

type APIConfig struct {
	ServiceName string        `json:"service,omitempty"`
	BaseURI     string        `json:"baseUri,omitempty"`
	Port        int           `json:"port,omitempty"`
	Commands    []JRPCCommand `json:"-"`
	Version     string        `json:"version,omitempty"`
	Hostname    string        `json:"host,omitempty"`   // Which host service is bound to - if blank defaults to os.Hostname(), used for consul connection
	Consul      string        `json:"consul,omitempty"` // where consul host is located. If blank no consul integration made: its host and port
}
