package model

type RPCCommandResponse struct {
	Code        int         `json:"code"`
	Error       string      `json:"error"`
	Value       interface{} `json:"value"`
	RawResponse bool        `json:"-"`
}
