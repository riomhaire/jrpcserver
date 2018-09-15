package usecases

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/riomhaire/jrpcserver/model/jrpcerror"
)

func PingCommand(metadata map[string]string, payload io.ReadCloser) (interface{}, jrpcerror.JrpcError) {
	return "pong", jrpcerror.JrpcError{}
}

func PongCommand(metadata map[string]string, payload io.ReadCloser) (interface{}, jrpcerror.JrpcError) {
	return "ping", jrpcerror.JrpcError{}
}

func EchoCommand(metadata map[string]string, payload io.ReadCloser) (interface{}, jrpcerror.JrpcError) {
	for name, value := range metadata {
		fmt.Printf("%s = %s\n", name, value)
	}

	data, _ := ioutil.ReadAll(payload)
	msg := make(map[string]interface{})

	json.Unmarshal(data, &msg)

	return msg, jrpcerror.JrpcError{}
}

func ListCommandsCommand(metadata map[string]string, payload io.ReadCloser) (interface{}, jrpcerror.JrpcError) {
	names := make([]string, 0)
	for _, cmd := range Commands {
		names = append(names, cmd.Name)
	}
	return names, jrpcerror.JrpcError{}
}
