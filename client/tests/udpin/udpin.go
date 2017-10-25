package udpin

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/netsec-ethz/conn-tester/client/tests"
	"github.com/netsec-ethz/conn-tester/lib/common"
	"github.com/netsec-ethz/conn-tester/lib/httputils"
	"github.com/netsec-ethz/conn-tester/lib/udputils"
	"github.com/netsec-ethz/conn-tester/server/requestmessages"
)

type UDPInTest struct {
	Host    string        `json:"host"`
	MyPort  string        `json:"my_port"`
	Timeout time.Duration `json:"timeout"`
}

func (t *UDPInTest) GetTestName() string {
	return "UDP in reachability test"
}

func (t *UDPInTest) GetTestDescription() string {
	return fmt.Sprintf("This test checks if this machine can be accessed over UDP port %s from internet", t.MyPort)
}

func (t *UDPInTest) Run() *tests.TestResult {
	// Create http client
	httpClient, err := httputils.CreateNewClientWrapper(t.Timeout * time.Second)
	if err != nil {
		return &tests.TestResult{Success: false, Message: err.Error()}
	}

	udpServer, err := udputils.NewUdpServer(t.MyPort)
	if err != nil {
		return &tests.TestResult{Success: false, Message: err.Error()}
	}
	defer udpServer.Close()

	nonce := common.GenerateNonce(64)
	go udpServer.HandleRequests(func(request string) string {
		if request == nonce {
			return "success"
		} else {
			return "fail"
		}
	})
	defer udpServer.Close()

	success, err := httpClient.SendCommand(t.Host, "POST",
		requestmessages.TCPTestRequest{InPort: t.MyPort, Timeout: t.Timeout, Nonce: nonce})

	if err != nil {
		return &tests.TestResult{Success: false, Message: err.Error()}
	}

	return &tests.TestResult{Success: success}
}

func Create(params *json.RawMessage) tests.Test {
	var testInterface tests.Test
	var t UDPInTest

	err := json.Unmarshal(*params, &t)
	if err != nil {
		fmt.Println("Error reding JSON params!")
		// TODO: Handle error
	}

	testInterface = &t

	return testInterface
}
