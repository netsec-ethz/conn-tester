package requestmessages

import (
	"time"
)

type TCPTestRequest struct {
	InPort  string        `json:"port"`
	Timeout time.Duration `json:"timeout"`
	Nonce   string        `json:"nonce"`
}

type UDPTestRequest struct {
	InPort  string        `json:"port"`
	Timeout time.Duration `json:"timeout"`
	Nonce   string        `json:"nonce"`
}
