package model

import (
	"io"

	"github.com/riomhaire/jrpcserver/model/jrpcerror"
)

type JRPCCommand struct {
	Name    string
	Command func(map[string]string, io.ReadCloser) (interface{}, jrpcerror.JrpcError)
}