package tcpout

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/netsec-ethz/conn-tester/client/tests"
	"github.com/netsec-ethz/conn-tester/lib/tcputils"
)

type TCPOutTest struct {
	Host             string        `json:"host"`
	Port             string        `json:"port"`
	Request          string        `json:"request"`
	Timeout          time.Duration `json:"timeout"`
	ExpectedResponse string        `json:"expected_response"`
	CompareResponse  bool          `json:"compare_response"`
}

func (t *TCPOutTest) GetTestName() string {
	return "TCP destination reachability test"
}

func (t *TCPOutTest) GetTestDescription() string {
	return fmt.Sprintf("This test checks if TCP server %s on port %s can be accessed from network", t.Host, t.Port)
}

func (t *TCPOutTest) Run() *tests.TestResult {
	socket, err := tcputils.NewTcpClient(t.Host, t.Port, t.Timeout*time.Second)
	if err != nil {
		return &tests.TestResult{Success: false, Message: err.Error()}
	}
	defer socket.Close()

	if t.CompareResponse {
		match, err := socket.SendRequestCompareResponse([]byte(t.Request), []byte(t.ExpectedResponse))
		if match {
			return &tests.TestResult{Success: true}
		} else {
			if err == nil {
				return &tests.TestResult{Success: false, Message: "Mismatched response"}
			} else {
				return &tests.TestResult{Success: false, Message: err.Error()}
			}
		}
	} else {
		err := socket.SendData([]byte(t.Request))
		if err == nil {
			return &tests.TestResult{Success: true}
		} else {
			return &tests.TestResult{Success: false, Message: err.Error()}
		}
	}
}

func Create(params *json.RawMessage) []tests.Test {
	var t TCPOutTest

	err := json.Unmarshal(*params, &t)
	if err != nil {
		fmt.Println("Error reding JSON params!")
		// TODO: Handle error
	}

	return []tests.Test{&t}
}
