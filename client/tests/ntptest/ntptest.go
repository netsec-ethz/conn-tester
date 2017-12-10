package ntptest

import (
	"encoding/json"
	"fmt"

	"github.com/beevik/ntp"

	"github.com/netsec-ethz/conn-tester/client/tests"
)

type NtpTest struct {
	NtpServer string `json:"ntp_server"`
}

func (ntpt *NtpTest) GetTestName() string {
	return "NTP server test"
}

func (ntpt *NtpTest) GetTestDescription() string {
	return fmt.Sprintf("This test checks if NTP server %s is accessible from network", ntpt.NtpServer)
}

func (ntpt *NtpTest) Run() *tests.TestResult {
	response, err := ntp.Query(ntpt.NtpServer)
	var testRes tests.TestResult
	if err != nil {
		testRes.Success = false
		testRes.Message = err.Error()
	} else {
		testRes.Success = true
		fmt.Printf("Time from NTP server received: %s\n", response.Time)
	}

	return &testRes
}

func Create(params *json.RawMessage) []tests.Test {
	var ntp NtpTest

	err := json.Unmarshal(*params, &ntp)
	if err != nil {
		fmt.Println("Error!")
		// TODO: Handle error
	}

	return []tests.Test{&ntp}
}
