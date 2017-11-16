package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"

	"github.com/netsec-ethz/conn-tester/client/tests"
	"github.com/netsec-ethz/conn-tester/client/tests/httptest"
	"github.com/netsec-ethz/conn-tester/client/tests/ntptest"
	"github.com/netsec-ethz/conn-tester/client/tests/tcpin"
	"github.com/netsec-ethz/conn-tester/client/tests/tcpout"
	"github.com/netsec-ethz/conn-tester/client/tests/udpin"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	configFilePath = kingpin.Flag("config", "Name of configuration file").Required().String()
	output         = kingpin.Flag("output_result", "Should output test result").Bool()
	outputFilePath = kingpin.Flag("output_path", "Location to output path").Default("result.json").String()
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	fmt.Println("--- Starting client application ---")

	kingpin.Parse()

	factory := tests.CreateTestFactory()
	factory.AddTest("ntp_test", ntptest.Create)
	factory.AddTest("http_test", httptest.Create)
	factory.AddTest("tcp_out", tcpout.Create)
	factory.AddTest("tcp_in", tcpin.Create)
	factory.AddTest("udp_in", udpin.Create)

	testList, err := tests.LoadTests(factory, *configFilePath)
	if err != nil {
		fmt.Println(err)
	}

	for _, config := range testList.TestList {
		tests.RunTest(config)
	}

	if *output {
		fmt.Printf("Saving test result in %s \n", *outputFilePath)

		resultJson, _ := json.Marshal(testList)
		err = ioutil.WriteFile(*outputFilePath, resultJson, 0644)
	}

}
