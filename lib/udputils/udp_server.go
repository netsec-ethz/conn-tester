package udputils

import (
	"fmt"
	"net"
)

type UdpServer struct {
	server        *net.UDPConn
	rcvBufferSize uint32
}

type UdpPacketHandler func(request string) string

func (u *UdpServer) Close() {
	u.server.Close()
}

func (u *UdpServer) HandleRequests(handler UdpPacketHandler) {
	go func() {
		buffer := make([]byte, u.rcvBufferSize)
		for {
			n, addr, err := u.server.ReadFromUDP(buffer)
			if err != nil {
				return // When server is stopped
			}
			resp := handler(string(buffer[:n]))
			u.server.WriteTo([]byte(resp), addr)
		}
	}()
}

func NewUdpServer(port string) (*UdpServer, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%s", port))
	if err != nil {
		return nil, err
	}

	server, err := net.ListenUDP("udp", udpAddr)

	udpServer := UdpServer{server: server, rcvBufferSize: 1024}

	return &udpServer, nil
}
