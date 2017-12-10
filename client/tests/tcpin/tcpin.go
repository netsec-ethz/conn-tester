package tcpin

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/netsec-ethz/conn-tester/client/tests"
	"github.com/netsec-ethz/conn-tester/lib/common"
	"github.com/netsec-ethz/conn-tester/lib/httputils"
	"github.com/netsec-ethz/conn-tester/lib/tcputils"
	"github.com/netsec-ethz/conn-tester/server/requestmessages"
)

type TCPInTest struct {
	Host    string        `json:"host"`
	MyPort  string        `json:"my_port"`
	Timeout time.Duration `json:"timeout"`
}

func (t *TCPInTest) GetTestName() string {
	return "TCP in reachability test"
}

func (t *TCPInTest) GetTestDescription() string {
	return fmt.Sprintf("This test checks if this machine can be accessed over TCP port %s from internet", t.MyPort)
}

func (t *TCPInTest) Run() *tests.TestResult {

	// Create http client
	httpClient, err := httputils.CreateNewClientWrapper(t.Timeout * time.Second)
	if err != nil {
		return &tests.TestResult{Success: false, Message: err.Error()}
	}

	// Start TCP server
	tcpServer, err := tcputils.NewTcpServer(t.MyPort, t.Timeout*time.Second)
	if err != nil {
		return &tests.TestResult{Success: false, Message: err.Error()}
	}
	defer tcpServer.Stop()

	// We define and setup tcp server
	nonce := common.GenerateNonce(64)
	go tcpServer.HandleRequests(func(connection *tcputils.TcpSocket) {
		// We need to read from connection first
		match, _ := common.CompareBytesWithStream([]byte(nonce), connection.Conn)

		if err != nil {
			err = connection.SendData([]byte(err.Error()))
		} else {
			if match {
				connection.SendData([]byte("success"))
			} else {
				connection.SendData([]byte("fail"))
			}
		}
		defer connection.Close()
	})

	success, err := httpClient.SendCommand(t.Host, "POST",
		requestmessages.TCPTestRequest{InPort: t.MyPort, Timeout: t.Timeout, Nonce: nonce})

	if err != nil {
		return &tests.TestResult{Success: false, Message: err.Error()}
	}

	return &tests.TestResult{Success: success}
}

func Create(params *json.RawMessage) []tests.Test {
	var t TCPInTest

	err := json.Unmarshal(*params, &t)
	if err != nil {
		fmt.Println("Error reding JSON params!")
		// TODO: Handle error
	}

	return []tests.Test{&t}
}
