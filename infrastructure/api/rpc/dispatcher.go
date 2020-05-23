package rpc

import (
	"io"

	"github.com/riomhaire/jrpcserver/model"
)

type Dispatcher struct {
	config       interface{}
	serverConfig *model.ServerConfig
}

func NewDispatcher(config interface{}) *Dispatcher {
	var dispatcher Dispatcher
	serverConfigAccessor, ok := config.(model.ServerConfigReader)

	if ok {
		server, _ := serverConfigAccessor.ReadServerConfig()
		dispatcher.serverConfig = server
	}
	dispatcher.config = config
	return &dispatcher

}

func (d *Dispatcher) Execute(method string, headers map[string]string, payload io.ReadCloser) *model.RPCCommandResponse {
	response := model.RPCCommandResponse{}
	found := false
	for _, cmd := range d.serverConfig.Commands {
		if cmd.Name == method {
			found = true
			result, err := cmd.Command(d.config, headers, payload)
			if err.Code != 0 {
				response.Code = err.Code
				response.Error = err.Error
				response.RawResponse = cmd.RawResponse
			} else {
				response.Code = 0
				response.Value = result
				response.RawResponse = cmd.RawResponse
			}
		}
	}
	if !found {
		response.Code = 404
		response.Error = "command not found"
	}
	return &response
}
