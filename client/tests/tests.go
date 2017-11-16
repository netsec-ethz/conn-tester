package tests

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fatih/color"
)

type TestResult struct {
	Success bool   `json:"success"`
	Message string `json:"description"`
}

type Test interface {
	GetTestName() string
	GetTestDescription() string
	Run() *TestResult
}

type TestConfiguration struct {
	Name   string           `json:"name"`
	Params *json.RawMessage `json:"params"`
	test   Test
	Result *TestResult `json:"result"`
}

type TestConfigurationList struct {
	TestList []*TestConfiguration `json:"tests"`
}

func LoadTests(testFactory *TestFactory, configFilePath string) (*TestConfigurationList, error) {
	configFile, err := os.Open(configFilePath)
	defer configFile.Close()
	if err != nil {
		return nil, err
	}

	// Load configuration from file
	var configuration TestConfigurationList
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&configuration)

	// Create tests based on configuration

	for _, testConfig := range configuration.TestList {

		testConfig.test = testFactory.CreateTest(testConfig.Name, testConfig.Params)

		if testConfig.test == nil {
			fmt.Printf("Unknown test <%s> \n", testConfig.Name)
		}
	}

	return &configuration, nil
}

func RunTest(tConfig *TestConfiguration) {
	if tConfig.test == nil {
		tConfig.Result = &TestResult{Success: false, Message: "Not a valid test"}
		color.Yellow("Skipping %s as it is not properly initialized\n\n", tConfig.Name)
		return
	}

	response := make(chan *TestResult, 1)

	go func() {
		fmt.Printf("Running: %s \n(%s)\n", tConfig.test.GetTestName(), tConfig.test.GetTestDescription())
		tConfig.Result = tConfig.test.Run()
		response <- tConfig.Result
	}()

	result := <-response
	if result.Success {
		color.Green("SUCCESS! \n\n")
	} else {
		color.Yellow("FAIL! details: %s\n\n", result.Message)
	}
}
