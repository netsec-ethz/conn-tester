package tcputils

import (
	"net"
	"time"

	"github.com/netsec-ethz/conn-tester/lib/common"
)

type TcpSocket struct {
	Conn    net.Conn
	Timeout time.Duration
}

func (t *TcpSocket) Close() {
	t.Conn.Close()
}

func (t *TcpSocket) SendData(request []byte) error {
	t.Conn.SetDeadline(time.Now().Add(t.Timeout))
	_, err := t.Conn.Write(request)
	if err != nil {
		return err
	}

	return nil
}

func (t *TcpSocket) SendRequestCompareResponse(request []byte, expectedResponse []byte) (bool, error) {
	t.SendData(request)
	t.Conn.SetDeadline(time.Now().Add(t.Timeout))
	if equal, err := common.CompareBytesWithStream(expectedResponse, t.Conn); err != nil {
		return false, err
	} else {
		return equal, nil
	}
}

func NewTcpClient(host string, port string, timeout time.Duration) (*TcpSocket, error) {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if err != nil {
		return nil, err
	}

	sock := TcpSocket{Conn: conn, Timeout: timeout}

	return &sock, nil
}
