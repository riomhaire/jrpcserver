package model

import (
	"io"

	"github.com/riomhaire/jrpcserver/model/jrpcerror"
)

type JRPCCommand struct {
	Name        string
	Command     func(APIConfig, map[string]string, io.ReadCloser) (interface{}, jrpcerror.JrpcError)
	RawResponse bool // If true then response is just marshalled otherwise rpc format
}
