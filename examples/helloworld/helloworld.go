package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/riomhaire/jrpcserver/infrastructure/api/rpc"
	"github.com/riomhaire/jrpcserver/model"
	"github.com/riomhaire/jrpcserver/model/jrpcerror"
)

// A simple example 'helloworld' program to show how use the framework
//
// to execute (after compiling):
//      helloworld -port=9999 -consul=consul:8500
//
func main() {

	name := flag.String("name", "helloworld", "name of service")
	path := flag.String("path", "/api/v1/helloworld", "Path to which to  <path>/<command>  points to action")
	consulHost := flag.String("consul", "", "consul host usually something like 'localhost:8500'. Leave blank if not required")
	port := flag.Int("port", 9999, "port to use")
	flag.Parse()

	config := rpc.APIConfig{
		ServiceName: *name,
		BaseURI:     *path,
		Port:        *port,
		Commands:    Commands(),
		Consul:      *consulHost,
	}
	rpc.StartAPI(config) // Start service -  wont return
}

func Commands() []model.JRPCCommand {
	commands := make([]model.JRPCCommand, 0)

	commands = append(commands, model.JRPCCommand{"example.helloworld", HelloWorldCommand})
	return commands
}

func HelloWorldCommand(metadata map[string]string, payload io.ReadCloser) (interface{}, jrpcerror.JrpcError) {
	data, err := ioutil.ReadAll(payload)
	if err != nil {
		return "", jrpcerror.JrpcError{500, err.Error()}
	} else {
		fmt.Println(string(data))
		response := fmt.Sprintf("Hello %v", string(data))

		return response, jrpcerror.JrpcError{}
	}
}
