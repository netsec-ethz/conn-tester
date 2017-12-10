package udpin

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
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

const port_range_regex = "\\[(\\d+)\\-(\\d+)\\]"

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

func Create(params *json.RawMessage) []tests.Test {
	var t UDPInTest

	err := json.Unmarshal(*params, &t)
	if err != nil {
		fmt.Println("Error reding JSON params!")
		// TODO: Handle error
	}

	// We need to check if port range is specified and create more instances
	re := regexp.MustCompile(port_range_regex)
	port_range := re.FindAllStringSubmatch(t.MyPort, 1)

	var startPort = 0
	var endPort = 0

	if len(port_range) == 1 {
		// Port range has been specified
		fmt.Printf("Received port range from %s to %s \n", port_range[0][1], port_range[0][2])
		startPort, _ = strconv.Atoi(port_range[0][1])
		endPort, _ = strconv.Atoi(port_range[0][2])

	} else {
		startPort, _ = strconv.Atoi(t.MyPort)
		endPort, _ = strconv.Atoi(t.MyPort)
	}

	var totalPorts = endPort - startPort + 1

	var generatedTests = make([]tests.Test, 0, totalPorts)
	for i := startPort; i <= endPort; i++ {
		newTest := t
		newTest.MyPort = strconv.Itoa(i)
		generatedTests = append(generatedTests, &newTest)
	}

	return generatedTests
}
