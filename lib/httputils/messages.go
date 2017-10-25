package httputils

type ServerResponseMessage struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
