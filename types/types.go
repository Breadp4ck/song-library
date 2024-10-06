package types

type Response struct {
	Message any `json:"message,omitempty"`
	Error   any `json:"error,omitempty"`
}

var ResponseOK = Response{Message: "ok"}
