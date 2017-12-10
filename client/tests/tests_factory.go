package tests

import (
	"encoding/json"
)

type InitFunc func(params *json.RawMessage) []Test

type TestFactory struct {
	generators map[string]InitFunc
}

func (tf *TestFactory) AddTest(testName string, initializer InitFunc) {
	tf.generators[testName] = initializer
}

func (tf *TestFactory) CreateTests(testName string, params *json.RawMessage) []Test {
	if init, exists := tf.generators[testName]; exists {
		return init(params)
	} else {
		return make([]Test, 0, 0)
	}
}

func CreateTestFactory() *TestFactory {
	var tf TestFactory
	tf.generators = make(map[string]InitFunc)
	return &tf
}
