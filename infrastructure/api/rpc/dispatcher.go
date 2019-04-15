package rpc

import (
	"io"

	"github.com/riomhaire/jrpcserver/model"
)

type Dispatcher struct {
	config model.APIConfig
}

func NewDispatcher(config model.APIConfig) *Dispatcher {
	return &Dispatcher{config}

}

func (d *Dispatcher) Execute(method string, headers map[string]string, payload io.ReadCloser) *model.RPCCommandResponse {
	response := model.RPCCommandResponse{}
	found := false
	for _, cmd := range d.config.Commands {
		if cmd.Name == method {
			found = true
			result, err := cmd.Command(d.config, headers, payload)
			if err.Code != 0 {
				response.Code = err.Code
				response.Error = err.Error
			} else {
				response.Code = 0
				response.Value = result
			}
		}
	}
	if !found {
		response.Code = 404
		response.Error = "command not found"
	}
	return &response
}
