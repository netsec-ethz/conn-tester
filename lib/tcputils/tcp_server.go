package tcputils

import (
	"fmt"
	"net"
	"time"
)

type TcpConnectionHandler func(connection *TcpSocket)

type TcpServer struct {
	server   net.Listener
	Timeout  time.Duration
	stopChan chan bool
}

func (t *TcpServer) Stop() {
	t.stopChan <- true
}

func (t *TcpServer) HandleRequests(handler TcpConnectionHandler) error {
	connChannel := acceptClients(t.server)
	var err error
	for stop := false; stop != true; {
		select {
		case <-t.stopChan:
			stop = true
		case <-time.After(t.Timeout):
			err = fmt.Errorf("Timed out! No requests received")
			stop = true
		case conn := <-connChannel:
			socket := &TcpSocket{Conn: conn, Timeout: t.Timeout}
			go handler(socket)
		}
	}
	fmt.Println("Stopping request handler")
	defer t.server.Close()

	return err
}

func acceptClients(listener net.Listener) chan net.Conn {
	ch := make(chan net.Conn)
	go func() {
		for {
			client, err := listener.Accept()
			if err != nil {
				break
			}
			ch <- client
		}
	}()
	return ch
}

func NewTcpServer(port string, timeout time.Duration) (*TcpServer, error) {
	server, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return nil, err
	}

	tcpServer := TcpServer{server: server, Timeout: timeout}
	tcpServer.stopChan = make(chan bool, 1)

	return &tcpServer, nil
}
