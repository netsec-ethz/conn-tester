package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/netsec-ethz/conn-tester/lib/httputils"
	"github.com/netsec-ethz/conn-tester/lib/tcputils"
	"github.com/netsec-ethz/conn-tester/lib/udputils"
	"github.com/netsec-ethz/conn-tester/server/requestmessages"
)

func index(hostAddress string, params *json.RawMessage) (bool, error) {
	return true, nil
}

func udpTest(hostAddress string, params *json.RawMessage) (bool, error) {
	var command requestmessages.UDPTestRequest

	err := json.Unmarshal(*params, &command)
	if err != nil {
		return false, err
	}

	udpClient, err := udputils.NewUdpClient(hostAddress, command.InPort, time.Duration(command.Timeout)*time.Second)
	if err != nil {
		return false, err
	}
	defer udpClient.Close()

	equal, err := udpClient.SendRequestUntilResponse([]byte(command.Nonce), []byte("success"))
	if err != nil {
		return false, err
	}
	if !equal {
		return false, fmt.Errorf("Received reply did not match")
	}

	return true, nil
}

func tcpTest(hostAddress string, params *json.RawMessage) (bool, error) {
	var command requestmessages.TCPTestRequest

	err := json.Unmarshal(*params, &command)
	if err != nil {
		return false, err
	}

	tcpCilent, err := tcputils.NewTcpClient(hostAddress, command.InPort, time.Duration(command.Timeout)*time.Second)
	if err != nil {
		return false, err
	}
	defer tcpCilent.Close()

	equal, err := tcpCilent.SendRequestCompareResponse([]byte(command.Nonce), []byte("success"))
	if err != nil {
		return false, err
	}

	if !equal {
		return false, fmt.Errorf("Replies do not match")
	}

	return true, nil
}

func registerHandlers(rh *httputils.RequestMultiplexer) {
	rh.RegisterHandler("/", "GET", index)
	rh.RegisterHandler("/tcp-test", "POST", tcpTest)
	rh.RegisterHandler("/udp-test", "POST", udpTest)
}
