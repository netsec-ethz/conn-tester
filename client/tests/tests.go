package tests

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fatih/color"
)

type TestResult struct {
	Success bool
	Message string
}

type Test interface {
	GetTestName() string
	GetTestDescription() string
	Run() *TestResult
}

type TestConfiguration struct {
	Name   string           `json:"name"`
	Params *json.RawMessage `json:"params"`
}

type testConfigurationFile struct {
	TestList []TestConfiguration `json:"tests"`
}

func LoadTests(testFactory *TestFactory, configFilePath string) ([]Test, error) {
	configFile, err := os.Open(configFilePath)
	defer configFile.Close()
	if err != nil {
		return nil, err
	}

	// Load configuration from file
	var configuration testConfigurationFile
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&configuration)

	// Create tests based on configuration
	var testsToRun []Test
	for _, testConfig := range configuration.TestList {
		newTest := testFactory.CreateTest(testConfig.Name, testConfig.Params)
		if newTest != nil {
			testsToRun = append(testsToRun, newTest)
		} else {
			fmt.Printf("Unknown test <%s> \n", testConfig.Name)
		}
	}

	return testsToRun, nil
}

func RunTest(test Test) {
	response := make(chan *TestResult, 1)

	go func() {
		fmt.Printf("Running: %s \n(%s)\n", test.GetTestName(), test.GetTestDescription())
		response <- test.Run()
	}()

	result := <-response
	if result.Success {
		color.Green("SUCCESS! \n\n")
	} else {
		color.Yellow("FAIL! details: %s\n\n", result.Message)
	}
}
