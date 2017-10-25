package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/netsec-ethz/conn-tester/client/tests"
	"github.com/netsec-ethz/conn-tester/client/tests/httptest"
	"github.com/netsec-ethz/conn-tester/client/tests/ntptest"
	"github.com/netsec-ethz/conn-tester/client/tests/tcpin"
	"github.com/netsec-ethz/conn-tester/client/tests/tcpout"
	"github.com/netsec-ethz/conn-tester/client/tests/udpin"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	fmt.Println("--- Starting client application ---")

	if len(os.Args) < 2 {
		fmt.Println("Error! You must specify config file name!")
		os.Exit(-1)
	}

	factory := tests.CreateTestFactory()
	factory.AddTest("ntp_test", ntptest.Create)
	factory.AddTest("http_test", httptest.Create)
	factory.AddTest("tcp_out", tcpout.Create)
	factory.AddTest("tcp_in", tcpin.Create)
	factory.AddTest("udp_in", udpin.Create)

	testList, err := tests.LoadTests(factory, os.Args[1])
	if err != nil {
		fmt.Println(err)
	}

	for _, test := range testList {
		tests.RunTest(test)
	}
}
