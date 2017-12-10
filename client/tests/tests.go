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
	Name    string           `json:"name"`
	Params  *json.RawMessage `json:"params"`
	Tests   []Test           `json:"tests"`
	Results []*TestResult    `json:"results"`
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

		testConfig.Tests = testFactory.CreateTests(testConfig.Name, testConfig.Params)
		testConfig.Results = make([]*TestResult, 0, cap(testConfig.Tests))

		if len(testConfig.Tests) == 0 {
			fmt.Printf("Unknown name <%s>, test won't be created!\n", testConfig.Name)
		}
	}

	return &configuration, nil
}

func RunTest(tConfig *TestConfiguration) {
	if len(tConfig.Tests) == 0 {
		color.Yellow("Skipping %s as it is not properly initialized\n\n", tConfig.Name)
		return
	}

	done := make(chan bool, 1)

	go func() {
		for _, test := range tConfig.Tests {
			fmt.Printf("Running: %s \n(%s)\n", test.GetTestName(), test.GetTestDescription())
			result := test.Run()

			if result.Success {
				color.Green("SUCCESS! \n\n")
			} else {
				color.Yellow("FAIL! details: %s\n\n", result.Message)
			}

			tConfig.Results = append(tConfig.Results, result)
		}

		done <- true
	}()

	<-done

	fmt.Printf("Finished running [%s] \n", tConfig.Name)
}
