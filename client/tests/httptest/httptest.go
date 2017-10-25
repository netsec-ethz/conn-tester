package httptest

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/netsec-ethz/conn-tester/client/tests"
	"github.com/netsec-ethz/conn-tester/lib/httputils"
)

type HttpTest struct {
	Host             string        `json:"host"`
	Method           string        `json:"method"`
	Request          string        `json:"request"`
	ExpectedResponse string        `json:"expected_response"`
	Timeout          time.Duration `json:"timeout"`
}

func (h *HttpTest) GetTestName() string {
	return "HTTP out test"
}

func (h *HttpTest) GetTestDescription() string {
	return fmt.Sprintf("This test checks if %s can be accessed over HTTP", h.Host)
}

func (h *HttpTest) Run() *tests.TestResult {
	var testRes tests.TestResult

	cw, err := httputils.CreateNewClientWrapper(h.Timeout * time.Second)
	if err != nil {
		return &tests.TestResult{Success: false, Message: err.Error()}
	}
	success, err := cw.SendRequestCompareResponse(h.Host, h.Method, []byte(h.Request), []byte(h.ExpectedResponse))

	testRes.Success = success
	if success == true {
		return &tests.TestResult{Success: true}
	} else {
		if err != nil {
			return &tests.TestResult{Success: false, Message: err.Error()}
		} else {
			return &tests.TestResult{Success: false,
				Message: "Response from webpage doesn't match expected output",
			}
		}
	}
}

func Create(params *json.RawMessage) tests.Test {
	var testInterface tests.Test
	var h HttpTest

	err := json.Unmarshal(*params, &h)
	if err != nil {
		fmt.Println("Error reding JSON params!")
		// TODO: Handle error
	}

	testInterface = &h

	return testInterface
}
