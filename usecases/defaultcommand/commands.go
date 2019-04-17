package defaultcommand

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/riomhaire/jrpcserver/model"
	"github.com/riomhaire/jrpcserver/model/jrpcerror"
)

func PingCommand(config model.APIConfig, metadata map[string]string, payload io.ReadCloser) (interface{}, jrpcerror.JrpcError) {
	return "pong", jrpcerror.JrpcError{}
}

func PongCommand(config model.APIConfig, metadata map[string]string, payload io.ReadCloser) (interface{}, jrpcerror.JrpcError) {
	return "ping", jrpcerror.JrpcError{}
}

func EchoCommand(config model.APIConfig, metadata map[string]string, payload io.ReadCloser) (interface{}, jrpcerror.JrpcError) {
	for name, value := range metadata {
		fmt.Printf("%s = %s\n", name, value)
	}

	data, _ := ioutil.ReadAll(payload)
	msg := make(map[string]interface{})

	json.Unmarshal(data, &msg)

	return msg, jrpcerror.JrpcError{}
}

func ListCommandsCommand(config model.APIConfig, metadata map[string]string, payload io.ReadCloser) (interface{}, jrpcerror.JrpcError) {
	names := make([]string, 0)
	for _, cmd := range config.Commands {
		names = append(names, cmd.Name)
	}
	return names, jrpcerror.JrpcError{}
}

func VersionCommand(config model.APIConfig, metadata map[string]string, payload io.ReadCloser) (interface{}, jrpcerror.JrpcError) {
	return config.Version, jrpcerror.JrpcError{}
}

func InfoCommand(config model.APIConfig, metadata map[string]string, payload io.ReadCloser) (interface{}, jrpcerror.JrpcError) {
	return config, jrpcerror.JrpcError{}
}
