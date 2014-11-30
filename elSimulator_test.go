package main

import "testing"

func init() {
	//prepare
	elSimulatorConfig.baseDirectory = "baseDirectory"
}

type testbase struct {
	method   string
	context  string
	expected string
}

var paramTestBase = []testbase{
	{"method", context, "baseDirectory/file/method/"},
	{"method", "/file/notegal", "baseDirectory/file/notegal/method/"},
}

func TestNameFileParameter_Base(t *testing.T) {
	for i, testbase := range paramTestBase {
		var buffer NameFileParameter
		baseActual := buffer.Base(testbase.method, testbase.context)
		if baseActual != testbase.expected {
			t.Error("i", i, "actual", baseActual, "expected", testbase.expected)
		}

	}

}
