package defaultcommand

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/riomhaire/jrpcserver/model"
	"github.com/riomhaire/jrpcserver/model/jrpcerror"
)

func PingCommand(config interface{}, metadata map[string]string, payload io.ReadCloser) (interface{}, jrpcerror.JrpcError) {
	return "pong", jrpcerror.JrpcError{}
}

func PongCommand(config interface{}, metadata map[string]string, payload io.ReadCloser) (interface{}, jrpcerror.JrpcError) {
	return "ping", jrpcerror.JrpcError{}
}

func EchoCommand(config interface{}, metadata map[string]string, payload io.ReadCloser) (interface{}, jrpcerror.JrpcError) {
	for name, value := range metadata {
		fmt.Printf("%s = %s\n", name, value)
	}

	data, _ := ioutil.ReadAll(payload)
	msg := make(map[string]interface{})

	json.Unmarshal(data, &msg)

	return msg, jrpcerror.JrpcError{}
}

func ListCommandsCommand(config interface{}, metadata map[string]string, payload io.ReadCloser) (interface{}, jrpcerror.JrpcError) {
	names := make([]string, 0)
	serverConfigAccessor, ok := config.(model.ServerConfigReader)
	if !ok {
		return "unknown", jrpcerror.JrpcError{}
	}
	server, _ := serverConfigAccessor.ReadServerConfig()

	for _, cmd := range server.Commands {
		names = append(names, cmd.Name)
	}
	return names, jrpcerror.JrpcError{}
}

func VersionCommand(config interface{}, metadata map[string]string, payload io.ReadCloser) (interface{}, jrpcerror.JrpcError) {
	serverConfigAccessor, ok := config.(model.ServerConfigReader)
	if !ok {
		return "unknown", jrpcerror.JrpcError{}
	}
	server, _ := serverConfigAccessor.ReadServerConfig()
	return server.Version, jrpcerror.JrpcError{}
}

func InfoCommand(config interface{}, metadata map[string]string, payload io.ReadCloser) (interface{}, jrpcerror.JrpcError) {
	return config, jrpcerror.JrpcError{}
}

func HealthCommand(config interface{}, metadata map[string]string, payload io.ReadCloser) (interface{}, jrpcerror.JrpcError) {
	return "UP", jrpcerror.JrpcError{}
}
