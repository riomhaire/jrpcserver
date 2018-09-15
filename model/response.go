package model

type RPCCommandResponse struct {
	Code  int
	Error string
	Value interface{}
}
