package udputils

import (
	"fmt"
	"net"
	"time"

	"github.com/netsec-ethz/conn-tester/lib/common"
)

type UdpClient struct {
	conn          *net.UDPConn
	Timeout       time.Duration
	MaxRetransmit int
}

func (u *UdpClient) Close() {
	u.conn.Close()
}

func (u *UdpClient) SendRequestUntilResponse(request []byte, expectedResponse []byte) (bool, error) {
	for i := 0; i < u.MaxRetransmit; i++ {
		u.conn.SetDeadline(time.Now().Add(u.Timeout))
		_, err := u.conn.Write(request)
		if err != nil {
			fmt.Println(err)
			continue //Error sending, decrement retry count, no point in receiving
		}

		if responseReceived, _ := common.CompareBytesWithStream(expectedResponse, u.conn); responseReceived == true {
			return true, nil
		}
	}

	return false, fmt.Errorf("Timed out! Haven't received desired reponse in expected time")
}

func (u *UdpClient) SendRequestUntilStopped(request []byte, stop <-chan bool) (bool, error) {
	for i := 0; i < u.MaxRetransmit; i++ {
		u.conn.SetDeadline(time.Now().Add(u.Timeout))
		u.conn.Write(request)

		select {
		case <-time.After(u.Timeout):
		case <-stop:
			return true, nil
		}
	}

	return false, fmt.Errorf("Timed out!")
}

func NewUdpClient(host string, port string, connectTimeout time.Duration) (*UdpClient, error) {
	conn, err := net.DialTimeout("udp", net.JoinHostPort(host, port), connectTimeout*time.Second)
	udpConn := conn.(*net.UDPConn)
	if err != nil {
		return nil, err
	}

	udpClient := UdpClient{conn: udpConn, Timeout: 1 * time.Second, MaxRetransmit: 15}

	return &udpClient, nil
}
